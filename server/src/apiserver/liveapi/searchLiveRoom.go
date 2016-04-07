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

func SearchLiveRoom(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.SearchLiveStreamResponse{}
	req := gotye_protocol.SearchLiveStreamRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("SearchLiveRoom : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("SearchLiveRoom : not json = ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("SearchLiveRoom : req=", string(readdata))
	service.LiveRoomSearch(&resp, req.SessionId, req.Keyword)

end:
	resp.SetAccess("/live/SearchLiveRoom")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
