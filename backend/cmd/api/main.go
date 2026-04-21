package main

import "go.uber.org/fx"

func main() {
	fx.New(
		ConfigModule,
		MongoModule,
		RepositoryModule,
		ServiceModule,
		HandlerModule,
		ServerModule,
	).Run()
}
