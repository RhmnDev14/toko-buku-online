package usecase

import (
	"context"
	"fmt"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/helper"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/repository"
	"toko_buku_online/internal/service"

	"golang.org/x/crypto/bcrypt"
)

type AuhtUc interface {
	Login(ctx context.Context, payload dto.LoginReq) (dto.AuthResponseDto, error)
	Register(ctx context.Context, payload dto.RegisterReq) error
}

type authUc struct {
	authRepo   repository.AuthRepo
	log        logger.Logger
	jwtService service.JwtService
}

func NewAuthUc(authRepo repository.AuthRepo, log logger.Logger, jwtService service.JwtService) AuhtUc {
	return &authUc{
		authRepo:   authRepo,
		log:        log,
		jwtService: jwtService,
	}
}

func (u *authUc) Login(ctx context.Context, payload dto.LoginReq) (dto.AuthResponseDto, error) {
	u.log.Info("login in usecase", payload)

	user, err := u.authRepo.Login(ctx, payload.Email)
	if err != nil {
		return dto.AuthResponseDto{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return dto.AuthResponseDto{}, fmt.Errorf(constant.ErrorLogin)
	}

	return u.jwtService.CreateToken(user)
}
func (u *authUc) Register(ctx context.Context, payload dto.RegisterReq) error {
	u.log.Info("register in usecase", payload)

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("Error : ", err)
		return fmt.Errorf(constant.ErrorInternalSystem)
	}

	payload.Password = string(hash)
	return u.authRepo.Register(ctx, helper.RegisterToEntity(payload))
}
