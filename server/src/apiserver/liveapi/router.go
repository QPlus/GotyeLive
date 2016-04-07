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
	httpSvc.HandleFunc("/live/Login", httplib.MakeHandler(Login))
	httpSvc.HandleFunc("/live/Register", httplib.MakeHandler(Register))
	httpSvc.HandleFunc("/live/ModifyUserInfo", httplib.MakeHandler(ModifyUserInfo))
	httpSvc.HandleFunc("/live/ModifyUserHeadPic", httplib.MakeHandler(ModifyUserHeadPic))
	httpSvc.HandleFunc("/live/ModifyUserPwd", httplib.MakeHandler(ModifyUserPwd))
	httpSvc.HandleFunc("/live/GetUserHeadPic", httplib.MakeHandler(GetUserHeadPic))

	//my live room
	httpSvc.HandleFunc("/live/CreateLiveRoom", httplib.MakeHandler(CreateLiveRoom))
	httpSvc.HandleFunc("/live/ModifyMyLiveRoom", httplib.MakeHandler(ModifyMyLiveRoom))
	httpSvc.HandleFunc("/live/GetMyLiveRoomId", httplib.MakeHandler(GetMyLiveRoomId))
	httpSvc.HandleFunc("/live/GetMyLiveRoom", httplib.MakeHandler(GetMyLiveRoom))

	//follow live room
	httpSvc.HandleFunc("/live/FollowLiveRoom", httplib.MakeHandler(FollowLiveRoom))

	//list live room
	httpSvc.HandleFunc("/live/GetLiveRoomList", httplib.MakeHandler(GetLiveRoomList))

	//push live stream
	httpSvc.HandleFunc("/live/PushLiveStream", httplib.MakeHandler(PushLiveStream))

	//search live room
	httpSvc.HandleFunc("/live/SearchLiveRoom", httplib.MakeHandler(SearchLiveRoom))

	httpSvc.ListenAndServe()
}
