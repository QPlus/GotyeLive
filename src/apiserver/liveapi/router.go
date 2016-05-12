package liveapi

import (
	"github.com/futurez/litego/httplib"
)

func StartHttpServer(host string, port int) {
	cfg := httplib.Config{
		Host: host,
		Port: port,
	}
	httpSvc := httplib.NewServer(cfg)

	//user info
	httpSvc.HandleFunc("/live/Login", Login)
	httpSvc.HandleFunc("/live/Register", Register)
	httpSvc.HandleFunc("/live/ModifyUserInfo", ModifyUserInfo)
	httpSvc.HandleFunc("/live/ModifyUserHeadPic", ModifyUserHeadPic)
	httpSvc.HandleFunc("/live/ModifyUserPwd", ModifyUserPwd)
	httpSvc.HandleFunc("/live/GetUserHeadPic", GetUserHeadPic)
	httpSvc.HandleFunc("/live/AuthCode", AuthCode)

	//my live room
	httpSvc.HandleFunc("/live/CreateLiveRoom", CreateLiveRoom)
	httpSvc.HandleFunc("/live/ModifyMyLiveRoom", ModifyMyLiveRoom)
	httpSvc.HandleFunc("/live/GetMyLiveRoomId", GetMyLiveRoomId)
	httpSvc.HandleFunc("/live/GetMyLiveRoom", GetMyLiveRoom)

	//follow live room
	httpSvc.HandleFunc("/live/FollowLiveRoom", FollowLiveRoom)

	//list live room
	httpSvc.HandleFunc("/live/GetLiveRoomList", GetLiveRoomList)

	//push live stream
	httpSvc.HandleFunc("/live/PushLiveStream", PushLiveStream)

	//play live stream
	httpSvc.HandleFunc("/live/PlayLiveStream", PlayLiveStream)

	//get live room player number.
	httpSvc.HandleFunc("/live/GetLiveroomNumber", GetLiveroomNumber)

	//search live room
	httpSvc.HandleFunc("/live/SearchLiveRoom", SearchLiveRoom)

	//refresh domain
	httpSvc.HandleFunc("/admin/RefreshPlayUrls", RefreshPlayUrls)

	//pay system
	//httpSvc.HandleFunc("/pay/ChargeRMB", ChargeRMB)
	httpSvc.HandleFunc("/pay/PayQinCoin", PayQinCoin)
	httpSvc.HandleFunc("/pay/GetPayAccount", GetPayAccount)

	httpSvc.ListenAndServe()
}
