package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"swagger/api"
	_ "swagger/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server.
// @termsOfService https://github.com/fujiawei-dev

// @contact.name Rustle Karl
// @contact.url https://github.com/fujiawei-dev
// @contact.email fu.jiawei@outlook.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes http https
// @host localhost:8080
// @BasePath /api/v1
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	v1.POST("/admin/auth", api.Auth)
	v1.POST("/accounts", api.Add)
	v1.GET("/bottles:id", api.One)
	v1.GET("/bottles", api.List)

	examples := v1.Group("/examples")
	{
		examples.GET("ping", api.Ping)
		examples.GET("calc", api.Calc)
		examples.GET("groups/:group_id/accounts/:account_id", api.PathParams)
		examples.GET("header", api.Header)
		examples.GET("securities", api.Securities)
		examples.GET("attribute", api.Attribute)
	}

	_ = r.Run(":8080")
}
