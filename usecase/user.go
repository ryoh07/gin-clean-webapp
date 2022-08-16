package usecase

import (
	"context"

	"github.com/ryoh07/gin-clean-webapp/common/dto"
	"github.com/ryoh07/gin-clean-webapp/domain/entity"
	"github.com/ryoh07/gin-clean-webapp/domain/repository"
	"github.com/ryoh07/gin-clean-webapp/transaction"
)

type UserInputPort interface {
	GetUser(userId string) (*dto.User, error)
	UpdateUser(ctx context.Context, user *dto.User) error
	CreateUser(ctx context.Context, user *dto.User) error
	DeleteUser(ctx context.Context, userId string) error
}

type userInteractor struct {
	repository.IUserRepository
	transaction.Transaction
}

func NewUserInteractor(repo repository.IUserRepository, tx transaction.Transaction) UserInputPort {
	return &userInteractor{repo, tx}
}

// ユーザの取得
func (s *userInteractor) GetUser(userId string) (*dto.User, error) {

	userE, err := s.IUserRepository.GetUser(userId)
	if err != nil {
		return nil, err
	}

	// 取得したエンティティをDTOに変換
	userD := dto.NewUser(userE.Id, userE.Name, userE.Icon, userE.SelfIntroduction)

	return userD, nil
}

func (s *userInteractor) UpdateUser(ctx context.Context, user *dto.User) error {
	userE := entity.NewUser(user.Id, user.Name, user.Icon, user.SelfIntroduction)

	err := s.Transaction.DoInTx(ctx, func(ctx context.Context) error {
		return s.IUserRepository.UpdateUser(ctx, userE)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *userInteractor) CreateUser(ctx context.Context, user *dto.User) error {
	userE := entity.NewInUser(user.Name, user.Icon, user.SelfIntroduction)

	err := s.Transaction.DoInTx(ctx, func(ctx context.Context) error {
		return s.IUserRepository.CreateUser(ctx, userE)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *userInteractor) DeleteUser(ctx context.Context, userId string) error {
	err := s.Transaction.DoInTx(ctx, func(ctx context.Context) error {
		return s.IUserRepository.DeleteUser(ctx, userId)
	})
	return err
}
