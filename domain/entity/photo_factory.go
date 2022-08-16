package entity

import (
	"github.com/ryoh07/gin-clean-webapp/common"
)

func NewPhoto(userId string, photo string, contents string, createdAt string) *Photo {
	return &Photo{
		PhotoId:   string([]rune(common.NewUuid())[:8]),
		UserId:    userId,
		Photo:     photo,
		Contents:  contents,
		CreatedAt: createdAt,
	}
}

func NewPhotoUser(userId string, name string, icon string, selfIntroduction string) *PhotoUser {
	return &PhotoUser{
		UserId:           userId,
		Name:             name,
		Icon:             icon,
		SelfIntroduction: selfIntroduction,
	}
}

func NewInPhoto(photo *Photo, items []*Item, tags []*Tag) *InPhoto {
	return &InPhoto{
		Photo: photo,
		Items: items,
		Tags:  tags,
	}
}

func NewLike(id uint, photoId string, userId string) *Like {
	return &Like{
		Id:      id,
		PhotoId: photoId,
		UserId:  userId,
	}
}

func NewItem(id uint, name string, price int, image string) *Item {
	return &Item{
		Id:    id,
		Name:  name,
		Price: price,
		Image: image,
	}
}

func NewTag(id uint, name string) *Tag {
	return &Tag{
		Id:   id,
		Name: name,
	}
}
