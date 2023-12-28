package serviceauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/DrumPatiphon/go-rest-api-service/config"
	"github.com/DrumPatiphon/go-rest-api-service/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type serviceAuth struct {
	mapClaims *serviceMapClaim
	cfg       config.IJwtConfig
}

type ServiceAdmin struct {
	*serviceAuth
}

type serviceMapClaim struct {
	Claims               *users.UserClaims `json:"claims"` // Know As Payload
	jwt.RegisteredClaims                   // Required inDocument
}

type IServiceAuth interface {
	SignToken() string
}
type IserviceAdmin interface {
	SignToken() string
}

func (a *serviceAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims) // Sign Token with Payload
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}

func (a *ServiceAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims) // Sign Token with Payload
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*serviceMapClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &serviceMapClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("singing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token fomat is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*serviceMapClaim); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*serviceMapClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &serviceMapClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("singing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token fomat is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token has expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*serviceMapClaim); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return &jwt.NumericDate{time.Unix(t, 0)} // Convert time.seconds to time.unix
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &serviceAuth{
		cfg: cfg,
		mapClaims: &serviceMapClaim{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecommerceshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewServiceAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IServiceAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknow tokenType")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IServiceAuth {
	return &serviceAuth{
		cfg: cfg,
		mapClaims: &serviceMapClaim{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecommerceshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AcessExpriresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IServiceAuth {
	return &serviceAuth{
		cfg: cfg,
		mapClaims: &serviceMapClaim{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "ecommerceshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IServiceAuth {
	return &ServiceAdmin{
		serviceAuth: &serviceAuth{
			cfg: cfg,
			mapClaims: &serviceMapClaim{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "ecommerceshop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}

}
