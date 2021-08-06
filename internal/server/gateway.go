package server

// TODO Implement queue
import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.com/eco_system/gateway/api/go/gen/accounts/v1"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"gitlab.com/eco_system/gateway/internal/common/log"
	"gitlab.com/eco_system/gateway/internal/common/validator"
	"gitlab.com/eco_system/gateway/internal/k8s"
	"gitlab.com/eco_system/gateway/internal/k8s/enum"
	"gitlab.com/eco_system/gateway/internal/observer"
	"gitlab.com/eco_system/gateway/internal/services"
)

// GatewayServer encapsulates GatewayServer and Observer
type GatewayServer struct {
	sync.Mutex
	Name        string
	production  bool
	Server      *http.Server
	services    map[string]*k8s.APIService
	logger      *logrus.Logger
	k8sClient   *k8s.Client
	accountSVC  accounts.AccountsServiceClient
	redisClient *redis.Client
}

var defaultGatewayOptions = options{
	logger: log.NewLogger(),
	addr:   ":5000",
}

var OmitHeaders = map[string]bool{
	"content-length": true,
}

var OmitWebSocketHeaders = map[string]bool{
	"connection": true,
	"upgrade":    true,
}

func essentialMiddleware(gs *GatewayServer) func(*gin.Engine) {
	return func(r *gin.Engine) {
		r.Use(gin.Recovery())
		r.Use(requestLogger(gs.logger))
		r.Use(WebhookRequests)
		// Health route
		r.GET("/hz", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "success",
			})
		})
		r.NoRoute(notFoundMW)
	}
}

// forwardHeaders packages http headers into headers-bin and forwards the metadata
func forwardHeaders(_ context.Context, r *http.Request) metadata.MD {
	md := metadata.Pairs()
	for k, v := range r.Header {
		header := strings.ToLower(k)
		if OmitHeaders[header] {
			continue
		}
		if len(v) > 0 {
			md.Set(header, v[0])
		}
	}
	return md
}

// newGrpcMux creates a new mux that handles grpc calls.
func (gs *GatewayServer) newGrpcMux(ctx context.Context) *runtime.ServeMux {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
		runtime.WithMetadata(forwardHeaders),
		runtime.WithErrorHandler(ProtoErrorWithLogger(gs.logger, gs.production)),
	)
	return mux
}

// newRouter returns a new gin Engine which handles routing for gateway
func newRouter(attachMiddleware ...func(r *gin.Engine)) (*gin.Engine, error) {
	r := gin.New()

	if attachMiddleware == nil {
		return r, nil
	}
	for _, m := range attachMiddleware {
		m(r)
	}
	return r, nil
}

func transformSubProtocolHeader(header string) string {
	tokens := strings.SplitN(header, "Bearer,", 2)

	if len(tokens) < 2 {
		return ""
	}

	return fmt.Sprintf("Bearer %v", strings.Trim(tokens[1], " "))
}

