package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	Title       string    `bson:"title,omitempty" json:"title,omitempty"`
	Description string    `bson:"description,omitempty" json:"description,omitempty"`
	Keywords    []string  `bson:"keywords,omitempty" json:"keywords,omitempty"`
	Content     string    `bson:"content,omitempty" json:"content,omitempty"`
	UpdateAt    time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt    time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
