package services

import (
	"context"
	"database/sql"
	db "echo/db/repository"
	"echo/tui/messages"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

var hashingCost int = 13

type UserService struct {
	userRepo db.UserRepository
}

func NewUserService(repo db.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (us UserService) SignIn(username, password string) tea.Msg {
	user, err := us.userRepo.GetUserByUsername(context.Background(), username)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return messages.AuthFailedMsg{
				Reason:      "We don't know u yet :/, Sign-Up to fix that",
				DebugReason: "database error: " + err.Error(),
			}
		}

		return messages.AuthFailedMsg{
			Reason:      "Oops, something went wrong. Try again ^^",
			DebugReason: "database error: " + err.Error(),
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return messages.AuthFailedMsg{
			Reason:      "Wrong username or password",
			DebugReason: "password verification failed for some reason: " + err.Error(),
		}
	}

	return messages.AuthSuccessMsg{
		User: user,
	}
}

func (us UserService) SignUp(username string, password string) tea.Msg {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashingCost)
	if err != nil {
		return messages.AuthFailedMsg{
			Reason:      "Oops, something went wrong. Try again ^^",
			DebugReason: "failed to hash the sighup password for some odd reason: " + err.Error(),
		}
	}

	log.Debugf("username: %s, hashedPassword: %s", username, string(hashedPassword))
	user, err := us.userRepo.CreateUser(context.Background(), username, string(hashedPassword))

	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// PostgreSQL error code for unique_violation is "23505"
			if pgErr.Code == "23505" {

				return messages.AuthFailedMsg{
					Reason:      "Username already taken unfortunately :/, or Sign-In if u already have an account",
					DebugReason: "database error (unique constraint violation): " + err.Error(),
				}
			}
		}

		return messages.AuthFailedMsg{
			Reason:      "Oops, something went wrong. Try again ^^",
			DebugReason: "database error: " + err.Error(),
		}
	}

	return messages.AuthSuccessMsg{
		User: user,
	}
}
