package service

import (
	"encoding/base64"
	"gotye_protocol"

	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

func UserRegister(resp *gotye_protocol.RegisterResponse, req *gotye_protocol.RegisterRequest) {
	if len(req.AuthCode) == 0 {
		logger.Warn("UserRegister : authcode is null. ")
		resp.SetStatus(gotye_protocol.API_AUTHCODE_ERROR)
		return
	}

	if !SP_phoneCode.Check(req.Phone, req.AuthCode) {
		logger.Warn("UserRegister : authcode is error, authcode=", req.AuthCode)
		resp.SetStatus(gotye_protocol.API_AUTHCODE_ERROR)
		return
	}
	SP_phoneCode.Delete(req.Phone)

	user_id := DBCreateUserAccount(req.Phone, req.Passwd)
	if user_id < 0 {
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		logger.Warn("UserRegister : create user error!")
		return
	}

	resp.SetStatus(gotye_protocol.API_SUCCESS)
	logger.Info("UserRegister : Success. phone=", req.Phone)
}

func UserLogin(resp *gotye_protocol.LoginResponse, req *gotye_protocol.LoginRequest) {
	user_id, headPicId, nickname, sex, status_code := DBCheckUserAccount(req.Account, req.Passwd)

	resp.SetStatus(status_code)
	if status_code == gotye_protocol.API_SUCCESS {
		resp.NickName = nickname
		resp.LiveRoomID = DBGetLiveroomIdByUserId(user_id)
		resp.HeadPicId = headPicId
		resp.Sex = sex

		//判断是否已经登录过.
		resp.SessionID = SP_sessionMgr.addSession(user_id, resp.LiveRoomID, resp.NickName)
		logger.Info("UserLogin success nickname=", resp.NickName)
	} else {
		logger.Warn("UserLogin failed. account=", req.Account, ", pwd=", req.Passwd)
	}
}

func UserInfoModify(resp *gotye_protocol.ModifyUserInfoResponse, req *gotye_protocol.ModifyUserInfoRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Info("UserInfoModify : get session data failed.")
		return
	}
	sd.UpdateTick()

	var err error
	if (req.Sex == 1 || req.Sex == 2) || len(req.Address) > 0 {
		if err = DBModifyUserInfo(sd.user_id, req.Sex, req.Address); err != nil {
			resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
			logger.Warn("UserInfoModify : user_id=", sd.user_id, ", sex=", req.Sex, ", addr=", req.Address)
			return
		}
	}

	if len(req.NickName) > 0 {
		if err = DBModifyUserNickName(sd.user_id, req.NickName); err != nil {
			resp.SetStatus(gotye_protocol.API_USERNAME_EXISTS_ERROR)
			logger.Warn("UserInfoModify : user_id=", sd.user_id, ", nickname=", req.NickName)
			return
		}
	}
	logger.Info("UserInfoModify : user_id=", sd.user_id, ", nickname=", req.NickName, ", sex=", req.Sex, ", address=", req.Address)
	resp.SetStatus(gotye_protocol.API_SUCCESS)
}

func UserHeadPicModify(resp *gotye_protocol.ModifyUserHeadPicResponse, req *gotye_protocol.ModifyUserHeadPicRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Info("UserHeadPicModify : get session data failed.")
		return
	}
	sd.UpdateTick()

	logger.Debug("UserHeadPicModify : headPicLen=", len(req.HeadPic))
	headPic, err := base64.StdEncoding.DecodeString(req.HeadPic)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_DECODE_HEAD_PIC_ERROR)
		logger.Info("UserHeadPicModify : decode err =", err.Error())
		return
	}

	resp.HeadPicId, err = DBModifyUserHeadPic(sd.user_id, headPic)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		logger.Info("UserHeadPicModify : update err = ", err)
		return
	}
	resp.SetStatus(gotye_protocol.API_SUCCESS)
}

func GetHeadPicById(id int64) ([]byte, error) {
	logger.Info("GetHeadPicById : id=", id)
	return DBGetUserHeadPic(id)
}

func UserPwdModify(resp *gotye_protocol.ModifyUserPwdResponse, req *gotye_protocol.ModifyUserPwdRequest) {
	if !util.CheckPhone(req.Phone) {
		resp.SetFormatStatus(gotye_protocol.API_PARAM_ERROR, "invalid phone.")
		logger.Warn("UserPwdModify : invalid phone=", req.Phone)
		return
	}
	err := DBModifyUserPwd(req.Phone, req.Passwd)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_PHONE_NOT_EXISTS_ERROR)
	} else {
		resp.SetStatus(gotye_protocol.API_SUCCESS)
	}
}
