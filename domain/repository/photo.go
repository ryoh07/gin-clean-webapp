package repository

import (
	"context"

	"github.com/ryoh07/gin-clean-webapp/domain/entity"
	"github.com/ryoh07/gin-clean-webapp/domain/service"
)

type IPhotoRepository interface {
	GetPhoto(photoId string) (*entity.Photo, error)
	GetPhotoUser(photoId string) (*entity.PhotoUser, error)
	CreatePhoto(ctx context.Context, inphoto *entity.InPhoto) error
	UpdatePhoto(ctx context.Context, inphoto *entity.InPhoto) error
	//	DeletePhoto(photo *entity.Photo) error
	FindPhotoCard(opts *service.PhotoCardOpts) ([]*entity.PhotoCard, int64, error)

	FindItem(photoId string, userId string) ([]*entity.Item, error)
	// DeleteItem(id uint) error

	FindTag(photoId string) ([]*entity.Tag, error)

	CreateLike(ctx context.Context, photoId string, userId string) error
	DeleteLike(ctx context.Context, photoId string, userId string) error
	GetLikeCounts(photoId string) (int, error)
}
