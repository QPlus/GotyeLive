package service

import (
	"gotye_protocol"

	"github.com/futurez/litego/logger"
)

const (
	DEFAULT_LIST_COUNT = 5
)

func GetAllLiveRoomList(resp *gotye_protocol.GetAllLiveRoomListResponse, req *gotye_protocol.GetLiveRoomListRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		logger.Warn("GetLiveRoomList : get session data failed.")
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		return
	}
	sd.UpdateTick()
	resp.Type = req.Type

	if req.Refresh == 1 {
		sd.allLastId = 0
	}

	if req.Count == 0 {
		req.Count = DEFAULT_LIST_COUNT
	}

	logger.Info("GetAllLiveRoomList start allLastId=", sd.allLastId)
	var err error
	sd.allLastId, err = DBGetAllLiveRoomList(resp, sd.allLastId, req.Count)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	} else {
		resp.SetStatus(gotye_protocol.API_SUCCESS)
	}
	logger.Info("GetAllLiveRoomList end allLastId=", sd.allLastId)
	for i := range resp.List {
		resp.List[i].IsFollow = DBIsFollowLiveRoom(sd.user_id, resp.List[i].LiveRoomId)
	}
}

func GetFcousLiveRoomList(resp *gotye_protocol.GetFcousLiveRoomListResponse, req *gotye_protocol.GetLiveRoomListRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		logger.Warn("GetLiveRoomList : get session data failed.")
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		return
	}
	sd.UpdateTick()
	resp.Type = req.Type

	if req.Refresh == 1 {
		sd.fcousLastId = 0
		sd.bfcousOnline = true
	}

	count := req.Count
	if count == 0 {
		count = DEFAULT_LIST_COUNT
	}

	var err error
	if sd.bfcousOnline {
		logger.Info("GetFcousLiveRoomList : get online fcous list, account=", sd.account)
		sd.fcousLastId, err = DBGetOnlineFocusLiveRoomList(resp, sd.user_id, sd.fcousLastId, count)
		if err != nil {
			resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
			logger.Warn("GetFcousLiveRoomList : online list, err=", err.Error())
			return
		} else {
			resp.SetStatus(gotye_protocol.API_SUCCESS)
			if len(resp.OnlineList) >= count {
				logger.Infof("GetFcousLiveRoomList : get online list full, account=%d, len=%d", sd.account, len(resp.OnlineList))
				return
			} else {
				logger.Infof("GetFcousLiveRoomList : get online finished. and start get offline list.")
				sd.fcousLastId = 0
				count -= len(resp.OnlineList)
				sd.bfcousOnline = false
			}
		}
	}
	logger.Infof("GetFcousLiveRoomList : start get offline")

	sd.fcousLastId, err = DBGetOfflineFocusLiveRoomList(resp, sd.user_id, sd.fcousLastId, count)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		logger.Warn("GetFcousLiveRoomList : online list, err=", err.Error())
		return
	} else {
		resp.SetStatus(gotye_protocol.API_SUCCESS)
	}
	logger.Infof("GetFcousLiveRoomList : account=%s, onlineLen=%d, offlineLen=%d", sd.account, len(resp.OnlineList), len(resp.OfflineList))
}
