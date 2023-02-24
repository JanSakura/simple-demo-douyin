package userLoginRegister

import (
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/userLoginRegister"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userRegisterResponse struct {
	models.ResponseStatus
	*userLoginRegister.LoginResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Query("username")
	rawVal, _ := c.Get("password")
	password, ok := rawVal.(string)
	if !ok {
		c.JSON(http.StatusOK, userRegisterResponse{
			ResponseStatus: models.ResponseStatus{
				StatusCode: 1,
				StatusMsg:  "密码解析出错",
			},
		})
		return
	}
	userRegisterResp, err := userLoginRegister.PostUserLogin(username, password)

	if err != nil {
		c.JSON(http.StatusOK, userRegisterResponse{
			ResponseStatus: models.ResponseStatus{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, userRegisterResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 0},
		LoginResponse:  userRegisterResp,
	})
}
