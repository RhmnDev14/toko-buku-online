package usecase

import (
	"context"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/helper"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/repository"
)

type CategoryUc interface {
	CreateCategory(ctx context.Context, payload dto.CategoryReq) error
	GetCategories(ctx context.Context) ([]entity.Category, error)
	UpdateCategory(ctx context.Context, id int, payload dto.CategoryReq) error
	DeleteCategory(ctx context.Context, payload int) error
}

type categoryUc struct {
	log          logger.Logger
	categoryRepo repository.CategoryRepo
}

func NewCategoryUc(log logger.Logger, categoryRepo repository.CategoryRepo) CategoryUc {
	return &categoryUc{
		log:          log,
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUc) CreateCategory(ctx context.Context, payload dto.CategoryReq) error {
	u.log.Info("Create category in uc", payload)
	return u.categoryRepo.CreateCategory(ctx, helper.CategoryToEntity(payload))
}

func (u *categoryUc) GetCategories(ctx context.Context) ([]entity.Category, error) {
	u.log.Info("get categories in uc", nil)
	return u.categoryRepo.GetCategories(ctx)
}

func (u *categoryUc) UpdateCategory(ctx context.Context, id int, payload dto.CategoryReq) error {
	u.log.Info("update category in uc", id)
	return u.categoryRepo.UpdateCategory(ctx, id, helper.CategoryToEntity(payload))
}

func (u *categoryUc) DeleteCategory(ctx context.Context, payload int) error {
	u.log.Info("Delete category in uc", payload)
	return u.categoryRepo.DeleteCategory(ctx, payload)
}
