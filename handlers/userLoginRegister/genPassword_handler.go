package userLoginRegister

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
)

func SHA256(s string) string {
	//SHA256 口令
	auth := sha256.New()
	auth.Write([]byte(s))
	return hex.EncodeToString(auth.Sum(nil))
}

func SHAPassword(c *gin.Context) {
	password := c.Query("password")
	if password == "" {
		password = c.PostForm("password")
	}
	salt := "salt" //简单盐值
	c.Set("password", SHA256(password+salt))
	c.Next()
}
