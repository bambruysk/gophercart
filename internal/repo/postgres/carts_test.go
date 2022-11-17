package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"gophercart/internal/models"
	"gophercart/internal/repo/postgres/testhelpers"
)

type TestStorager interface {
	Repo
}

type StorageTestSuite struct {
	suite.Suite
	TestStorager
	container *testhelpers.TestDatabase
}

func (sts *StorageTestSuite) SetupTest() {

	storageContainer := testhelpers.NewTestDatabase(sts.T())

	opts := &Options{
		Host:         storageContainer.Host(),
		Port:         storageContainer.Port(sts.T()),
		DatabaseName: "postgres",
		MigrationOptions: &MigrationOptions{
			Enable:  true,
			Path:    "./migrate",
			Version: 1,
		},
		auth: &AuthOptions{
			User:     "postgres",
			Password: "postgres",
		},
	}

	logger := logrus.StandardLogger()

	store, err := NewStorage(context.Background(), logger, opts)
	require.NoError(sts.T(), err)

	sts.TestStorager = store
	sts.container = storageContainer
}

func (sts *StorageTestSuite) TearDownTest() {
	sts.container.Close(sts.T())
}

func TestStorageTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
		return
	}

	t.Parallel()
	suite.Run(t, new(StorageTestSuite))
}

func (sts *StorageTestSuite) Test_storage_Create() {
	type args struct {
		ctx    context.Context
		userID models.User
		good   models.Good
	}

	now := time.Now()
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "Positive",
			args: args{
				ctx:    context.Background(),
				userID: models.User(uuid.New()),
				good: models.Good{
					ID:        uuid.New(),
					Count:     0,
					CreatedAt: now,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		sts.Run(tt.name, func() {
			s := sts.TestStorager
			if err := s.Create(tt.args.ctx, tt.args.userID, tt.args.good); (err != nil) != tt.wantErr {
				sts.T().Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
