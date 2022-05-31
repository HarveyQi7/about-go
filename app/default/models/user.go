package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Username    string             `json:"username" bson:"username"`
	PhonePrefix string             `json:"phonePrefix" bson:"phonePrefix"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	Password    string             `json:"password" bson:"password"`
	Authorities UserAuthority      `json:"authorities" bson:"authorities"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type UserAuthority struct {
	AuthId  primitive.ObjectID `json:"authId" bson:"authId"`
	Role    string             `json:"role" bson:"role"`
	Level   int64              `json:"level" bson:"level"`
	UserMgm int64              `json:"userMgm" bson:"userMgm"`
	AddAuth bool               `json:"addAuth" bson:"addAuth"`
	UpdAuth bool               `json:"updAuth" bson:"updAuth"`
	DelAuth bool               `json:"delAuth" bson:"delAuth"`
}

func (user User) ColName() string {
	return "User"
}
