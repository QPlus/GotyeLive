// my_liveroom.go
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

func ModifyMyLiveRoom(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.ModifyMyLiveRoomResponse{}
	req := gotye_protocol.ModifyMyLiveRoomRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("ModifyMyLiveRoom : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("ModifyMyLiveRoom : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("ModifyMyLiveRoom : req=", string(readdata))
	service.ModifyMyLiveRoom(&resp, &req)

end:
	resp.SetAccess("/live/ModifyMyLiveRoom")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