func (gs *GatewayServer) authorization(svc *k8s.APIService, cb gin.HandlerFunc) gin.HandlerFunc {
	if svc == nil {
		gs.logger.Errorf("authorization: svc is nil")
		return func(c *gin.Context) {
			c.String(500, "Internal Server Error")
			c.Abort()
		}
	}
	required := true
	r := svc.Authentication.Required

	switch svc.Authentication.Required.(type) {
	case bool:
		b, ok := r.(bool)
		if ok {
			r = b
		}
	case string:
		if s, ok := r.(string); ok {
			b, err := strconv.ParseBool(s)
			if err != nil {
				gs.logger.Warnf("[%v] Error parsing %v - authentication.required. Falling back to true", svc.ServiceName, required)
			}
			required = b
		}
	}

	return func(c *gin.Context) {
		if !required {
			cb(c)
			return
		}
		em := gin.H{
			"error": "Request is malformed",
		}
		url := c.Request.URL
		// TODO: Revise logic exclude
		for _, e := range svc.Authentication.Exclude {
			if e == url.String() {
				cb(c)
				return
			}
		}

		token := ""
		// Websocket proxy will verify the header
		// there's no way of introspecting it or having a middleware
		// Code extracted from github.com/tmc/grpc-websocket-proxy
		if websocket.IsWebSocketUpgrade(c.Request) {
			if swsp := c.Request.Header.Get("Sec-WebSocket-Protocol"); swsp != "" {
				token = transformSubProtocolHeader(swsp)
			}
		} else {
			token = c.Request.Header.Get("Authorization")
		}
		if token == "" {
			c.JSON(401, em)
			c.Abort()
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// TODO: Include scope
		ip := c.ClientIP()
		md := metadata.Pairs("x-forwarded-for", ip)
		ctx = metadata.NewOutgoingContext(ctx, md)
		res, err := gs.accountSVC.Introspect(ctx, &accounts.IntrospectRequest{Token: token})
		if err != nil {
			gs.logger.Errorf("authentication: %v", err)
			c.JSON(401, em)
			c.Abort()
			return
		}

		// add subject to header
		if res.Active {
			c.Request.Header.Add("x-subject", res.Sub)
			cb(c)
			return
		}
		c.JSON(401, gin.H{
			"active": false,
		})
	}
}

func (gs *GatewayServer) applyHTTP(r *gin.Engine, svc *k8s.APIService, target string, route string) error {
	if r == nil {
		return InvalidParams("applyHttp: gin r router is nil")
	}
	if target == "" {
		return InvalidParams("applyHttp: target is empty")
	}
	// TODO: Test connectivity
	rp := gs.authorization(svc, ReverseProxy(target))
	// Routes all requests to service
	r.GET(route, rp)
	r.POST(route, rp)
	r.Any(fmt.Sprintf("%v/*any", route), rp)
	return nil
}

func (gs *GatewayServer) applyGrpc(r *gin.Engine, svc *k8s.APIService, serviceName string, conn *grpc.ClientConn, target string, route string) error {
	var err error
	var svcDesc *ServiceDesc
	if r == nil {
		return InvalidParams("r router is nil")
	}
	if serviceName == "" {
		return InvalidParams("serviceName is empty")
	}
	// Check if service exists in proto definition
	for k, s := range GetProtos() {
		// Splits service name i.e api.auth.authservice ["api", "auth", "authservice"]
		split := strings.Split(k, ".")
		// Formats it to dash i.e api-auth-authservice
		dash := strings.ToLower(strings.Join(split, "-"))
		// Formats it to underscore i.e api_auth_authservice
		underscore := strings.ToLower(strings.Join(split, "_"))
		if serviceName == dash || serviceName == underscore {
			svcDesc = &s
			break
		}
	}
	if svcDesc == nil {
		return NotFound("No gw proto generated or found for " + serviceName)
	}
	// Create connection
	if conn == nil {
		// TLS will be handled by sidecar proxy. I.E. Istio's Sidecar
		conn, err = grpc.DialContext(context.Background(), target, grpc.WithInsecure())
		if err != nil {
			return nil
		}
	}
	// creates a new mux
	mux := gs.newGrpcMux(context.Background())
	// pass mux to grpc handlers
	err = svcDesc.Handler(context.Background(), mux, conn)
	if err != nil {
		return fmt.Errorf("mux Handler: %v", err)
	}
	handler := gs.authorization(svc, func(c *gin.Context) {
		wsproxy.WebsocketProxy(mux,
			wsproxy.WithForwardedHeaders(func(header string) bool {
				if OmitWebSocketHeaders[strings.ToLower(header)] {
					return false
				}
				return true
			}),
			wsproxy.WithLogger(gs.logger),
		).ServeHTTP(c.Writer, c.Request)
	})

	r.GET(route, handler)
	r.POST(route, handler)
	r.Any(fmt.Sprintf("%v/*any", route), handler)
	return nil
}

// applyRoutes creates a new replaces the server's http handler with
// newly populated routes
func (gs *GatewayServer) applyRoutes() {
	// Create new Router
	r, err := newRouter(
		essentialMiddleware(gs),
	)
	if err != nil {
		gs.logger.Errorf("applyRoutes: %v", err)
	}

	// Iterate services mapping them to gin router
	for _, svc := range gs.services {
		validate := validator.Instance()
		// Ensures required fields are populated
		err := validate.Struct(svc)
		if err != nil {
			gs.logger.Error(err)
			continue
		}
		// Iterate exposed ports of services
		// TODO: Test if routes is are alive
		for _, port := range svc.Ports {
			target := fmt.Sprintf("%v.svc.cluster.local:%v", svc.DNSPath, port.Port)
			route := fmt.Sprintf("/%v%v", svc.APIversion, svc.Path)

			lf := logrus.Fields{
				"protocol":       "grpc",
				"route":          route,
				"service":        svc.ServiceName,
				"authentication": svc.Authentication.Required,
				"apiVersion":     svc.APIversion,
			}

			switch port.Name {
			case "http":
				err := gs.applyHTTP(r, svc, target, route)
				if err != nil {
					gs.logger.WithFields(lf).Error(err)
				}
				if gs.production {
					gs.logger.WithFields(lf).Info("Route -->")
				}
			case "grpc":
				err := gs.applyGrpc(r, svc, svc.ServiceName, svc.GRPCClientConn, target, route)
				if err != nil {
					gs.logger.WithFields(lf).Error(err)
				}
				if gs.production {
					gs.logger.WithFields(lf).Info("Route -->")
				}
			}
		}
	}
	gs.Lock()
	defer gs.Unlock()
	// Apply routes to server handler
	gs.Server.Handler = r
}

// updateServices updates gateway services
func (gs *GatewayServer) updateServices(service *k8s.APIService) {
	if service == nil {
		gs.logger.Error("updateServices: service is nil")
		return
	}
	if service.DNSPath == "" {
		gs.logger.Error("updateServices: service.DNSPath is empty")
		return
	}
	// Create map if services is not assigned
	if gs.services == nil {
		gs.services = map[string]*k8s.APIService{}
	}
	gs.Lock()
	defer gs.Unlock()
	gs.services[service.DNSPath] = service
}

func (gs *GatewayServer) deleteServices(service *k8s.APIService) {
	if service == nil {
		gs.logger.Error("deleteServices: service is nil")
		return
	}
	if service.DNSPath == "" {
		gs.logger.Error("deleteServices: service.DNSPath is empty")
		return
	}
	// Create map if services is not assigned
	if gs.services == nil {
		return
	}
	gs.Lock()
	defer gs.Unlock()
	delete(gs.services, service.DNSPath)
}

// fetchAllServices fetches all services from cluster
func (gs *GatewayServer) fetchAllServices() error {
	// Get K8S Services in cluster
	svc, err := gs.k8sClient.GetServices("default")
	if err != nil {
		return fmt.Errorf("fetchAllServices: GetServices: %v", err)
	}
	for _, d := range svc {
		o, err := gs.k8sClient.CoreAPI().Admission().UnmarshalK8SObject(d)
		if err != nil {
			gs.logger.Errorf("fetchAllServices: %v", err)
			continue
		}
		// Filter admission request if incoming request does not have api-service labeled.
		if strings.ToLower(o.Metadata.Labels.ResourceType) != string(enum.LabelAPIService) {
			continue
		}
		s, err := gs.k8sClient.CoreAPI().APIServices().ObjectToAPI(o)
		if err != nil {
			gs.logger.Errorf("fetchAllServices: %v", err)
			continue
		}
		gs.updateServices(s)
	}
	return nil
}

// apply updates apiservices to gatway services
func (gs *GatewayServer) apply(ar *k8s.AdmissionRegistration) {
	s, err := gs.k8sClient.CoreAPI().APIServices().ObjectToAPI(&ar.Request.Object)
	if err != nil {
		gs.logger.Error(err)
	}
	gs.updateServices(s)
	gs.applyRoutes()
}

func (gs *GatewayServer) delete(ar *k8s.AdmissionRegistration) {
	s, err := gs.k8sClient.CoreAPI().APIServices().ObjectToAPI(&ar.Request.OldObject)
	if err != nil {
		gs.logger.Error(err)
	}
	gs.deleteServices(s)
	gs.applyRoutes()
}

// directAdmission streamlines a series operations
// such as parsing AdmissionRequest, filtering
// and routing to its necessary
func (gs *GatewayServer) directAdmission(d []byte) {
	ar, err := gs.k8sClient.CoreAPI().Admission().Unmarshal(d)
	if err != nil {
		gs.logger.Error(err)
		return
	}
	// Filter admission request if incoming request is not of K8S Service Object.
	if strings.ToLower(ar.Request.Kind.Kind) != string(enum.K8SServiceObject) {
		gs.logger.Printf("Request request %v is not service")
		return
	}
	// Filter admission request if incoming request does not have api-service labeled.
	if strings.ToLower(ar.Request.OldObject.Metadata.Labels.ResourceType) != string(enum.LabelAPIService) &&
		strings.ToLower(ar.Request.Object.Metadata.Labels.ResourceType) != string(enum.LabelAPIService) {
		return
	}
	switch ar.Request.Operation {
	case string(enum.Create):
		fallthrough
	case string(enum.Update):
		gs.apply(ar)
	case string(enum.Delete):
		gs.delete(ar)
	}
}

func (gs *GatewayServer) subscribe(ctx context.Context) {
	if gs.redisClient == nil {
		return
	}
	sub := gs.redisClient.Subscribe(string(RGateway))
	_, err := sub.Receive()
	if err != nil {
		gs.logger.Errorf("gateway: subscribe: %v", err)
		return
	}

	subCh := sub.Channel()

	go func() {
		defer func() {
			sub.Close()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-subCh:
				if !ok {
					return
				}
				gs.directAdmission([]byte(msg.Payload))
			}
		}
	}()
}

