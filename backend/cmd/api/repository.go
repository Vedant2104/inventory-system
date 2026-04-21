package main

import (
	mongorepo "github.com/Vedant2104/inventory-system/internals/adapters/repository/mongo"
	"github.com/Vedant2104/inventory-system/internals/ports"
	"go.uber.org/fx"
)

var RepositoryModule = fx.Options(
	fx.Provide(
		fx.Annotate(
			mongorepo.NewProductRepository,
			fx.ParamTags(`name:"productCollection"`),
			fx.As(new(ports.ProductRepository)),
		),
		fx.Annotate(
			mongorepo.NewProductCategoryRepository,
			fx.ParamTags(`name:"categoryCollection"`),
			fx.As(new(ports.ProductCategoryRepository)),
		),
	),
)
