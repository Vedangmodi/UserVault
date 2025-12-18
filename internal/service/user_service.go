package service

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"

	"uservault/internal/models"
	"uservault/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id int64) (*models.User, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]*models.User, error)
	UpdateUser(ctx context.Context, id int64, req models.UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	repo      repository.UserRepository
	validator *validator.Validate
}

func NewUserService(repo repository.UserRepository, v *validator.Validate) UserService {
	return &userService{
		repo:      repo,
		validator: v,
	}
}

func parseDOB(dobStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dobStr)
}

func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	dob, err := parseDOB(req.DOB)
	if err != nil {
		return nil, err
	}
	u, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &models.User{
		ID:   int64(u.ID),
		Name: u.Name,
		DOB:  u.Dob.Time.Format("2006-01-02"),
		Age:  models.CalculateAge(u.Dob.Time, now),
	}, nil
}

func (s *userService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	u, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &models.User{
		ID:   int64(u.ID),
		Name: u.Name,
		DOB:  u.Dob.Time.Format("2006-01-02"),
		Age:  models.CalculateAge(u.Dob.Time, now),
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int32) ([]*models.User, error) {
	users, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	result := make([]*models.User, 0, len(users))
	for _, u := range users {
		user := &models.User{
			ID:   int64(u.ID),
			Name: u.Name,
			DOB:  u.Dob.Time.Format("2006-01-02"),
			Age:  models.CalculateAge(u.Dob.Time, now),
		}
		result = append(result, user)
	}
	return result, nil
}

func (s *userService) UpdateUser(ctx context.Context, id int64, req models.UpdateUserRequest) (*models.User, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, err
	}
	dob, err := parseDOB(req.DOB)
	if err != nil {
		return nil, err
	}
	u, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &models.User{
		ID:   int64(u.ID),
		Name: u.Name,
		DOB:  u.Dob.Time.Format("2006-01-02"),
		Age:  models.CalculateAge(u.Dob.Time, now),
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}


