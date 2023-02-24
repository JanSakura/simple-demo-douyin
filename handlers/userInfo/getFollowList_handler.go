package userInfo

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/userInfo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type followListResponse struct {
	models.ResponseStatus
	*userInfo.FollowList
}

func GetFollowListHandler(c *gin.Context) {
	newProxyQueryFollowList(c).Do()
}

type proxyQueryFollowList struct {
	*gin.Context
	userId int64
	*userInfo.FollowList
}

func newProxyQueryFollowList(context *gin.Context) *proxyQueryFollowList {
	return &proxyQueryFollowList{Context: context}
}

func (p *proxyQueryFollowList) Do() {
	var err error
	if err = p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	if err = p.prepareData(); err != nil {
		p.SendError(err.Error())
		return
	}
	p.SendOk("请求成功")
}

func (p *proxyQueryFollowList) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId
	return nil
}

func (p *proxyQueryFollowList) prepareData() error {
	list, err := userInfo.QueryFollowList(p.userId)
	if err != nil {
		return err
	}
	p.FollowList = list
	return nil
}

func (p *proxyQueryFollowList) SendError(msg string) {
	p.JSON(http.StatusOK, followListResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *proxyQueryFollowList) SendOk(msg string) {
	p.JSON(http.StatusOK, followListResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 0, StatusMsg: msg},
		FollowList:     p.FollowList,
	})
}
