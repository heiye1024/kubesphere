package gateway

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"

	gatewayclient "kubesphere.io/kubesphere/pkg/virtualization/gateway/client"
	"kubesphere.io/kubesphere/pkg/virtualization/gateway/dto"
	"kubesphere.io/kubesphere/pkg/virtualization/gateway/middleware"
)

// Router wires HTTP routes to handlers.
type Router struct {
	K8s   ctrlclient.Client
	Authn middleware.Authenticator
	Audit gin.HandlerFunc
	Store gatewayclient.VirtualizationStore
}

// NewRouter constructs a Router backed by Kubernetes clients.
func NewRouter(cfg *rest.Config) (*Router, error) {
	k8sClient, err := gatewayclient.NewDynamicClient(cfg)
	if err != nil {
		return nil, err
	}
	authn := middleware.NewJWTImpersonationAuthenticator(cfg)
	store := gatewayclient.NewVirtualizationStore(k8sClient)
	return &Router{
		K8s:   k8sClient,
		Authn: authn,
		Audit: middleware.AuditTrail(),
		Store: store,
	}, nil
}

func allowCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		headers := []string{
			"Authorization",
			"Content-Type",
			"Impersonate-User",
			"Impersonate-Group",
			"Impersonate-UID",
			"X-Request-Id",
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-Deny-Reason")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// RegisterRoutes adds virtualization endpoints to gin.Engine.
func RegisterRoutes(engine *gin.Engine, router *Router) {
	engine.Use(allowCORS())

	basePath := "/kapis/virtualization.kubesphere.io/v1beta1"
	engine.GET(basePath+"/openapi", func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.OpenAPISchema(basePath))
	})

	group := engine.Group(basePath + "/projects/:namespace")
	group.Use(router.Audit, router.Authn.GinMiddleware())

	group.GET("/vms", router.listVMs)
	group.POST("/vms", router.createVM)
	group.POST("/vms/:name:powerOn", router.powerOn)
	group.POST("/vms/:name:powerOff", router.powerOff)
	group.POST("/vms/:name:migrate", router.migrate)
	group.POST("/vms/:name:console", router.console)

	group.GET("/disks", router.listDisks)
	group.POST("/disks", router.createDisk)

	group.GET("/nets", router.listNets)
	group.POST("/nets", router.createNet)

	group.GET("/snapshots", router.listSnapshots)
	group.POST("/snapshots", router.createSnapshot)

	group.GET("/templates", router.listTemplates)
	group.POST("/templates", router.createTemplate)
}
