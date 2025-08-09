package persistence

import "go.uber.org/fx"

// PersistenceModule 仓储层模块
var PersistenceModule = fx.Options(
	fx.Provide(
		NewUserRepository,
	),
)