package handler

import (
	"context"
	"strconv"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/usecase"
	"toko_buku_online/toko_buku_online/api/gen/go/toko/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	toko.UnimplementedAuthServiceServer
	uc  usecase.AuhtUc
	log logger.Logger
}

func NewAuthHandler(uc usecase.AuhtUc, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		uc:  uc,
		log: log,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *toko.LoginRequest) (*toko.LoginResponse, error) {
	h.log.Info("Login in handler", req)

	paylod := dto.LoginReq{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	res, err := h.uc.Login(ctx, paylod)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	id := strconv.Itoa(res.UserId)
	return &toko.LoginResponse{
		Token:  res.Token,
		UserId: id,
		Role:   res.Role,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *toko.RegisterRequest) (*toko.RegisterRespone, error) {
	h.log.Info("Register in handler", req)

	dtoReq := dto.RegisterReq{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	}

	err := h.uc.Register(ctx, dtoReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &toko.RegisterRespone{
		Message: constant.Succes,
	}, nil
}
