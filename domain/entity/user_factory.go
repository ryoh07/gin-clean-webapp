package entity

import "github.com/ryoh07/gin-clean-webapp/common"

func NewUser(id string, name string, icon string, selfIntroduction string) *User {
	return &User{
		Id:               id,
		Name:             name,
		Icon:             icon,
		SelfIntroduction: selfIntroduction,
	}
}

func NewInUser(name string, icon string, selfIntroduction string) *User {
	return &User{
		Id:               string([]rune(common.NewUuid())[:8]),
		Name:             name,
		Icon:             icon,
		SelfIntroduction: selfIntroduction,
	}
}
