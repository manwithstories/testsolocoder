package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Pagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page < 1 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		if pageSize > 100 {
			pageSize = 100
		}

		c.Set("page", page)
		c.Set("page_size", pageSize)
		c.Next()
	}
}

func ParseUintID(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param(paramName)
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.Set(paramName, uint(0))
		} else {
			c.Set(paramName, uint(id))
		}
		c.Next()
	}
}
