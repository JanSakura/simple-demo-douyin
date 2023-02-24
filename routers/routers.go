package routers

import (
	"github.com/JanSakura/simple-demo-douyin/Init"
	"github.com/JanSakura/simple-demo-douyin/handlers"
	"github.com/JanSakura/simple-demo-douyin/handlers/comment"
	"github.com/JanSakura/simple-demo-douyin/handlers/userInfo"
	"github.com/JanSakura/simple-demo-douyin/handlers/userLoginRegister"
	"github.com/JanSakura/simple-demo-douyin/handlers/video"
	"github.com/gin-gonic/gin"
)

// 路由
func router() *gin.Engine {
	r := gin.Default()
	r.Static("static", "./static")
	topGroup := r.Group("/douyin/")
	{ //basic interface
		topGroup.GET("feed/", video.FeedVideoListHandler)                                                     //视频流接口
		topGroup.POST("user/register/", userLoginRegister.SHAPassword, userLoginRegister.UserRegisterHandler) //用户注册接口
		topGroup.POST("user/login/", userLoginRegister.SHAPassword, userLoginRegister.UserLoginHandler)       //用户登录
		topGroup.GET("user/", handlers.JWTAuthHandler(), userInfo.UserInfoHandler)                            //用户信息
		topGroup.POST("publish/action/", handlers.JWTAuthHandler(), video.PublishVideoHandler)                //视频投稿
		topGroup.GET("publish/list/", handlers.NoAuthToGetUserId(), video.GetVideoListHandler)                //发布列表
		//互动接口
		topGroup.POST("favorite/action/", handlers.JWTAuthHandler(), video.PostFavorHandler)         //点赞
		topGroup.GET("favorite/list/", handlers.NoAuthToGetUserId(), video.GetFavorVideoListHandler) //喜欢列表
		topGroup.POST("comment/action/", handlers.JWTAuthHandler(), comment.PostCommentHandler)      //评论操作
		topGroup.GET("comment/list/", handlers.JWTAuthHandler(), comment.GetCommentListHandler)      //视频评论列表
		//社交接口
		topGroup.POST("relation/action/", handlers.JWTAuthHandler(), userInfo.PostRelationActionHandler)   //关系操作,关注取关
		topGroup.GET("relation/follow/list/", handlers.NoAuthToGetUserId(), userInfo.GetFollowListHandler) //用户关注列表
		topGroup.GET("relation/follower/list/", handlers.NoAuthToGetUserId(), userInfo.GetFollowerHandler) //用户粉丝列表
		topGroup.GET("relation/friend/list/", handlers.NoAuthToGetUserId())                                //用户好友列表
		//社交--消息接口
		topGroup.GET("message/chat/")    //聊天记录
		topGroup.POST("message/action/") //发送消息操作
	}
	return r
}

func InitRouter() *gin.Engine {
	Init.InitDB()
	r := router()
	return r
}
