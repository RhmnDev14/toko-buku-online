package middleware

import (
	"context"
	"fmt"
	"strings"
	"toko_buku_online/internal/config"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthMiddleware interface {
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

func (a *authMiddleware) Middleware(method string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		a.log.Info("gRPC RequireToken: Checking token", nil)

		// Ambil metadata dari context
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			a.log.Error("Metadata missing", nil)
			return nil, fmt.Errorf(constant.ErrorInternalSystem)
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			a.log.Error("Missing token", nil)
			return nil, fmt.Errorf(constant.ErrorInternalSystem)
		}

		tokenHeader := strings.TrimPrefix(authHeaders[0], "Bearer ")
		if tokenHeader == "" {
			a.log.Error("Empty token", nil)
			return nil, fmt.Errorf(constant.ErrorInternalSystem)
		}

		claims, err := a.jwtService.DecodeToken(tokenHeader)
		if err != nil {
			a.log.Error("Token decode failed", err)
			return nil, fmt.Errorf(constant.ErrorInternalSystem)
		}

		if claims.Role == constant.User && method == constant.ORDER || method == constant.GETBOOK {
			return nil, fmt.Errorf(constant.ErrorDontPermission)
		}

		return handler(ctx, req)
	}
}
