package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ctrl "sigs.k8s.io/controller-runtime"

	"kubesphere.io/kubesphere/pkg/virtualization/gateway"
)

func main() {
	cfg, err := ctrl.GetConfig()
	if err != nil {
		log.Fatalf("failed to load kubeconfig: %v", err)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	router, err := gateway.NewRouter(cfg)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	gateway.RegisterRoutes(engine, router)

	srv := &http.Server{
		Addr:              ":9444",
		Handler:           engine,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("virtualization gateway listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("gateway failed: %v", err)
	}
}
