package liveapi

import (
	"apiserver/service"
	"encoding/json"
	"gotye_protocol"
	"io/ioutil"
	"net/http"

	"github.com/futurez/litego/httplib"
	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

func AuthCode(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.AuthCodeResponse{}
	req := gotye_protocol.AuthCodeRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("AuthCode : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("AuthCode : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	if !util.CheckPhone(req.Phone) {
		logger.Warn("AuthCode : phone is invalid. ", string(readdata))
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "phone invalid.")
		goto end
	}

	logger.Info("AuthCode : phone=", req.Phone)
	service.RequestAuthCode(&resp, &req)

end:
	resp.SetAccess("/live/AuthCode")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
