package helper

import (
	"context"
	"fmt"
	"time"
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
