package jsonapi

type ResponseMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Data     interface{}     `json:"data"`
	Messages ResponseMessage `json:"messages"`
	Errors   ResponseMessage `json:"errors"`
}

func NewSuccessResponse(result interface{}) Response {
	return Response{
		Data:     result,
		Messages: ResponseMessage{},
		Errors:   ResponseMessage{},
	}
}

func NewSuccessResponseWithMessage(result interface{}, code int, message string) Response {
	return Response{
		Data:     result,
		Messages: ResponseMessage{code, message},
		Errors:   ResponseMessage{},
	}
}

func NewErrorResponse(code int, message string) Response {
	return Response{
		Data:     nil,
		Messages: ResponseMessage{},
		Errors:   ResponseMessage{code, message},
	}
}
