package usecase

import (
	"context"

	"github.com/ryoh07/gin-clean-webapp/common/dto"
	"github.com/ryoh07/gin-clean-webapp/domain/entity"
	"github.com/ryoh07/gin-clean-webapp/domain/repository"
	"github.com/ryoh07/gin-clean-webapp/domain/service"
	"github.com/ryoh07/gin-clean-webapp/transaction"
)

type PhotoInputPort interface {
	GetPhoto(photoId string) (*dto.PhotoPage, error)
	CreatePhoto(ctx context.Context, inphoto *dto.InPhoto) error
	UpdatePhoto(ctx context.Context, inphoto *dto.InPhoto) error
	GetPhotoCardList(queryCondition *dto.QueryCondition) (*dto.PhotoCardList, error)
	GetMyPhotoCardList(userId string) (*dto.PhotoCardList, error)
	GetLikePhotoCardList(userId string) (*dto.PhotoCardList, error)
	CreateLike(ctx context.Context, photoId string, userId string) error
	DeleteLike(ctx context.Context, photoId string, userId string) error
}

type photoInteractor struct {
	repository.IPhotoRepository
	transaction.Transaction
}

func NewPhotoInteractor(repo repository.IPhotoRepository, tx transaction.Transaction) PhotoInputPort {
	return &photoInteractor{repo, tx}
}

// 写真詳細ページの取得
func (s *photoInteractor) GetPhoto(photoId string) (*dto.PhotoPage, error) {

	photo, err := s.IPhotoRepository.GetPhoto(photoId)
	if err != nil {
		return nil, err
	}
	// 写真ユーザを取得
	photoUser, err := s.IPhotoRepository.GetPhotoUser(photoId)
	if err != nil {
		return nil, err
	}
	// いいね数取得
	likes, err := s.IPhotoRepository.GetLikeCounts(photoId)
	if err != nil {
		return nil, err
	}
	// アイテム取得しDTO変換
	items, err := s.IPhotoRepository.FindItem(photoId, "")
	if err != nil {
		return nil, err
	}
	itemsDto := make([]*dto.Item, 0)
	for _, v := range items {
		itemsDto = append(itemsDto, dto.NewItem(v.Id, v.Name, v.Price, v.Image))
	}
	// アイテム数取得
	itemCounts := (len(itemsDto))

	// タグ取得しDTO変換
	tags, err := s.IPhotoRepository.FindTag(photoId)
	if err != nil {
		return nil, err
	}
	tagsDto := make([]*dto.Tag, 0)
	for _, v := range tags {
		tagsDto = append(tagsDto, dto.NewTag(v.Id, v.Name))
	}

	// 取得したエンティティをDTOに変換
	photoPage := dto.NewPhotoPage(
		photo.PhotoId, photo.Photo, photo.Contents, photo.CreatedAt, likes,
		photoUser.UserId, photoUser.Name, photoUser.Icon, photoUser.SelfIntroduction,
		itemCounts, itemsDto, tagsDto)

	return photoPage, nil
}

