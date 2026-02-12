package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/MrTomatePNG/projeto-m/internal/database"
	"github.com/MrTomatePNG/projeto-m/internal/dto"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameTaken = errors.New("username already exists")
	ErrEmailTaken    = errors.New("email already exists")
)

type UserService struct {
	db  *database.Queries
	ctx context.Context
}

func NewUserService(db *database.Queries) *UserService {
	ctx := context.Background()
	return &UserService{
		db:  db,
		ctx: ctx,
	}
}

func (s *UserService) RegisterUser(username, email, password string) (*dto.UserDto, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("unreachble password: %s", err.Error())
	}

	_, err = s.db.GetUserByEmail(s.ctx, email)
	if err == nil {
		return nil, fmt.Errorf("esse email ja existe")
	}

	user, err := s.db.CreateUser(s.ctx, database.CreateUserParams{
		Username: username,
		Email:    email,
		Password: string(hashPassword),
	})

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Printf("PgError ConstraintName: %s\n", pgErr.ConstraintName)

			if pgErr.ConstraintName == "users_username_key" {
				return nil, ErrUsernameTaken
			}
			if pgErr.ConstraintName == "users_email_key" {
				return nil, ErrEmailTaken
			}
		}
		return nil, err
	}

	return &dto.UserDto{
		ID:       user.ID,
		Username: user.Username,
		Bio:      user.Bio.String,
	}, nil
}

func (s *UserService) GetUserByEmail(email string) (*database.User, error) {
	user, err := s.db.GetUserByEmail(s.ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("email user not found: %s", email)
		}
		return nil, fmt.Errorf("internal server error")
	}

	return &user, nil
}
