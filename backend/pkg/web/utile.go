package web

import (
	"cc.allio/fusion/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"io"
	"mime/multipart"
)

// ParseNumberForQuery 从Get请求中获取给定的参数Integer类型，如果为""或则错误则返回参数orElse给定的数值
func ParseNumberForQuery(c *gin.Context, param string, orElse int) int {
	value := c.Query(param)
	return utils.ToStringInt(value, orElse)
}

// ParseBoolForQuery 从Get请求中获取给定的参数bool类型，如果为""或则错误则返回参数orElse给定的数值
func ParseBoolForQuery(c *gin.Context, param string, orElse bool) bool {
	value := c.Query(param)
	return lo.If[bool](value == "", orElse).ElseF(func() bool {
		return utils.ToStringBool(value)
	})
}

// ParseNumberForPath 从Get请求的路径中获取参数并转换为int
func ParseNumberForPath(c *gin.Context, param string, orElse int) int {
	value := c.Param(param)
	return utils.ToStringInt(value, orElse)
}

// ReadMultipartFile read specific multipart file then close file
func ReadMultipartFile(file multipart.File) ([]byte, error) {
	body, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return body, nil
}
