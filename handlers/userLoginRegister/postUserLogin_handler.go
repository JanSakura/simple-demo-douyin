package userLoginRegister

import (
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/userLoginRegister"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userLoginResponse struct {
	models.ResponseStatus
	*userLoginRegister.LoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	raw, _ := c.Get("password")
	password, ok := raw.(string) //检测password类型是否为string
	if !ok {
		c.JSON(http.StatusOK, userLoginResponse{
			ResponseStatus: models.ResponseStatus{
				StatusCode: 1,
				StatusMsg:  "密码解析错误",
			},
		})
	}
	//获得用户登录响应
	userLoginResp, err := userLoginRegister.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, userLoginResponse{
			ResponseStatus: models.ResponseStatus{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, userLoginResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 0},
		LoginResponse:  userLoginResp,
	})
}
