package userInfo

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/userInfo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type followerListResponse struct {
	models.ResponseStatus
	*userInfo.FollowerList
}

func GetFollowerHandler(c *gin.Context) {
	newProxyQueryFollowerHandler(c).Do()
}

type proxyQueryFollowerHandler struct {
	*gin.Context

	userId int64

	*userInfo.FollowerList
}

func newProxyQueryFollowerHandler(context *gin.Context) *proxyQueryFollowerHandler {
	return &proxyQueryFollowerHandler{Context: context}
}

func (p *proxyQueryFollowerHandler) Do() {
	var err error
	if err = p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	if err = p.prepareData(); err != nil {
		if errors.Is(err, userInfo.ErrUserNotExist) {
			p.SendError(err.Error())
		} else {
			p.SendError("准备数据出错")
		}
		return
	}
	p.SendOk("成功")
}

func (p *proxyQueryFollowerHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId
	return nil
}

func (p *proxyQueryFollowerHandler) prepareData() error {
	list, err := userInfo.QueryFollowerList(p.userId)
	if err != nil {
		return err
	}
	p.FollowerList = list
	return nil
}

func (p *proxyQueryFollowerHandler) SendError(msg string) {
	p.JSON(http.StatusOK, followerListResponse{
		ResponseStatus: models.ResponseStatus{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func (p *proxyQueryFollowerHandler) SendOk(msg string) {
	p.JSON(http.StatusOK, followerListResponse{
		ResponseStatus: models.ResponseStatus{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FollowerList: p.FollowerList,
	})
}
