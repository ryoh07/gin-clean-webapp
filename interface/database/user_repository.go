package database

import (
	"bytes"
	"context"
	"net/url"

	"github.com/go-xorm/xorm"
	"github.com/ryoh07/gin-clean-webapp/common"
	"github.com/ryoh07/gin-clean-webapp/domain/entity"
	"github.com/ryoh07/gin-clean-webapp/domain/repository"
	"github.com/ryoh07/gin-clean-webapp/interface/aws"
)

type UserRepository struct {
	*xorm.Engine
}

func NewUserRepository(db *xorm.Engine) repository.IUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUser(userId string) (*entity.User, error) {

	user := entity.User{}
	_, err := r.Table("users").
		Where("id = ?", userId).
		Get(&user)

	if err != nil {
		return nil, err
	}

	// 画像を取得しエンコード
	awsS3 := aws.NewAwsS3()
	var imgByte *bytes.Buffer
	imgByte, _ = awsS3.TestS3Downloader(user.Icon)
	user.Icon = common.Encode(imgByte)

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	// トランザクションオブジェクトをコンテキストから取得する
	session := GetTx(ctx)
	// ユーザアップデート
	_, err := session.Table("users").
		Where("id = ?", user.Id).
		Update(user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	// トランザクションオブジェクトをコンテキストから取得する
	session := GetTx(ctx)

	// 画像をS3に保存後、保存パスをモデルに格納
	var urlStr string
	var u *url.URL
	awsS3 := aws.NewAwsS3()

	urlStr, _ = awsS3.S3Uploader(common.Decode(user.Icon), common.NewUuid(), "png")
	u, _ = url.Parse(urlStr)
	user.Icon = u.Path

	// ユーザインサート
	_, err := session.Table("users").Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId string) error {
	// トランザクションオブジェクトをコンテキストから取得する
	session := GetTx(ctx)

	user := entity.User{}
	_, err := session.Table("users").
		Where("id = ?", userId).
		Delete(user)
	if err != nil {
		return err
	}
	return nil
}
