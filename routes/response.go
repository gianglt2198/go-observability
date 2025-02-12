package routes

func ErrorResponse(err error) *map[string]interface{} {

	res := make(map[string]interface{})
	res["success"] = false
	res["error"] = err.Error()

	return &res
}

func SuccessResponse(data interface{}) *map[string]interface{} {

	res := make(map[string]interface{})
	res["success"] = true
	res["data"] = data

	return &res
}
