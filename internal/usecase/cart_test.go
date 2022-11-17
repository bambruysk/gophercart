package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"gophercart/internal/models"
	"gophercart/internal/usecase/mocks"
)

func Test_processor_Add(t *testing.T) {
	type fields struct {
		logger    *logrus.Logger
		repo      *mocks.Repo
		warehouse *mocks.Warehouse
	}
	type args struct {
		ctx    context.Context
		userID models.User
		good   models.Good
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantRepoReq bool
		wantErr     bool
		wantCount   int
	}{
		{
			name: "Test",
			fields: fields{
				logger:    logrus.StandardLogger(),
				repo:      mocks.NewRepo(t),
				warehouse: mocks.NewWarehouse(t),
			},
			args: args{
				ctx:    context.Background(),
				userID: models.User(uuid.New()),
				good: models.Good{
					ID:        uuid.New(),
					Count:     1,
					CreatedAt: time.Time{},
				},
			},
			wantCount:   100,
			wantErr:     false,
			wantRepoReq: true,
		},
		{
			name: "Test",
			fields: fields{
				logger:    logrus.StandardLogger(),
				repo:      mocks.NewRepo(t),
				warehouse: mocks.NewWarehouse(t),
			},
			args: args{
				ctx:    context.Background(),
				userID: models.User(uuid.New()),
				good: models.Good{
					ID:        uuid.New(),
					Count:     1000,
					CreatedAt: time.Time{},
				},
			},
			wantCount:   100,
			wantRepoReq: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &processor{
				logger:    tt.fields.logger,
				repo:      tt.fields.repo,
				warehouse: tt.fields.warehouse,
			}
			if tt.wantRepoReq {
				tt.fields.repo.EXPECT().
					Create(tt.args.ctx, tt.args.userID, tt.args.good).
					Return(errors.New("my erros"))
			}
			tt.fields.warehouse.EXPECT().
				CheckGood(tt.args.ctx, tt.args.good).
				Return(tt.wantCount, nil)

			if err := p.Add(tt.args.ctx, tt.args.userID, tt.args.good); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
