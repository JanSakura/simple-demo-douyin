package userInfo

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/userInfo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PostRelationActionHandler(c *gin.Context) {
	newProxyPostFollowAction(c).Do()
}

type proxyPostFollowAction struct {
	*gin.Context

	userId     int64
	followId   int64
	actionType int
}

func newProxyPostFollowAction(context *gin.Context) *proxyPostFollowAction {
	return &proxyPostFollowAction{Context: context}
}

func (p *proxyPostFollowAction) Do() {
	var err error
	if err = p.prepareNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	if err = p.startAction(); err != nil {
		//当错误为model层发生的，那么就是重复键值的插入了
		if errors.Is(err, userInfo.ErrIvdAct) || errors.Is(err, userInfo.ErrIvdFolUsr) {
			p.SendError(err.Error())
		} else {
			p.SendError("请勿重复关注")
		}
		return
	}
	p.SendOk("操作成功")
}

func (p *proxyPostFollowAction) prepareNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId

	//解析需要关注的id
	followId := p.Query("to_user_id")
	parseInt, err := strconv.ParseInt(followId, 10, 64)
	if err != nil {
		return err
	}
	p.followId = parseInt

	//解析action_type
	actionType := p.Query("action_type")
	parseInt, err = strconv.ParseInt(actionType, 10, 32)
	if err != nil {
		return err
	}
	p.actionType = int(parseInt)
	return nil
}

func (p *proxyPostFollowAction) startAction() error {
	err := userInfo.PostFollowAction(p.userId, p.followId, p.actionType)
	if err != nil {
		return err
	}
	return nil
}

func (p *proxyPostFollowAction) SendError(msg string) {
	p.JSON(http.StatusOK, models.ResponseStatus{StatusCode: 1, StatusMsg: msg})
}

func (p *proxyPostFollowAction) SendOk(msg string) {
	p.JSON(http.StatusOK, models.ResponseStatus{StatusCode: 1, StatusMsg: msg})
}
