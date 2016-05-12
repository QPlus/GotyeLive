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

func PlayLiveStream(w http.ResponseWriter, r *http.Request) {
	var status int
	req := gotye_protocol.PlayLiveStreamRequest{}
	resp := gotye_protocol.PlayLiveStreamResponse{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("PlayLiveStream : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("PlayLiveStream : not json = ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("PlayLiveStream : req=", string(readdata))
	status = service.PlayLiveStream(req.SessionId, req.LiveroomId, req.Status)
	resp.SetStatus(status)

end:
	resp.SetAccess("/live/PushLiveStream")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
