package pvz_test

import (
	"context"
	"testing"

	"github.com/ereminiu/pvz/internal/pb/api"
	"github.com/ereminiu/pvz/internal/transport/grpc/handler/pvz"
	mock_usecases "github.com/ereminiu/pvz/internal/usecases/mocks"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestServer_RefundList(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockPVZUsecases, page, limit int, orderBy string, pattern map[string]string)
	testTable := []struct {
		name  string
		page  int
		limit int
		call  mockBehavior
		Error error
	}{
		{
			name:  "OK",
			page:  1,
			limit: 5,
			call: func(ctx context.Context, s *mock_usecases.MockPVZUsecases, page, limit int, orderBy string, pattern map[string]string) {
				s.EXPECT().GetRefunds(ctx, page, limit, orderBy, pattern).Return(nil, nil)
			},
			Error: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			uc := mock_usecases.NewMockPVZUsecases(c)
			tc.call(ctx, uc, tc.page, tc.limit, "", nil)

			pvzClient := pvz.New(uc)
			_, err := pvzClient.RefundList(ctx, &api.RefundListRequest{
				Page:    int32(tc.page),
				Limit:   int32(tc.limit),
				OrderBy: "",
				Pattern: nil,
			})

			assert.ErrorIs(t, err, tc.Error)
		})
	}
}

func TestServer_HistoryList(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_usecases.MockPVZUsecases, page, limit int, orderBy string, pattern map[string]string)
	testTable := []struct {
		name  string
		page  int
		limit int
		call  mockBehavior
		Error error
	}{
		{
			name:  "OK",
			page:  1,
			limit: 5,
			call: func(ctx context.Context, s *mock_usecases.MockPVZUsecases, page, limit int, orderBy string, pattern map[string]string) {
				s.EXPECT().GetHistory(ctx, page, limit, orderBy, pattern).Return(nil, nil)
			},
			Error: nil,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx := context.Background()

			uc := mock_usecases.NewMockPVZUsecases(c)
			tc.call(ctx, uc, tc.page, tc.limit, "", nil)

			pvzClient := pvz.New(uc)
			_, err := pvzClient.HistoryList(ctx, &api.HistoryListRequest{
				Page:    int32(tc.page),
				Limit:   int32(tc.limit),
				OrderBy: "",
				Pattern: nil,
			})

			assert.ErrorIs(t, err, tc.Error)
		})
	}
}
