package gotye_protocol

import (
	"fmt"
)

//All API respone must inherit this struct.
type ApiResponse struct {
	Access string `json:"access"`
	Status int    `json:"status"`
	Desc   string `json:"desc"`
}

func (r *ApiResponse) SetAccess(access string) {
	r.Access = access
}

func (r *ApiResponse) SetStatus(status int) {
	r.Status = status
	r.Desc = ApiStatus[r.Status]
}

func (r *ApiResponse) SetFormatStatus(status int, val string) {
	r.Status = status
	r.Desc = fmt.Sprintf("%s", val)
}
