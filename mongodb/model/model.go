package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	// 用户名
	Username string `bson:"username" json:"username"`

	// 密码
	Password string `bson:"password" json:"password,omitempty"`

	// 电子邮件
	Email string `bson:"email" json:"email"`

	// 所属部门
	Department *primitive.ObjectID `bson:"department" json:"-"`

	// 权限组
	Roles []primitive.ObjectID `bson:"roles" json:"roles,omitempty"`

	// 称呼
	Name string `bson:"name" json:"name"`

	// 头像
	Avatar string `bson:"avatar" json:"avatar"`

	// 飞书 OpenID
	Feishu bson.M `json:"feishu" bson:"feishu"`

	// 标签
	Labels map[string]string `bson:"labels" json:"labels"`

	// 状态
	Status bool `bson:"status" json:"status"`

	// 会话次数
	Sessions int64 `bson:"sessions" json:"sessions"`

	// 最近一次登录时间
	LastTime time.Time `bson:"last_time" json:"last_time"`

	// 创建时间
	CreateTime time.Time `bson:"create_time" json:"-"`

	// 更新时间
	UpdateTime time.Time `bson:"update_time" json:"-"`
}

type Record struct {
	Time  time.Time `bson:"time"`
	Email string    `bson:"email" faker:"email"`
	Name  string    `bson:"name" faker:"name"`
	Ip    string    `bson:"ip" faker:"ipv4"`
}