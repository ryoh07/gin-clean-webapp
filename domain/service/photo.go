package service

// 投稿写真一覧取得関数の引数を示す構造体
// 未指定とゼロ値を区別するため、型はポインタにする
type PhotoCardOpts struct {
	UserId         *string
	LikeUserId     *string
	QueryCondition *QueryCondition
}

func NewPhotoCardOpts(userId *string, likeUserId *string, queryCondition *QueryCondition) *PhotoCardOpts {
	return &PhotoCardOpts{
		UserId:         userId,
		LikeUserId:     likeUserId,
		QueryCondition: queryCondition,
	}
}

// 投稿写真一覧の検索条件
// 未指定とゼロ値を区別するため、型はポインタにする
type QueryCondition struct {
	Keyword *string `form:"keyword"`
	TagId   *string `form:"tag_id"`
	Sort    *string `form:"sort"`
	PageNo  *int    `form:"page_no"`
}

func NewQueryCondition(keyword *string, tagId *string, sort *string, pageNo *int) *QueryCondition {
	return &QueryCondition{
		Keyword: keyword,
		TagId:   tagId,
		Sort:    sort,
		PageNo:  pageNo,
	}
}
