package service

import (
	"fmt"
	"time"
	"toko_buku_online/internal/config"
	"toko_buku_online/internal/constant"
	"toko_buku_online/internal/dto"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtService interface {
	CreateToken(user entity.User) (dto.AuthResponseDto, error)
	DecodeToken(tokenString string) (*entity.Claim, error)
}

type jwtService struct {
	cfgToken config.TokenConfig
	log      logger.Logger
}

func (j *jwtService) CreateToken(user entity.User) (dto.AuthResponseDto, error) {
	j.log.Info("create access token in service", user)

	jti := uuid.NewString()
	claims := entity.Claim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfgToken.IssuerName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfgToken.JwtExpiresTime)),
		},
		UserId: int(user.ID),
		Role:   user.Role,
		Jti:    jti,
	}

	token := jwt.NewWithClaims(j.cfgToken.JwtSigningMethod, claims)

	ss, err := token.SignedString(j.cfgToken.JwtSignatureKy)
	if err != nil {
		j.log.Error("Error : ", err)
		return dto.AuthResponseDto{}, fmt.Errorf(constant.ErrorInternalSystem)
	}

	return dto.AuthResponseDto{
		Token:  ss,
		UserId: claims.UserId,
		Role:   claims.Role,
		Jti:    claims.Jti,
	}, nil
}

func (j *jwtService) DecodeToken(tokenString string) (*entity.Claim, error) {
	j.log.Info("decode token in service", nil)

	token, err := jwt.ParseWithClaims(tokenString, &entity.Claim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return j.cfgToken.JwtSignatureKy, nil
	})
	if err != nil {
		j.log.Error("Error : ", err)
		return nil, fmt.Errorf(constant.ErrorInternalSystem)
	}

	claims, ok := token.Claims.(*entity.Claim)
	if !ok || !token.Valid {
		j.log.Error("Error : invalid token", nil)
		return nil, fmt.Errorf(constant.ErrorInternalSystem)
	}

	return claims, nil
}

func NewJwtService(cfgToken config.TokenConfig, log logger.Logger) JwtService {
	return &jwtService{
		cfgToken: cfgToken,
		log:      log,
	}
}
