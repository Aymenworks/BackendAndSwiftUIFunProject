package errors

var (
	UserNameEmpty        = newBadRequest("auth_username_invalid", "The username should not be empty")
	PasswordEmpty        = newBadRequest("auth_password_empty", "The password should not be empty")
	IncorrectCredentials = newBadRequest("auth_incorrect_credentials", "Incorrect credentials")
)
