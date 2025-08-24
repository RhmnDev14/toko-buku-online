package handler

import (
	"context"
	"strconv"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/usecase"
	"toko_buku_online/toko_buku_online/api/gen/go/toko/v1"

	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	toko.UnimplementedAuthServiceServer
	uc usecase.AuhtUc
}

func NewAuthHandler(uc usecase.AuhtUc) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *toko.LoginRequest) (*toko.LoginResponse, error) {
	// panggil usecase.Login dengan username
	res, err := h.uc.Login(ctx, req.GetUsername())
	if err != nil {
		return nil, err
	}

	id := strconv.Itoa(res.UserId)

	// mapping DTO ke response protobuf
	return &toko.LoginResponse{
		Token:  res.Token,
		UserId: id,
		Role:   res.Role,
	}, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *toko.RegisterRequest) (*emptypb.Empty, error) {
	// mapping request protobuf ke DTO
	dtoReq := dto.RegisterReq{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	}

	// panggil usecase.Register
	err := h.uc.Register(ctx, dtoReq)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
