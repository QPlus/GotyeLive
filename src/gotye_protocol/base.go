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
	r.Desc = fmt.Sprintf("%s %s", ApiStatus[r.Status], val)
}

//func CheckAuthValid(sessionId, accessToken string, vResponse *ApiResponse) (userId int, valid bool) {
//	if sessionId == "" {
//		vResponse.SetFormatCode(API_PARAM_ERROR, "session id is empty")
//		return
//	}

//	if accessToken == "" {
//		vResponse.SetFormatCode(API_PARAM_ERROR, "access token is empty")
//		return
//	}

//	if !utils.IsAccessTokenValid(sessionId, accessToken) {
//		vResponse.SetCode(API_UNAUTHORIZED_ERROR)
//		return
//	}

//	gUserId, gErr := model.GetSession(sessionId)
//	if gErr != nil {
//		vResponse.SetCode(API_SESSION_EXPIRED_ERROR)
//		return
//	}

//	userId = gUserId
//	valid = true
//	return
//}
