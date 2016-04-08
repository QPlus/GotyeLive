package liveapi

import (
	"apiserver/service"
	"encoding/json"
	"gotye_protocol"
	"io/ioutil"
	"net/http"

	"github.com/futurez/litego/httplib"
	"github.com/futurez/litego/logger"
)

func FollowLiveRoom(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.FollowLiveRoomResponse{}
	req := gotye_protocol.FollowLiveRoomRequest{}

	var status int

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("FollowLiveRoom : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("FollowLiveRoom : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("FollowLiveRoom : req=", string(readdata))
	status = service.FollowLiveRoom(req.SessionId, req.LiveRoomId, req.IsFollow)
	if status != gotye_protocol.API_SUCCESS {
		logger.Warn("FollowLiveRoom : status=", status)
	} else {
		logger.Info("FollowLiveRoom : success ", string(readdata))
	}
	resp.SetStatus(status)

end:
	resp.SetAccess("/live/FollowLiveRoom")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
