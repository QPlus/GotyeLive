http://162.gotlive.com.cn/live_api.html

https://livevip.com.cn/live/admin/index

http://www.gotye.com.cn/dev2/docs/live/start.html


1. 申请三个域名提供给gotye  
	视频推流加速域名: rtmp://rtmppublisher.livevip.com.cn/live 
	视频播放加速域名: rtmp://liveplayer.livevip.com.cn/live 
	视频下载加速域名: http://liveplayer.livevip.com.cn/live 
	Access Secret: 8aa76e4e008e4a8b82db68f289c8ead0 

2. IOS SDK
	GotyeLiveCore : 登录和验证(我要提供哪些参数)
	1). init sdk （AccessSecret, CompanyID） [已提供]
	
	2). 验证直播室(gotyeRoomID(直播室ID), 
				 password(房间密码),
				 bindAccount(估计是gotye的账号),
				 nickname(终端用户昵称)
		 		 是否是第三方房间号）
		
	GotyeLiveChat
	1). login room (account, nickname)
	
	GotyeLivePlayer
	GotyeLivePublisher
	
3.  Server
	1). AccessToken 		
		我的后台
		{
			"scope" : "app",
			"username" : "zhouxueshi@gotye.com.cn"
			"password" : "123456"
		}

		主播室口令授权(roomid, roomkey) 分别用于哪些情况
		{
			account: 
			nickname:
		}
			
			
	2). 账户信息
		GetAppUser (是开发者账号，不是app用户的账号)	
		ModifyAppUser
		
		GetClientSdk (开发者APP)
		SetClientSdk 
				
	3). 播放地址管理
		GetVideoUrls (获取app推流，播放地址)
		{
			"uploadUrl" : 推流地址
			"rtmpUrl"   : 播放地址
			"httpUrl"   : 下载地址(点播地址)
			
		}
			
		GetClientUrls (获取html5的分享地址)
		{
			"modeChatUrl" : 分享
		}
			
			
	
	3). 主播室管理
	a.	CreateRoom (创建一个主播室是视频还是聊天室)
		{
			"roomName"    : 主播室名称
			"anchorPwd"   : 主播登录密码
			"assistPwd"   : 助理登录密码
			"userPwd"	  : 用户登录密码
			"anchorDesc"  : "主播描述"
			"contentDesc"  : "内容描述"
			"thirdRoomID" : //是否用第三方roomID
		}
			
		{
			"roomId" :    //房间ID
			"appUserID" : //APPID
			"roomName" : 
		    "anchorPwd": "000000",
		    "userPwd": "222222",
		    "assistPwd": "111111",
		    "anchorDesc": "just test",
		    "contentDesc": "今天我们来讲讲骑行西藏是一种怎样的体验",
		    "dateCreate": 1445235210660,//创建时间
		    "thirdRoomId": "55555",
		}
	
	b.	GetRooms (获取主播室列表)
	
	c.  GetRoom
		{
			一些基本信息
		}

	d.	ModifyRoom
		{
			"roomId" : 100041, //必填
			"roomName" : "",
			"enableRecordFlag": 1,
    		"permanentPlayFlag": 1,
			"anchorPwd" : //主播
			"assistPwd" : //助理
			"userPwd" : //观众
			"anchorDesc" : //主播描述
			"contenDesc" : //主播室描述
		}
		
	e.	DeleteRoom (删除聊天室)
		{
			"roomId" : 858588
		}
		
			
	4). 直播上下文管理	
	a.  GetLiveContext (room/app)
		{
			"roomId" : 111	
		}

		{
			"publisherAccount" : 主讲人账号
			"share": null, // 当前共享数据
    		"recordingStatus": 0, // 1-录制中 0-停录中
    		"talkerAccount": null, // 当前给麦人
    		"playUserCount": 0 // 当前播放视频人数
		}	
	
	b.  SetLiveStatus (room/app)
		{
			"roomId"  : 111,
			"timeout" : 10,
			"status"  : 1,
		}
	
	5). 获取在线人数统计结果
		/GetPlayUserBillingStati 
	
	
	
	
	
	
		
				