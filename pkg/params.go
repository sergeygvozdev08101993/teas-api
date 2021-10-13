package pkg

var (
	errBadRequest = &Error{"bad_request", 400, "Bad request", "Request body is not well-formed."}
	errNotAcceptable = &Error{"not_acceptable", 406, "Not Acceptable", "Accept header must be set to 'application/vnd.api+json'."}
	errUnsupportedType = &Error{"unsupported_media_type", 415, "Unsupported Media Type", "Content-Type header must be set to 'application/vnd.api+json'."}
	errInternalServer = &Error{"internal_server_error", 500, "Internal Server Error", "Something went wrong"}
)

type Error struct {
	Id string `json:"id"`
	Status int `json:"status"`
	Title string `json:"title"`
	Detail string `json:"detail"`
}

type Errors struct {
	Errors []*Error `json:"errors"`
}

type Tea struct {
	Name string `json:"name"`
	Category string `json:"category"`
}

type TeaResource struct {
	Status int `json:"status"`
	Data Tea `json:"data"`
}

type TeasCollection struct {
	Data []Tea `json:"data"`
}
