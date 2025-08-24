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

type CategoryHandler struct {
	toko.UnimplementedCategoryServiceServer
	log        logger.Logger
	uc         usecase.CategoryUc
	middleware middleware.AuthMiddleware
}

func NewCategoryHandler(log logger.Logger, uc usecase.CategoryUc, middleware middleware.AuthMiddleware) *CategoryHandler {
	return &CategoryHandler{
		log:        log,
		uc:         uc,
		middleware: middleware,
	}
}

func (h *CategoryHandler) CreateCategory(ctx context.Context, req *toko.CategoryRequest) (*toko.CategoryResponse, error) {
	h.log.Info("Create category in handler", req)

	payload := dto.CategoryReq{
		Name: req.GetName(),
	}

	err := h.uc.CreateCategory(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.CategoryResponse{
		Message: constant.Succes,
	}, nil
}
func (h *CategoryHandler) GetCategories(ctx context.Context, req *toko.Empty) (*toko.CategoryResponseList, error) {
	h.log.Info("Get categories in handler", req)

	data, err := h.uc.GetCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var result []*toko.CategoryResponseData
	for _, d := range data {
		result = append(result, &toko.CategoryResponseData{
			Id:   int32(d.ID),
			Name: d.Name,
		})
	}

	return &toko.CategoryResponseList{
		Data: result,
	}, nil
}
func (h *CategoryHandler) UpdateCategory(ctx context.Context, req *toko.CategoryUpdateRequest) (*toko.CategoryResponse, error) {
	h.log.Info("Create category in handler", req)

	payload := dto.CategoryReq{
		Name: req.GetName(),
	}

	ids, err := helper.GetIdFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = h.uc.UpdateCategory(ctx, id, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.CategoryResponse{
		Message: constant.Succes,
	}, nil
}
func (h *CategoryHandler) DeleteCategory(ctx context.Context, req *toko.Empty) (*toko.CategoryResponse, error) {
	h.log.Info("Create category in handler", nil)

	ids, err := helper.GetIdFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = h.uc.DeleteCategory(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.CategoryResponse{
		Message: constant.Succes,
	}, nil
}
