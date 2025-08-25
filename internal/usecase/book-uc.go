package usecase

import (
	"context"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/helper"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/repository"
)

type BookUc interface {
	CreateBook(ctx context.Context, payload dto.BookReq) error
	GetBooks(ctx context.Context, param entity.Param) ([]entity.Book, int64, error)
	GetBookById(ctx context.Context, payload int) (entity.Book, error)
	UpdateBook(ctx context.Context, id int, payload dto.BookReq) error
	DeleteBook(ctx context.Context, payload int) error
}

type bookUc struct {
	log      logger.Logger
	bookRepo repository.BookRepo
}

func NewBookUc(log logger.Logger, bookRepo repository.BookRepo) BookUc {
	return &bookUc{
		log:      log,
		bookRepo: bookRepo,
	}
}

func (u *bookUc) CreateBook(ctx context.Context, payload dto.BookReq) error {
	u.log.Info("create book in uc", payload)
	return u.bookRepo.CreateBook(ctx, helper.BookToEntity(payload))
}
func (u *bookUc) GetBooks(ctx context.Context, param entity.Param) ([]entity.Book, int64, error) {
	u.log.Info("get books in uc", nil)
	return u.bookRepo.GetBooks(ctx, param.Search, param.SortBy, param.SortType, param.Page, param.Limit, param.Filter)
}
func (u *bookUc) GetBookById(ctx context.Context, payload int) (entity.Book, error) {
	u.log.Info("get book by id in uc", payload)
	return u.bookRepo.GetBookById(ctx, payload)
}
func (u *bookUc) UpdateBook(ctx context.Context, id int, payload dto.BookReq) error {
	u.log.Info("update book in uc", id)
	return u.bookRepo.UpdateBook(ctx, id, helper.BookToEntity(payload))
}
func (u *bookUc) DeleteBook(ctx context.Context, payload int) error {
	u.log.Info("delete book in uc", payload)
	return u.bookRepo.DeleteBook(ctx, payload)
}
