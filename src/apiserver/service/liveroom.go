package service

import (
	"gotye_protocol"
	"gotye_protocol/gotye_sdk"

	"github.com/futurez/litego/logger"
)

const (
	DefaultAnchorPwd = "openAnchor"
	DefaultAssistPwd = "openAssist"
	DefaultUserPwd   = "openuser"
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
	var (
		err      error
		apptoken string
		result   *gotye_sdk.CreateRoomResponse
	)

sucLoop:
	for i := 0; i < 2; i++ {
		apptoken, err = GotyeAccessAppToken()
		if err != nil {
			logger.Error("CreateLiveRoom : AccessToken Failed, ", err.Error())
			goto errLoop
		}

		result, err = GotyeCreateRoom(apptoken,
			req.LiveRoomName,
			req.LiveAnchorPwd,
			req.LiveAssistPwd,
			req.LiveUserPwd,
			req.LiveRoomDesc,
			req.LiveRoomTopic)
		if err != nil {
			logger.Error("CreateLiveRoom : create room error=", err.Error())
			goto errLoop
		}

		switch result.Status {
		case gotye_sdk.API_SUCCESS:
			logger.Info("CreateLiveRoom : success create liveroom_id= ", result.Entity.RoomId)
			break sucLoop

		case gotye_sdk.API_INVALID_TOKEN_ERROR:
			if i == 0 {
				GotyeClearAccessToken()
				logger.Info("CreateLiveRoom : invalid token error, and accesstoken again.")
			} else {
				logger.Error("CreateLiveRoom : why access new token, but return invalid")
				goto errLoop
			}

		default:
			logger.Error("CreateLiveRoom : create room status=%d", result.Status)
			goto errLoop
		}
	}

	err = DBCreateLiveroom(sd.user_id,
		result.Entity.RoomId,
		result.Entity.RoomName,
		result.Entity.AnchorDesc,
		result.Entity.ContentDesc,
		result.Entity.AnchorPwd,
		result.Entity.AssistPwd,
		result.Entity.UserPwd)
	if err != nil {
		logger.Error("CreateLiveRoom : ", err.Error())
		goto errLoop
	} else {
		sd.liveroom_id = result.Entity.RoomId
		resp.LiveRoomId = result.Entity.RoomId
		resp.SetStatus(gotye_protocol.API_SUCCESS)
		return
	}

errLoop:
	resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
	return
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

sucLoop:
	for i := 0; i < 2; i++ {
		apptoken, err := GotyeAccessAppToken()
		if err != nil {
			logger.Error("ModifyMyLiveRoom : AccessToken Failed, ", err.Error())
			resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
			return
		}

		status := GotyeModifyRoom(apptoken, req.LiveRoomID, req.LiveRoomName,
			req.LiveRoomAnchorPwd, DefaultAssistPwd, req.LiveUserPwd, req.LiveRoomDesc, req.LiveRoomTopic)

		switch status {
		case gotye_sdk.API_SUCCESS:
			logger.Info("ModifyMyLiveRoom : success modify liveroom_id= .", req.LiveRoomID)
			break sucLoop

		case gotye_sdk.API_INVALID_TOKEN_ERROR:
			if i == 0 {
				GotyeClearAccessToken()
				logger.Info("ModifyMyLiveRoom : invalid token error, and accesstoken again.")
			} else {
				logger.Error("ModifyMyLiveRoom : why access new token, but return invalid")
				resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
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
		logger.Info("GetMyLiveRoom : account=", sd.account, "not exist liveroom")
		resp.SetStatus(gotye_protocol.API_LIVEROOM_NOT_EXISTS_ERROR)
		return
	}

	resp.LiveRoomId,
		resp.LiveRoomName,
		resp.LiveRoomDesc,
		resp.LiveRoomTopic,
		resp.LiveAnchorPwd,
		resp.LiveUserPwd,
		ok = DBGetLiveRoomByUserId(sd.user_id)
	if !ok {
		logger.Warnf("GetMyLiveRoom : why get user_id=%d liveroom_id=%d failed.", sd.user_id, sd.liveroom_id)
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
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
