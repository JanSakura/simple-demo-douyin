package handlers

import (
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NoAuthToGetUserId() gin.HandlerFunc { //特意表明是一个handler
	return func(c *gin.Context) {
		rawId := c.Query("user_id")
		if rawId == "" {
			rawId = c.PostForm("user_id")
		}
		//用户不存在
		if rawId == "" {
			c.JSON(http.StatusOK, models.ResponseStatus{StatusCode: http.StatusUnauthorized, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		userId, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, models.ResponseStatus{StatusCode: http.StatusUnauthorized, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
		}
		c.Set("user_id", userId)
		c.Next()
	}
}
