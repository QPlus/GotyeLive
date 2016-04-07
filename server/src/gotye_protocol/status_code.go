package gotye_protocol

const (
	//common
	API_SUCCESS               = 10000
	API_SERVER_ERROR          = 10001
	API_PARAM_ERROR           = 10002
	API_UNAUTHORIZED_ERROR    = 10003
	API_EXPIRED_SESSION_ERROR = 10004
	API_READ_REQDATA_ERROR    = 10005

	API_ACCOUNT_NOT_EXISTS_ERROR = 10100
	API_LOGIN_PASSWORD_ERROR     = 10101
	API_ACCOUNT_EXISTS_ERROR     = 10102
	API_PHONE_EXISTS_ERROR       = 10103
	API_PHONE_NOT_EXISTS_ERROR   = 10104
	API_EMAIL_EXISTS_ERROR       = 10105

	API_END_LIVEROOM_ERROR              = 10200
	API_LIVEROOM_ID_NOT_EXIST_ERROR     = 10201
	API_LIVEROOM_ID_INVALID_ERROR       = 10202
	API_REPECT_PASSWORD_LIVEROOM_ERROR  = 10203
	API_LIVEROOM_NOT_EXISTS_ERROR       = 10204
	API_INVALID_LIVEROOM_NAME_ERROR     = 10205
	API_INVALID_PASSWORD_LIVEROOM_ERROR = 10206

	API_DECODE_HEAD_PIC_ERROR = 10301
)

//API_READ_REQDATA_ERROR:   "read reqdata error",
//

var ApiStatus = map[int]string{
	API_SUCCESS:                     "success",
	API_SERVER_ERROR:                "server intra error",
	API_PARAM_ERROR:                 "input param error",
	API_UNAUTHORIZED_ERROR:          "unauthorized error",
	API_EXPIRED_SESSION_ERROR:       "expired session error",
	API_READ_REQDATA_ERROR:          "read reqdata error",
	API_ACCOUNT_NOT_EXISTS_ERROR:    "account not exists error",
	API_LOGIN_PASSWORD_ERROR:        "login password error",
	API_ACCOUNT_EXISTS_ERROR:        "account exists error",
	API_PHONE_EXISTS_ERROR:          "phone exists error",
	API_EMAIL_EXISTS_ERROR:          "email exists error",
	API_END_LIVEROOM_ERROR:          "end liveroom error",
	API_LIVEROOM_ID_NOT_EXIST_ERROR: "liveroom id not exist error",
	API_LIVEROOM_NOT_EXISTS_ERROR:   "liveroom not exist error",
	API_DECODE_HEAD_PIC_ERROR:       "decode head pic error",
}
