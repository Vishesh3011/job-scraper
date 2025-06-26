package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/repository"
	"job-scraper.go/internal/types"
)

type UserService interface {
	CreateUser(*models.UserInput) (*models.User, error)
	UpdateUser(*models.UserInput) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
}

type userService struct {
	context.Context
	*repository.Queries
}

func newUserService(ctx context.Context, q *repository.Queries) UserService {
	return userService{
		Context: ctx,
		Queries: q,
	}
}

func (u userService) CreateUser(ui *models.UserInput) (*models.User, error) {
	loc, err := json.Marshal(ui.Locations)
	if err != nil {
		return nil, err
	}

	keywords, err := json.Marshal(ui.Keywords)
	if err != nil {
		return nil, err
	}

	if err := u.Queries.CreateUser(u.Context, repository.CreateUserParams{
		Name:     ui.Name,
		Email:    *ui.Email,
		Location: loc,
		Keywords: keywords,
	}); err != nil {
		return nil, err
	}

	jsu, err := u.Queries.GetUserByEmail(u.Context, *ui.Email)
	if err != nil {
		return nil, err
	}

	return models.NewUser(&jsu), nil
}

func (u userService) UpdateUser(ui *models.UserInput) (*models.User, error) {
	loc, err := json.Marshal(ui.Locations)
	if err != nil {
		return nil, err
	}

	keywords, err := json.Marshal(ui.Keywords)
	if err != nil {
		return nil, err
	}

	if err := u.Queries.UpdateUser(u.Context, repository.UpdateUserParams{
		Name:     ui.Name,
		Email:    *ui.Email,
		Location: loc,
		Keywords: keywords,
	}); err != nil {
		return nil, err
	}

	jsu, err := u.Queries.GetUserByEmail(u.Context, *ui.Email)
	if err != nil {
		return nil, err
	}

	return models.NewUser(&jsu), nil
}

func (u userService) GetUserByEmail(email string) (*models.User, error) {
	jsu, err := u.Queries.GetUserByEmail(u.Context, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrRecordNotFound
		}
		return nil, err
	}

	return models.NewUser(&jsu), nil
}
