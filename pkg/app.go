package pkg

import (
	route "organization_management/pkg/api/routes"
	database "organization_management/pkg/database/mongodb"
	util "organization_management/pkg/utils"

	"github.com/gin-gonic/gin"
)

func StartApplication() {
	port := util.EnvPort()
	router := gin.New()

	// apply middleware
	router.Use(gin.Logger())

	//run database
	database.ConnectDB()

	// apply routes
	public := router.Group("/api")
	{
		route.AuthRoutes(public)
	}

	// run the server
	router.Run(":" + port)
}
