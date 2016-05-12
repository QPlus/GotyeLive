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

type Req_RefreshPlayUrls struct {
	Secret string `json:"secret"`
}

type Resp_RefreshPlayUrls struct {
	gotye_protocol.ApiResponse
}

func RefreshPlayUrls(w http.ResponseWriter, r *http.Request) {
	req := Req_RefreshPlayUrls{}
	resp := Resp_RefreshPlayUrls{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("RefreshPlayUrls : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("RefreshPlayUrls : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	if req.Secret != "gotyeopenlive" {
		logger.Warn("RefreshPlayUrls : secret error ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	resp.SetStatus(service.RefreshPlayUrls())

end:
	resp.SetAccess("/admin/RefreshPlayUrls")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
