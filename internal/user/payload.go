package user

import "time"

type GetByIDResponse struct {
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"gte=1,lte=255"`
	Age  int    `json:"age" validate:"gte=1,lte=120"`
}