func (s *photoInteractor) CreatePhoto(ctx context.Context, inPhoto *dto.InPhoto) error {

	photo := entity.NewPhoto(inPhoto.UserId, inPhoto.Photo, inPhoto.Contents, "")
	items := make([]*entity.Item, 0)
	for _, v := range inPhoto.Items {
		items = append(items, entity.NewItem(v.Id, v.Name, v.Price, v.Image))
	}
	tags := make([]*entity.Tag, 0)
	for _, v := range inPhoto.Tags {
		tags = append(tags, entity.NewTag(v.Id, v.Name))
	}
	err := s.Transaction.DoInTx(ctx, func(ctx context.Context) error {
		return s.IPhotoRepository.CreatePhoto(ctx, &entity.InPhoto{
			Photo: photo,
			Items: items,
			Tags:  tags,
		})
	})
	return err
}
func (s *photoInteractor) UpdatePhoto(ctx context.Context, inPhoto *dto.InPhoto) error {
	photo := entity.NewPhoto(inPhoto.UserId, inPhoto.Photo, inPhoto.Contents, "")
	items := make([]*entity.Item, 0)
	for _, v := range inPhoto.Items {
		items = append(items, entity.NewItem(v.Id, v.Name, v.Price, v.Image))
	}
	tags := make([]*entity.Tag, 0)
	for _, v := range inPhoto.Tags {
		tags = append(tags, entity.NewTag(v.Id, v.Name))
	}
	err := s.Transaction.DoInTx(ctx, func(ctx context.Context) error {
		return s.IPhotoRepository.UpdatePhoto(ctx, &entity.InPhoto{
			Photo: photo,
			Items: items,
			Tags:  tags,
		})
	})
	return err
}
func (s *photoInteractor) GetPhotoCardList(dtoQC *dto.QueryCondition) (*dto.PhotoCardList, error) {

	serviceQC := service.NewQueryCondition(dtoQC.Keyword, dtoQC.TagId,
		dtoQC.Sort, dtoQC.PageNo)
	photoCardListE, photoCounts, err := s.IPhotoRepository.FindPhotoCard(&service.PhotoCardOpts{
		QueryCondition: serviceQC,
	})
	if err != nil {
		return nil, err
	}
	// 取得結果をDTOに変換
	PhotoCardListD := make([]*dto.PhotoCard, 0)
	for _, v := range photoCardListE {
		PhotoCardListD = append(PhotoCardListD, dto.NewPhotoCard(v.Id, v.Photo, v.Likes, v.CreatedAt, v.UserId, v.Name, v.Icon))
	}

	return dto.NewPhotoCardList(photoCounts, PhotoCardListD), nil
}

// マイページ-投稿写真
func (s *photoInteractor) GetMyPhotoCardList(userId string) (*dto.PhotoCardList, error) {
	photoCardListE, photoCounts, err := s.IPhotoRepository.FindPhotoCard(&service.PhotoCardOpts{
		UserId: &userId,
	})
	if err != nil {
		return nil, err
	}
	// 取得結果をDTOに変換
	PhotoCardListD := make([]*dto.PhotoCard, 0)
	for _, v := range photoCardListE {
		PhotoCardListD = append(PhotoCardListD, dto.NewPhotoCard(v.Id, v.Photo, v.Likes, v.CreatedAt, v.UserId, v.Name, v.Icon))
	}
	return dto.NewPhotoCardList(photoCounts, PhotoCardListD), nil
}

// マイページ-いいねした写真一覧
func (s *photoInteractor) GetLikePhotoCardList(userId string) (*dto.PhotoCardList, error) {
	photoCardListE, photoCounts, err := s.IPhotoRepository.FindPhotoCard(&service.PhotoCardOpts{
		LikeUserId: &userId,
	})
	if err != nil {
		return nil, err
	}
	// 取得結果をDTOに変換
	PhotoCardListD := make([]*dto.PhotoCard, 0)
	for _, v := range photoCardListE {
		PhotoCardListD = append(PhotoCardListD, dto.NewPhotoCard(v.Id, v.Photo, v.Likes, v.CreatedAt, v.UserId, v.Name, v.Icon))
	}

	return dto.NewPhotoCardList(photoCounts, PhotoCardListD), nil
}

// いいね付与
func (s *photoInteractor) CreateLike(ctx context.Context, photoId string, userId string) error {
	err := s.Transaction.DoInTx(ctx, func(ctx context.Context) error {
		return s.IPhotoRepository.CreateLike(ctx, photoId, userId)
	})
	return err
}

// いいね削除
func (s *photoInteractor) DeleteLike(ctx context.Context, photoId string, userId string) error {
	err := s.Transaction.DoInTx(ctx, func(ctx context.Context) error {
		return s.IPhotoRepository.DeleteLike(ctx, photoId, userId)
	})
	return err
}
