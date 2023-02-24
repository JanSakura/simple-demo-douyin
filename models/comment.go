package models

import "time"

type Comment struct {
	Id         int64     `json:"id"`                   //评论ID
	UserInfoId int64     `json:"-"`                    //用于一对多关系的用户id
	VideoId    int64     `json:"-"`                    //一对多，视频对评论
	User       UserInfo  `json:"user" gorm:"-"`        //评论的用户信息
	Content    string    `json:"content"`              //内容
	CreatedAt  time.Time `json:"-"`                    //创建时间
	CreateDate string    `json:"create_date" gorm:"-"` //创建日期
}
