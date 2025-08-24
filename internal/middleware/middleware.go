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
	Require(ctx context.Context, method string) error
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

func (a *authMiddleware) Require(ctx context.Context, method string) error {
	a.log.Info("gRPC RequireToken: Checking token", method)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf(constant.ErrorInternalSystem)
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return fmt.Errorf(constant.ErrorInternalSystem)
	}

	tokenHeader := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if tokenHeader == "" {
		return fmt.Errorf(constant.ErrorInternalSystem)
	}

	claims, err := a.jwtService.DecodeToken(tokenHeader)
	if err != nil {
		return fmt.Errorf(constant.ErrorInternalSystem)
	}

	if claims.Role == constant.User && (method != constant.ORDER || method != constant.GETBOOK) {
		return fmt.Errorf(constant.ErrorDontPermission)
	}

	return nil
}
