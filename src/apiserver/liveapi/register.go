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

func Register(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.RegisterResponse{}
	req := gotye_protocol.RegisterRequest{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("Register : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("Register : reqdata not json ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	if len(req.Account) == 0 {
		logger.Warn("Register : account is null. ", string(readdata))
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "account is null.")
		goto end
	}

	if len(req.Phone) == 0 || !util.CheckPhone(req.Phone) {
		logger.Warn("Register : phone is null. ", string(readdata))
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "phone invalid.")
		goto end
	}

	if len(req.Email) == 0 || !util.ChechEmail(req.Email) {
		logger.Warn("Register : email is null. ", string(readdata))
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "email invalid.")
		goto end
	}

	if len(req.Password) < 6 {
		logger.Warn("Register : pwd so short ", string(readdata))
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "pwd so short.")
		goto end
	}

	logger.Info("Register : account=", req.Account, ",phone=", req.Phone, ",email=", req.Email)
	service.UserRegister(&resp, &req)

end:
	resp.SetAccess("/live/Register")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
