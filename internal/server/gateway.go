package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/isaiahwong/gateway-go/internal/k8s"
	"github.com/isaiahwong/gateway-go/internal/k8s/enum"
	"github.com/isaiahwong/gateway-go/internal/observer"
	"github.com/isaiahwong/gateway-go/internal/util/log"
)

// GatewayServer encapsulates GatewayServer and Observer
type GatewayServer struct {
	Server   *http.Server
	services map[string]*k8s.APIService
	opts     *gatewayOptions
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
func Logger(l log.Logger) GatewayOption {
	return func(o *gatewayOptions) {
		o.logger = l
	}
}

// K8SClient sets k8s client for GatewayServer.
// Though there isn't a generic type interface :(
func K8SClient(k *k8s.Client) GatewayOption {
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

// updateServices updates gateway services
func (gs *GatewayServer) updateServices(service *k8s.APIService) {
	if service == nil {
		gs.opts.logger.Error("updateServices: service is nil")
		return
	}
	if service.DNSPath == "" {
		gs.opts.logger.Error("updateServices: service.DNSPath is empty")
		return
	}
	// Create map if services is not assigned
	if gs.services == nil {
		gs.services = map[string]*k8s.APIService{}
	}
	gs.services[service.DNSPath] = service
}

// applyRoutes creates a new replaces the server's http handler with
// newly populated routes
func (gs *GatewayServer) applyRoutes() {
	// Create new Router
	r, err := newRouter(
		httpMiddleware,
		// gatewayMiddleware(gs.gw),
	)
	if err != nil {
		gs.opts.logger.Errorf("applyRoutes: %v", err)
	}
	// Create a new http router
	for _, svc := range gs.services {
		for _, port := range svc.Ports {
			gs.opts.logger.Println(svc.DNSPath, port.Port)
			switch port.Name {
			case "http":
				path := fmt.Sprintf("/%v%v", svc.APIversion, svc.Path)
				rp := reverseProxy(fmt.Sprintf("http://%v:%v%v", svc.DNSPath, port.Port, path))
				// Routes GET root path
				r.GET(path, rp)
				// Routes all requests to service
				r.Any(fmt.Sprintf("%v/*any", path), rp)
			}
		}
	}

	// Apply routes to server handler
	gs.Server.Handler = r
}

// fetchAllServices fetches all services from cluster
func (gs *GatewayServer) fetchAllServices() error {
	// Get K8S Services in cluster
	svcs, err := gs.opts.k8sClient.GetServices("default")
	if err != nil {
		return err
	}
	for _, d := range svcs {
		o, err := gs.opts.k8sClient.CoreAPI().Admission().UnmarshalObject(d)
		if err != nil {
			return err
		}
		// Filter admission request if incoming request does not have api-service labeled.
		if strings.ToLower(o.Metadata.Labels.ResourceType) != string(enum.LabelAPIService) {
			continue
		}
		s, err := gs.opts.k8sClient.CoreAPI().APIServices().ObjectToAPI(o)
		if err != nil {
			return err
		}
		gs.updateServices(s)
	}
	return nil
}

// directAdmission streamlines a series operations
// such as parsing AdmissionRequest, filtering
// and routing to its necessary
func (gs *GatewayServer) directAdmission(d []byte) {
	ar, err := gs.opts.k8sClient.CoreAPI().Admission().Unmarshal(d)
	if err != nil {
		gs.opts.logger.Error(err)
		return
	}
	// Filter admission request if incoming request is not of K8S Service Object.
	if strings.ToLower(ar.Request.Kind.Kind) != string(enum.K8SServiceObject) {
		gs.opts.logger.Printf("Admission request %v is not service")
		return
	}
	// Filter admission request if incoming request does not have api-service labeled.
	if strings.ToLower(ar.Request.Object.Metadata.Labels.ResourceType) != string(enum.LabelAPIService) {
		return
	}
	switch ar.Request.Operation {
	case string(enum.Create):
		gs.create(ar)
	}
}

// create adds apiservices to gatway services
func (gs *GatewayServer) create(ar *k8s.AdmissionRegistration) {
	s, err := gs.opts.k8sClient.CoreAPI().APIServices().ObjectToAPI(&ar.Request.Object)
	if err != nil {
		gs.opts.logger.Error(err)
	}
	gs.updateServices(s)
	gs.applyRoutes()
}

type ServiceDesc struct {
	ServiceName    string
	PackageSVC     string
	Package        string
	CurrentPackage string
	Handler        func(context.Context, *gwruntime.ServeMux, *grpc.ClientConn) error
}

func gatewayMiddleware(gw *gwruntime.ServeMux) func(*gin.Engine) {
	return func(r *gin.Engine) {
		if gw == nil {
			return
		}
		// Proxies to gateway services
		r.Any("/v1/*any", gin.WrapF(gw.ServeHTTP))
	}
}

func httpMiddleware(r *gin.Engine) {
	r.Use(gin.Recovery())
	// Health route
	r.GET("/hz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "success",
		})
	})
	r.NoRoute(notFound)
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

// func newGateway(ctx context.Context, opts []gwruntime.ServeMuxOption) (*gwruntime.ServeMux, error) {
// 	mux := gwruntime.NewServeMux(opts...)
// 	conn, err := grpc.DialContext(ctx, "payment-service.default:50051", grpc.WithInsecure())
// 	if err != nil {
// 		return nil, err
// 	}
// 	for k, s := range GetProtos() {
// 		err = f(ctx, mux, conn)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return mux, nil
// }

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
		opts:     &opts,
		services: map[string]*k8s.APIService{},
		Server:   s,
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
