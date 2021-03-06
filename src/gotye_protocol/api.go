package gotye_protocol

//http://127.0.0.1:8080/live/login

//api : /live/Login
/*
{
    "account" : "zhangsan"/"13512324432",
    "password": "123456"
}
*/
type LoginRequest struct {
	Account string `json:"account"`
	Passwd  string `json:"password"`
}

/*
{
    "access"   : "/live/Login"
    "status"   :
    "desc"     :
    "account"  : "zhangsan"
    "nickName" : "aaaaaa"
    "liveRoomId" : //如果这个用户有liveRoomId,就返回，如果没有返回0
    "sessionId": "sasassasaasasasasas" (32个字节的字符串)
}
*/
type LoginResponse struct {
	ApiResponse
	NickName   string `json:"nickName"`
	LiveRoomID int64  `json:"liveRoomId"`
	SessionID  string `json:"sessionId"`
	HeadPicId  int64  `json:"headPicId"`
	Sex        int8   `json:"sex"` //1:male, 2: female
}

//api : /live/AuthCode
/*
{
    "phone"    : "13512023289",
}
*/
type AuthCodeRequest struct {
	Phone string `json:"phone"`
}

/*
{
    "access"   : "/live/AuthCode"
    "status"   :
    "desc"     :
}
*/
type AuthCodeResponse struct {
	ApiResponse
}

//api : /live/Register
/*
{
    "phone"    : "13512023289",
    "password" : "123456"
    "authCode" : "",
}
*/
type RegisterRequest struct {
	Phone    string `json:"phone"`
	Passwd   string `json:"password"`
	AuthCode string `json:"authCode"`
}

/*
{
    "access"   : "/live/Register"
    "status"   :
    "desc"     :
}
*/
type RegisterResponse struct {
	ApiResponse
}

//api : live/ModifyUserPwd
/*
{
    "phone" : "",
    "password": "123456"
}
*/
type ModifyUserPwdRequest struct {
	Phone  string `json:"phone"`
	Passwd string `json:"password"`
}

/*
{
    "access"   : "/live/ModifyUserPwd"
}
*/
type ModifyUserPwdResponse struct {
	ApiResponse
}

//api : /live/ModifyUserInfo
/*
{
    "sessionId" : "asssssssssssssssssss"
	"nickName"  : 修改用户昵称
}
*/
type ModifyUserInfoRequest struct {
	SessionID string `json:"sessionId"`
	NickName  string `json:"nickName"`
	Sex       int8   `json:"sex"` //1:male, 2: female
	Address   string `json:"address"`
}

type ModifyUserInfoResponse struct {
	ApiResponse
}

//api : /live/ModifyUserHeadPic
/*
{
    "sessionId" : "asssssssssssssssssss"
	"headPic"   : 修改用户头像 (图片base64之后成字符串再传入)
}
*/
type ModifyUserHeadPicRequest struct {
	SessionID string `json:"sessionId"`
	HeadPic   string `json:"headPic"`
}

type ModifyUserHeadPicResponse struct {
	ApiResponse
	HeadPicId int64 `json:"headPicId"`
}

//api : /live/GetUserHeadPic
type GetUserHeadPicResponse struct {
	ApiResponse
}

//api :　/live/GetLiveRoomList
const (
	ALL_LIVE_ROOM_LIST = iota + 1
	FOCUS_LIVE_ROOM_LIST
)

/*
{
    "sessionId"     : "asssssssssssssssssss"
    "type"    : 1(全部), 2(关注)
    "refresh" : 1(刷新), 0(获取下一页)
    "count"   : 可以填充一次性刷新几个，如果不填，默认是5个
}
*/
type GetLiveRoomListRequest struct {
	SessionID string `json:"sessionId"`
	Type      int    `json:"type"`
	Refresh   int    `josn:"refresh"`
	Count     int    `json:"count"`
}

/*
 {
     "access"   : "/live/GetLiveRoomList"
     "status"   :
     "desc"     :
     "type"     : 1(全部), 2(关注)
     "lastId"   : 这次获取到的节点
     "list"     : [
                {
                    "liveRoomId" :
                    "liveRoomanchorPwd": 主播密码
    	            "liveRoomUserPwd": 观看直播的用户密码
    	            "liveRoomName" : 直播室名称
    	            "liveRoomDesc" : 直播室描述
    	            "liveRoomTopic": 演讲的题目
    	            "anchorName"   : 主播昵称
    	            "anchorIcon"   : 主播头像
    	            "followCount"  : 被关注量
                }
                ]
 }
*/
type GetAllLiveRoomListResponse struct {
	ApiResponse
	Type int            `json:"Type"`
	List []LiveRoomInfo `json:"list"`
}

