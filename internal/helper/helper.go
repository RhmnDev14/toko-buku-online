package helper

import (
	"context"
	"fmt"
	"time"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/entity"

	"google.golang.org/grpc/metadata"
)

func RegisterToEntity(dto dto.RegisterReq) entity.User {
	return entity.User{
		Name:      dto.Name,
		Email:     dto.Email,
		Password:  dto.Password,
		Role:      dto.Role,
		CreatedAt: time.Now(),
	}
}

func CategoryToEntity(dto dto.CategoryReq) entity.Category {
	return entity.Category{
		Name: dto.Name,
	}
}

func GetIdFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata tidak tersedia")
	}

	id := md.Get("id")
	if len(id) == 0 {
		return "", fmt.Errorf("authorization metadata tidak ditemukan")
	}

	return id[0], nil
}

func BookToEntity(dto dto.BookReq) entity.Book {
	return entity.Book{
		Title:       dto.Title,
		Author:      dto.Author,
		Price:       dto.Price,
		Stock:       dto.Stock,
		Year:        dto.Year,
		CategoryID:  uint(dto.CategoryId),
		ImageBase64: dto.ImageBase64,
	}
}

func GetParamFromMetadata(ctx context.Context) (string, string, string, string, string, string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", "", "", "", "", "", fmt.Errorf("metadata tidak tersedia")
	}

	src := md.Get("search")
	if len(src) == 0 {
		return "", "", "", "", "", "", fmt.Errorf("search metadata tidak ditemukan")
	}
	sortBy := md.Get("sortBy")
	if len(sortBy) == 0 {
		return "", "", "", "", "", "", fmt.Errorf("sortBy metadata tidak ditemukan")
	}
	sortType := md.Get("sortType")
	if len(sortType) == 0 {
		return "", "", "", "", "", "", fmt.Errorf("sortType metadata tidak ditemukan")
	}
	filter := md.Get("filter")
	if len(filter) == 0 {
		return "", "", "", "", "", "", fmt.Errorf("filter metadata tidak ditemukan")
	}
	page := md.Get("page")
	if len(page) == 0 {
		return "", "", "", "", "", "", fmt.Errorf("page metadata tidak ditemukan")
	}
	limit := md.Get("limit")
	if len(limit) == 0 {
		return "", "", "", "", "", "", fmt.Errorf("limit metadata tidak ditemukan")
	}

	return src[0], sortBy[0], sortType[0], filter[0], page[0], limit[0], nil
}

func OrderItemToEntity(dto dto.OrderItem) entity.OrderItem {
	return entity.OrderItem{
		BookID:   uint(dto.BookId),
		Quantity: dto.Quantity,
		Price:    dto.Price,
	}
}

func TotalAmount(o dto.Order) float64 {
	var total float64
	for _, item := range o.OrderItems {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func OrderToEntity(ctx context.Context, dto dto.Order) entity.Order {
	total := TotalAmount(dto)
	userID := ctx.Value(constant.UserIDKey).(int)
	return entity.Order{
		UserID:     uint(userID),
		TotalPrice: total,
		Status:     constant.PENDING,
		CreatedAt:  time.Now(),
	}
}
