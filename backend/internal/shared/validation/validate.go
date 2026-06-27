package validation

func ValidateStruct(data interface{}) error {
	return Validator().Struct(data)
}
