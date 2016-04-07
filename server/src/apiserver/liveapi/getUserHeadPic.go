package liveapi

import (
	"apiserver/service"
	"gotye_protocol"
	"net/http"
	"strconv"

	"github.com/futurez/litego/httplib"
	"github.com/futurez/litego/logger"
)

// http://xxx:xx/live/GetUserHeadPic?id=
func GetUserHeadPic(w http.ResponseWriter, r *http.Request) {
	resp := gotye_protocol.GetUserHeadPicResponse{}
	r.ParseMultipartForm(32 << 20)
	imageId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		logger.Warn("GetUserHeadPic : failed.")
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		resp.SetAccess("/live/GetUserHeadPic")
		httplib.HttpResponseJson(w, http.StatusOK, resp)
	}

	if imageId == 0 {
		logger.Warn("GetUserHeadPic : failed.")
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		resp.SetAccess("/live/GetUserHeadPic")
		httplib.HttpResponseJson(w, http.StatusOK, resp)
		return
	}

	headPic := service.GetHeadPicById(int64(imageId))
	logger.Infof("GetUserHeadPic : headPicId=%d, headPic=%d", imageId, len(headPic))
	httplib.HttpResponseImage(w, headPic)
}
