package entity

type User struct {
	Id               string `xorm:"id"`
	Name             string
	Icon             string
	SelfIntroduction string
}
