package route

import (
	controller "organization_management/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/signup", controller.RegisterUser())
	routerGroup.POST("/signin", controller.LoginUser())
	routerGroup.POST("/refresh-token", controller.RefreshToken())
}