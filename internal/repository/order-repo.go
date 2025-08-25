package repository

import (
	"context"
	"fmt"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"

	"gorm.io/gorm"
)

type OrderRepo interface {
	Begin() *gorm.DB
	CreateOrder(ctx context.Context, tx *gorm.DB, payload entity.Order) (entity.Order, error)
	CreateOrderItem(ctx context.Context, tx *gorm.DB, payload entity.OrderItem) error
	PayOrder(ctx context.Context, payload int) error
}

type orderRepo struct {
	log logger.Logger
	db  *gorm.DB
}

func NewOrderRepo(log logger.Logger, db *gorm.DB) OrderRepo {
	return &orderRepo{
		log: log,
		db:  db,
	}
}

func (r *orderRepo) Begin() *gorm.DB {
	r.log.Info("Begin transaction", nil)
	return r.db.Begin()
}

func (r *orderRepo) CreateOrder(ctx context.Context, tx *gorm.DB, payload entity.Order) (entity.Order, error) {
	r.log.Info("create order in repo", payload)

	conn := tx
	if conn == nil {
		conn = r.db
	}

	err := conn.WithContext(ctx).Create(&payload).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return entity.Order{}, fmt.Errorf(constant.ErrorServerCreate)
	}
	return payload, nil
}

func (r *orderRepo) CreateOrderItem(ctx context.Context, tx *gorm.DB, payload entity.OrderItem) error {
	r.log.Info("create order item in repo", payload)

	conn := tx
	if conn == nil {
		conn = r.db
	}

	err := conn.WithContext(ctx).Create(&payload).Error
	if err != nil {
		r.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorServerCreate)
	}
	return nil
}

func (r *orderRepo) PayOrder(ctx context.Context, payload int) error {
	r.log.Info("pay order in repo", payload)

	result := r.db.WithContext(ctx).
		Model(&entity.Order{}).
		Where("id = ?", payload).
		Updates(map[string]interface{}{
			"status": constant.PAID,
		})

	if result.Error != nil {
		r.log.Error("Error : ", result.Error)
		return fmt.Errorf(constant.ErrorServerUpdate)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf(constant.ErrorDataNotFound)
	}

	return nil
}
