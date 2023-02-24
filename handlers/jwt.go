package handlers

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

// 用于签名的字符串
var jwtKey = []byte("DouyinKey") //jwt的秘钥
// myClaims 创建Claim
type myClaims struct {
	UserId               int64
	jwt.RegisteredClaims //最新版的StandardClaims已废弃
}

// GenerateToken 生成JWT:封装生成Token
func GenerateToken(user models.UserLogin) (string, error) {
	//Token过期时间，如24小时
	const TokenExpireDuration = time.Hour * 24
	//创建自己的声明
	claims := &myClaims{
		UserId: user.UserInfoId, //自定义字段
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()), //DefaultClaims是int64格式,即time.now().Unix()
			Issuer:    "douyin_demo",                  //签发人
			Subject:   "GenToken",                     //主题
		}}
	//使用指定的签名方法创建签名对象,如SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//使用指定的secret签名，获得完整编码后的字符串token
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseToken 解析给定的JWT字符串
func ParseToken(tokenStr string) (*myClaims, error) {
	//解析token,如果是自定义Claims则使用ParseWithClaims
	token, err := jwt.ParseWithClaims(tokenStr, &myClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	//有错误，先处理返回
	if err != nil {
		return nil, err
	}
	//对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*myClaims); ok && token.Valid { //token校验
		return claims, nil //校验成功
	}
	return nil, errors.New("invalid token") //校验失败
}

// AuthHandler 认证，获取Token
//func AuthHandler(uid int64) {
//	var user models.UserInfo
//
//}

// JWTAuthHandler 基于JWT的Token认证校验中间件,通过则设置user_id;基础性的Handler，所以特意写的和其他功能性Handler不一样
func JWTAuthHandler() gin.HandlerFunc { //只是标注是HandlerFunc类型
	return func(c *gin.Context) {
		//Token假定放在Header的Authorization,不过不知道是不是就不采用了,直接获取
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, models.ResponseStatus{
				StatusCode: http.StatusUnauthorized,
				StatusMsg:  "用户不存在",
			})
			c.Abort() //阻止向下执行Handler
			return
		}
		//校验token
		tokenMsg, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, models.ResponseStatus{
				StatusCode: http.StatusForbidden,
				StatusMsg:  "token错误",
			})
			c.Abort()
			return
		}
		//token超时
		//通过ExpiresAt值判断是否过期，如果没有设置，则返回true，没有过期
		if !tokenMsg.VerifyExpiresAt(time.Now(), false) {
			c.JSON(http.StatusOK, models.ResponseStatus{
				StatusCode: http.StatusPaymentRequired,
				StatusMsg:  "token过期",
			})
			c.Abort()
			return
		}
		//将当前请求的userId设置到请求上下文的c的user_id
		c.Set("user_id", tokenMsg.UserId)
		//后续处理函数可通过c.Get("user_id")获取当前请求的用户信息
		c.Next()
	}
}
