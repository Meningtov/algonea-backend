package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetIntQueryParam(c *gin.Context, name string) (uint64, bool) {
	queryParam := c.Query(name)
	if queryParam == "" {
		return 0, false
	}
	assetId, err := strconv.ParseUint(queryParam, 10, 0)
	if err != nil {
		return 0, false
	}
	return assetId, true
}
