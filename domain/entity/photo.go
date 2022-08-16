package entity

// 投稿写真リスト取得
// 投稿写真取得
type Photo struct {
	PhotoId   string `xorm:"id"`
	UserId    string
	Photo     string
	Contents  string
	CreatedAt string `xorm:"created"`
}

type PhotoUser struct {
	UserId           string `xorm:"id"`
	Name             string
	Icon             string
	SelfIntroduction string
}

// 投稿写真と付随データのインサートエンティティ
type InPhoto struct {
	Photo *Photo
	Items []*Item
	Tags  []*Tag
}

// 投稿写真のカード
type PhotoCard struct {
	Id        string `xorm:"id"`
	Photo     string
	Likes     int
	CreatedAt string
	UserId    string
	Name      string
	Icon      string
}

// 投稿写真一覧の検索条件
// 未指定とゼロ値を区別するため、型はポインタにする
type QueryCondition struct {
	Keyword *string `form:"keyword"`
	TagId   *string `form:"tag_id"`
	Sort    *string `form:"sort"`
	PageNo  *int    `form:"page_no"`
}

// 投稿写真一覧取得関数の引数を示す構造体
// 未指定とゼロ値を区別するため、型はポインタにする
type PhotoCardOpts struct {
	UserId         *string
	LikeUserId     *string
	QueryCondition *QueryCondition
}

// いいね
type Like struct {
	Id      uint `xorm:"autoincr"`
	PhotoId string
	UserId  string
}

// アイテム
type Item struct {
	Id    uint   `xorm:"autoincr"`
	Name  string `xorm:"varchar(256)"`
	Price int    `xorm:"varchar(256)"`
	Image string `xorm:"varchar(256)"`
}

// タグ
type Tag struct {
	Id   uint   `xorm:"autoincr"`
	Name string `xorm:"varchar(256)"`
}
