package router

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// RouterRegistry 路由注册器
type RouterRegistry struct {
	routers []Router
}

// RouterRegistryParams 路由注册器参数
type RouterRegistryParams struct {
	fx.In

	Routers []Router `group:"routers"`
}

// NewRouterRegistry 创建路由注册器
func NewRouterRegistry(params RouterRegistryParams) *RouterRegistry {
	return &RouterRegistry{
		routers: params.Routers,
	}
}

// RegisterAllRoutes 注册所有路由
func (r *RouterRegistry) RegisterAllRoutes(app *fiber.App) {
	// 为每个路由器创建对应的路由组
	for _, router := range r.routers {
		prefix := router.GetPrefix()
		if prefix != "" {
			group := app.Group(prefix)
			// 在路由组中注册路由
			router.RegisterRoutes(group)
		} else {
			// 直接在app上注册路由
			router.RegisterRoutes(app)
		}
	}
}
