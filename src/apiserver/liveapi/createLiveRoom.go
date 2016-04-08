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

func CreateLiveRoom(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.CreateLiveRoomResponse{}
	req := gotye_protocol.CreateLiveRoomRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("CreateLiveRoom : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("CreateLiveRoom : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("CreateLiveRoom : req = ", string(readdata))
	service.CreateLiveRoom(&resp, &req)

end:
	resp.SetAccess("/live/CreateLiveRoom")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
