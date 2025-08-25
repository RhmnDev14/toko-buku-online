package repository

import (
	"context"
	"fmt"
	"net/url"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"

	"gorm.io/gorm"
)

type BookRepo interface {
	CreateBook(ctx context.Context, payload entity.Book) error
	GetBooks(ctx context.Context, search, sortBy, sortType string, page, limit, filter int) ([]entity.Book, int64, error)
	GetBookById(ctx context.Context, payload int) (entity.Book, error)
	UpdateBook(ctx context.Context, id int, payload entity.Book) error
	DeleteBook(ctx context.Context, payload int) error
}

type bookRepo struct {
	log logger.Logger
	db  *gorm.DB
}

func NewBookRepo(log logger.Logger, db *gorm.DB) BookRepo {
	return &bookRepo{
		log: log,
		db:  db,
	}
}

func (r *bookRepo) CreateBook(ctx context.Context, payload entity.Book) error {
	r.log.Info("create book in repo", payload)

	err := r.db.WithContext(ctx).Create(&payload).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorServerCreate)
	}
	return nil
}
func (r *bookRepo) GetBooks(ctx context.Context, search, sortBy, sortType string, page, limit, filter int) ([]entity.Book, int64, error) {
	r.log.Info("get books in repo", map[string]interface{}{
		"search": search,
		"page":   page,
		"limit":  limit,
		"filter": filter,
	})

	query := r.db.WithContext(ctx).Model(&entity.Book{})

	decodename, err := url.QueryUnescape(search)
	if err != nil {
		r.log.Error("Error : ", err)
		return nil, 0, fmt.Errorf(constant.ErrorInternalSystem)
	}

	if search != "" {
		query = query.Where("title ILIKE ? OR author ILIKE ? ", "%"+decodename+"%", "%"+decodename+"%")
	}

	if filter != 0 {
		query = query.Where("category_id = ?", filter)
	}

	var totalRecords int64
	offset := (page - 1) * limit

	err = query.Count(&totalRecords).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return nil, 0, fmt.Errorf(constant.ErrorInternalSystem)
	}

	var books []entity.Book
	if sortBy == "" {
		sortBy = "id"
	}
	if sortType == "" {
		sortType = "asc"
	}

	err = query.Order(fmt.Sprintf("%s %s", sortBy, sortType)).
		Limit(limit).
		Offset(offset).
		Find(&books).Error

	if err != nil {
		r.log.Error("Error :", err)
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf(constant.ErrorDataNotFound)
		}
		return nil, 0, fmt.Errorf(constant.ErrorServerGet)
	}

	return books, totalRecords, nil
}
func (r *bookRepo) GetBookById(ctx context.Context, payload int) (entity.Book, error) {
	r.log.Info("get book by id in repo", payload)

	var book entity.Book
	err := r.db.WithContext(ctx).Where("id = ?", payload).First(&book).Error
	if err != nil {
		r.log.Error("Error : ", err)
		if err == gorm.ErrRecordNotFound {
			return book, fmt.Errorf(constant.ErrorDataNotFound)
		}
		return book, fmt.Errorf(constant.ErrorServerGet)
	}
	return book, nil
}
func (r *bookRepo) UpdateBook(ctx context.Context, id int, payload entity.Book) error {
	r.log.Info("update book in repo", id)

	err := r.db.WithContext(ctx).Where("id = ?", id).Updates(&payload).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorServerUpdate)
	}
	return nil
}
func (r *bookRepo) DeleteBook(ctx context.Context, payload int) error {
	r.log.Info("delete book in repo", payload)

	var book entity.Book
	err := r.db.WithContext(ctx).Where("id = ?", payload).Delete(&book).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorServerUpdate)
	}
	return nil
}
