package dto

type RegisterUser struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=6"`
	Name     string `json:"name" xml:"name" form:"name" validate:"required"`
}
