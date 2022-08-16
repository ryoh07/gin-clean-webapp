package database

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/go-xorm/xorm"
	"github.com/ryoh07/gin-clean-webapp/domain/entity"
	"github.com/ryoh07/gin-clean-webapp/domain/repository"
	"github.com/ryoh07/gin-clean-webapp/domain/service"
	"github.com/ryoh07/gin-clean-webapp/interface/aws"

	"github.com/ryoh07/gin-clean-webapp/common"

	"xorm.io/builder"
)

type PhotoRepository struct {
	*xorm.Engine
}

func NewPhotoRepository(db *xorm.Engine) repository.IPhotoRepository {
	return &PhotoRepository{db}
}

func (r *PhotoRepository) GetPhoto(photoId string) (*entity.Photo, error) {
	photo := entity.Photo{}
	_, err := r.Table("photos").
		Where("id = ?", photoId).
		Get(&photo)

	if err != nil {
		return nil, err
	}
	// 画像ダウンロードし、bace64にエンコード
	// エンコードの値をエンティティに格納する。
	awsS3 := aws.NewAwsS3()
	var imgByte *bytes.Buffer
	imgByte, _ = awsS3.TestS3Downloader(photo.Photo)
	photo.Photo = common.Encode(imgByte)

	return &photo, nil
}

func (r *PhotoRepository) GetPhotoUser(photoId string) (*entity.PhotoUser, error) {

	photoUser := entity.PhotoUser{}

	query := builder.
		Select(" u.id,u.name,u.icon,u.self_introduction").
		From("users u").
		Join("INNER", "photos p", "u.id = p.user_id").
		Where(builder.Eq{"p.id": photoId})

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}
	engine := r.SQL(sql, args...)

	_, err = engine.Get(&photoUser)

	if err != nil {
		return nil, err
	}
	// 画像ダウンロードし、bace64にエンコード
	// エンコードの値をエンティティに格納する。
	awsS3 := aws.NewAwsS3()
	var imgByte *bytes.Buffer
	imgByte, _ = awsS3.TestS3Downloader(photoUser.Icon)
	photoUser.Icon = common.Encode(imgByte)

	return &photoUser, nil
}

// 投稿写真一覧を取得
// ユーザが投稿した一覧
// いいねをした一覧
func (r *PhotoRepository) FindPhotoCard(opts *service.PhotoCardOpts) ([]*entity.PhotoCard, int64, error) {
	subQuery, _, err := builder.
		Select("photo_id,user_id,count(*) as likes").
		From("photo_likes").
		GroupBy("photo_id").
		ToSQL()
	if err != nil {
		return nil, 0, err
	}
	query := builder.Dialect("MYSQL").
		Select("Distinct SQL_CALC_FOUND_ROWS p.id,p.photo,ifnull(l.likes,0) as likes,p.created_at,u.id as user_id,u.name,u.icon").
		From("photos p").
		Join("INNER", "users u", "p.user_id = u.id").
		Join("INNER", "photo_tags pt", "p.id = pt.photo_id").
		Join("LEFT", fmt.Sprintf("(%s) AS l", subQuery), "p.id = l.photo_id")

	// 指定ユーザの投稿写真一覧を取得
	if opts.UserId != nil {
		query = query.Where(builder.Eq{"u.id": *opts.UserId})
	}

	// いいねした投稿写真一覧を取得
	if opts.LikeUserId != nil {
		query = query.Where(builder.Eq{"l.user_id": *opts.LikeUserId})
	}

	// クエリパラメータ条件
	if opts.QueryCondition != nil {
		// キーワード検索
		if opts.QueryCondition.Keyword != nil {
			query = query.Where(builder.Like{"p.contents", *opts.QueryCondition.Keyword})
		}
		// タグ検索
		if opts.QueryCondition.TagId != nil {
			query = query.Where(builder.Eq{"t.id": *opts.QueryCondition.TagId})
		}
		// ソート条件
		if opts.QueryCondition.Sort != nil {
			if *opts.QueryCondition.Sort == "like" {
				query = query.OrderBy("l.likes DESC")
			}
		} else {
			query = query.OrderBy("p.created_at DESC")
		}
		// ページネーション
		if opts.QueryCondition.PageNo != nil {
			query.Limit(60, (*opts.QueryCondition.PageNo-1)*60)
		} else {
			query = query.Limit(60, 0)
		}
	} else {
		query = query.OrderBy("p.created_at DESC")
		query = query.Limit(60, 0)
	}

	// 検索結果取得
	sql, args, err := query.ToSQL()
	engine := r.SQL(sql, args...)
	if err != nil {
		return nil, 0, err
	}
	PhotoCards := []*entity.PhotoCard{}
	err = engine.Where("user_id = ?", "a").Find(&PhotoCards)
	if err != nil {
		return nil, 0, err
	}

	// 件数取得
	var count int64
	cntSql, args, err := builder.
		Select("FOUND_ROWS()").From("DUAL").ToSQL()
	if err != nil {
		return nil, 0, err
	}
	_, err = r.SQL(cntSql, args...).Get(&count)
	if err != nil {
		return nil, 0, err
	}

	// 画像を取得しエンコード
	awsS3 := aws.NewAwsS3()
	var imgByte *bytes.Buffer
	for _, v := range PhotoCards {
		imgByte, _ = awsS3.TestS3Downloader(v.Photo)
		v.Photo = common.Encode(imgByte)
		imgByte, _ = awsS3.TestS3Downloader(v.Icon)
		v.Icon = common.Encode(imgByte)
	}
	return PhotoCards, count, nil
}

