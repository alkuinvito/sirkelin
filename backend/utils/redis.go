package utils

import (
	"github.com/alkuinvito/sirkelin/initializers"
	"github.com/gin-gonic/gin"
	"time"
)

func SetBlacklist(c *gin.Context, jti string) error {
	return initializers.Rdb.SetEx(c, jti, nil, time.Hour).Err()
}

func CheckBlacklist(c *gin.Context, jti string) bool {
	return initializers.Rdb.Exists(c, jti).Val() == 1
}
