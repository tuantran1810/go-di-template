package repositories

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
	"github.com/tuantran1810/go-di-template/internal/repositories/mysql"
)

func (s *UserAttributeRepositoryTestSuite) getTestData(t *testing.T) ([]entities.User, []entities.UserAttribute) {
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
		}, []entities.UserAttribute{
			{
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "test1",
				Value:     "test1",
			},
			{
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    1,
				Key:       "test2",
				Value:     "test2",
			},
			{
				CreatedAt: now,
				UpdatedAt: now,
				UserID:    2,
				Key:       "test3",
				Value:     "test3",
			},
		}
}

func (s *UserAttributeRepositoryTestSuite) createTestData(t *testing.T, userRepository *UserRepository, attStore *UserAttributeRepository) {
	t.Helper()

	users, atts := s.getTestData(t)

	if _, err := userRepository.CreateMany(context.Background(), nil, users); err != nil {
		t.Errorf("failed to create data: %v", err)
		return
	}
	if _, err := attStore.CreateMany(context.Background(), nil, atts); err != nil {
		t.Errorf("failed to create data: %v", err)
		return
	}
}

func (s *UserAttributeRepositoryTestSuite) setup(t *testing.T, port int) (*UserRepository, *UserAttributeRepository, error) {
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
		return nil, nil, err
	}

	userTransformer := entities.NewExtendedDataTransformer(&userTransformer{})
	attTransformer := entities.NewExtendedDataTransformer(&userAttributeTransformer{})
	return &UserRepository{
			GenericRepository: mysql.NewGenericRepository(r, userTransformer),
		}, &UserAttributeRepository{
			GenericRepository: mysql.NewGenericRepository(r, attTransformer),
		}, nil
}

func (s *UserAttributeRepositoryTestSuite) cleanup(t *testing.T, store *UserRepository) {
	t.Helper()

	if err := store.DB().Exec("DROP TABLE IF EXISTS `test`.`users`").Error; err != nil {
		t.Logf("failed to cleanup data: %v\n", err)
		return
	}

	if err := store.DB().Exec("DROP TABLE IF EXISTS `test`.`user_attributes`").Error; err != nil {
		t.Logf("failed to cleanup data: %v\n", err)
		return
	}
}

type UserAttributeRepositoryTestSuite struct {
	suite.Suite
	userRepository *UserRepository
	attStore       *UserAttributeRepository
	container      *mysqlModule.MySQLContainer
}

func (s *UserAttributeRepositoryTestSuite) SetupSuite() {
	t := s.T()
	if err := os.Setenv("TZ", "UTC"); err != nil {
		t.Errorf("failed to set time zone: %v", err)
		return
	}

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

	userRepository, attStore, err := s.setup(t, port.Int())
	s.Require().NoError(err)
	s.userRepository = userRepository
	s.Require().NotNil(s.userRepository)
	s.attStore = attStore
	s.Require().NotNil(s.attStore)
	if err := userRepository.DB().Exec("SET @@global.time_zone = '+00:00'").Error; err != nil {
		t.Errorf("failed to set time zone: %v", err)
		return
	}
}

func (s *UserAttributeRepositoryTestSuite) TearDownSuite() {
	t := s.T()
	s.cleanup(t, s.userRepository)

	if err := testcontainers.TerminateContainer(s.container); err != nil {
		t.Errorf("failed to terminate container: %v", err)
		return
	}
}

func (s *UserAttributeRepositoryTestSuite) SetupTest() {
	t := s.T()
	s.Require().NoError(s.userRepository.AutoMigrate(context.Background()))

	if err := s.userRepository.DB().Exec("TRUNCATE TABLE `test`.`users`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	s.Require().NoError(s.attStore.AutoMigrate(context.Background()))

	if err := s.attStore.DB().Exec("TRUNCATE TABLE `test`.`user_attributes`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	s.createTestData(t, s.userRepository, s.attStore)
}

func (s *UserAttributeRepositoryTestSuite) TearDownTest() {
	t := s.T()
	if err := s.userRepository.DB().Exec("TRUNCATE TABLE `test`.`users`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}

	if err := s.attStore.DB().Exec("TRUNCATE TABLE `test`.`user_attributes`").Error; err != nil {
		t.Errorf("failed to cleanup data: %v\n", err)
		return
	}
}

func (s *UserAttributeRepositoryTestSuite) TestUserAttributeRepository_GetByUserID() {
	t := s.T()
	store := s.attStore
	_, atts := s.getTestData(t)
	now := atts[0].CreatedAt

	tests := []struct {
		name    string
		userID  uint
		want    []entities.UserAttribute
		wantErr bool
	}{
		{
			name:   "user1",
			userID: 1,
			want: []entities.UserAttribute{
				{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "test1",
					Value:     "test1",
				},
				{
					ID:        2,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    1,
					Key:       "test2",
					Value:     "test2",
				},
			},
			wantErr: false,
		},
		{
			name:   "user2",
			userID: 2,
			want: []entities.UserAttribute{
				{
					ID:        3,
					CreatedAt: now,
					UpdatedAt: now,
					UserID:    2,
					Key:       "test3",
					Value:     "test3",
				},
			},
			wantErr: false,
		},
		{
			name:    "not found",
			userID:  10,
			want:    []entities.UserAttribute{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.GetByUserID(context.TODO(), nil, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserAttributeRepository.GetByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserAttributeRepository.GetByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAttributeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserAttributeRepositoryTestSuite))
}
