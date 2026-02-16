package interceptor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/Molov30/go-microservices/services/gateway/internal/model"
)

type JWTCLaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTInterceptor struct {
	jwtSecret []byte
}

var publicMethods = map[string]bool{
	"/gateway.Gateway/Register": true,
	"/gateway.Gateway/Login":    true,
	"/gateway.Gateway/Refresh":  true,
}

func NewJWTInterceptor(jwtSecret string) *JWTInterceptor {
	return &JWTInterceptor{
		jwtSecret: []byte(jwtSecret),
	}
}

func (i *JWTInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		userID, err := i.extractUserIDFromJWT(ctx)
		if err != nil {
			return nil, err
		}

		ctxWithUserID := context.WithValue(ctx, model.UserIDKey, userID)

		return handler(ctxWithUserID, req)
	}
}

func (i *JWTInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if publicMethods[info.FullMethod] {
			return handler(srv, ss)
		}

		userID, err := i.extractUserIDFromJWT(ss.Context())
		if err != nil {
			return err
		}

		ctxWithUserID := context.WithValue(ss.Context(), model.UserIDKey, userID)
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctxWithUserID,
		}

		return handler(srv, wrappedStream)
	}
}

func (i *JWTInterceptor) extractUserIDFromJWT(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authHeaders := md.Get("Authorization")
	if len(authHeaders) == 0 {
		return 0, status.Errorf(codes.Unauthenticated, "authorization header is not provided")
	}

	token := authHeaders[0]
	if !strings.HasPrefix(token, "Bearer ") {
		return 0, status.Errorf(codes.Unauthenticated, "invalid token format")
	}

	accessToken := strings.TrimPrefix(token, "Bearer ")
	if accessToken == "" {
		return 0, status.Errorf(codes.Unauthenticated, "access token is empty")
	}

	userID, err := i.parseAndValidateJWT(accessToken)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return userID, nil
}

func (i *JWTInterceptor) parseAndValidateJWT(tokenString string) (uint64, error) {
	claimsNew := JWTCLaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claimsNew, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return i.jwtSecret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, errors.New("token is not valid")
	}

	claims, ok := token.Claims.(*JWTCLaims)
	if !ok {
		return 0, errors.New("failed to extract claims")
	}

	if claims.UserID == 0 {
		return 0, errors.New("user_id not found in token")
	}

	return claims.UserID, nil
}

func (i *JWTInterceptor) ValidateTokenWithoutAuthService(tokenString string) (uint64, error) {
	return i.parseAndValidateJWT(tokenString)
}

func GetUserIDFromContext(ctx context.Context) (uint64, error) {
	userID, ok := ctx.Value(model.UserIDKey).(uint64)
	if !ok {
		return 0, errors.New("user_id not found in context")
	}
	return userID, nil
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
