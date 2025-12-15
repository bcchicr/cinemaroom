package grpc

import (
	"context"
	"strings"

	authn "github.com/bcchicr/cinemaroom/services/authn/generated/cinemaroom/authn/v1"
	v1 "github.com/bcchicr/cinemaroom/services/authn/generated/cinemaroom/common/v1"
	"github.com/bcchicr/cinemaroom/services/authn/internal/application/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	authn.UnimplementedAuthnServiceServer

	registrationHandler handlers.RegistrationHandler
	loginHandler        handlers.LoginHandler
	logoutHandler       handlers.LogoutHandler
	refreshHandler      handlers.RefreshHandler
	editHandler         handlers.EditHandler
}

func NewServer(
	registrationHandler handlers.RegistrationHandler,
	loginHandler handlers.LoginHandler,
	logoutHandler handlers.LogoutHandler,
	refreshHandler handlers.RefreshHandler,
	editHandler handlers.EditHandler,
	errorHandlerInterceptor grpc.UnaryServerInterceptor,
) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(errorHandlerInterceptor),
	)

	authn.RegisterAuthnServiceServer(s, &server{
		registrationHandler: registrationHandler,
		loginHandler:        loginHandler,
		logoutHandler:       logoutHandler,
		refreshHandler:      refreshHandler,
		editHandler:         editHandler,
	})
	return s
}

func (s *server) Register(ctx context.Context, req *authn.RegisterRequest) (*authn.RegisterResponse, error) {
	command := handlers.RegisterCommand{
		Login:    req.Login,
		Email:    req.Email,
		Password: req.Password,
	}
	access, refresh, err := s.registrationHandler.Handle(ctx, command)
	if err != nil {
		return nil, err
	}

	response := &authn.RegisterResponse{
		Tokens: &authn.Tokens{
			JwtToken: &authn.JwtToken{
				Value:     access.Value(),
				ExpiresAt: timestamppb.New(access.ExpiresAt()),
			},
			RefreshToken: &authn.RefreshToken{
				Value:     refresh.Value,
				ExpiresAt: timestamppb.New(refresh.ExpiresAt),
			},
		},
	}

	return response, nil
}

func (s *server) Login(ctx context.Context, req *authn.LoginRequest) (*authn.LoginResponse, error) {
	var Identifier string
	var IdentifierType handlers.LoginIdentifierType

	switch req.Identifier.(type) {
	case *authn.LoginRequest_Email:
		Identifier = req.GetEmail()
		IdentifierType = handlers.LoginIdentifierEmail
	case *authn.LoginRequest_Login:
		Identifier = req.GetLogin()
		IdentifierType = handlers.LoginIdentifierLogin
	}

	command := handlers.LoginCommand{
		Identifier:     Identifier,
		IdentifierType: IdentifierType,
		Password:       req.Password,
	}
	access, refresh, err := s.loginHandler.Handle(ctx, command)
	if err != nil {
		return nil, err
	}

	response := &authn.LoginResponse{
		Tokens: &authn.Tokens{
			JwtToken: &authn.JwtToken{
				Value:     access.Value(),
				ExpiresAt: timestamppb.New(access.ExpiresAt()),
			},
			RefreshToken: &authn.RefreshToken{
				Value:     refresh.Value,
				ExpiresAt: timestamppb.New(refresh.ExpiresAt),
			},
		},
	}

	return response, nil
}

func (s *server) Refresh(ctx context.Context, req *authn.RefreshRequest) (*authn.RefreshResponse, error) {

	command := handlers.RefreshCommand{
		RefreshTokenString: req.RefreshToken,
	}

	access, refresh, err := s.refreshHandler.Handle(ctx, command)
	if err != nil {
		return nil, err
	}

	response := &authn.RefreshResponse{
		Tokens: &authn.Tokens{
			JwtToken: &authn.JwtToken{
				Value:     access.Value(),
				ExpiresAt: timestamppb.New(access.ExpiresAt()),
			},
			RefreshToken: &authn.RefreshToken{
				Value:     refresh.Value,
				ExpiresAt: timestamppb.New(refresh.ExpiresAt),
			},
		},
	}

	return response, nil
}

func (s *server) Logout(ctx context.Context, req *authn.LogoutRequest) (*authn.LogoutResponse, error) {
	command := handlers.LogoutCommand{
		JwtTokenString:     req.JwtToken,
		RefreshTokenString: req.RefreshToken,
	}

	accountID, err := s.logoutHandler.Handle(ctx, command)
	if err != nil {
		return nil, err
	}

	response := &authn.LogoutResponse{
		AccountId: &v1.Uuid{
			Value: accountID.String(),
		},
	}
	return response, nil
}
func (s *server) Edit(ctx context.Context, req *authn.EditRequest) (*authn.EditResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("Authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := authHeaders[0]
	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, prefix)

	command := handlers.EditCommand{
		AccessTokenString: tokenString,
		Login:             req.Login,
		Email:             req.Email,
		Password:          req.Password,
	}

	accountID, err := s.editHandler.Handle(ctx, command)
	if err != nil {
		return nil, err
	}

	response := &authn.EditResponse{
		AccountId: &v1.Uuid{
			Value: accountID.String(),
		},
	}

	return response, nil
}
