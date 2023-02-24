package models

import "time"

type Message struct {
	MessageId int64     `json:"message_id"` //消息id
	Id        int64     //用户ID
	TargetId  int64     //发送的目标用户ID
	CreatedAt time.Time //创建时间
	Context   string    //创建内容
}
