package model

import (
	"github.com/songquanpeng/go-api-starter/common"
)

func InitDB() {
	err := common.InitDB()
	if err != nil {
		panic(err)
	}
	AutoMigrate()
}

func AutoMigrate() {
	err := common.DB.AutoMigrate(
		&ParseTask{},
		&ParseResult{},
		&User{},
		&Session{},
		&CodeExchange{},
	)
	if err != nil {
		panic(err)
	}
}

func CloseDB() error {
	return common.CloseDB()
}
