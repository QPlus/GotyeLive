package service

import (
	"encoding/base64"
	"gotye_protocol"

	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

func UserLogin(resp *gotye_protocol.LoginResponse, req *gotye_protocol.LoginRequest) {
	user_id, headPicId, account, nickname, sex, status_code := DBCheckUserAccount(req.Account, req.Passwd)

	resp.SetStatus(status_code)
	if status_code == gotye_protocol.API_SUCCESS {
		resp.Account = account
		resp.NickName = nickname
		resp.LiveRoomID = DBGetLiveroomIdByUserId(user_id)
		resp.HeadPicId = headPicId
		resp.Sex = sex

		//判断是否已经登录过.
		resp.SessionID = SP_sessionMgr.addSession(user_id, resp.LiveRoomID, resp.Account, resp.NickName)
		logger.Info("UserLogin success. account=", resp.Account, ", nickname=", resp.NickName)
	} else {
		logger.Warn("UserLogin failed. account=", req.Account, ", pwd=", req.Passwd)
	}
}

func UserRegister(resp *gotye_protocol.RegisterResponse, req *gotye_protocol.RegisterRequest) {
	if DBIsAccountExists(req.Account) {
		resp.SetStatus(gotye_protocol.API_ACCOUNT_EXISTS_ERROR)
		logger.Info("UserRegister : account = ", req.Account, " is exists.")
		return
	}

	if DBIsPhoneExists(req.Phone) {
		resp.SetStatus(gotye_protocol.API_PHONE_EXISTS_ERROR)
		logger.Info("UserRegister : phone = ", req.Phone, " is exists.")
		return
	}

	if DBIsEmailExists(req.Email) {
		resp.SetStatus(gotye_protocol.API_EMAIL_EXISTS_ERROR)
		logger.Info("UserRegister : email = ", req.Email, " is exists.")
		return
	}

	user_id := DBCreateUserAccount(req.Account, req.Phone, req.Email, req.Password)
	if user_id < 0 {
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		logger.Warn("UserRegister : create user error!")
		return
	}

	resp.SetStatus(gotye_protocol.API_SUCCESS)
	logger.Info("UserRegister : Success. account=", req.Account, ",phone=", req.Phone, ",email=", req.Email)
}

func UserInfoModify(resp *gotye_protocol.ModifyUserInfoResponse, req *gotye_protocol.ModifyUserInfoRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionID)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Info("UserInfoModify : get session data failed.")
		return
	}
	sd.UpdateTick()

	err := DBModifyUserInfo(sd.user_id, req.NickName, req.Sex, req.Address)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	}
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
