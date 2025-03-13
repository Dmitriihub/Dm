package handlers

import (
	"context"
	"fmt"
	"newproject/internal/userService"
	"newproject/internal/web/users"
)

type UserHandler struct {
	Service *userService.UserService
}

// DeleteUsersId implements users.StrictServerInterface.
func (h *UserHandler) DeleteUsersId(ctx context.Context, request users.DeleteUsersIdRequestObject) (users.DeleteUsersIdResponseObject, error) {
	id := request.Id

	err := h.Service.DeleteUserByID(uint(id))
	if err != nil {
		fmt.Println("Ошибка удаления пользователя:", err)
		return nil, fmt.Errorf("не удалось удалить пользователя с id %d: %w", id, err)
	}

	return users.DeleteUsersId204Response{}, nil
}

// PatchUsersId implements users.StrictServerInterface.
func (h *UserHandler) PatchUsersId(ctx context.Context, request users.PatchUsersIdRequestObject) (users.PatchUsersIdResponseObject, error) {
	id := request.Id
	userRequest := request.Body

	// Проверяем, что Email и Password не равны nil
	if userRequest.Email == nil && userRequest.Name == nil && userRequest.Password == nil {
		return nil, fmt.Errorf("email, name, and password must be provided at least one field to update")
	}

	// Формируем объект для обновления
	userToUpdate := userService.User{
		Email:    *userRequest.Email,
		Password: *userRequest.Password,
		Name:     *userRequest.Name, // Учитываем поле Name
	}

	// Вызов метода сервиса для обновления пользователя
	updatedUser, err := h.Service.UpdateUserByID(id, userToUpdate)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	// Формируем ответ
	response := users.PatchUsersId200JSONResponse{
		Id:       ptr(int(updatedUser.ID)),
		Email:    ptr(updatedUser.Email),
		Password: ptr(updatedUser.Password),
		Name:     ptr(updatedUser.Name), // Возвращаем Name
	}

	return response, nil
}

func NewUserHandler(service *userService.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) GetUsers(ctx context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	allUsers, err := h.Service.GetAllUsers()
	if err != nil {
		return nil, err
	}

	response := users.GetUsers200JSONResponse{}
	for _, usr := range allUsers {
		user := users.User{
			Id:       ptr(int(usr.ID)),
			Email:    ptr(usr.Email),
			Password: ptr(usr.Password),
			Name:     ptr(usr.Name), // Добавляем Name в ответ
		}
		response = append(response, user)
	}

	return response, nil
}

func ptr[T any](v T) *T {
	return &v
}

func (h *UserHandler) PostUsers(ctx context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	userRequest := request.Body

	// Проверяем, что Email и Password не равны nil
	if userRequest.Email == nil || userRequest.Password == nil || userRequest.Name == nil {
		return nil, fmt.Errorf("email, name, and password must be provided")
	}

	newUser := userService.User{
		Email:    *userRequest.Email,
		Password: *userRequest.Password,
		Name:     *userRequest.Name, // Передаем Name при создании
	}

	createdUser, err := h.Service.CreateUser(newUser)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	response := users.PostUsers201JSONResponse{
		Id:       ptr(int(createdUser.ID)),
		Email:    ptr(createdUser.Email),
		Password: ptr(createdUser.Password),
		Name:     ptr(createdUser.Name), // Возвращаем Name в ответе
	}

	return response, nil
}
