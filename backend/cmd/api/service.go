package main

import (
	"github.com/Vedant2104/inventory-system/internals/service"
	"go.uber.org/fx"
)

var ServiceModule = fx.Options(
	fx.Provide(
		service.NewProductCategoryService,
		service.NewProductService,
	),
)