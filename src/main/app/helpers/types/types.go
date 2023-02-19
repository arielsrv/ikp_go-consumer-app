package types

func String(value string) *string {
	return &value
}
func StringValue(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}
