package lib

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func ResponseSuccess(data interface{}, envelopes ...string) Response {
	if len(envelopes) > 0 {
		data = map[string]any{
			envelopes[0]: data,
		}
	}

	return Response{
		Status: "success",
		Data:   data,
	}
}

func ResponseError(err error) Response {
	return Response{
		Status:  "error",
		Message: err.Error(),
	}
}
