package dto

type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	Token string `json:"token"`
	User  string `json:"user"`
}
