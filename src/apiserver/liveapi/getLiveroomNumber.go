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

func GetLiveroomNumber(w http.ResponseWriter, r *http.Request) {
	var status int
	req := gotye_protocol.GetLiveroomNumberRequest{}
	resp := gotye_protocol.GetLiveroomNumberResponse{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("GetLiveroomNumber : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("GetLiveroomNumber : not json = ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("GetLiveroomNumber : req=", string(readdata))
	resp.Number, status = service.GetLiveroomNumber(req.SessionId, req.LiveroomId)
	resp.LiveroomId = req.LiveroomId
	resp.SetStatus(status)

end:
	resp.SetAccess("/live/PushLiveStream")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
