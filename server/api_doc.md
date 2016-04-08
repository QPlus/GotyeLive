#QPlus Live Server API 文档

##返回码列表

###请求相关

|返回码                   |描述               |值   |备注        |
|-------------------------|-------------------|-----|------------|
|API_SUCCESS              |"成功"             |10000|请求成功    |
|API_SERVER_ERROR         |"系统异常"         |10001|几乎不会发生|
|API_PARAM_ERROR	      |"参数错误"         |10002|几乎不会发生|
|API_EXPIRED_SESSION_ERROR|"登录过期,重新登录"|10003|离最后一次访问服务器超过一个小时过期|

###回复相关

####账号相关

|返回码                      |描述          |值   |备注      |
|----------------------------|--------------|-----|----------|
|API_ACCOUNT_NOT_EXISTS_ERROR|"账号不存在"  |10100|登录时    |
|API_LOGIN_PASSWORD_ERROR    |"登录密码错误"|10101|登录时    |
|API_ACCOUNT_EXISTS_ERROR	 |"账号已注册"	|10102|注册时    |
|API_PHONE_EXISTS_ERROR		 |"手机已注册"  |10103|注册时    |
|API_PHONE_NOT_EXISTS_ERROR  |"手机不存在"  |10104|修改密码时|
|API_EMAIL_EXISTS_ERROR      |"邮箱已注册"  |10105|注册时    |
|API_EMAIL_NOT_EXISTS_ERROR  |"邮箱不存在"	|10106|暂时未用  |

####直播室相关

|返回码                             |描述            |值   |备注               |
|-----------------------------------|----------------|-----|-------------------|
|API_LIVEROOM_ID_NOT_EXIST_ERROR    |"直播室ID不存在"|10200|推流,修改,搜索时   |
|API_REPECT_PASSWORD_LIVEROOM_ERROR |"直播室密码重复"|10201|修改直播室,不会发生|
|API_LIVEROOM_NOT_EXISTS_ERROR      |"直播室不存在"  |10202|获取自己聊天室时   |
|API_INVALID_LIVEROOM_NAME_ERROR    |"直播室名称非法"|10203|修改直播室时	   |
|API_INVALID_PASSWORD_LIVEROOM_ERROR|"直播室密码非法"|10204|一般不会发生       |

####图片相关

|返回码		              |描述          |值   |备注  						 |
|-------------------------|--------------|-----|-----------------------------|
|API_DECODE_HEAD_PIC_ERROR|"头像解码错误"|10300|base64解压头像时,一般不会发生|


##请求回复说明

所有的API，请求返回头部格式的内容：

```
{
  "access"   : <api接口名称>
  "status"   : <状态码>
  "desc"     : <描述>
}
```

例如：

```
{
	"access"  : "live/Login"
	"status"  : 10000,
	"desc"    : "成功"
}
```

```
{
	"access"  : "live/Register"
	"status"  : 10102,
	"desc"    : "账号已注册"
}
```

##以下API都是POST method

```
用户登录API : /live/Login

request
{
	"account" : 账号/手机号/邮箱
    "password": "123456"
}

response
{
    "access"     : "/live/Login"
    "status"     :
    "desc"       :
    "account"    : "zhangsan"
    "nickName"   : "aaaaaa"
    "liveRoomId" : 如果这个用户有liveRoomId,就返回，如果没有返回0
    "sessionId"  : (32个字节的字符串)
	"headPicId"  : 头像ID
	"sex"		 : 1 男，2 女
}
```

```
用户注册API : /live/Register

request
{
    "account"  : "zhangsan",
    "phone"    : "13512023289",
    "email"    : "example@qq.com"
    "password" : "123456"
}

response
{
    "access"   : "/live/Register"
    "status"   :
    "desc"     :
}
```

```
用户密码修改API : live/ModifyUserPwd

request
{
    "phone" : "",
    "password": "123456"
}

response
{
    "access" : "/live/ModifyUserPwd"
	"status"   :
    "desc"     :
}
```

```
用户个人信息修改API : /live/ModifyUserInfo

request
{
    "sessionId" : ""
	"nickName"  : 如果不修改填""
	"sex"		: 如果不修改填0
	"address"	: 如果不修改填""
}

response
{
    "access" : "/live/ModifyUserInfo"
	"status"   :
    "desc"     :
}
```

```
用户头像修改API : /live/ModifyUserHeadPic

request
{
    "sessionId" : "asssssssssssssssssss"
	"headPic"   : 修改用户头像 (图片base64之后成字符串再传入)
}

response
{
    "access" : "/live/ModifyUserHeadPic"
	"status"   :
    "desc"     :
}
```

