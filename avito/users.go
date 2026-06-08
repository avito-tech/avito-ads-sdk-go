package avito

import (
	"context"
	"fmt"
)

// UserService — методы управления пользователями аккаунта.
type UserService struct{ baseService }

// List возвращает список пользователей аккаунта.
//
// GET v1/account/{accountID}/users
func (s *UserService) List(ctx context.Context) ([]User, error) {
	var out struct {
		Users []User `json:"users"`
	}
	if _, err := s.tr.do(ctx, "GET", s.accountPath("users"), nil, &out); err != nil {
		return nil, err
	}
	return out.Users, nil
}

// Add добавляет пользователя в аккаунт.
//
// POST v1/account/{accountID}/add-user
func (s *UserService) Add(ctx context.Context, userID int64, role string) error {
	if !contains(validUserRoles, role) {
		return validationErrorf("недопустимая роль пользователя: %q", role)
	}
	body := map[string]any{"userId": userID, "role": role}
	_, err := s.tr.do(ctx, "POST", s.accountPath("add-user"), body, nil)
	return err
}

// SetRole изменяет роль пользователя.
//
// POST v1/account/{accountID}/set-user-role
func (s *UserService) SetRole(ctx context.Context, userID int64, role string) error {
	if !contains(validUserRoles, role) {
		return validationErrorf("недопустимая роль пользователя: %q", role)
	}
	body := map[string]any{"userId": userID, "role": role}
	_, err := s.tr.do(ctx, "POST", s.accountPath("set-user-role"), body, nil)
	return err
}

// Delete удаляет пользователя из аккаунта.
//
// DELETE v1/account/{accountID}/delete-user/{userID}
func (s *UserService) Delete(ctx context.Context, userID int64) error {
	path := s.accountPath(fmt.Sprintf("delete-user/%d", userID))
	_, err := s.tr.do(ctx, "DELETE", path, nil, nil)
	return err
}
