package registry

import (
	"github.com/gin-gonic/gin"
)

// TBD: Inject Serivce Here

// Router is an interface to register router handlers to base router
type Router interface {
	RegisterRoutes(base *gin.RouterGroup)
}

var routerFactories []func() Router

// RegisterRouter registers a router to the routers registry
func RegisterRouter(router func() Router) {
	routerFactories = append(routerFactories, router)
}

// Routers returns the registered routers from the registry
func Routers() []func() Router {
	return routerFactories
}
