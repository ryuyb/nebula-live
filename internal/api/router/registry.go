package router

import "github.com/gofiber/fiber/v2"

// Router 定义了路由设置的接口
type Router interface {
	SetupRoutes(v1 *fiber.Router)
}

// RouterRegistry 管理所有的路由器
type RouterRegistry struct {
	routers []Router
}

// NewRouterRegistry 创建一个新的路由器注册表
func NewRouterRegistry(routers ...Router) *RouterRegistry {
	return &RouterRegistry{
		routers: routers,
	}
}

// RegisterRouter 注册一个新的路由器
func (r *RouterRegistry) RegisterRouter(router Router) {
	r.routers = append(r.routers, router)
}

// SetupAllRoutes 设置所有注册的路由
func (r *RouterRegistry) SetupAllRoutes(v1 *fiber.Router) {
	for _, router := range r.routers {
		router.SetupRoutes(v1)
	}
}
