//go:build unit
// +build unit

package auth_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/ereminiu/pvz/internal/handler"
	myerrors "github.com/ereminiu/pvz/internal/my_errors"
	"github.com/ereminiu/pvz/internal/usecases"
	mock_usecases "github.com/ereminiu/pvz/internal/usecases/mocks"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestSignIn(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockAuthUsecases, username, password string)

	type SingInInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	testTable := []struct {
		name                 string
		input                SingInInput
		call                 mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: SingInInput{
				Username: "Alina Rin",
				Password: "qwerty128",
			},
			call: func(ctx context.Context, s *mock_usecases.MockAuthUsecases, username, password string) {
				s.EXPECT().SignIn(ctx, username, password).Return("", nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Invalid Password",
			input: SingInInput{
				Username: "Alina Rin",
				Password: "qwerty256",
			},
			call: func(ctx context.Context, s *mock_usecases.MockAuthUsecases, username, password string) {
				s.EXPECT().SignIn(ctx, username, password).Return("", myerrors.ErrIncorrectPassword)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			auth := mock_usecases.NewMockAuthUsecases(c)
			tc.call(ctx, auth, tc.input.Username, tc.input.Password)

			uc := &usecases.Usecases{
				AuthUsecases: auth,
			}

			handler := handler.New(ctx, uc)

			router := http.NewServeMux()
			router.HandleFunc("POST /sign-in", handler.SignIn)

			body, err := sonic.MarshalString(tc.input)
			if err != nil {
				panic(err)
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(body))

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
		})
	}
}
