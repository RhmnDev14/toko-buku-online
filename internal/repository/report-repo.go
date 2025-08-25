package repository

import (
	"context"
	"fmt"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"

	"gorm.io/gorm"
)

type ReportRepo interface {
	GetSalesReport(ctx context.Context) (entity.Sales, error)
	GetBestSeller(ctx context.Context) ([]entity.BestBook, error)
	GetPriceReport(ctx context.Context) (entity.PriceBook, error)
}

type reportRepo struct {
	log logger.Logger
	db  *gorm.DB
}

func NewReportRepo(log logger.Logger, db *gorm.DB) ReportRepo {
	return &reportRepo{
		log: log,
		db:  db,
	}
}

func (r *reportRepo) GetSalesReport(ctx context.Context) (entity.Sales, error) {
	r.log.Info("Get sales report in repo", ctx)

	var sales entity.Sales
	err := r.db.WithContext(ctx).
		Table("order_items").
		Select("SUM(price * quantity) as omset, SUM(quantity) as total_buku_terjual").
		Scan(&sales).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return entity.Sales{}, fmt.Errorf(constant.ErrorServerGet)
	}

	return sales, nil
}

func (r *reportRepo) GetBestSeller(ctx context.Context) ([]entity.BestBook, error) {
	r.log.Info("get best seller in repo", ctx)

	var books []entity.BestBook
	err := r.db.WithContext(ctx).
		Table("order_items oi").
		Select("b.id, b.title").
		Joins("JOIN books b ON oi.book_id = b.id").
		Group("b.id, b.title").
		Order("SUM(oi.quantity) DESC").
		Scan(&books).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return nil, fmt.Errorf(constant.ErrorServerGet)
	}

	return books, nil
}

func (r *reportRepo) GetPriceReport(ctx context.Context) (entity.PriceBook, error) {
	r.log.Info("get price report in repo", ctx)
	var pb entity.PriceBook
	err := r.db.WithContext(ctx).
		Table("books").
		Select("MAX(price) as max, MIN(price) as min, AVG(price) as avg").
		Scan(&pb).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return entity.PriceBook{}, fmt.Errorf(constant.ErrorServerGet)
	}
	return pb, nil
}
