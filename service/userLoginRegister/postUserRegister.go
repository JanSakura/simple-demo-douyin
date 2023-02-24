package userLoginRegister

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/dao"
	"github.com/JanSakura/simple-demo-douyin/handlers"
	"github.com/JanSakura/simple-demo-douyin/models"
)

const (
	MaxUserLen     = 100
	MaxPasswordLen = 30
	MinPasswordLen = 8
)

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
type queryUserLoginFlow struct {
	username string
	password string
	data     *LoginResponse
	userId   int64
	token    string
}

// QueryUserLogin 查询用户是否存在，并返回token和id
func QueryUserLogin(username, password string) (*LoginResponse, error) {
	return NewQueryUserLoginFlow(username, password).Do()
}

func NewQueryUserLoginFlow(username, password string) *queryUserLoginFlow {
	return &queryUserLoginFlow{username: username, password: password}
}

func (q *queryUserLoginFlow) Do() (*LoginResponse, error) {
	//对参数进行合法性验证
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	//准备好数据
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	//打包最终数据
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *queryUserLoginFlow) checkNum() error {
	if q.username == "" {
		return errors.New("用户名为空")
	}
	if len(q.username) > MaxUserLen {
		return errors.New("用户名长度超出限制")
	}
	if q.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (q *queryUserLoginFlow) prepareData() error {
	userLoginDAO := dao.NewUserLoginDao()
	var login models.UserLogin
	//准备好userid
	err := userLoginDAO.InquiryUserLogin(q.username, q.password, &login)
	if err != nil {
		return err
	}
	q.userId = login.UserInfoId

	//颁发token
	token, err := handlers.GenerateToken(login)
	if err != nil {
		return err
	}
	q.token = token
	return nil
}

func (q *queryUserLoginFlow) packData() error {
	q.data = &LoginResponse{
		UserId: q.userId,
		Token:  q.token,
	}
	return nil
}