// OnNotify receives events when being triggered
func (gs *GatewayServer) OnNotify(e observer.Event) {
	if e.Data == nil || len(e.Data) == 0 {
		return
	}
	gs.directAdmission(e.Data)
}

func (gs *GatewayServer) Gracefully(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Close all grpc client connection
	for _, s := range gs.services {
		if s.GRPCClientConn != nil {
			s.GRPCClientConn.Close()
		}
	}
	if err := gs.Server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		gs.logger.Errorf("Gracefully: Error shutting down %v server: %v\n", gs.Name, err)
	}
	gs.logger.Infof("%v closed\n", gs.Name)
}

// Run executes GatewayServer
func (gs *GatewayServer) Run(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)

	go func() {
		defer cancel()
		gs.logger.Infof("Running %v on [%v] - Production: %v\n", gs.Name, gs.Server.Addr, gs.production)
		if err := gs.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			gs.logger.Fatalf("Gateway Server: %s\n", err)
		}
	}()

	go gs.subscribe(runCtx)

	return nil
}

// NewGatewayServer returns a new gin server
func NewGatewayServer(opt ...Option) (*GatewayServer, error) {
	var ac accounts.AccountsServiceClient
	var err error

	opts := defaultGatewayOptions
	for _, o := range opt {
		o(&opts)
	}

	if !opts.accountsDisabled {
		// Initialize AccountsService
		ac, err = services.NewAccountsClient(
			services.WithAddress(opts.accountsAddr),
			services.WithTimeout(opts.accountsTimeout),
		)
		if err != nil {
			return nil, err
		}
	}

	// Initialize Http Server
	s := &http.Server{
		Addr: opts.addr,
	}
	// Initialize GatewayServer
	gs := &GatewayServer{
		services:    map[string]*k8s.APIService{},
		Server:      s,
		Name:        "Gateway Server",
		logger:      opts.logger,
		k8sClient:   opts.k8sClient,
		production:  opts.production,
		redisClient: opts.redisClient,
		accountSVC:  ac,
	}
	// Initialize K8SClient if nil
	if opts.k8sClient != nil {
		err := gs.fetchAllServices()
		if err != nil {
			return nil, err
		}
	}
	gs.applyRoutes()
	return gs, nil
}
