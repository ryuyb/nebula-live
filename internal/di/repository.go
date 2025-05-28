package di

import (
	"go.uber.org/fx"
	"nebulaLive/internal/repository"
)

// RepositoryModule 提供存储库相关的依赖项
var RepositoryModule = fx.Options(
	fx.Provide(
		repository.NewClient,
		repository.NewUserRepository,
	),
)
