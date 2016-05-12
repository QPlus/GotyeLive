package service

import (
	"gotye_protocol"

	"github.com/futurez/litego/logger"
)

func ChargeRMB(resp *gotye_protocol.ChargeRMBResponse, req *gotye_protocol.ChargeRMBRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionId)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Warn("ChargeRMB : get session data failed.")
		return
	}
	sd.UpdateTick()

	if req.RMB <= 0 {
		resp.SetStatus(gotye_protocol.API_PARAM_ERROR)
		logger.Error("ChargeRMB : user=", sd.nickname, ", invalid rmb=", req.RMB)
		return
	}

	var err error
	resp.QinCoin, err = dbChargeRMB(sd.user_id, req.RMB)
	if err != nil {
		resp.SetStatus(gotye_protocol.API_CHARGE_RMB_ERROR)
		return
	}
	resp.SetStatus(gotye_protocol.API_SUCCESS)
}

func PayQinCoin(resp *gotye_protocol.PayQinCoinResponse, req *gotye_protocol.PayQinCoinRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionId)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Warn("PayQinCoin : get session data failed.")
		return
	}
	sd.UpdateTick()

	//get anchorUserId account user_id
	anchorUserId := DBGetUserIdByNickname(req.AnchorAccount)
	if anchorUserId == 0 {
		resp.SetStatus(gotye_protocol.API_USERNAME_NOT_EXISTS_ERROR)
		logger.Warn("PayQinCoin : not this account=", req.AnchorAccount)
		return
	}

	errorno := dbUpdateJiaCoin(anchorUserId, sd.user_id, req.QinCoin)
	switch errorno {
	case 0:
		resp.SetStatus(gotye_protocol.API_LACK_OF_BALANCE_ERROR)
	case 1:
		resp.SetStatus(gotye_protocol.API_SUCCESS)
	default:
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
	}
	logger.Info("PayQinCoin : userid=", sd.user_id, "pay qincoin=", req.QinCoin, " to anchoruserid=", anchorUserId, "(",
		req.AnchorAccount, "), errorno=", errorno)
}

func GetPayAccount(resp *gotye_protocol.GetPayAccountResponse, req *gotye_protocol.GetPayAccountRequest) {
	sd, ok := SP_sessionMgr.readSession(req.SessionId)
	if !ok {
		resp.SetStatus(gotye_protocol.API_EXPIRED_SESSION_ERROR)
		logger.Warn("GetPayAccount : get session data failed.")
		return
	}
	sd.UpdateTick()

	db := SP_MysqlDbPool.GetDBConn()
	err := db.QueryRow("SELECT `qin_coin`,`jia_coin`,`level`,`xp` FROM tbl_pay_account WHERE user_id=?", sd.user_id).
		Scan(&resp.QinCoin, &resp.JiaCoin, &resp.Level, &resp.XP)

	if err != nil {
		logger.Error("GetPayAccount : err=", err.Error())
		resp.SetStatus(gotye_protocol.API_SERVER_ERROR)
		return
	}
	logger.Info("GetPayAccount : suc qincoin=", resp.QinCoin, ",jiacoin=", resp.JiaCoin, ",level=", resp.Level, ",xp=", resp.XP)
	resp.SetStatus(gotye_protocol.API_SUCCESS)
	return
}

func dbUpdateJiaCoin(anchorUserId int64, vistorUserId int64, qinCoin int) int {
	db := SP_MysqlDbPool.GetDBConn()
	tx, err := db.Begin()
	if err != nil {
		logger.Error("dbUpdateJiaCoin : err=", err.Error())
		return -1
	}
	defer tx.Commit()
	//	defer (
	////		tx.Exec("UNLOCK TABLES")

	//	)
	////	tx.Exec("LOCK TABLE tbl_pay_account WRITE")

	var errorno int
	err = tx.QueryRow("CALL pay_qin_coin(?, ?, ?)", anchorUserId, vistorUserId, qinCoin).Scan(&errorno)
	if err != nil {
		logger.Error("dbUpdateJiaCoin : CALL charge_rmb err=", err.Error())
		return -1
	}
	logger.Info("dbUpdateJiaCoin : errorno=", errorno)
	return errorno
}

func dbChargeRMB(user_id int64, rmb int) (int, error) {
	db := SP_MysqlDbPool.GetDBConn()
	tx, err := db.Begin()
	if err != nil {
		logger.Error("dbChargeRMB : err=", err.Error())
		return 0, err
	}

	defer tx.Commit()
	/* (
		//tx.Exec("UNLOCK TABLES")

	)*/

	//tx.Exec("LOCK TABLE tbl_pay_account WRITE")

	var total_qin_coin int
	err = tx.QueryRow("CALL charge_rmb(?, ?)", user_id, rmb).Scan(&total_qin_coin)
	if err != nil {
		logger.Error("dbChargeRMB : CALL charge_rmb err=", err.Error())
		return 0, err
	}
	logger.Info("dbChargeRMB : user_id=", user_id, ", rmb=", rmb)
	return total_qin_coin, nil
}
