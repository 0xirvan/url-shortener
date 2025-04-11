package validation

type Register struct {
	Name     string `validate:"required,min=3,max=100"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=100,password"`
}

type Login struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=100"`
}

type Logout struct {
	RefreshToken string `validate:"required"`
}

type ForgotPassword struct {
	Email string `validate:"required,email"`
}

type GoogleLogin struct {
	Name          string `validate:"required,max=100"`
	Email         string `validate:"required,email"`
	VerifiedEmail bool   `validate:"required,boolean"`
}

type RefreshToken struct {
	RefreshToken string `validate:"required"`
}

type Token struct {
	Token string `validate:"required"`
}
