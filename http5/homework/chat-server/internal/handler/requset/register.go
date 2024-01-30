package requset

type RegisterRequest struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}
