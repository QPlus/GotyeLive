package service

import (
	"gotye_protocol"
	"strconv"

	"github.com/futurez/litego/logger"
)

func LiveRoomSearch(resp *gotye_protocol.SearchLiveStreamResponse, sessionId, keyword string) {
	sd, ok := SP_sessionMgr.readSession(sessionId)
	if !ok {
		logger.Warn("LiveRoomSearch : get session data failed.")
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		return
	}
	sd.UpdateTick()
	logger.Infof("LiveRoomSearch : sessionId=%s, keyword=%s", sessionId, keyword)

	liveroomId, err := strconv.Atoi(keyword)
	if err != nil {
		logger.Warnf("LiveRoomSearch : keyword=%d, err=%s", keyword, err.Error())
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		return
	}

	err = DBGetLiveRoomByLiveroomId(resp, int64(liveroomId), sd.user_id)
	if err != nil {
		logger.Warn("LiveRoomSearch : keyword=", keyword)
		resp.SetStatus(gotye_protocol.API_LIVEROOM_ID_NOT_EXIST_ERROR)
		return
	}
	resp.FollowCount = DBGetFollowCount(resp.LiveRoomId)
	resp.IsFollow = DBIsFollowLiveRoom(sd.user_id, resp.LiveRoomId)
	resp.IsPlay = DBIsOnlineLiveRoom(resp.LiveRoomId)
	resp.SetStatus(gotye_protocol.API_SUCCESS)
	return
}
