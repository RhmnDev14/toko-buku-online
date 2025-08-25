package middleware

import (
	"context"
	"fmt"
	"strings"
	"toko_buku_online/internal/config"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/service"

	"google.golang.org/grpc/metadata"
)

type AuthMiddleware interface {
	Require(ctx context.Context, method string) (context.Context, error)
}
type authMiddleware struct {
	jwtService service.JwtService
	log        logger.Logger
	cfg        config.TokenConfig
}

func NewGRPCAuthMiddleware(jwt service.JwtService, log logger.Logger, cfg config.TokenConfig) AuthMiddleware {
	return &authMiddleware{
		jwtService: jwt,
		log:        log,
		cfg:        cfg,
	}
}

func (a *authMiddleware) Require(ctx context.Context, method string) (context.Context, error) {
	a.log.Info("gRPC RequireToken: Checking token", method)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, fmt.Errorf(constant.ErrorInternalSystem)
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return ctx, fmt.Errorf(constant.ErrorInternalSystem)
	}

	tokenHeader := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if tokenHeader == "" {
		return ctx, fmt.Errorf(constant.ErrorInternalSystem)
	}

	claims, err := a.jwtService.DecodeToken(tokenHeader)
	if err != nil {
		return ctx, fmt.Errorf(constant.ErrorInternalSystem)
	}

	if claims.Role == constant.User && (method != constant.ORDER || method != constant.GETBOOK) {
		return ctx, fmt.Errorf(constant.ErrorDontPermission)
	}

	ctx = context.WithValue(ctx, constant.UserIDKey, claims.UserId)

	return ctx, nil
}
