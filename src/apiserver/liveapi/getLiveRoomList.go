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

func GetLiveRoomList(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.ApiResponse{}
	req := gotye_protocol.GetLiveRoomListRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("GetLiveRoomList : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("GetLiveRoomList : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("GetLiveRoomList : req=", string(readdata))

	if req.Type == gotye_protocol.ALL_LIVE_ROOM_LIST {
		GetAllLiveRoomList(w, &req)
		return
	} else if req.Type == gotye_protocol.FOCUS_LIVE_ROOM_LIST {
		GetFcousLiveRoomList(w, &req)
		return
	} else {
		logger.Warn("GetLiveRoomList : unknown type=", req.Type)
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
	}

end:
	resp.SetAccess("liveapi/GetLiveRoomList")
	logger.Info("GetLiveRoomList : resp=", resp)
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}

func GetAllLiveRoomList(w http.ResponseWriter, req *gotye_protocol.GetLiveRoomListRequest) {
	logger.Info("GetAllLiveRoomList : ")
	resp := gotye_protocol.GetAllLiveRoomListResponse{}
	resp.SetAccess("liveapi/GetLiveRoomList")

	service.GetAllLiveRoomList(&resp, req)

	logger.Info("GetAllLiveRoomList : resp=", resp)
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}

func GetFcousLiveRoomList(w http.ResponseWriter, req *gotye_protocol.GetLiveRoomListRequest) {
	logger.Info("GetFcousLiveRoomList : ")
	resp := gotye_protocol.GetFcousLiveRoomListResponse{}
	resp.SetAccess("liveapi/GetLiveRoomList")

	service.GetFcousLiveRoomList(&resp, req)
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
