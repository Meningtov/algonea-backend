package handler

var InternalServerError = Error{
	Message: "Internal server error",
	Code:    "INTERNAL_SERVER_ERROR",
}

var MissingQueryParam = Error{
	Message: "Missing query param",
	Code:    "MISSING_QUERY_PARAM",
}

var MissingPathParam = Error{
	Message: "Missing path param",
	Code:    "MISSING_PATH_PARAM",
}

type Error struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}
