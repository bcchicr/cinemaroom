package api

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/bcchicr/cinemaroom/services/authn/app"
	"github.com/bcchicr/cinemaroom/services/authn/config"
	api "github.com/bcchicr/cinemaroom/services/authn/generated/cinemaroom/authn/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type registeredUserType struct {
	Email    string
	Login    string
	Password string
}

var registeredUser registeredUserType

func randomEmail() string {
	return fmt.Sprintf(
		"user%d@example.com",
		rand.Int63(),
	)
}

func randomLogin() string {
	return fmt.Sprintf(
		"user%d",
		rand.Int63(),
	)
}

func mustContext(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func startTestServer(container *app.Container, t *testing.T) (addr string, stop func()) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	s := container.GrpcServer

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Log(err)
		}
	}()

	return lis.Addr().String(), func() {
		lis.Close()
	}
}

func newClient(t *testing.T) (clinet api.AuthnServiceClient, cleanup func()) {
	cfg, err := config.NewConfig()
	if err != nil {
		t.Fatalf("cannot get config: %v", err)
	}

	container, err := app.NewContainer(cfg)
	if err != nil {
		t.Fatalf("failed to build container: %v", err)
	}

	addr, stopServer := startTestServer(container, t)

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		stopServer()
		container.Close()
		t.Fatal(err)
	}

	cleanup = func() {
		conn.Close()
		stopServer()
		container.Close()
	}

	return api.NewAuthnServiceClient(conn), cleanup
}

