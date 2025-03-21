package models

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Tasks    []Task `json:"tasks"`
}

type Task struct {
	ID     uint   `json:"id"`
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
	UserID uint   `json:"user_id"`
	Title  string `json:"title"`
}

type NewTaskRequest struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
	UserId uint   `json:"user_id"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}
