package service

import (
	"gotye_protocol"
	"gotye_protocol/gotye_sdk"

	"github.com/futurez/litego/logger"
)

func RefreshPlayUrls() int {
	liveroomIds := DBGetLiveroomIds()

	total := len(liveroomIds)
	if total <= 0 {
		return gotye_protocol.API_SUCCESS
	}

	limit := (total + PreRefreshNum) / PreRefreshNum * PreRefreshNum
	logger.Debug("RefreshPlayUrls : total=", total, ",limit=", limit)

	for i := 0; i < limit; i += PreRefreshNum {
		var ids []int64
		if i != (limit - PreRefreshNum) {
			ids = liveroomIds[i : i+PreRefreshNum]
			logger.Debug("RefreshPlayUrls : last i=", i, ",ids=", ids)
		} else {
			ids = liveroomIds[i:total]
			logger.Debug("RefreshPlayUrls : pre i=", i, ",ids=", ids)
		}

		resp, err := GotyeGetRoomsLiveInfo(ids...)
		if err != nil {
			logger.Error("UpdateNum : GotyeGetLiveContext Failed, ", err.Error())
			continue
		}

		if resp.Status != gotye_sdk.API_SUCCESS {
			logger.Error("UpdateNum : GotyeGetLiveContext status=", resp.Status)
			continue
		}

		for _, entity := range resp.Entities {
			DBUpdateLiveroomUrls(entity.RoomId, entity.PlayRtmpUrls[0], entity.PlayHlsUrls[0], entity.PlayFlvUrls[0])
		}
	}
	return gotye_protocol.API_SUCCESS
}
