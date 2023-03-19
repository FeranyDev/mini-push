package util

func DefaultString(value, defaultValue string) string {
	if value != "" || value != *new(string) {
		return value
	} else {
		return defaultValue
	}
}

func DefaultInt(value, defaultValue int) int {
	if value != *new(int) {
		return value
	} else {
		return defaultValue
	}
}
