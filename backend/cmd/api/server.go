package main

import (
	"context"
	"log"
	"net/http"

	httprepo "github.com/Vedant2104/inventory-system/internals/adapters/http"
	"github.com/Vedant2104/inventory-system/internals/infrastructure/logger"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type ServiceParams struct {
	fx.In

	Config          *Config
	ProductHandler  *httprepo.ProductHandler
	CategoryHandler *httprepo.ProductCategoryHandler
	Logger          zerolog.Logger
	LC              fx.Lifecycle
}

func NewHTTPServer(p ServiceParams) *http.ServeMux {
	router := http.NewServeMux()
	httprepo.RegisterProductHandler(router, p.ProductHandler)
	httprepo.RegisterProductCategoryHandler(router, p.CategoryHandler)

	handler := httprepo.RequestLogger(p.Logger)(httprepo.EnableCORS()(router))

	server := &http.Server{
		Addr:    ":" + p.Config.Port,
		Handler: handler,
	}

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting HTTP server on port", p.Config.Port)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Failed to start HTTP server", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return router
}

var ServerModule = fx.Options(
	fx.Provide(
		logger.GetLogger,
		NewHTTPServer,
	),
	fx.Invoke(func(*http.ServeMux) {}),
)
