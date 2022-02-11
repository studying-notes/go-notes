package api

import (
	"github.com/gin-gonic/gin"
)

// 用法示例

// @Summary ping example
// @Description do ping
// @Tags example
// @Accept json
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /examples/ping [get]
func Ping(c *gin.Context) {}

// @Summary calc example
// @Description plus
// @Tags example
// @Accept json
// @Produce json
// @Param val1 query int true "used for calc"
// @Param val2 query int true "used for calc"
// @Success 200 {integer} string "请求结果"
// @Failure 400 {string} string "操作失败"
// @Router /examples/calc [get]
func Calc(c *gin.Context) {}

// @Summary path params example
// @Description path params
// @Tags example
// @Accept json
// @Produce json
// @Param group_id path int true "Group ID"
// @Param account_id path int true "Account ID"
// @Success 200 {string} string "请求结果"
// @Failure 400 {string} string "操作失败"
// @Router /examples/groups/{group_id}/accounts/{account_id} [get]
func PathParams(c *gin.Context) {}

// @Summary custom header example
// @Description custom header
// @Tags example
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} string "请求结果"
// @Failure 400 {string} string "操作失败"
// @Router /examples/header [get]
func Header(c *gin.Context) {}

// @Summary security example
// @Description security
// @Tags example
// @Accept json
// @Produce json
// @Param Authorization header string true "Authentication header"
// @Success 200 {string} string "请求结果"
// @Failure 400 {string} string "操作失败"
// @Security ApiKeyAuth
// @Security OAuth2Implicit[admin, write]
// @Router /examples/securities [get]
func Securities(c *gin.Context) {}

// @Summary attribute example
// @Description attribute
// @Tags example
// @Accept json
// @Produce json
// @Param enumstring query string false "string enums" Enums(A, B, C)
// @Param enumint query int false "int enums" Enums(1, 2, 3)
// @Param enumnumber query number false "int enums" Enums(1.1, 1.2, 1.3)
// @Param string query string false "string valid" minlength(5) maxlength(10)
// @Param int query int false "int valid" mininum(1) maxinum(10)
// @Param default query string false "string default" default(A)
// @Success 200 {string} string "请求结果"
// @Failure 400 {string} string "操作失败"
// @Router /examples/attribute [get]
func Attribute(c *gin.Context) {}
