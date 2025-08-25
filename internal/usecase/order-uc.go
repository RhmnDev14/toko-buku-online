package usecase

import (
	"context"
	"fmt"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/helper"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/repository"
)

type OrderUc interface {
	CreateOrder(ctx context.Context, payload dto.Order) error
	PayOrder(ctx context.Context, payload int) error
	GetOrders(ctx context.Context) ([]entity.Order, error)
	DetailOrder(ctx context.Context, payload int) (entity.DetailOrder, error)
}

type orderUc struct {
	log       logger.Logger
	orderRepo repository.OrderRepo
	bookRepo  repository.BookRepo
}

func NewOrderUc(log logger.Logger, orderRepo repository.OrderRepo, bookRepo repository.BookRepo) OrderUc {
	return &orderUc{
		log:       log,
		orderRepo: orderRepo,
		bookRepo:  bookRepo,
	}
}

func (u *orderUc) CreateOrder(ctx context.Context, payload dto.Order) error {
	u.log.Info("create order in uc", payload)

	tx := u.orderRepo.Begin()
	order, err := u.orderRepo.CreateOrder(ctx, tx, helper.OrderToEntity(ctx, payload))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, v := range payload.OrderItems {
		book, err := u.bookRepo.GetBookById(ctx, v.BookId)
		if err != nil {
			if err == fmt.Errorf(constant.ErrorDataNotFound) {
				return fmt.Errorf("buku dengan id %d tidak ditemukan", v.BookId)
			}
			return err
		}

		if book.Stock < v.Quantity {
			return fmt.Errorf("buku dengan id %d tidak memiliki stock yang mencukupi", v.BookId)
		}

		oi := helper.OrderItemToEntity(v)
		oi.OrderID = order.ID

		err = u.orderRepo.CreateOrderItem(ctx, tx, oi)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (u *orderUc) PayOrder(ctx context.Context, payload int) error {
	u.log.Info("pay order in uc", payload)
	return u.orderRepo.PayOrder(ctx, payload)
}

func (u *orderUc) GetOrders(ctx context.Context) ([]entity.Order, error) {
	u.log.Info("get orders in uc", ctx)
	return u.orderRepo.GetOrders(ctx)
}

func (u *orderUc) DetailOrder(ctx context.Context, payload int) (entity.DetailOrder, error) {
	u.log.Info("detail order in uc", payload)

	order, err := u.orderRepo.GetOrder(ctx, payload)
	if err != nil {
		return entity.DetailOrder{}, err
	}
	orderItems, err := u.orderRepo.GetOrderItem(ctx, int(order.ID))
	if err != nil {
		return entity.DetailOrder{}, err
	}
	return entity.DetailOrder{
		Order:     order,
		OrderItem: orderItems,
	}, nil
}