/*
 {
     "access"   : "/live/GetLiveRoomList"
     "status"   :
     "desc"     :
     "type"     : 1(全部), 2(关注)
     "lastId"   : 这次获取到的节点
     "onlineList"     : [
                {
                    "liveRoomId" :
                    "liveRoomanchorPwd": 主播密码
    	            "liveRoomUserPwd": 观看直播的用户密码
    	            "liveRoomName" : 直播室名称
    	            "liveRoomDesc" : 直播室描述
    	            "liveRoomTopic": 演讲的题目
    	            "anchorName"   : 主播昵称
    	            "anchorIcon"   : 主播头像
    	            "followCount"  : 被关注量
                }
                ],
    "offlineList"     : [
    {
        "liveRoomId" :
        "liveRoomanchorPwd": 主播密码
        "liveRoomUserPwd": 观看直播的用户密码
        "liveRoomName" : 直播室名称
        "liveRoomDesc" : 直播室描述
        "liveRoomTopic": 演讲的题目
        "anchorName"   : 主播昵称
        "anchorIcon"   : 主播头像
        "followCount"  : 被关注量
    }
    ]
 }
*/
type GetFcousLiveRoomListResponse struct {
	ApiResponse
	Type        int            `json:"Type"`
	OnlineList  []LiveRoomInfo `json:"onlineList"`
	OfflineList []LiveRoomInfo `json:"offlineList"`
}

type LiveRoomInfo struct {
	LiveRoomId    int64  `json:"liveRoomId"`
	LiveAnchorPwd string `json:"liveRoomAnchorPwd"` //主播密码
	LiveUserPwd   string `json:"liveRoomUserPwd"`   //观看直播的用户密码
	LiveRoomName  string `json:"liveRoomName"`
	LiveRoomDesc  string `json:"liveRoomDesc"`
	LiveRoomTopic string `json:"liveRoomTopic"`
	AnchorName    string `json:"anchorName"`
	HeadPicId     int64  `json:"headPicId"`
	IsFollow      int8   `json:"isFollow"` //1 : 关注, 0: 未关注
	FollowCount   int    `json:"followCount"`
	PlayerCount   int    `json:"playerCount"`
	PlayRtmpUrl   string `json:"playRtmpUrl"`
	PlayHlsUrl    string `json:"playHlsUrl"`
	PlayFlvUrl    string `json:"playFlvUrl"`
}

//api : /live/GetMyLiveRoomId
/*
{
    "sessionId" : ""
}
*/
type GetMyLiveRoomIdRequest struct {
	SessionID string `json:"sessionId"`
}

/*
{
    "access"   : " /live/GetMyLiveRoomId"
    "status"   :
    "desc"     :
    "LiveRoomId": 1234567
}
*/
type GetMyLiveRoomIdResponse struct {
	ApiResponse
	LiveRoomId int64 `json:"LiveRoomId"`
}

//api : /live/GetMyLiveRoom
/*
{
    "sessionId" : ""
}
*/
type GetMyLiveRoomRequest struct {
	SessionID string `json:"sessionId"`
}

/*
{
    "access"   : " /live/GetMyLiveRoom"
    "status"   :
    "desc"     :
    "liveRoomId" :
    "liveRoomanchorPwd": 主播密码
    "liveRoomUserPwd": 观看直播的用户密码
    "liveRoomName" : 直播室名称
    "liveRoomDesc" : 直播室描述
    "liveRoomTopic": 演讲的题目
    "anchorName"   : 主播昵称
    "anchorIcon"   : 主播头像
    "followCount"   : 被关注量

}
*/
type GetMyLiveRoomResponse struct {
	ApiResponse
	LiveRoomInfo
}

//api : /live/ModifyMyLiveRoom
/*
{
	"sessionId"  :
	"liveRoomId" :        可为0,
	"liveRoomAnchorPwd" : 主播密码
	"liveRoomUserPwd"   : 观看直播的用户密码
	"liveRoomName"      :
	"liveRoomDesc"      :
	"liveRoomTopic"     :
}
*/
type ModifyMyLiveRoomRequest struct {
	SessionID         string `json:"sessionId"`
	LiveRoomID        int64  `json:"liveRoomId"`        //可为0,
	LiveRoomAnchorPwd string `json:"liveRoomAnchorPwd"` //主播密码
	LiveUserPwd       string `json:"liveRoomUserPwd"`   //观看直播的用户密码
	LiveRoomName      string `json:"liveRoomName"`
	LiveRoomDesc      string `json:"liveRoomDesc"`
	LiveRoomTopic     string `json:"liveRoomTopic"`
}

/*
{
    "access"   : " /live/ModifyMyLiveRoom"
    "status"   :
}
*/
type ModifyMyLiveRoomResponse struct {
	ApiResponse
}

