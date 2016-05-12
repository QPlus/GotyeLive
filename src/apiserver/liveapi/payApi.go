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

func ChargeRMB(w http.ResponseWriter, r *http.Request) {
	req := gotye_protocol.ChargeRMBRequest{}
	resp := gotye_protocol.ChargeRMBResponse{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("ChargeRMB : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("ChargeRMB : not json = ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("ChargeRMB : req=", string(readdata))
	service.ChargeRMB(&resp, &req)

end:
	resp.SetAccess("/pay/ChargeRMB")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}

func PayQinCoin(w http.ResponseWriter, r *http.Request) {
	req := gotye_protocol.PayQinCoinRequest{}
	resp := gotye_protocol.PayQinCoinResponse{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("PayQinCoin : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("PayQinCoin : not json = ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("PayQinCoin : req=", string(readdata))
	service.PayQinCoin(&resp, &req)

end:
	resp.SetAccess("/pay/PayQinCoin")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}

func GetPayAccount(w http.ResponseWriter, r *http.Request) {
	req := gotye_protocol.GetPayAccountRequest{}
	resp := gotye_protocol.GetPayAccountResponse{}

	readdata, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Warn("GetPayAccount : ", err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	err = json.Unmarshal(readdata, &req)
	if err != nil {
		logger.Warn("GetPayAccount : not json = ", string(readdata))
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		goto end
	}

	logger.Info("GetPayAccount : req=", string(readdata))
	service.GetPayAccount(&resp, &req)

end:
	resp.SetAccess("/pay/GetPayAccount")
	httplib.HttpResponseJson(w, http.StatusOK, &resp)
}
