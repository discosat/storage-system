package dam

import "reflect"

func FilterOutEmptyFields(req ImageRequest) map[string]interface{} {
	filteredResult := make(map[string]interface{})
	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)

	for i := 0; i < val.NumField(); i++ {
		formVal := val.Field(i)
		formType := typ.Field(i).Tag.Get("form")

		if formType == "" {
			formType = typ.Field(i).Name
		}

		if formVal.Kind() == reflect.Ptr && !formVal.IsNil() {
			filteredResult[formType] = formVal.Elem().Interface()
		} else if formVal.Kind() != reflect.Ptr && !formVal.IsZero() {
			filteredResult[formType] = formVal.Interface()
		}
	}
	return filteredResult
}