//api : /live/CreateLiveRoom
/*
{
    "sessionId"    :
	"liveRoomAnchorPwd" : 主播密码
    "liveRoomAssistPwd" : 助理密码
	"liveRoomUserPwd"   : 观看密码
	"liveRoomName"      : 直播室名称
	"liveRoomDesc"      : 可不填
	"liveRoomTopic"     : 可不填

}
*/
type CreateLiveRoomRequest struct {
	SessionID     string `json:"sessionId"`
	LiveAnchorPwd string `json:"liveRoomAnchorPwd"` //主播密码
	LiveAssistPwd string `json:"liveRoomAssistPwd"` //助理密码
	LiveUserPwd   string `json:"liveRoomUserPwd"`   //观看密码
	LiveRoomName  string `json:"liveRoomName"`
	LiveRoomDesc  string `json:"liveRoomDesc"`
	LiveRoomTopic string `json:"liveRoomTopic"`
}

/*
{
    "access"    : " /live/CreateLiveRoom"
    "status"    :
    "liveRoomId":
}
*/
type CreateLiveRoomResponse struct {
	ApiResponse
	LiveRoomId int64 `json:"liveRoomId"`
}

//api : /live/FollowLiveRoom
/*
{
    "sessionId"  :
	"liveRoomId" :
    "isFollow"   : 1 关注, 0 取消

}
*/
type FollowLiveRoomRequest struct {
	SessionId  string `json:"sessionId"`
	LiveRoomId int64  `json:"liveRoomId"`
	IsFollow   int    `json:"isFollow"`
}

/*
{
    "access"    : " /live/CreateLiveRoom"
    "status"    :
    "desc"      :
}
*/
type FollowLiveRoomResponse struct {
	ApiResponse
}

//api : /live/PushLiveStream
/*
{
	"sessionId"  :
	"liveRoomId" :
    "status"     : // 1. push, 0. stop
	"timeout"    : //超时时间,默认60秒. 如果在此时间内没有设置状态为1，直播状态变为0
}
*/
type PushLiveStreamRequest struct {
	SessionId  string `json:"sessionId"`
	LiveRoomId int64  `json:"liveRoomId"`
	Status     int    `json:"status"`  // 1. push, 0. stop
	Timeout    int    `json:"timeout"` //超时时间,默认60秒. 如果在此时间内没有设置状态为1，直播状态变为0
}

/*
{
    "access"    : " /live/PushLiveStream"
    "status"    :
    "desc"
}
*/
type PushLiveStreamResponse struct {
	ApiResponse
}

/*
api : /live/PlayLiveStream
{
    "sessionId"  :
	"liveroomId" :
	"status"     :
}
*/
type PlayLiveStreamRequest struct {
	SessionId  string `json:"sessionId"`
	LiveroomId int64  `json:"liveroomId"`
	Status     int    `json:"status"` //1. start, 0, stop
}

type PlayLiveStreamResponse struct {
	ApiResponse
}

/*
 api : /live/GetLiveroomNumber
*/
type GetLiveroomNumberRequest struct {
	SessionId  string `json:"sessionId"`
	LiveroomId int64  `json:"liveroomId"`
}

type GetLiveroomNumberResponse struct {
	ApiResponse
	LiveroomId int64 `json:"liveroomId"`
	Number     int   `json:"number"`
}

//api : /live/SearchLiveStream
type SearchLiveStreamRequest struct {
	SessionId string `json:"sessionId"`
	Keyword   string `json:"keyword"`
}

type SearchLiveStreamResponse struct {
	ApiResponse
	LiveRoomInfo
	IsPlay int8 `json:"isPlay"` //1 : 直播, 0: 未直播
}

//api : /live/GetAdLiveRoom
type GetAdLiveRoomListRequest struct {
	SessionId string `json:"sessionId"`
}

type GetAdLiveRoomListResponse struct {
	ApiResponse
	LiveRoomInfo
}

/*
 api : /pay/ChargeRMB
*/
type ChargeRMBRequest struct {
	SessionId string `json:"sessionId"`
	RMB       int    `json:"rmb"`
}

type ChargeRMBResponse struct {
	ApiResponse
	QinCoin int `json:"qinCoin"`
}

/*
 api : /pay/PayQinCoin
*/
type PayQinCoinRequest struct {
	SessionId     string `json:"sessionId"`
	QinCoin       int    `json:"qinCoin"` //支付亲元
	AnchorAccount string `json:"anchorAccount"`
}

type PayQinCoinResponse struct {
	ApiResponse
	QinCoin int `json:"qinCoin"` //剩余多少亲元
}

/*
 api : /pay/GetPayAccount
*/
type GetPayAccountRequest struct {
	SessionId string `json:"sessionId"`
}

type GetPayAccountResponse struct {
	ApiResponse
	QinCoin int `json:"qinCoin"` //剩余亲元
	JiaCoin int `json:"jiaCoin"` //收入加元
	Level   int `json:"level"`   //用户等级
	XP      int `json:"xp"`      //用户经验值
}
