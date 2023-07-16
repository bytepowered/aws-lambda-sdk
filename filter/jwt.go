package filter

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bytepowered/aws-lambda-sdk"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
)

const (
	JwtEnvSecret = "AUTH_JWT_SECRET"
	JwtEnvIssuer = "AUTH_JWT_ISSUER"
	JwtEnvExpsec = "AUTH_JWT_EXPSEC"
)

const (
	ctxKeyJwtClaims = "@internal.aws.lambda.sdk.jwt.claims"
)

// JWTSign 生成JWT字符串
func JWTSign(sub string, id string) (string, error) {
	expsec := lambda.RequiredIntEnv(JwtEnvExpsec, 7*24*3600)
	now := time.Now()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        id,
		Issuer:    lambda.RequiredEnv(JwtEnvIssuer),
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(expsec))),
		IssuedAt:  jwt.NewNumericDate(now),
	}).SignedString([]byte(lambda.RequiredEnv(JwtEnvSecret)))
	if err != nil {
		return "", fmt.Errorf("JWT-SIGN-TOKEN: %w", err)
	}
	return token, nil
}

// JWTParse 解析JWT字符串
func JWTParse(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(lambda.RequiredEnv(JwtEnvSecret)), nil
	}, jwt.WithIssuer(lambda.RequiredEnv(JwtEnvIssuer)))
	if token.Valid {
		return token, nil
	}
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return nil, fmt.Errorf("JWT-MALFORMED: %w", err)
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return nil, fmt.Errorf("JWT-SIGN-INVALID: %w", err)
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return nil, fmt.Errorf("JWT-EXPIRED: %w", err)
	} else {
		return nil, fmt.Errorf("JWT-UNKNOWN: %w", err)
	}
}

func JWTFilter(next lambda.HandleFunc) lambda.HandleFunc {
	return func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (*events.APIGatewayV2HTTPResponse, error) {
		str := lambda.HeaderAuthorization(&req)
		if !lambda.CheckNotEmpty(str) {
			log.Println("ERROR: jwt token not found")
			return lambda.SendInvalidArgs("authorization")
		}
		token, err := JWTParse(str)
		if err != nil {
			log.Println("ERROR: jwt token invalid, error:", err, ", token:", str)
			return lambda.SendInvalidToken()
		}
		return next(context.WithValue(ctx, ctxKeyJwtClaims, token.Claims), req)
	}
}

func JWTLoadClaims(ctx context.Context) (claims jwt.Claims, ok bool) {
	value := ctx.Value(ctxKeyJwtClaims)
	claims, ok = value.(jwt.Claims)
	return
}