package repository

import (
	"context"
	"fmt"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"

	"gorm.io/gorm"
)

type CategoryRepo interface {
	CreateCategory(ctx context.Context, payload entity.Category) error
	GetCategories(ctx context.Context) ([]entity.Category, error)
	UpdateCategory(ctx context.Context, id int, payload entity.Category) error
	DeleteCategory(ctx context.Context, payload int) error
}

type categoryRepo struct {
	log logger.Logger
	db  *gorm.DB
}

func NewCategoryRepo(log logger.Logger, db *gorm.DB) CategoryRepo {
	return &categoryRepo{
		log: log,
		db:  db,
	}
}

func (r *categoryRepo) CreateCategory(ctx context.Context, payload entity.Category) error {
	r.log.Info("Create category in repo", payload)

	err := r.db.WithContext(ctx).Create(&payload).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorServerCreate)
	}
	return nil
}
func (r *categoryRepo) GetCategories(ctx context.Context) ([]entity.Category, error) {
	r.log.Info("Get categories in repo", nil)

	var categories []entity.Category
	err := r.db.WithContext(ctx).Find(&categories).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return nil, fmt.Errorf(constant.ErrorServerGet)
	}
	return categories, nil
}
func (r *categoryRepo) UpdateCategory(ctx context.Context, id int, payload entity.Category) error {
	r.log.Info("Update category in repo", payload)

	updates := map[string]interface{}{}

	if payload.Name != "" {
		updates["name"] = payload.Name
	}

	if len(updates) == 0 {
		r.log.Info("No fields to update, skip update", payload)
		return nil
	}

	err := r.db.WithContext(ctx).
		Model(&entity.Category{}).
		Where("id = ?", id).
		Updates(updates).Error

	if err != nil {
		r.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorServerUpdate)
	}

	return nil
}

func (r *categoryRepo) DeleteCategory(ctx context.Context, payload int) error {
	r.log.Info("Update category in repo", payload)

	var category entity.Category
	err := r.db.WithContext(ctx).Where("id = ?", payload).Delete(&category).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorServerUpdate)
	}
	return nil
}
