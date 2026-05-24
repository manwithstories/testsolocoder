package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	SortBy     string `json:"sort_by"`
	SortOrder  string `json:"sort_order"`
	Keyword    string `json:"keyword"`
}

func GetPagination(c *gin.Context) Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	keyword := c.Query("keyword")

	return Pagination{
		Page:      page,
		PageSize:  pageSize,
		Offset:    (page - 1) * pageSize,
		Limit:     pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Keyword:   keyword,
	}
}

func GetTotalPages(total int64, pageSize int) int {
	if pageSize == 0 {
		return 0
	}
	pages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		pages++
	}
	return pages
}
