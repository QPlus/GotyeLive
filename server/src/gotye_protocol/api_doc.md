response status code :
//common
API_SUCCESS                  = 10000
API_SERVER_ERROR             = 10001
API_PARAM_ERROR              = 10002
API_EXPIRED_SESSION_ID_ERROR = 10004    //sessionid过期,这时应该重新申请
API_LIVEROOM_NOT_EXISTS_ERROR = 10005

// api : /live/Login
API_ACCOUNT_NOT_EXISTS_ERROR = 10100
API_LOGIN_PASSWORD_ERROR     = 10101

// api : /live/Register
API_ACCOUNT_EXISTS_ERROR     = 10200  //账号已经存在
API_PHONE_EXISTS_ERROR       = 10201  //手机号已经存在
API_EMAIL_EXISTS_ERROR       = 10202  //邮箱已经存在

// api : /live/GetLiveRoomList
API_NOT_PUBLISHING_LIVEROOM_ERROR = 10300  //当前没有正在直播的用户
API_NOT_FOLLOW_LIVEROOM_ERROR = 10301  //如果你没有关注任何直播用户



以下API都是POST method

api : /live/Login
request:
{
    "account" : "账号/邮箱/手机号码",
    "password": "123456"
}

response:
{
    "access"   : "/live/Login"
    "status"   :
    "desc"     :
    "account"  : "zhangsan"
    "nickName" : "aaaaaa"
    "liveRoomId" : //如果这个用户有liveRoomId,就返回，如果没有返回0
    "sessionId": "sasassasaasasasasas" (32个字节的字符串)
}

api : /live/Register
request:
{
    "username" : "zhangsan",
    "phone"    : "13512023289",
    "email"    : "example@gotye.com.cn"
    "password" : "123456"
}

response:
{
    "access"   : "/live/Register"
    "status"   :
    "desc"     :
    "account"  : "zhangsan"
    "nickName" : "zhangsan" //刚注册成功的,nickName等于account
    "sessionId": "sasassasaasasasasas" (32个字节的字符串)
}

api : /live/modifyUserInfo
request:
{
    "sessionId" : "asssssssssssssssssss"
	"nickName"  : 修改用户昵称,如果不修改传""
	"userIcon"  : 修改用户头像,如果不修改传""
}

response:
{
    "access"   : "/live/Register"
    "status"   :
    "desc"     :
}

api :　/live/GetLiveRoomList
request
{
    "sessionId"     : "asssssssssssssssssss"
    "liveRoomType"  : 1(全部), 2(关注)
    "index"         : 从哪一个开始获取,第一次获取填0,下一次填充服务器返回值
    "count"         : 可以填充一次性刷新几个，如果不填，默认是5个
}

response
{
    "access"   : "/live/GetLiveRoomList"
    "status"   :
    "desc"     :
    "liveRoomType" : 1(全部), 2(关注)
    "index"    : 5     //获取到的位置,用于下一次获取的index
    "count"    : 获取了几个
    "total"    : 123    //如果已经获得的LiveRoom已经等于Total,就代表已经获取结束。
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

api : /live/GetMyLiveRoomId
request:
{
    "sessionId" : ""
}

response
{
    "access"   : " /live/GetMyLiveRoomId"
    "status"   : API_LIVEROOM_NOT_EXISTS_ERROR
    "desc"     :
    "LiveRoomId": 1234567
}

api : /live/GetMyLiveRoom
request :
{
    "sessionId" : ""
}

response :
{
    "access"   : " /live/GetMyLiveRoom"
    "status"   : API_LIVEROOM_NOT_EXISTS_ERROR
    "desc"     :
    "liveRoomId" :
    "liveRoomanchorPwd": 主播密码
    "liveRoomUserPwd": 观看直播的用户密码
    "liveRoomName" : 直播室名称
    "liveRoomDesc" : 直播室描述
    "liveRoomTopic": 演讲的题目
    "anchorName"   : 主播昵称
    "anchorIcon"   : 主播头像
    "focusCount"   : 被关注量

}



api : /live/ModifyMyLiveRoom
request:
{
	"sessionId"  :
	"liveRoomId" :        可为0,
	"liveRoomAnchorPwd" : 主播密码
	"liveRoomUserPwd"   : 观看直播的用户密码
	"liveRoomName"      :
	"liveRoomDesc"      :
	"liveRoomTopic"     :
}

response:
{
    "access" : " /live/ModifyMyLiveRoom"
    "status" : API_LIVEROOM_NOT_EXISTS_ERROR
    "desc"   : 
}

api : /live/CreateLiveRoom
request:
{
    "sessionId"    :
	"liveRoomanchorPwd" : 主播密码
	"liveRoomUserPwd"   : 观看直播的用户密码
	"liveRoomName"      : 直播室名称
	"liveRoomDesc"      : 可不填
	"liveRoomTopic"     : 可不填

}

response:
{
    "access"    : " /live/CreateLiveRoom"
    "status"    : API_SERVER_ERROR / API_PARAM_ERROR
    "liveRoomId":
}

api : /live/FollowLiveRoom
request:
{
    "sessionId"  :
	"liveRoomId" :
    "isFollow"   : 1 关注, 0 取消

}

response
{
    "access"    : " /live/CreateLiveRoom"
    "status"    :
    "desc"
}


api : /live/PushLiveStream
request:
{    
	"sessionId"  :
	"liveRoomId" :
    "status"     : 1在推流, 0停止推流
	"timeout"    : 超时时间,默认60秒. 如果在此时间内没有设置状态为1，直播状态变为0
}

response:
{
    "access"    : " /live/PushLiveStream"
    "status"    :
    "desc"
}
