package errors

var (
	TipNameInvalid = NewBadRequest("tip_name_invalid", "The tip name is invalid")
)
