package response

type GlobalResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"response"`
}

func (gr GlobalResponse) GetGlobalResponse(status int, message string, data interface{}) GlobalResponse {
	gr.Message = message
	gr.Status = status
	gr.Data = data
	return gr
}
