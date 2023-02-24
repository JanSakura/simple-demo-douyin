package dao

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"sync"
)

// UserLoginDAO 抽象的User空结构体，用于核心的密码，基本只是在登录注册等状态校验中用
type UserLoginDAO struct {
}

var (
	userLoginDAO  *UserLoginDAO
	userLoginOnce sync.Once
)

// NewUserLoginDao Dao表示方法，DAO用于变量，新增用户登录
func NewUserLoginDao() *UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDAO = new(UserLoginDAO)
	})
	return userLoginDAO
}

func (u *UserLoginDAO) InquiryUserLogin(userName string, password string, login *models.UserLogin) error {
	if login == nil {
		return errNullPtr
	}
	models.GlobalDB.Where("username=? and password=?", userName, password).First(login)
	if login.Id == 0 {
		return errors.New("用户不存在，或账号、密码有错误")
	}
	return nil
}

func (u *UserLoginDAO) IsUserExistByUsername(userName string) bool {
	var userLogin models.UserLogin
	models.GlobalDB.Where("username=?", userName).First(&userLogin)
	if userLogin.Id == 0 {
		return false
	}
	return true
}