```
获取直播室列表API :　/live/GetLiveRoomList

request
{
    "sessionId": 
    "type"     : 1(全部), 2(关注)
    "refresh"  : 1(刷新), 0(获取下一页)
    "count"    : 可以填充一次性刷新几个，如果不填，默认是5个
}

response(all)
{
	"access"   : "/live/GetLiveRoomList"
	"status"   :
	"desc"     :
	"type"     : 1(全部)
	"list"     : [
			{
				"liveRoomId" :
				"liveRoomanchorPwd": 主播密码
				"liveRoomUserPwd": 观看直播的用户密码
				"liveRoomName" : 直播室名称
				"liveRoomDesc" : 直播室描述
				"liveRoomTopic": 演讲的题目
				"anchorName"   : 主播昵称
				"headPicId"    : 主播头像ID
				"isFollow"     : 1 关注, 0 未关注
				"followCount"  : 被关注量					
				"playerCount"  : 当前观看人数
			},
		]
}

response(fcous)
{
     "access"   : "/live/GetLiveRoomList"
     "status"   :
     "desc"     :
     "type"     : 2(关注)
     "onlineList" : [	当前正在直播的list
            {
                "liveRoomId" :
				"liveRoomanchorPwd": 主播密码
				"liveRoomUserPwd": 观看直播的用户密码
				"liveRoomName" : 直播室名称
				"liveRoomDesc" : 直播室描述
				"liveRoomTopic": 演讲的题目
				"anchorName"   : 主播昵称
				"headPicId"    : 主播头像ID
				"isFollow"     : 1 关注, 0 未关注
				"followCount"  : 被关注量					
				"playerCount"  : 当前观看人数
            },
		],
    "offlineList" : [	当前未在直播的list
			{
                "liveRoomId" :
				"liveRoomanchorPwd": 主播密码
				"liveRoomUserPwd": 观看直播的用户密码
				"liveRoomName" : 直播室名称
				"liveRoomDesc" : 直播室描述
				"liveRoomTopic": 演讲的题目
				"anchorName"   : 主播昵称
				"headPicId"    : 主播头像ID
				"isFollow"     : 1 关注, 0 未关注
				"followCount"  : 被关注量					
				"playerCount"  : 当前观看人数
            },
    ]
}
```

```
获取用户自己直播室账号API : /live/GetMyLiveRoomId

request
{
    "sessionId" : ""
}

response
{
    "access"   : " /live/GetMyLiveRoomId"
    "status"   :
    "desc"     :
    "LiveRoomId": 1234567
}
```

```
获取用户自己直播室API : /live/GetMyLiveRoom

request
{
    "sessionId" : ""
}

response
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
	"headPicId"    : 主播头像ID
	"isFollow"     : 1 关注, 0 未关注
	"followCount"  : 被关注量					
	"playerCount"  : 当前观看人数
}
```

```
修改用户自己直播室API : /live/ModifyMyLiveRoom

request
{
	"sessionId"  :
	"liveRoomId" :        可为0,
	"liveRoomAnchorPwd" : 主播密码
	"liveRoomUserPwd"   : 观看直播的用户密码
	"liveRoomName"      : 直播室名称
	"liveRoomDesc"      : 直播室描述
	"liveRoomTopic"     : 直播的主题
}

response
{	
	"access"   : " /live/ModifyMyLiveRoom"
    "status"   :
    "desc"     :
}
```

```
创建直播室API : /live/CreateLiveRoom

request
{
    "sessionId"    :
	"liveRoomAnchorPwd" : 主播密码
    "liveRoomAssistPwd" : 助理密码
	"liveRoomUserPwd"   : 观看密码
	"liveRoomName"      : 直播室名称
	"liveRoomDesc"      : 直播室描述
	"liveRoomTopic"     : 直播室主题

}

response
{
    "access"    : " /live/CreateLiveRoom"
    "status"   :
    "desc"     :
    "liveRoomId":
}
```

```
关注和取消其他直播室API : /live/FollowLiveRoom

request
{
    "sessionId"  :
	"liveRoomId" :
    "isFollow"   : 1 关注, 0 取消

}

response
{
    "access"    : " /live/CreateLiveRoom"
    "status"    :
    "desc"      :
}
```

```
开启和关闭直播API : /live/PushLiveStream

request
{
	"sessionId"  :
	"liveRoomId" :
    "status"     : 1 push, 0 stop
	"timeout"    : 超时时间,默认60秒. 如果在此时间内没有设置状态为1，直播状态变为0
}

response
{
    "access"    : " /live/PushLiveStream"
    "status"    :
    "desc"
}
```

```
搜索直播室API : /live/SearchLiveStream

request
{
	"sessionId" :
	"keyword"   : 目前只支持liveroomId
}

response
{
    "access"   : " /live/SearchLiveStream"
    "status"   :
    "desc"     :
    "liveRoomId" :
	"liveRoomanchorPwd": 主播密码
	"liveRoomUserPwd": 观看直播的用户密码
	"liveRoomName" : 直播室名称
	"liveRoomDesc" : 直播室描述
	"liveRoomTopic": 演讲的题目
	"anchorName"   : 主播昵称
	"headPicId"    : 主播头像ID
	"isFollow"     : 1 关注, 0 未关注
	"followCount"  : 被关注量					
	"playerCount"  : 当前观看人数
}
```


##以下API都是Get method

```
获取用户头像API : /live/GetUserHeadPic

request
/live/GetUserHeadPic?id=headPicId

response
	如果成功，直接返回图片数据

	如果失败，返回
	{
		"access": "/live/GetUserHeadPic",
		"status": 10002,
		"desc"  : "参数错误"
	}
```




