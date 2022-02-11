package api

import "github.com/gin-gonic/gin"

// @Summary Auth admin
// @Description get admin info
// @Tags accounts,admin
// @Accept json
// @Produce json
// @Success 200 {object} response.Admin
// @Failure 400 {object} HTTPError
// @Security ApiKeyAuth
// @Router /admin/auth [post]
func Auth(c *gin.Context) {}

// @Summary Show a bottle
// @Description get string by ID
// @ID get-string-by-int
// @Tags bottles
// @Accept json
// @Produce json
// @Param id path int true "Bottle ID"
// @Success 200 {object} response.Bottle
// @Failure 400 {object} HTTPError
// @Router /bottles/{id} [get]
func One(c *gin.Context) {}

// @Summary List bottles
// @Description get bottles
// @Tags bottles
// @Accept json
// @Produce json
// @Success 200 {array} response.Bottle
// @Failure 400 {object} HTTPError
// @Router /bottles [get]
func List(c *gin.Context) {}

// @Summary Add an account
// @Description add by json account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body request.AddAccount true "Add account"
// @Success 200 {object} response.Account
// @Failure 400 {object} HTTPError
// @Router /accounts [post]
func Add(c *gin.Context) {}
