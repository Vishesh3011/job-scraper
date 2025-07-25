package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"job-scraper.go/internal/models"
	"job-scraper.go/internal/repository"
	"job-scraper.go/internal/types"
	"job-scraper.go/internal/utils"
	"log/slog"
)

type UserService interface {
	CreateUser(*models.UserInput) (*models.User, error)
	UpdateUser(*models.UserInput) (*models.User, error)
	GetUserByEmail(string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}

type userService struct {
	context.Context
	*repository.Queries
	key    string
	logger *slog.Logger
}

func newUserService(ctx context.Context, q *repository.Queries, key string, logger *slog.Logger) UserService {
	return userService{
		Context: ctx,
		Queries: q,
		key:     key,
		logger:  logger,
	}
}

func (u userService) CreateUser(ui *models.UserInput) (*models.User, error) {
	cookie, err := utils.EncryptStr(ui.Cookie, u.key)
	if err != nil {
		u.logger.Error(utils.PrepareLogMsg("Failed to encrypt cookie"))
		return nil, err
	}

	token, err := utils.EncryptStr(ui.CsrfToken, u.key)
	if err != nil {
		u.logger.Error(utils.PrepareLogMsg("Failed to encrypt CSRF token"))
		return nil, err
	}

	user := models.NewUser(ui.Name, ui.Email, ui.Locations, ui.Keywords, cookie, token)
	if err := u.Queries.CreateUser(u.Context, user.ToCreateUserParam(token, cookie)); err != nil {
		u.logger.Error(utils.PrepareLogMsg("Failed to create user in database"), slog.Any("error", err))
		return nil, err
	}

	jsu, err := u.Queries.GetUserByID(u.Context, user.Id)
	if err != nil {
		u.logger.Error(utils.PrepareLogMsg("Failed to retrieve user after creation"), slog.Any("error", err))
		return nil, err
	}

	return models.ToUser(&jsu), nil
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

	cookie, err := utils.EncryptStr(ui.Cookie, u.key)
	if err != nil {
		return nil, err
	}

	token, err := utils.EncryptStr(ui.CsrfToken, u.key)
	if err != nil {
		return nil, err
	}

	if err := u.Queries.UpdateUser(u.Context, repository.UpdateUserParams{
		Name:      ui.Name,
		Email:     utils.ToSQLNullStr(ui.Email),
		Location:  loc,
		Keywords:  keywords,
		Cookie:    cookie,
		CsrfToken: token,
	}); err != nil {
		return nil, err
	}

	jsu, err := u.Queries.GetUserByEmail(u.Context, utils.ToSQLNullStr(ui.Email))
	if err != nil {
		return nil, err
	}

	return models.ToUser(&jsu), nil
}

func (u userService) GetUserByEmail(email string) (*models.User, error) {
	jsu, err := u.Queries.GetUserByEmail(u.Context, utils.ToSQLNullStr(&email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrRecordNotFound
		}
		return nil, err
	}

	return models.ToUser(&jsu), nil
}

func (u userService) GetAllUsers() ([]models.User, error) {
	users, err := u.Queries.GetAllUsers(u.Context)
	if err != nil {
		return nil, err
	}
	var result []models.User
	for _, user := range users {
		result = append(result, *models.ToUser(&user))
	}
	return result, nil
}
