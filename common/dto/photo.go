package dto

type PhotoCardList struct {
	PhotoCounts   int64        `json:"photo_counts"`
	PhotoCardList []*PhotoCard `json:"photo_list"`
}

func NewPhotoCardList(counts int64, photoCardList []*PhotoCard) *PhotoCardList {

	return &PhotoCardList{
		PhotoCounts:   counts,
		PhotoCardList: photoCardList,
	}
}

type PhotoCard struct {
	PhotoId   string    `json:"id"`
	Photo     string    `json:"photo"`
	Likes     int       `json:"likes"`
	CreatedAt string    `json:"created_at"`
	User      PhotoUser `json:"user"`
}

func NewPhotoCard(id string, photo string, likes int, createdAt string, userId string, name string, icon string) *PhotoCard {

	return &PhotoCard{
		PhotoId:   id,
		Photo:     photo,
		Likes:     likes,
		CreatedAt: createdAt,
		User: PhotoUser{
			Id:   userId,
			Name: name,
			Icon: icon,
		},
	}
}

// 写真ページ
type PhotoPage struct {
	Id         string    `json:"id"`
	Photo      string    `json:"photo"`
	Contents   string    `json:"contents"`
	CreatedAt  string    `json:"created_at"`
	Likes      int       `json:"likes"`
	User       PhotoUser `json:"user"`
	ItemCounts int       `json:"item_counts"`
	Items      []*Item   `json:"items"`
	Tags       []*Tag    `json:"tags"`
}

// 写真詳細エンティティを生成
func NewPhotoPage(id string, photo string, contents string, createdAt string,
	likes int, userId string, userName string, icon string, selfIntroduction string, itemCounts int, items []*Item, tags []*Tag) *PhotoPage {

	return &PhotoPage{
		Id:        id,
		Photo:     photo,
		Contents:  contents,
		CreatedAt: createdAt,
		Likes:     likes,
		User: PhotoUser{
			Id:               userId,
			Name:             userName,
			Icon:             icon,
			SelfIntroduction: selfIntroduction,
		},
		ItemCounts: itemCounts,
		Items:      items,
		Tags:       tags,
	}
}

type PhotoUser struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Icon             string `json:"icon"`
	SelfIntroduction string `json:"self_introduction"`
}

// new
type InPhoto struct {
	Id       string `json:"id"`
	Photo    string `json:"photo"`
	Contents string `json:"contents"`
	UserId   string `json:"user_id"`
	Items    []Item `json:"items"`
	Tags     []Tag  `json:"tags"`
}

// アイテム
type Item struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

func NewItem(id uint, name string, price int, image string) *Item {
	return &Item{
		Id:    id,
		Name:  name,
		Price: price,
		Image: image,
	}
}

// タグ
type Tag struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func NewTag(id uint, name string) *Tag {
	return &Tag{
		Id:   id,
		Name: name,
	}
}

// 検索条件
type QueryCondition struct {
	Keyword *string `form:"keyword"`
	TagId   *string `form:"tag_id"`
	Sort    *string `form:"sort"`
	PageNo  *int    `form:"page_no"`
}
