package errors

var (
	PathKeyInvalid = newBadRequest("path_key_invalid", "The path key is invalid")
)