// 投稿写真と付随データのインサート
func (r *PhotoRepository) CreatePhoto(ctx context.Context, inPhoto *entity.InPhoto) error {

	// トランザクションオブジェクトをコンテキストから取得する
	session := GetTx(ctx)

	// 画像をS3に保存後、保存パスをモデルに格納
	var urlStr string
	var u *url.URL
	awsS3 := aws.NewAwsS3()

	urlStr, _ = awsS3.S3Uploader(common.Decode(inPhoto.Photo.Photo), common.NewUuid(), "png")
	u, _ = url.Parse(urlStr)
	inPhoto.Photo.Photo = u.Path

	for _, v := range inPhoto.Items {
		urlStr, _ = awsS3.S3Uploader(common.Decode(v.Image), common.NewUuid(), "png")
		u, _ := url.Parse(urlStr)
		v.Image = u.Path
	}

	// 投稿写真インサート
	_, err := session.Table("photos").Insert(inPhoto.Photo)
	if err != nil {
		return err
	}

	// アイテムインサート
	for _, item := range inPhoto.Items {
		_, err = session.Table("items").Insert(item)
		if err != nil {
			return err
		}

		// アイテムIDとフォトIDを関連付ける
		_, err = session.Query(fmt.Sprintf(
			"INSERT INTO photo_items (item_id, photo_id) VALUES (%d ,'%s')",
			item.Id,
			inPhoto.Photo.PhotoId,
		))
		if err != nil {
			return err
		}
	}

	for _, tag := range inPhoto.Tags {

		// タグが存在していなければインサート
		exist, err := session.Table("tags").
			Where("name = ?", tag.Name).
			Exist()
		if err != nil {
			return err
		}
		if !exist {
			_, err = session.Table("tags").Insert(tag)
			if err != nil {
				return err
			}
		}

		// タグIDとフォトIDを関連付ける
		_, err = session.Query(fmt.Sprintf(
			"INSERT INTO photo_tags (tag_id, photo_id) "+
				"SELECT id , '%s' "+
				"FROM tags "+
				"WHERE name = '%s' ",
			inPhoto.Photo.PhotoId,
			tag.Name))
		if err != nil {
			return err
		}
	}
	return nil
}

