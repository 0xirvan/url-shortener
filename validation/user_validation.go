package validation

type CreateUser struct {
	Name     string `validate:"required,min=3,max=100"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=100,password"`
	Role     string `validate:"required,oneof=user admin"`
}

type UpdateUser struct {
	Name     string `validate:"omitempty,min=3,max=100"`
	Email    string `validate:"omitempty,email"`
	Password string `validate:"omitempty,min=8,max=100,password"`
}

type UpdatePassOrVerifyUser struct {
	Password      string `validate:"omitempty,min=8,max=100,password"`
	VerifiedEmail bool   `validate:"omitempty,boolean"`
}

type QueryUser struct {
	Page   int    `validate:"omitempty,number"`
	Limit  int    `validate:"omitempty,number"`
	Search string `validate:"omitempty,max=100"`
}
