package model

import (
	"time"
)

type User struct {
	Id        int       `json:"id" xorm:"not null pk autoincr INT(11)"`
	Name      string    `json:"name" xorm:"VARCHAR(45)"`
	Age       string    `json:"age" xorm:"VARCHAR(45)"`
	CreatedAt time.Time `json:"created_at" xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
	UpdatedAt time.Time `json:"updated_at" xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}
