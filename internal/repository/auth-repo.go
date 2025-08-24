package repository

import (
	"context"
	"fmt"
	"strings"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"

	"gorm.io/gorm"
)

type AuthRepo interface {
	Login(ctx context.Context, payload string) (entity.User, error)
	Register(ctx context.Context, user entity.User) error
}

type authRepo struct {
	log logger.Logger
	db  *gorm.DB
}

func NewAuthRepo(log logger.Logger, db *gorm.DB) AuthRepo {
	return &authRepo{
		log: log,
		db:  db,
	}
}

func (r *authRepo) Login(ctx context.Context, payload string) (entity.User, error) {
	r.log.Info("login in repo", payload)

	var user entity.User
	err := r.db.WithContext(ctx).
		Where("email = ? ", payload).
		First(&user).Error
	if err != nil {
		r.log.Error("Error : ", err)
		if err == gorm.ErrRecordNotFound {
			return entity.User{}, fmt.Errorf(constant.ErrorLogin)
		}
		return entity.User{}, fmt.Errorf(constant.ErrorServerGet)
	}

	return user, nil
}

func (r *authRepo) Register(ctx context.Context, user entity.User) error {
	r.log.Info("register in repository", user)

	err := r.db.WithContext(ctx).
		Create(&user).Error
	if err != nil {
		r.log.Error("Error : ", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") &&
			strings.Contains(err.Error(), "idx_users_email") {
			return fmt.Errorf(constant.ErrorEmailHasBeenUsed)
		}
		return fmt.Errorf(constant.ErrorServerCreate)
	}

	return nil
}
