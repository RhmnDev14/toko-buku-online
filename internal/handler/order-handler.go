package handler

import (
	"context"
	"strconv"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/helper"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/middleware"
	"toko_buku_online/internal/usecase"
	"toko_buku_online/toko_buku_online/api/gen/go/toko/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	toko.UnimplementedOrderServiceServer
	log        logger.Logger
	uc         usecase.OrderUc
	middleware middleware.AuthMiddleware
}

func NewOrderHandler(log logger.Logger, uc usecase.OrderUc, middleware middleware.AuthMiddleware) *OrderHandler {
	return &OrderHandler{
		log:        log,
		uc:         uc,
		middleware: middleware,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *toko.OrderRequest) (*toko.OrderResponse, error) {
	h.log.Info("Create category in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.CREATE)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var items []dto.OrderItem
	for _, oi := range req.GetOrderItems() {
		items = append(items, dto.OrderItem{
			BookId:   int(oi.GetBookId()),
			Price:    oi.GetPrice(),
			Quantity: int(oi.GetQuantity()),
		})
	}

	payload := dto.Order{
		OrderItems: items,
	}

	err = h.uc.CreateOrder(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.OrderResponse{
		Message: constant.Succes,
	}, nil
}
func (h *OrderHandler) PayOrder(ctx context.Context, req *toko.EmptyOrder) (*toko.OrderResponse, error) {
	h.log.Info("Create category in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.PUT)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ids, err := helper.GetIdFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	payload, err := strconv.Atoi(ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = h.uc.PayOrder(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.OrderResponse{
		Message: constant.Succes,
	}, nil
}
