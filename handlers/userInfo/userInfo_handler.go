package userInfo

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/dao"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type userResponse struct {
	models.ResponseStatus
	User *models.UserInfo `json:"user"`
}

// UserInfoHandler 获取用户信息
func UserInfoHandler(c *gin.Context) {
	p := newProxyUserInfo(c)
	//得到上层中间件根据token解析的userId
	rawId, ok := c.Get("user_id")
	if !ok {
		p.UserInfoError("解析userId出错")
		return
	}
	err := p.DoQueryUserInfoByUserId(rawId)
	if err != nil {
		p.UserInfoError(err.Error())
	}
}

// 注意小写，在router处隐藏
type proxyUserInfo struct {
	c *gin.Context
}

func newProxyUserInfo(c *gin.Context) *proxyUserInfo {
	return &proxyUserInfo{c: c}
}

func (p *proxyUserInfo) DoQueryUserInfoByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("解析userId失败")
	}
	//由于得到userinfo不需要处理dao层的数据，所以直接调用dao层的接口
	userinfoDAO := dao.NewUserInfoDao()

	var userInfo models.UserInfo
	err := userinfoDAO.InquiryUserInfoByID(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

func (p *proxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, userResponse{
		ResponseStatus: models.ResponseStatus{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func (p *proxyUserInfo) UserInfoOk(user *models.UserInfo) {
	p.c.JSON(http.StatusOK, userResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 0},
		User:           user,
	})
}
