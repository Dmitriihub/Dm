package handlers

import (
	"context"
	"fmt"

	"newproject/internal/taskService"
	"newproject/internal/web/tasks"
	// Импортируем наш сервис
)

type Handler struct {
	Service *taskService.TaskService
}

// DeleteTasksId implements tasks.StrictServerInterface.
func (h *Handler) DeleteTasksId(ctx context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	id := request.Id // Извлекаем ID из запроса

	// Удаляем задачу через сервис
	err := h.Service.DeleteTaskByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("не удалось найти задачу с id %d: %w", id, err)
	}

	return tasks.DeleteTasksId204Response{}, nil // Возвращаем успешный ответ
}

// PatchTasksId implements tasks.StrictServerInterface.
func (h *Handler) PatchTasksId(ctx context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	id := request.Id            // Извлекаем ID из запроса
	taskRequest := request.Body // Получаем тело запроса

	taskToUpdate := taskService.Task{
		Task:   taskRequest.Task,
		IsDone: taskRequest.IsDone,
	}

	updatedTask, err := h.Service.UpdateTaskByID(uint(id), taskToUpdate)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти задачу с id %d: %w", id, err)
	}

	response := tasks.PatchTasksId200JSONResponse{
		Id:     int64(updatedTask.ID),
		Task:   updatedTask.Task,
		IsDone: updatedTask.IsDone,
	}

	return response, nil
}

// GetTasks implements tasks.StrictServerInterface.
func (h *Handler) GetTasks(ctx context.Context, request tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	response := tasks.GetTasks200JSONResponse{}
	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:     int64(tsk.ID),
			Task:   tsk.Task,
			IsDone: tsk.IsDone,
		}
		response = append(response, task)
	}

	return response, nil
}

// PostTasks implements tasks.StrictServerInterface.
func (h *Handler) PostTasks(ctx context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	// Логируем тело запроса
	fmt.Printf("Request body: %+v\n", request.Body)

	// Проверяем, что поля Task и IsDone не пустые
	if request.Body.Task == "" || request.Body.IsDone {
		return nil, fmt.Errorf("task and isDone fields are required")
	}

	// Создаем задачу для создания
	taskToCreate := taskService.Task{
		Task:   request.Body.Task,   // Здесь мы уже работаем с обычным значением, не указателем
		IsDone: request.Body.IsDone, // Убедимся, что IsDone не nil
	}

	// Создаем задачу через сервис
	createdTask, err := h.Service.CreateTask(taskToCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	// Создаем ответ
	response := tasks.PostTasks201JSONResponse{
		Id:     int64(createdTask.ID), // Используем int, так как в структуре PostTasksResponseObject должно быть поле типа int
		Task:   createdTask.Task,
		IsDone: createdTask.IsDone,
	}

	return response, nil
}

// Нужна для создания структуры Handler на этапе инициализации приложения

func NewHandler(service *taskService.TaskService) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) GetTask(_ context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	// Получение всех задач из сервиса
	allTasks, err := h.Service.GetAllTasks()
	if err != nil {
		return nil, err
	}

	// Создаем переменную респон типа 200джейсонРеспонс
	// Которую мы потом передадим в качестве ответа
	response := tasks.GetTasks200JSONResponse{}

	// Заполняем слайс response всеми задачами из БД
	for _, tsk := range allTasks {
		task := tasks.Task{
			Id:     int64(tsk.ID),
			Task:   tsk.Task,
			IsDone: tsk.IsDone,
		}
		response = append(response, task)
	}

	// САМОЕ ПРЕКРАСНОЕ. Возвращаем просто респонс и nil!
	return response, nil
}

func (h *Handler) PostTask(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	// Распаковываем тело запроса напрямую, без декодера!
	taskRequest := request.Body
	// Обращаемся к сервису и создаем задачу
	taskToCreate := taskService.Task{
		Task:   taskRequest.Task,
		IsDone: taskRequest.IsDone,
	}
	createdTask, err := h.Service.CreateTask(taskToCreate)

	if err != nil {
		return nil, err
	}
	// создаем структуру респонс
	response := tasks.PostTasks201JSONResponse{
		Id:     int64(createdTask.ID),
		Task:   createdTask.Task,
		IsDone: createdTask.IsDone,
	}
	// Просто возвращаем респонс!
	return response, nil
}

func (h *Handler) PatchTaskId(ctx context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	id := request.Id            // Извлекаем ID из запроса
	taskRequest := request.Body // Получаем тело запроса

	taskToUpdate := taskService.Task{
		Task:   taskRequest.Task,
		IsDone: taskRequest.IsDone,
	}

	// Здесь нужно привести id к int64
	updatedTask, err := h.Service.UpdateTaskByID(uint(id), taskToUpdate)
	if err != nil {
		return nil, fmt.Errorf("не удалось найти задачу с id %d: %w", id, err)
	}

	response := tasks.PatchTasksId200JSONResponse{
		Id:     int64(updatedTask.ID), // Обратите внимание, что здесь `updatedTask.ID` теперь должен быть int64
		Task:   updatedTask.Task,
		IsDone: updatedTask.IsDone,
	}

	return response, nil
}

func (h *Handler) DeleteTaskId(ctx context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	// Извлекаем ID из запроса
	id := request.Id

	// Удаляем задачу через сервис
	err := h.Service.DeleteTaskByID(uint(id))
	if err != nil {
		return nil, fmt.Errorf("не удалось найти задачу с id %d: %w", id, err)
	}

	// Возвращаем успешный ответ
	return tasks.DeleteTasksId204Response{}, nil
}
