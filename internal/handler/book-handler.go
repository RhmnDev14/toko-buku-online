package handler

import (
	"context"
	"strconv"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/helper"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/middleware"
	"toko_buku_online/internal/usecase"
	"toko_buku_online/toko_buku_online/api/gen/go/toko/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookHanlder struct {
	toko.UnimplementedBookServiceServer
	log        logger.Logger
	uc         usecase.BookUc
	middleware middleware.AuthMiddleware
}

func NewBookHandler(log logger.Logger, uc usecase.BookUc, middleware middleware.AuthMiddleware) *BookHanlder {
	return &BookHanlder{
		log:        log,
		uc:         uc,
		middleware: middleware,
	}
}

func (h *BookHanlder) CreateBook(ctx context.Context, req *toko.BookRequest) (*toko.BookResponse, error) {
	h.log.Info("Create book in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.CREATE)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	payload := dto.BookReq{
		Title:       req.GetTitle(),
		Author:      req.GetAuthor(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
		Year:        int(req.GetYear()),
		CategoryId:  int(req.GetCategoryId()),
		ImageBase64: req.GetImageBase64(),
	}

	err = h.uc.CreateBook(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.BookResponse{
		Message: constant.Succes,
	}, nil
}

func (h *BookHanlder) GetBooks(ctx context.Context, req *toko.EmptyBook) (*toko.BookResponseList, error) {
	h.log.Info("Get categories in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.GETBOOK)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	src, sortBy, sortType, filter, page, limit, err := helper.GetParamFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if filter == "" {
		filter = "0"
	}
	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}
	fltr, err := strconv.Atoi(filter)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	pg, err := strconv.Atoi(page)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	lm, err := strconv.Atoi(limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	param := entity.Param{
		Search:   src,
		Filter:   fltr,
		SortBy:   sortBy,
		SortType: sortType,
		Page:     pg,
		Limit:    lm,
	}

	data, total, err := h.uc.GetBooks(ctx, param)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var result []*toko.BookResponseData
	for _, d := range data {
		result = append(result, &toko.BookResponseData{
			Id:          int32(d.ID),
			Title:       d.Title,
			Author:      d.Author,
			Price:       d.Price,
			Year:        int32(d.Year),
			Stock:       int32(d.Stock),
			CategoryId:  int32(d.CategoryID),
			ImageBase64: d.ImageBase64,
		})
	}

	totalPages := (total + int64(lm) - 1) / int64(lm)
	metaDataBooks := &toko.MetaData{
		TotalData: int32(total),
		TotalPage: int32(totalPages),
		Page:      int32(pg),
		Limit:     int32(lm),
	}
	return &toko.BookResponseList{
		Data: result,
		Meta: metaDataBooks,
	}, nil
}

func (h *BookHanlder) GetBookById(ctx context.Context, req *toko.EmptyBook) (*toko.BookResponseData, error) {
	h.log.Info("Create book in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.GETBOOK)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ids, err := helper.GetIdFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	book, err := h.uc.GetBookById(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.BookResponseData{
		Id:          int32(book.ID),
		Title:       book.Title,
		Author:      book.Author,
		Price:       book.Price,
		Stock:       int32(book.Stock),
		Year:        int32(book.Year),
		CategoryId:  int32(book.CategoryID),
		ImageBase64: book.ImageBase64,
	}, nil
}
func (h *BookHanlder) UpdateBook(ctx context.Context, req *toko.BookRequest) (*toko.BookResponse, error) {
	h.log.Info("Create book in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.PUT)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	payload := dto.BookReq{
		Title:       req.GetTitle(),
		Author:      req.GetAuthor(),
		Price:       req.GetPrice(),
		Stock:       int(req.GetStock()),
		Year:        int(req.GetYear()),
		CategoryId:  int(req.GetCategoryId()),
		ImageBase64: req.GetImageBase64(),
	}

	ids, err := helper.GetIdFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = h.uc.UpdateBook(ctx, id, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.BookResponse{
		Message: constant.Succes,
	}, nil
}
func (h *BookHanlder) DeleteBook(ctx context.Context, req *toko.EmptyBook) (*toko.BookResponse, error) {
	h.log.Info("Create book in handler", req)

	ctx, err := h.middleware.Require(ctx, constant.PUT)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ids, err := helper.GetIdFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, err := strconv.Atoi(ids)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = h.uc.DeleteBook(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.BookResponse{
		Message: constant.Succes,
	}, nil
}
