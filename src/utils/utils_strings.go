package utils

func IsEmpty(s string) bool {
	return s == ""
}

func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}
