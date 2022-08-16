package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/ryoh07/gin-clean-webapp/common/di"
)

func SetRoutes(engine *gin.Engine, db *xorm.Engine) {
	user := di.InitUser(db)
	engine.GET("/user/:id", user.GetUser)
	engine.POST("/user", user.CreateUser)
	engine.POST("/user/edit", user.UpdateUser)
	engine.DELETE("/user/:userid", user.DeleteUser)

	photo := di.InitPhoto(db)
	engine.GET("/photo", photo.GetPhotoCardList)
	engine.GET("/photo/:id", photo.GetPhoto)
	engine.POST("/photo", photo.CreatePhoto)
	engine.POST("/photo/edit", photo.UpdatePhoto)
	engine.GET("/photo/myphoto/:id", photo.GetMyPhotoCardList)
	engine.GET("/photo/mylike/:id", photo.GetLikePhotoCardList)
	engine.PUT("/photo/:photoid/:userid/like", photo.CreateLike)
	engine.DELETE("/photo/:photoid/:userid/like", photo.DeleteLike)
}
