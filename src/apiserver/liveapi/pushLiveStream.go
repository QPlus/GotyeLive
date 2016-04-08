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

func PushLiveStream(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.PushLiveStreamResponse{}
	req := gotye_protocol.PushLiveStreamRequest{}

	var status int

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("PushLiveStream : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("PushLiveStream : not json = ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("PushLiveStream : req=", string(readdata))
	status = service.PushingLiveStream(req.SessionId, req.LiveRoomId, req.Status, req.Timeout)
	resp.SetStatus(status)

end:
	resp.SetAccess("/live/PushLiveStream")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
