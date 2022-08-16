package dto

type User struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	SelfIntroduction string `json:"self_introduction"`
}

func NewUser(id string, name string, icon string, selfIntroduction string) *User {
	return &User{
		Id:               id,
		Name:             name,
		Icon:             icon,
		SelfIntroduction: selfIntroduction,
	}
}
