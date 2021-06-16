package errors

var (
	RefreshTokenEmpty       = newBadRequest("refresh_token_empty", "The refresh token should not be empty")
	UserNameEmpty           = newBadRequest("auth_username_invalid", "The username should not be empty")
	PasswordEmpty           = newBadRequest("auth_password_empty", "The password should not be empty")
	IncorrectCredentials    = newBadRequest("auth_incorrect_credentials", "Incorrect credentials")
	UsernameAlreadyAssigned = newBadRequest("auth_username_already_assigned", "Username is already assigned")
	TokenNotSet             = newUnauthorizedRequest("token_not_set", "Token is not set")
	TokenInvalid            = newUnauthorizedRequest("token_invalid", "Token invalid")
)
