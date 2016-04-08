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

func GetMyLiveRoom(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.GetMyLiveRoomResponse{}
	req := gotye_protocol.GetMyLiveRoomRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("PushLiveStream : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("PushLiveStream : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("GetMyLiveRoom : req=", req)
	service.GetMyLiveRoom(&resp, &req)

end:
	resp.SetAccess("/live/GetMyLiveRoom")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
