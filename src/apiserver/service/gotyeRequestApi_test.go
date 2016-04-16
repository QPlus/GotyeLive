package service

import (
	"gotye_protocol/gotye_sdk"
	"testing"
)

func TestGotyeGetRoomsLiveInfo(t *testing.T) {
	InitMysqlDbPool("192.168.1.141", "gotye_open_live", "appuser", "123456")

	ids := []int64{210336}

	//get live room play urls.
	liveInfo, err := GotyeGetRoomsLiveInfo(ids...)
	if err != nil {
		t.Error("CreateLiveRoom : ", err.Error())
		return
	}
	if liveInfo.Status != gotye_sdk.API_SUCCESS || len(liveInfo.Entities) == 0 {
		t.Error("CreateLiveRoom : GotyeGetLiveContext status=", liveInfo.Status)
		return
	}

	//	//insert to tbl_liveroom_urls
	//	for _, entity := range liveInfo.Entities {
	//		err = DBInsertLiveroomUrls(entity.RoomId, entity.PlayRtmpUrls, entity.PlayHlsUrls, entity.PlayFlvUrls)
	//		if err != nil {
	//			logger.Error("CreateLiveRoom : ", err.Error())
	//		}
	//	}
}
