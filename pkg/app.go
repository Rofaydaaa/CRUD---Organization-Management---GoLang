package pkg

import (
	middleware "organization_management/pkg/api/middleware"
	route "organization_management/pkg/api/routes"
	database "organization_management/pkg/database/mongodb"
	util "organization_management/pkg/utils"
	redis "organization_management/pkg/database/redis"

	"github.com/gin-gonic/gin"
)

func StartApplication() {
	port := util.EnvPort()
	router := gin.New()

	// apply middleware
	router.Use(gin.Logger())

	//run database
	database.ConnectDB()
	redis.InitRedis()

	// apply routes
	public := router.Group("/api")
	{
		route.AuthRoutes(public)
	}
	protected := router.Group("/api")
	{
		protected.Use(middleware.JwtAuthMiddleware())
		route.OrganizationRoutes(protected)
		route.ProtectedUderRoutes(protected)
	}

	// run the server
	router.Run(":" + port)
}
