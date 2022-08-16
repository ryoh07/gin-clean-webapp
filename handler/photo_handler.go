package handler

import (
	"context"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryoh07/gin-clean-webapp/common/dto"
	"github.com/ryoh07/gin-clean-webapp/usecase"
)

type IPhotoHandler interface {
	// IDから投稿写真を取得
	GetPhoto(c *gin.Context)
	// 投稿写真を作成
	CreatePhoto(c *gin.Context)
	// 投稿写真を更新
	UpdatePhoto(c *gin.Context)
	// 検索条件より投稿写真一覧を取得
	GetPhotoCardList(c *gin.Context)
	// ログインユーザの投稿写真一覧を取得
	GetMyPhotoCardList(c *gin.Context)
	// ログインユーザのいいねした写真一覧を取得
	GetLikePhotoCardList(c *gin.Context)
	// 写真のいいねを作成
	CreateLike(c *gin.Context)
	// 写真のいいねを削除
	DeleteLike(c *gin.Context)
}

type PhotoHandler struct {
	usecase.PhotoInputPort
}

func NewPhotoHandler(srv usecase.PhotoInputPort) IPhotoHandler {
	return &PhotoHandler{srv}
}

func (h *PhotoHandler) GetPhoto(c *gin.Context) {
	photoId := c.Param("id")
	photo, _ := h.PhotoInputPort.GetPhoto(photoId)
	c.JSON(http.StatusOK, photo)
}

func (h *PhotoHandler) CreatePhoto(c *gin.Context) {
	ctx := context.Background()
	// リクエストにはIDがないので上手くマッピングできるかテスト
	inPhoto := dto.InPhoto{}
	if err := c.ShouldBindJSON(&inPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.PhotoInputPort.CreatePhoto(ctx, &inPhoto)
}

func (h *PhotoHandler) UpdatePhoto(c *gin.Context) {
	ctx := context.Background()
	// リクエストにはphotoがないので上手くマッピングできるかテスト
	inPhoto := dto.InPhoto{}
	if err := c.ShouldBindJSON(&inPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.PhotoInputPort.UpdatePhoto(ctx, &inPhoto)
}
func (h *PhotoHandler) GetPhotoCardList(c *gin.Context) {
	query := dto.QueryCondition{}
	c.ShouldBindQuery(&query)
	photo, err := h.PhotoInputPort.GetPhotoCardList(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, photo)
}
func (h *PhotoHandler) GetMyPhotoCardList(c *gin.Context) {
	userId := c.Param("id")
	photo, err := h.PhotoInputPort.GetMyPhotoCardList(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, photo)
}
func (h *PhotoHandler) GetLikePhotoCardList(c *gin.Context) {
	userId := c.Param("id")
	photo, err := h.PhotoInputPort.GetLikePhotoCardList(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, photo)
}
func (h *PhotoHandler) CreateLike(c *gin.Context) {
	ctx := context.Background()
	photoId := c.Param("photoid")
	userId := c.Param("userid")

	err := h.PhotoInputPort.CreateLike(ctx, photoId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
func (h *PhotoHandler) DeleteLike(c *gin.Context) {
	ctx := context.Background()
	photoId := c.Param("photoid")
	userId := c.Param("userid")

	err := h.PhotoInputPort.DeleteLike(ctx, userId, photoId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