func TestAuthn_Register(t *testing.T) {
	client, cleanup := newClient(t)
	defer cleanup()

	someEmail := randomEmail()
	someLogin := randomLogin()
	somePassword := "password123"

	t.Run("happy path", func(t *testing.T) {
		resp, err := client.Register(mustContext(t), &api.RegisterRequest{
			Login:    someLogin,
			Email:    &someEmail,
			Password: somePassword,
		})
		if err != nil {
			t.Fatalf("register error: %v", err)
		}

		if resp.Tokens == nil {
			t.Fatal("tokens is nil")
		}
		if resp.Tokens.JwtToken.Value == "" {
			t.Fatal("jwt token is empty")
		}
		if resp.Tokens.RefreshToken.Value == "" {
			t.Fatal("refresh token is empty")
		}

		registeredUser = registeredUserType{
			Email:    someEmail,
			Login:    someLogin,
			Password: somePassword,
		}
	})

	t.Run("email already taken", func(t *testing.T) {
		_, err := client.Register(mustContext(t), &api.RegisterRequest{
			Login:    randomLogin(),
			Email:    &someEmail,
			Password: "password123",
		})

		if err == nil {
			t.Fatalf("expected register error, got success")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})

	t.Run("login already taken", func(t *testing.T) {
		email := randomEmail()

		_, err := client.Register(mustContext(t), &api.RegisterRequest{

			Login:    someLogin,
			Email:    &email,
			Password: "password123",
		})

		if err == nil {
			t.Fatalf("expected register error, got success")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})
}

func TestAuthn_Login(t *testing.T) {
	client, cleanup := newClient(t)
	defer cleanup()

	tests := []struct {
		name string
		req  *api.LoginRequest
	}{
		{
			name: "login by email",
			req: &api.LoginRequest{
				Identifier: &api.LoginRequest_Email{
					Email: registeredUser.Email,
				},
				Password: registeredUser.Password,
			},
		},
		{
			name: "login by login",
			req: &api.LoginRequest{
				Identifier: &api.LoginRequest_Login{
					Login: registeredUser.Login,
				},
				Password: registeredUser.Password,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.Login(mustContext(t), tt.req)
			if err != nil {
				t.Fatalf("login error: %v", err)
			}
			if resp.Tokens.JwtToken.Value == "" {
				t.Fatal("jwt token is empty")
			}
			if resp.Tokens.RefreshToken.Value == "" {
				t.Fatal("refresh token is empty")
			}
		})
	}

	t.Run("failed login by email", func(t *testing.T) {
		email := randomEmail()

		_, err := client.Login(mustContext(t), &api.LoginRequest{
			Identifier: &api.LoginRequest_Email{
				Email: email,
			},
			Password: "non-existent-password",
		})
		if err == nil {
			t.Fatalf("expected login error: %v", err)
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})

	t.Run("failed login by login", func(t *testing.T) {
		login := randomLogin()

		_, err := client.Login(mustContext(t), &api.LoginRequest{
			Identifier: &api.LoginRequest_Login{
				Login: login,
			},
			Password: "non-existent-password",
		})
		if err == nil {
			t.Fatalf("expected login error: %v", err)
		}
		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})
}

func TestAuthn_Refresh(t *testing.T) {
	client, cleanup := newClient(t)
	defer cleanup()

	t.Run("happy path", func(t *testing.T) {
		loginResp, err := client.Login(mustContext(t), &api.LoginRequest{
			Identifier: &api.LoginRequest_Login{
				Login: registeredUser.Login,
			},
			Password: registeredUser.Password,
		})
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		refreshResp, err := client.Refresh(mustContext(t), &api.RefreshRequest{
			RefreshToken: loginResp.Tokens.RefreshToken.Value,
		})
		if err != nil {
			t.Fatalf("refresh failed: %v", err)
		}

		if refreshResp.Tokens.JwtToken.Value == "" {
			t.Fatal("jwt token is empty")
		}
		if refreshResp.Tokens.RefreshToken.Value == "" {
			t.Fatal("refresh token is empty")
		}
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		_, err := client.Refresh(mustContext(t), &api.RefreshRequest{
			RefreshToken: "invalid_token",
		})
		if err == nil {
			t.Fatalf("expected refresh error, got success")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})
}

func TestAuthn_Edit(t *testing.T) {
	client, cleanup := newClient(t)
	defer cleanup()

	t.Run("happy path", func(t *testing.T) {
		loginResp, err := client.Login(mustContext(t), &api.LoginRequest{
			Identifier: &api.LoginRequest_Login{
				Login: registeredUser.Login,
			},
			Password: registeredUser.Password,
		})
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		ctx := metadata.AppendToOutgoingContext(
			mustContext(t),
			"Authorization",
			"Bearer "+loginResp.Tokens.JwtToken.Value,
		)

		login := randomLogin()
		resp, err := client.Edit(ctx, &api.EditRequest{
			Login: &login,
		})
		if err != nil {
			t.Fatalf("edit error: %v", err)
		}

		if resp.AccountId.Value == "" {
			t.Fatal("account id is empty")
		}

		registeredUser.Login = login
	})

	t.Run("same login", func(t *testing.T) {
		loginResp, err := client.Login(mustContext(t), &api.LoginRequest{
			Identifier: &api.LoginRequest_Login{
				Login: registeredUser.Login,
			},
			Password: registeredUser.Password,
		})
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		ctx := metadata.AppendToOutgoingContext(
			mustContext(t),
			"Authorization",
			"Bearer "+loginResp.Tokens.JwtToken.Value,
		)

		login := registeredUser.Login
		_, err = client.Edit(ctx, &api.EditRequest{
			Login: &login,
		})
		if err == nil {
			t.Fatal("expected unauthenticated error")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})

	t.Run("same email", func(t *testing.T) {
		loginResp, err := client.Login(mustContext(t), &api.LoginRequest{
			Identifier: &api.LoginRequest_Login{
				Login: registeredUser.Login,
			},
			Password: registeredUser.Password,
		})
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		ctx := metadata.AppendToOutgoingContext(
			mustContext(t),
			"Authorization",
			"Bearer "+loginResp.Tokens.JwtToken.Value,
		)

		email := registeredUser.Email
		_, err = client.Edit(ctx, &api.EditRequest{
			Email: &email,
		})
		if err == nil {
			t.Fatal("expected unauthenticated error")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})

	t.Run("unauthenticated edit", func(t *testing.T) {
		client, cleanup := newClient(t)
		defer cleanup()

		login := "new_login"
		_, err := client.Edit(mustContext(t), &api.EditRequest{
			Login: &login,
		})
		if err == nil {
			t.Fatal("expected unauthenticated error")
		}

		st, ok := status.FromError(err)
		if !ok {
			t.Fatalf("expected grpc status error, got: %v", err)
		}

		if st.Code() != codes.Unauthenticated {
			t.Fatalf("unexpected code: %v", st.Code())
		}
	})

}
