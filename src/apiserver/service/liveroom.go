package service

import (
	"gotye_protocol"
	"gotye_protocol/gotye_sdk"

	"github.com/futurez/litego/logger"
)

const (
	DefaultAnchorPwd = "anchorpwd"
	DefaultAssistPwd = "assistpwd"
	DefaultUserPwd   = "userpwd"
)

func CreateLiveRoom(resp *gotye_protocol.CreateLiveRoomResponse,
	req *gotye_protocol.CreateLiveRoomRequest) {

	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Info("CreateLiveRoom : get session data failed.")
		return
	}
	sd.UpdateTick()

	//check is have liveroom_id
	if sd.liveroom_id != 0 {
		resp.SetStatus(gotye_protocol.API_SUCCESS)
		resp.LiveRoomId = sd.liveroom_id
		logger.Warn("CreateLiveRoom : already have liveroom_id=", resp.LiveRoomId)
		return
	}

	//check three password
	if len(req.LiveAnchorPwd) == 0 {
		req.LiveAnchorPwd = DefaultAnchorPwd
	}
	if len(req.LiveAssistPwd) == 0 {
		req.LiveAssistPwd = DefaultAssistPwd
	}
	if len(req.LiveUserPwd) == 0 {
		req.LiveUserPwd = DefaultUserPwd
	}

	// create live room
	result, err := GotyeCreateRoom(req.LiveRoomName, req.LiveAnchorPwd, req.LiveAssistPwd,
		req.LiveUserPwd, req.LiveRoomDesc, req.LiveRoomTopic)
	if err != nil {
		logger.Error("CreateLiveRoom : create room error=", err.Error())
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	}
	if result.Status != gotye_sdk.API_SUCCESS {
		logger.Error("CreateLiveRoom : failed, status=", result.Status)
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	}
	//get live room play urls.
	playRtmpUrl, playHlsUrl, playFlvUrl := GotyeGetLiveroomUrl(result.Entity.RoomId)

	//insert to tbl_liverooms
	err = DBCreateLiveroom(sd.user_id, result.Entity.RoomId, result.Entity.RoomName, result.Entity.AnchorDesc,
		result.Entity.ContentDesc, result.Entity.AnchorPwd, result.Entity.AssistPwd, result.Entity.UserPwd,
		playRtmpUrl, playHlsUrl, playFlvUrl)
	if err != nil {
		logger.Error("CreateLiveRoom : ", err.Error())
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	}

	sd.liveroom_id = result.Entity.RoomId
	resp.LiveRoomId = result.Entity.RoomId
	resp.SetStatus(gotye_protocol.API_SUCCESS)
	logger.Info("CreateLiveRoom : success create liveroom_id= ", result.Entity.RoomId)
}

