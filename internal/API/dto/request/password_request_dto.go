package request

type UpdatePasswordDTO struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type AdminUpdatePasswordDTO struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
