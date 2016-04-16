// api.go
package gotye_sdk

var (
	HttpDomain           = "https://livevip.com.cn/liveApi/"
	HttpUrlAccessToken   = HttpDomain + "AccessToken"
	HttpUrlCreateRoom    = HttpDomain + "CreateRoom"
	HttpUrlModifyRoom    = HttpDomain + "ModifyRoom"
	HttpGetLiveContext   = HttpDomain + "GetLiveContext"
	HttpGetRoomsLiveInfo = HttpDomain + "GetRoomsLiveInfo"
)

type AccessTokenAppRequest struct {
	Scope    string `json:"scope"`    //scope = "app"
	UserName string `json:"username"` //developer account
	Password string `json:"password"` //developer password
}

type AccessTokenRoomRequest struct {
	Scope     string `json:"scope"`     //scope = "room"
	RoomId    int64  `json:"roomId"`    //liveroom_id
	Password  string `json:"password"`  //liveroom password
	SecretKey string `json:"secretKey"` //md5(roomId+password+accessSecret)
	NickName  string `json:"nickName"`  //
}

type AccessTokenResponse struct {
	ExpiresIn   int64  `json:"expiresIn"`   //有效时间,单位(秒)
	AccessToken string `json:"accessToken"` //返回的token
	Role        int    `json:"role"`        //角色: 1-后台用户, 2-主播端, 3-助理, 4-观众
}

type ApiResponse struct {
	AccessPath string `json:"accessPath"`
	Runtime    int    `json:"runtime"`
	Systime    int64  `json:"systime"`
	Status     int    `json:"status"`
}

type CreateRoomRequest struct {
	RoomName    string `json:"roomName"`
	AnchorPwd   string `json:"anchorPwd"`   //主播登录密码
	AssistPwd   string `json:"assistPwd"`   //助理登录密码
	UserPwd     string `json:"userPwd"`     //用户登录密码
	AnchorDesc  string `json:"anchorDesc"`  //"主播描述"
	ContentDesc string `json:"contentDesc"` //"内容描述"
	ThirdRoomId int64  `json:"thirdRoomID"` //是否用第三方roomID
}

type CreateRoomResponse struct {
	ApiResponse
	Entity CreateRoomEntity `json:"entity"`
}

type CreateRoomEntity struct {
	RoomId      int64  `json:"roomId"`      //live room id
	AppUserId   int    `json:"appUserId"`   //app  id
	RoomName    string `json:"roomName"`    //主播室名称
	AnchorPwd   string `json:"anchorPwd"`   //主播密码
	UserPwd     string `json:"userPwd"`     //观众密码
	AssistPwd   string `json:"assistPwd"`   //助理密码
	AnchorDesc  string `json:"anchorDesc"`  //主播室描述
	ContentDesc string `json:"contentDesc"` //演讲主题
	DateCreate  int64  `json:"dateCreate"`  //创建时间
	ThirdRoomId int64  `json:"thirdRoomId"`
}

type ModifyRoomRequest struct {
	RoomId            int64  `json:"roomId"`
	RoomName          string `json:"roomName"`
	EnableRecordFlag  int    `json:"enableRecordFlag"`
	PermanentPlayFlag int    `json:"permanentPlayFlag"`
	StartPlayTime     int    `json:"startPlayTime"`
	AnchorPwd         string `json:"anchorPwd"`
	AssistPwd         string `json:"assistPwd"`
	UserPwd           string `json:"userPwd"`
	AnchorDesc        string `json:"anchorDesc"`
	ContentDesc       string `json:"contentDesc"`
}

type ModifyRoomResponse struct {
	ApiResponse
	AffectedRows int `json:"affectedRows"`
}

type GetLiveContextRequest struct {
	RoomId int64 `json:"roomId"`
}

type GetLiveContextResponse struct {
	ApiResponse
	Entity LiveContext `json:"entity"`
}

type LiveContext struct {
	RecordingStatus int `json:"recordingStatus"`
	PlayUserCount   int `josn:"playUserCount"`
}

type GetRoomsLiveInfoRequest struct {
	RoomIds []int64 `json:"roomIds"`
}

type GetRoomsLiveInfoResponse struct {
	ApiResponse
	Entities []LiveRoomInfo `json:"entities"`
}

type LiveRoomInfo struct {
	RoomId        int64    `json:"roomId"`
	PlayUserCount int      `json:"playUserCount"`
	StreamStatus  int      `json:"streamStatus"`
	PlayRtmpUrls  []string `json:"playRtmpUrls"`
	PlayHlsUrls   []string `json:"playHlsUrls"`
	PlayFlvUrls   []string `json:"playFlvUrls"`
}
