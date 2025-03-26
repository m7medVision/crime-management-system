package dto

type ErrorDTO struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}
