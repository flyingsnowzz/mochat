package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体，与原始 PHP 项目保持一致
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Success 返回成功响应，code=0
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
		Msg:  "success",
	})
}

// SuccessMsg 返回成功响应，仅包含消息
func SuccessMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Data: nil,
		Msg:  msg,
	})
}

// Fail 返回失败响应，HTTP 状态码 200，业务错误码由 code 指定
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}

// FailWithHTTP 返回失败响应，可自定义 HTTP 状态码和业务错误码
func FailWithHTTP(c *gin.Context, httpCode int, code int, msg string) {
	c.JSON(httpCode, Response{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}

// PageResult 返回分页数据响应
func PageResult(c *gin.Context, list interface{}, total int64, currentPage int, pageSize int) {
	totalPage := 0
	if pageSize > 0 {
		totalPage = int(math.Ceil(float64(total) / float64(pageSize)))
	}
	Success(c, gin.H{
		"list":        list,
		"total":       total,
		"currentPage": currentPage,
		"pageSize":    pageSize,
		"page": gin.H{
			"perPage":     pageSize,
			"total":       total,
			"totalPage":   totalPage,
			"currentPage": currentPage,
		},
	})
}
