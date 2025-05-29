package auth

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
	Name     string `json:"name" validate:"required,gte=1"`
	Age      int    `json:"age" validate:"required,gte=1, lte=120"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
