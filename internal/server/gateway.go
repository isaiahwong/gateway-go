package server

// TODO Implement queue
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/isaiahwong/gateway-go/internal/k8s"
	"github.com/isaiahwong/gateway-go/internal/k8s/enum"
	"github.com/isaiahwong/gateway-go/internal/observer"
	"github.com/isaiahwong/gateway-go/internal/util/log"
	"github.com/isaiahwong/gateway-go/internal/util/validator"
	"github.com/isaiahwong/gateway-go/protogen"
)

// GatewayServer encapsulates GatewayServer and Observer
type GatewayServer struct {
	Name      string
	Server    *http.Server
	services  map[string]*k8s.APIService
	logger    log.Logger
	k8sClient *k8s.Client
}

type gatewayOptions struct {
	logger    log.Logger
	k8sClient *k8s.Client
}

var defaultGatewayOption = gatewayOptions{
	logger: log.NewLogger(),
}

// GatewayOption sets options for GatewayServer.
type GatewayOption func(*gatewayOptions)

// Logger sets logger for gateway
func WithLogger(l log.Logger) GatewayOption {
	return func(o *gatewayOptions) {
		o.logger = l
	}
}

// K8SClient sets k8s client for GatewayServer.
// Though there isn't a generic type interface :(
func WithK8SClient(k *k8s.Client) GatewayOption {
	return func(o *gatewayOptions) {
		o.k8sClient = k
	}
}

// OnNotify receives events when being triggered
func (gs *GatewayServer) OnNotify(e observer.Event) {
	if e.Data == nil || len(e.Data) == 0 {
		return
	}
	gs.directAdmission(e.Data)
}

func essentialMiddleware(gs *GatewayServer) func(*gin.Engine) {
	return func(r *gin.Engine) {
		r.Use(gin.Recovery())
		r.Use(requestLogger(gs.logger))
		r.Use(WebhookRequests)
		r.Use(authMW(&gs.services))
		// Health route
		r.GET("/hz", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "success",
			})
		})
		r.NoRoute(notFoundMW)
	}
}

// forwardAllHeaders packages http headers into headers-bin and forwards the metadata
func forwardAllHeaders(_ context.Context, r *http.Request) metadata.MD {
	headers := map[string]string{}
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[strings.ToLower(k)] = v[0]
		}
	}
	h, err := json.Marshal(headers)
	if err != nil {
		return nil
	}
	md := metadata.New(map[string]string{"headers-bin": string(h)})
	return md
}

// newGrpcMux creates a new mux that handles grpc calls.
func newGrpcMux(ctx context.Context) *runtime.ServeMux {
	runtime.HTTPError = HTTPError
	runtime.OtherErrorHandler = OtherErrorHandler
	mux := runtime.NewServeMux(
		runtime.WithMetadata(forwardAllHeaders),
	)
	return mux
}

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

func applyHTTP(r *gin.Engine, path string, route string) error {
	if r == nil {
		return InvalidParams("applyHttp: gin r router is nil")
	}
	if path == "" {
		return InvalidParams("applyHttp: path is empty")
	}
	// TODO: Test connectivity
	rp := ReverseProxy(path)
	// Routes all requests to service
	r.GET(route, rp)
	r.POST(route, rp)
	r.Any(fmt.Sprintf("%v/*any", route), rp)
	return nil
}

func applyGrpc(r *gin.Engine, serviceName string, conn *grpc.ClientConn, target string, route string) error {
	var err error
	var svc *protogen.ServiceDesc
	if r == nil {
		return InvalidParams("r router is nil")
	}
	if serviceName == "" {
		return InvalidParams("serviceName is empty")
	}
	// Check if service exists in proto definition
	for k, s := range protogen.GetProtos() {
		// Splits service name i.e api.auth.authservice ["api", "auth", "authservice"]
		split := strings.Split(k, ".")
		// Formats it to dash i.e api-auth-authservice
		dash := strings.ToLower(strings.Join(split, "-"))
		// Formats it to underscore i.e api_auth_authservice
		underscore := strings.ToLower(strings.Join(split, "-"))
		if serviceName == dash || serviceName == underscore {
			svc = &s
			break
		}
	}
	if svc == nil {
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
	mux := newGrpcMux(context.Background())
	// pass mux to grpc handlers
	svc.Handler(context.Background(), mux, conn)
	handler := func(c *gin.Context) {
		mux.ServeHTTP(c.Writer, c.Request)
	}
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
		// Test if routes is working
		for _, port := range svc.Ports {
			switch port.Name {
			case "http":
				route := fmt.Sprintf("/%v%v", svc.APIversion, svc.Path)
				path := fmt.Sprintf("%v.svc.cluster.local:%v", svc.DNSPath, port.Port)
				err := applyHTTP(r, path, route)
				if err != nil {
					gs.logger.Error(err)
				}
			case "grpc":
				target := fmt.Sprintf("%v.svc.cluster.local:%v", svc.DNSPath, port.Port)
				path := fmt.Sprintf("/%v%v", svc.APIversion, svc.Path)
				err := applyGrpc(r, svc.ServiceName, svc.GRPCClientConn, target, path)
				if err != nil {
					gs.logger.Error(err)
				}
			}
		}
	}
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
	delete(gs.services, service.DNSPath)
}

// fetchAllServices fetches all services from cluster
func (gs *GatewayServer) fetchAllServices() error {
	// Get K8S Services in cluster
	svcs, err := gs.k8sClient.GetServices("default")
	if err != nil {
		return fmt.Errorf("fetchAllServices: GetServices: %v", err)
	}
	for _, d := range svcs {
		o, err := gs.k8sClient.CoreAPI().Admission().UnmarshalK8SObject(d)
		if err != nil {
			return fmt.Errorf("fetchAllServices: %v", err)
		}
		// Filter admission request if incoming request does not have api-service labeled.
		if strings.ToLower(o.Metadata.Labels.ResourceType) != string(enum.LabelAPIService) {
			continue
		}
		s, err := gs.k8sClient.CoreAPI().APIServices().ObjectToAPI(o)
		if err != nil {
			return fmt.Errorf("fetchAllServices: %v", err)
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
		gs.logger.Printf("Admission request %v is not service")
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

// NewGatewayServer returns a new gin server
func NewGatewayServer(port string, opt ...GatewayOption) (*GatewayServer, error) {
	opts := defaultGatewayOption
	for _, o := range opt {
		o(&opts)
	}
	// Initialize Http Server
	s := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
	}
	// Initialize GatewayServer
	gs := &GatewayServer{
		services:  map[string]*k8s.APIService{},
		Server:    s,
		Name:      "Gateway Server",
		logger:    opts.logger,
		k8sClient: opts.k8sClient,
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
