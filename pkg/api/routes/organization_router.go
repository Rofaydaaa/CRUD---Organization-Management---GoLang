package route

import (
	controller "organization_management/pkg/controllers"
	"github.com/gin-gonic/gin"
)

func OrganizationRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/organization", controller.CreateOrganization())
	routerGroup.GET("/organization/:organization_id", controller.ReadOrganization())
	routerGroup.GET("/organization", controller.ReadAllOrganizations())
	routerGroup.PUT("/organization/:organization_id", controller.UpdateOrganization())
	routerGroup.DELETE("/organization/:organization_id", controller.DeleteOrganization())
	routerGroup.POST("/organization/:organization_id/invite", controller.InviteUserToOrganization())
}