package errors

const defaultErrorMessage = "An error happened"

var (
	UnknownError             = newUnknownError()
	NotFound                 = newNotFoundError()
	InvalidParameter         = newBadRequest("invalid_parameter", "Invalid parameter")
	InvalidContentTypeHeader = newBadRequest("invalid_content_type_header", "Invalid content-type header")
)
