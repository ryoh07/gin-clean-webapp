package di

import (
	"github.com/go-xorm/xorm"
	"github.com/ryoh07/gin-clean-webapp/handler"
	"github.com/ryoh07/gin-clean-webapp/interface/database"
	"github.com/ryoh07/gin-clean-webapp/usecase"
)

func InitUser(db *xorm.Engine) handler.IUserHandler {
	tx := database.NewTransaction(db)
	r := database.NewUserRepository(db)
	s := usecase.NewUserInteractor(r, tx)
	return handler.NewUserHandler(s)
}

func InitPhoto(db *xorm.Engine) handler.IPhotoHandler {
	tx := database.NewTransaction(db)
	r := database.NewPhotoRepository(db)
	s := usecase.NewPhotoInteractor(r, tx)
	return handler.NewPhotoHandler(s)
}
