package helper

import (
	"time"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/entity"
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
