package errors

var (
	TipNameInvalid = newBadRequest("tip_name_invalid", "The tip name is invalid")
)