func ModifyMyLiveRoom(resp *gotye_protocol.ModifyMyLiveRoomResponse,
	req *gotye_protocol.ModifyMyLiveRoomRequest) {

	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Info("ModifyMyLiveRoom : get session data failed.")
		return
	}
	sd.UpdateTick()

	if req.LiveRoomID == 0 {
		req.LiveRoomID = sd.liveroom_id
	}

	if sd.liveroom_id == 0 || req.LiveRoomID != sd.liveroom_id {
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "why liveroomid not equal.")
		logger.Warnf("ModifyMyLiveRoom : sd.liveroom_id=%d, req.liveroom_id=%d.", sd.liveroom_id, req.LiveRoomID)
		return
	}

	status := GotyeModifyRoom(req.LiveRoomID, req.LiveRoomName,
		req.LiveRoomAnchorPwd, DefaultAssistPwd, req.LiveUserPwd, req.LiveRoomDesc, req.LiveRoomTopic)

	switch status {
	case gotye_sdk.API_SUCCESS:
		{
			logger.Info("ModifyMyLiveRoom : success modify liveroom_id= .", req.LiveRoomID)
			err := DBModifyLiveRoomInfo(req.LiveRoomID, req.LiveRoomName, req.LiveRoomAnchorPwd, DefaultAssistPwd,
				req.LiveUserPwd, req.LiveRoomDesc, req.LiveRoomTopic)
			if err != nil {
				logger.Error("ModifyLiveRoom : ", err.Error())
				resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
				return
			}
			logger.Info("ModifyLiveRoom : Success.")
			resp.SetStatus(gotye_protocol.API_SUCCESS)
			return
		}

	case gotye_sdk.API_INVALID_LIVEROOM_ID_ERROR:
		fallthrough
	case gotye_sdk.API_NOT_EXISTS_LIVEROOM_ID_ERROR:
		logger.Error("ModifyMyLiveRoom : invalid liveroom_id =", req.LiveRoomID)
		resp.SetStatus(gotye_protocol.API_LIVEROOM_ID_NOT_EXIST_ERROR)
		return

	case gotye_sdk.API_REPECT_PASSWORD_LIVEROOM_ERROR:
		logger.Error("ModifyMyLiveRoom : repect anchor pwd =", req.LiveRoomAnchorPwd)
		resp.SetStatus(gotye_protocol.API_REPECT_PASSWORD_LIVEROOM_ERROR)
		return

	case gotye_sdk.API_INVALID_PASSWORD_LIVEROOM_ERROR:
		logger.Error("ModifyMyLiveRoom : invalid password =", req.LiveRoomAnchorPwd)
		resp.SetStatus(gotye_protocol.API_INVALID_PASSWORD_LIVEROOM_ERROR)
		return

	case gotye_sdk.API_INVALID_LIVEROOM_NAME_ERROR:
		logger.Error("ModifyMyLiveRoom : invalid roomName =", req.LiveRoomName)
		resp.SetStatus(gotye_protocol.API_INVALID_LIVEROOM_NAME_ERROR)
		return

	case gotye_sdk.API_NULL_LIVEROOM_ID_ERROR:
		logger.Error("ModifyMyLiveRoom : why null liveroom id =", req.LiveRoomID)
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return

	default:
		logger.Error("CreateLiveRoom : unknown status=%d", status)
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	}
}

func FollowLiveRoom(sessionId string, liveRoomId int64, isFollow int) int {
	sd, ok := SP_sessionMgr.readSession(sessionId)
	if !ok {
		logger.Info("ModifyMyLiveRoom : get session data failed.")
		return gotye_protocol.API_EXPIRED_SESSION_ERROR
	}
	sd.UpdateTick()

	var err error
	if isFollow == 1 {
		err = DBAddFollowLiveRoom(sd.user_id, liveRoomId)
	} else {
		err = DBDelFollowLiveRoom(sd.user_id, liveRoomId)
	}
	if err != nil {
		return gotye_protocol.API_SERVER_ERROR
	}
	return gotye_protocol.API_SUCCESS
}

func PushingLiveStream(sessinId string, liveRoomId int64, Status int, Timeout int) int {
	sd, ok := SP_sessionMgr.readSession(sessinId)
	if !ok {
		logger.Warn("PushingLiveStream : get session data failed.")
		return gotye_protocol.API_EXPIRED_SESSION_ERROR
	}
	sd.UpdateTick()

	if liveRoomId != sd.liveroom_id {
		logger.Warnf("PushingLiveStream : why liveroomId(%d) != sd.liveroomId(%d).", liveRoomId, sd.liveroom_id)
		return gotye_protocol.API_LIVEROOM_ID_NOT_EXIST_ERROR
	}

	if Status == 1 {
		SP_onlineLiveMgr.StartPushStream(liveRoomId, Timeout)
	} else {
		SP_onlineLiveMgr.StopPushStream(liveRoomId)
	}
	return gotye_protocol.API_SUCCESS
}