// 投稿写真の更新
func (r *PhotoRepository) UpdatePhoto(ctx context.Context, inPhoto *entity.InPhoto) error {

	// トランザクションオブジェクトをコンテキストから取得する
	session := GetTx(ctx)

	_, err := session.Table("photos").
		Where("id = ?", inPhoto.Photo.PhotoId).
		Update(inPhoto.Photo)
	if err != nil {
		return err
	}

	// 使用アイテム更新
	for _, item := range inPhoto.Items {
		// IDがゼロ値の時はインサート
		// ゼロ値以外の場合は更新する
		if item.Id == 0 {
			_, err = session.Table("items").Insert(item)
			if err != nil {
				return err
			}

			// アイテムIDとフォトIDを関連付ける
			_, err = session.Query(fmt.Sprintf(
				"INSERT INTO photo_items (item_id, photo_id) VALUES (%d , '%s')",
				item.Id,
				inPhoto.Photo.PhotoId,
			))
			if err != nil {
				return err
			}

		} else if item.Id != 0 {
			_, err = session.Table("items").
				Where("id = ?", item.Id).
				Update(item)
			if err != nil {
				return err
			}
		}
	}

	// タグ更新
	for _, tag := range inPhoto.Tags {

		// タグの重複チェック
		exist, err := session.Table("tags").
			Where("name = ?", tag.Name).
			Exist()
		if err != nil {
			return err
		}

		if !exist {
			// タグインサート
			_, err = session.Table("tags").Insert(tag)
			if err != nil {
				return err
			}
		}

		// タグIDとフォトIDを関連付ける
		// すでに関連付いてる場合は無視する
		_, err = session.Query(fmt.Sprintf(
			"INSERT INTO photo_tags (tag_id, photo_id) "+
				"SELECT id , %s "+
				"FROM tags "+
				"WHERE name = '%s' AND "+
				"NOT EXISTS( "+
				"SELECT 'X' "+
				"FROM photo_tags pt "+
				"INNER JOIN tags t "+
				"ON pt.tag_id = t.id "+
				"WHERE pt.photo_id = '%s' AND t.name = '%s') ",
			inPhoto.Photo.PhotoId,
			tag.Name,
			inPhoto.Photo.PhotoId,
			tag.Name))

		if err != nil {
			return err
		}
	}
	return nil
}

// いいねインサート
func (r *PhotoRepository) CreateLike(ctx context.Context, photoId string, userId string) error {

	// トランザクションオブジェクトをコンテキストから取得する
	session := GetTx(ctx)

	like := entity.Like{
		PhotoId: photoId,
		UserId:  userId,
	}
	_, err := session.Table("photo_likes").Insert(like)
	if err != nil {
		return err
	}
	return nil
}

// いいね削除
func (r *PhotoRepository) DeleteLike(ctx context.Context, photoId string, userId string) error {

	// トランザクションオブジェクトをコンテキストから取得する
	session := GetTx(ctx)

	like := entity.Like{}
	_, err := session.Table("photo_likes").
		Where("user_id = ?", userId).
		Where("photo_id = ?", photoId).
		Delete(like)
	if err != nil {
		return err
	}
	return nil
}

// アイテム取得
func (r *PhotoRepository) FindItem(photoId string, userId string) ([]*entity.Item, error) {

	query := builder.
		Select("i.id,i.name,i.price,i.image").
		From("items i").
		Join("INNER", "photo_items pi", "i.id = pi.item_id").
		Join("INNER", "photos p", "p.id = pi.photo_id")

	// 写真が持つアイテム取得
	if photoId != "" {
		query.Where(builder.Eq{"pi.photo_id": photoId})
	}
	// ユーザが投稿したアイテム取得
	if userId != "" {
		query.Where(builder.Eq{"p.user_id": userId})
	}

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}
	items := []*entity.Item{}
	r.SQL(sql, args...).Find(&items)

	if err != nil {
		return nil, err
	}
	// 画像を取得しエンコード
	awsS3 := aws.NewAwsS3()
	var imgByte *bytes.Buffer
	for _, v := range items {
		imgByte, _ = awsS3.TestS3Downloader(v.Image)
		v.Image = common.Encode(imgByte)
	}

	return items, nil
}

// タグ取得
func (r *PhotoRepository) FindTag(photoId string) ([]*entity.Tag, error) {

	query := builder.
		Select("t.id,t.name").
		From("tags t").
		Join("INNER", "photo_tags pt", "t.id = pt.tag_id").
		Where(builder.Eq{"pt.photo_id": photoId})

	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}
	tags := []*entity.Tag{}
	r.SQL(sql, args...).Find(&tags)

	if err != nil {
		return nil, err
	}
	return tags, nil
}

// いいね数取得
func (r *PhotoRepository) GetLikeCounts(photoId string) (int, error) {

	subQuery, _, err := builder.
		Select("photo_id,count(*) as likes").
		From("photo_likes").
		GroupBy("photo_id").
		ToSQL()
	if err != nil {
		return 0, err
	}
	query := builder.
		Select("ifnull(pl.likes,0) as likes").
		From("photos p").
		Join("INNER", fmt.Sprintf("(%s) AS pl", subQuery), "p.id = pl.photo_id").
		Where(builder.Eq{"pl.photo_id": photoId})

	sql, args, err := query.ToSQL()
	if err != nil {
		return 0, err
	}
	var count int
	r.SQL(sql, args...).Get(&count)

	if err != nil {
		return 0, err
	}
	return count, nil

}
