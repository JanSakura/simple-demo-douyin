package comment

import (
	"errors"
	"fmt"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/comment"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type proxyPostCommentHandler struct {
	*gin.Context

	videoId     int64
	userId      int64
	commentId   int64
	actionType  int64
	commentText string
}

func newProxyPostCommentHandler(context *gin.Context) *proxyPostCommentHandler {
	return &proxyPostCommentHandler{Context: context}
}

type postCommentResponse struct {
	models.ResponseStatus
	*comment.Response
}

func PostCommentHandler(c *gin.Context) {
	newProxyPostCommentHandler(c).Do()
}

func (p *proxyPostCommentHandler) Do() {
	//解析参数
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}

	//正式调用Service层
	commentRes, err := comment.PostComment(p.userId, p.videoId, p.commentId, p.actionType, p.commentText)
	if err != nil {
		p.SendError(err.Error())
		return
	}

	//成功返回
	p.SendOk(commentRes)
}

func (p *proxyPostCommentHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId

	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	p.videoId = videoId

	//根据actionType解析对应的可选参数
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	switch actionType {
	case comment.CREATE:
		p.commentText = p.Query("comment_text")
	case comment.DELETE:
		p.commentId, err = strconv.ParseInt(p.Query("comment_id"), 10, 64)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("未定义的行为%d", actionType)
	}
	p.actionType = actionType
	return nil
}

func (p *proxyPostCommentHandler) SendError(msg string) {
	p.JSON(http.StatusOK, postCommentResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 1, StatusMsg: msg}, Response: &comment.Response{}})
}

func (p *proxyPostCommentHandler) SendOk(comment *comment.Response) {
	p.JSON(http.StatusOK, postCommentResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 0},
		Response:       comment,
	})
}