func GetMyLiveRoom(resp *gotye_protocol.GetMyLiveRoomResponse, req *gotye_protocol.GetMyLiveRoomRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		logger.Info("GetMyLiveRoom : get session data failed.")
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		return
	}
	sd.UpdateTick()

	if sd.liveroom_id == 0 {
		logger.Info("GetMyLiveRoom : nickname=", sd.nickname, "not exist liveroom")
		resp.SetStatus(gotye_protocol.API_LIVEROOM_NOT_EXISTS_ERROR)
		return
	}

	liveroomInfo, ok := DBGetLiveRoomByUserId(sd.user_id)
	if !ok {
		logger.Warnf("GetMyLiveRoom : why get user_id=%d liveroom_id=%d failed.", sd.user_id, sd.liveroom_id)
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	}

	resp.LiveRoomId = liveroomInfo.LiveRoomId
	resp.LiveAnchorPwd = liveroomInfo.LiveAnchorPwd
	resp.LiveUserPwd = liveroomInfo.LiveUserPwd
	resp.LiveRoomName = liveroomInfo.LiveRoomName
	resp.LiveRoomDesc = liveroomInfo.LiveRoomDesc
	resp.LiveRoomTopic = liveroomInfo.LiveRoomTopic
	resp.AnchorName = liveroomInfo.AnchorName
	resp.PlayRtmpUrl = liveroomInfo.PlayRtmpUrl
	resp.PlayHlsUrl = liveroomInfo.PlayHlsUrl
	resp.PlayFlvUrl = liveroomInfo.PlayFlvUrl

	if resp.PlayRtmpUrl == "" || resp.PlayHlsUrl == "" || resp.PlayFlvUrl == "" {
		logger.Warn("GetMyLiveRoom : url is nil")
		resp.PlayRtmpUrl, resp.PlayHlsUrl, resp.PlayFlvUrl = GotyeGetLiveroomUrl(resp.LiveRoomId)
		DBUpdateLiveroomUrls(resp.LiveRoomId, resp.PlayRtmpUrl, resp.PlayHlsUrl, resp.PlayFlvUrl)
	}

	resp.FollowCount = DBGetFollowCount(sd.liveroom_id)
	resp.HeadPicId = DBGetHeadPicIdByUserId(sd.user_id)
	resp.AnchorName = sd.nickname
	resp.SetStatus(gotye_protocol.API_SUCCESS)
	return
}

func GetMyLiveRoomId(resp *gotye_protocol.GetMyLiveRoomIdResponse, req *gotye_protocol.GetMyLiveRoomIdRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		logger.Info("GetMyLiveRoom : get session data failed.")
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		return
	}
	sd.UpdateTick()

	resp.LiveRoomId = sd.liveroom_id
	resp.SetStatus(gotye_protocol.API_SUCCESS)
	logger.Info("GetMyLiveRoomId : user_id=%d, liveroom_id=%d.", sd.user_id, sd.liveroom_id)
}

func PlayLiveStream(sessinId string, liveroomId int64, Status int) int {
	sd, ok := SP_sessionMgr.readSession(sessinId)
	if !ok {
		logger.Warn("PushingLiveStream : get session data failed.")
		return gotye_protocol.API_EXPIRED_SESSION_ERROR
	}
	sd.UpdateTick()

	if Status == 1 {
		SP_onlineLiveMgr.StartPlayStream(liveroomId)
	} else {
		SP_onlineLiveMgr.StopPlayStream(liveroomId)
	}
	return gotye_protocol.API_SUCCESS
}

func GetLiveroomNumber(sessinId string, liveroomId int64) (int, int) {
	sd, ok := SP_sessionMgr.readSession(sessinId)
	if !ok {
		logger.Warn("PushingLiveStream : get session data failed.")
		return 0, gotye_protocol.API_EXPIRED_SESSION_ERROR
	}
	sd.UpdateTick()

	return SP_onlineLiveMgr.GetPlayCount(liveroomId), gotye_protocol.API_SUCCESS
}
