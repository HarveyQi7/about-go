package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authority struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Role      string             `json:"role" bson:"role"`
	Level     int64              `json:"level" bson:"level"`
	UserMgm   int64              `json:"userMgm" bson:"userMgm"`
	AddAuth   bool               `json:"addAuth" bson:"addAuth"`
	UpdAuth   bool               `json:"updAuth" bson:"updAuth"`
	DelAuth   bool               `json:"delAuth" bson:"delAuth"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func (a Authority) ColName() string {
	return "authority"
}
