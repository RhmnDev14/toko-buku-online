package usecase

import (
	"context"
	"fmt"
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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var totals []float64
	bookPrices := make(map[int]float64)

	for _, v := range payload.OrderItems {
		price, err := u.bookRepo.GetPrice(ctx, tx, v.BookId)
		if err != nil {
			tx.Rollback()
			return err
		}
		bookPrices[v.BookId] = price

		total := price * float64(v.Quantity)
		totals = append(totals, total)
	}

	totalAmount := helper.TotalAmount(totals)

	order, err := u.orderRepo.CreateOrder(ctx, tx, helper.OrderToEntity(ctx, payload, totalAmount))
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, v := range payload.OrderItems {
		book, err := u.bookRepo.GetBookById(ctx, v.BookId)
		if err != nil {
			tx.Rollback()
			return err
		}

		if book.Stock < v.Quantity {
			tx.Rollback()
			return fmt.Errorf("buku dengan id %d tidak memiliki stock yang mencukupi", v.BookId)
		}

		price := bookPrices[v.BookId]

		oi := helper.OrderItemToEntity(v)
		oi.OrderID = order.ID
		oi.Price = price * float64(v.Quantity)

		err = u.orderRepo.CreateOrderItem(ctx, tx, oi)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

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
