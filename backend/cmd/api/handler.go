package main

import (
	httprepo "github.com/Vedant2104/inventory-system/internals/adapters/http"
	"go.uber.org/fx"
)

var HandlerModule = fx.Options(
	fx.Provide(
		httprepo.NewProductHandler,
		httprepo.NewProductCategoryHandler,
	),
)
