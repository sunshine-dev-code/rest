package response

// Response define
type RestResponse struct {
	Code int16  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func NewOkRestResponse(data any) *RestResponse {
	return &RestResponse{0, "ok", data}
}

// error response
var (
	ErrUnknowResponse         = RestResponse{10000, "unknow_error", nil}
	ErrUnauthorizedResponse   = RestResponse{10001, "unauthorized", nil}
	ErrParamsResponse         = RestResponse{10002, "parameter_error", nil}
	ErrExistentResponse       = RestResponse{10003, "repeat_parameter", nil}
	ErrRecordNotFoundResponse = RestResponse{10004, "record_not_found", nil}
	ErrPermissionDenied       = RestResponse{10005, "permission_denied", nil}
)
