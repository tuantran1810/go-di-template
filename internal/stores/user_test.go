package stores

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	mysqlModule "github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/tuantran1810/go-di-template/internal/entities"
	"github.com/tuantran1810/go-di-template/internal/stores/mysql"
)

func (s *UserStoreTestSuite) getTestData(t *testing.T) []entities.User {
	t.Helper()
	now := time.Now().UTC().Truncate(time.Second)
	return []entities.User{
		{
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "user1",
		},
		{
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "user2",
		},
		{
			CreatedAt: now,
			UpdatedAt: now,
			Username:  "user3",
		},
	}
}

func (s *UserStoreTestSuite) createTestData(t *testing.T, store *UserStore) {
	t.Helper()

	if _, err := store.CreateMany(context.Background(), nil, s.getTestData(t)); err != nil {
		t.Errorf("failed to create data: %v", err)
		return
	}
}

func (s *UserStoreTestSuite) setup(t *testing.T, port int) (*UserStore, error) {
	t.Helper()

	config := mysql.RepositoryConfig{
		Username:  "root",
		Password:  "secret",
		Protocol:  "tcp",
		Address:   fmt.Sprintf("127.0.0.1:%d", port),
		Database:  "test",
		Params:    map[string]string{},
		Collation: "utf8mb4_general_ci",
		Loc:       time.Local,
		TLSConfig: "",

		Timeout:                 10 * time.Second,
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            10 * time.Second,
		AllowAllFiles:           false,
		AllowCleartextPasswords: false,
		AllowOldPasswords:       false,
		ClientFoundRows:         false,
		ColumnsWithAlias:        false,
		InterpolateParams:       false,
		MultiStatements:         false,
		ParseTime:               true,

		MaxOpenConns:           10,
		MaxIdleConns:           10,
		ConnMaxLifeTimeSeconds: 1800,
	}
	r := mysql.MustNewRepository(config)
	if err := r.Start(context.Background()); err != nil {
		return nil, err
	}

	transformer := entities.NewExtendedDataTransformer(&userTransformer{})
	return &UserStore{
		GenericStore: mysql.NewGenericStore(r, transformer),
	}, nil
}

func (s *UserStoreTestSuite) cleanup(t *testing.T, store *UserStore) {
	t.Helper()

	if err := store.DB().Exec("DROP TABLE IF EXISTS `test`.`users`").Error; err != nil {
		t.Logf("failed to cleanup data: %v\n", err)
		return
	}
}

type UserStoreTestSuite struct {
	suite.Suite
	store     *UserStore
	container *mysqlModule.MySQLContainer
}

func (s *UserStoreTestSuite) SetupSuite() {
	t := s.T()
	os.Setenv("TZ", "UTC")

	mysqlContainer, err := mysqlModule.Run(context.Background(),
		"mysql:lts",
		mysqlModule.WithDatabase("test"),
		mysqlModule.WithUsername("root"),
		mysqlModule.WithPassword("secret"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30*time.Second),
			wait.ForListeningPort("3306/tcp").WithStartupTimeout(30*time.Second),
		),
	)
	s.Require().NoError(err)

	port, err := mysqlContainer.MappedPort(context.Background(), "3306")
	s.Require().NoError(err)
	s.Require().NotNil(port)

	s.Require().NoError(err)
	s.container = mysqlContainer
	s.Require().NotNil(s.container)

	store, err := s.setup(t, port.Int())
	s.Require().NoError(err)
	s.store = store
	s.Require().NotNil(s.store)
	if err := store.DB().Exec("SET @@global.time_zone = '+00:00'").Error; err != nil {
		t.Errorf("failed to set time zone: %v", err)
		return
	}
}

func (s *UserStoreTestSuite) TearDownSuite() {
	t := s.T()
	cleanup(t, s.store)

	if err := testcontainers.TerminateContainer(s.container); err != nil {
		t.Errorf("failed to terminate container: %v", err)
		return
	}
}

func (s *UserStoreTestSuite) SetupTest() {
	t := s.T()
	s.Require().NoError(s.store.AutoMigrate(context.Background()))

	if err := s.store.DB().Exec("TRUNCATE TABLE `test`.`users`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	s.createTestData(t, s.store)
}

func (s *UserStoreTestSuite) TearDownTest() {
	t := s.T()
	if err := s.store.DB().Exec("TRUNCATE TABLE `test`.`users`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}
}

func (s *UserStoreTestSuite) TestUserStore_FindByUsername() {
	t := s.T()
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name     string
		username string
		want     *entities.User
		wantErr  bool
	}{
		{
			name:     "user1",
			username: "user1",
			want: &entities.User{
				ID:        1,
				CreatedAt: now,
				UpdatedAt: now,
				Username:  "user1",
			},
			wantErr: false,
		},
		{
			name:     "not found",
			username: "user10",
			want:     nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.store.FindByUsername(context.TODO(), nil, tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserStore.FindByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserStore.FindByUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserStoreTestSuite(t *testing.T) {
	suite.Run(t, new(UserStoreTestSuite))
}
