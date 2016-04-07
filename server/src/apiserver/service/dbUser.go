package service

import (
	"database/sql"
	"errors"
	"fmt"
	"gotye_protocol"
	"strings"

	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

func DBCheckUserAccount(username, password string) (userId, headPicId int64, account, nickname string, sex int8, status_code int) {
	db := SP_MysqlDbPool.GetDBConn()
	var pwd string
	err := db.QueryRow("SELECT user_id, account, nickname, pwd, headpic_id, sex FROM tbl_users WHERE account=? OR phone=? OR email=?",
		username, username, username).Scan(&userId, &account, &nickname, &pwd, &headPicId, &sex)
	switch {
	case err == sql.ErrNoRows:
		logger.Infof("DBCheckUserAccount : %s not exists.", username)
		status_code = gotye_protocol.API_ACCOUNT_NOT_EXISTS_ERROR
		return
	case err != nil:
		logger.Error("DBCheckUserAccount : ", err.Error())
		status_code = gotye_protocol.API_SERVER_ERROR
		return
	}

	if pwd != util.Md5Hash(password) {
		logger.Infof("DBCheckUserAccount : %s password error", username)
		status_code = gotye_protocol.API_LOGIN_PASSWORD_ERROR
		return
	}
	status_code = gotye_protocol.API_SUCCESS
	return
}

func DBIsAccountExists(account string) bool {
	logger.Info("DBIsAccountExists : account=", account)

	db := SP_MysqlDbPool.GetDBConn()
	var count int
	err := db.QueryRow("SELECT count(*) as count FROM tbl_users WHERE account=?", account).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		logger.Warn("DBIsPhoneExists : why not row.")
	case err != nil:
		logger.Error("DBIsAccountExists : ", err.Error())
	}
	return count != 0
}

func DBIsPhoneExists(phone string) bool {
	logger.Info("DBIsPhoneExists : phone=", phone)

	db := SP_MysqlDbPool.GetDBConn()
	var count int
	err := db.QueryRow("SELECT count(*) as count FROM tbl_users WHERE phone=?", phone).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		logger.Warn("DBIsPhoneExists : why not row.")
	case err != nil:
		logger.Error("DBIsPhoneExists : ", err.Error())
	}
	return count != 0
}

func DBIsEmailExists(email string) bool {
	logger.Info("DBIsEmailExists : email=", email)

	db := SP_MysqlDbPool.GetDBConn()
	var count int
	err := db.QueryRow("SELECT count(*) FROM tbl_users WHERE email=?", email).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		logger.Warn("DBIsEmailExists : why not row.")
	case err != nil:
		logger.Error("DBIsEmailExists : ", err.Error())
	}
	return count != 0
}

//create new user
func DBCreateUserAccount(account, email, phone, pwd string) int64 {
	db := SP_MysqlDbPool.GetDBConn()
	res, err := db.Exec("INSERT INTO tbl_users(account,phone,email,nickname,pwd) VALUES(?,?,?,?,?)",
		account, email, phone, account, util.Md5Hash(pwd))
	if err != nil {
		logger.Error("DBCreateUserAccount : ", err.Error())
		return -1
	}
	num, _ := res.LastInsertId()
	logger.Info("DBCreateUserAccount : LastInsertId=", num)
	return num
}

func DBModifyUserNickName(userid int64, nickname string) error {
	db := SP_MysqlDbPool.GetDBConn()
	res, err := db.Exec("UPDATE tbl_users SET nickname=? WHERE user_id=?", nickname, userid)
	if err != nil {
		logger.Error("DBModifyUserNickName : ", err.Error())
		return err
	}
	num, _ := res.RowsAffected()
	logger.Info("DBCreateUserAccount : RowsAffected=", num)
	return nil
}

func DBModifyUserInfo(userid int64, nickname string, sex int8, addr string) error {
	db := SP_MysqlDbPool.GetDBConn()
	var setValue []string
	if len(nickname) > 0 {
		setValue = append(setValue, fmt.Sprintf("nickname='%s'", nickname))
	}
	if sex == 1 || sex == 2 {
		setValue = append(setValue, fmt.Sprintf("sex=%d", sex))
	}
	if len(addr) > 0 {
		setValue = append(setValue, fmt.Sprintf("address='%s'", addr))
	}
	setData := strings.Join(setValue, ",")
	sql := fmt.Sprintf("UPDATE tbl_users SET %s WHERE user_id=%d", setData, userid)
	logger.Info("DBModifyUserInfo : sql=", sql)

	res, err := db.Exec(sql)
	if err != nil {
		logger.Error("DBModifyUserNickName : ", err.Error())
		return err
	}
	num, _ := res.RowsAffected()
	logger.Info("DBModifyUserInfo : RowsAffected=", num)
	return nil
}

func DBGetHeadPicIdByUserId(userid int64) int64 {
	db := SP_MysqlDbPool.GetDBConn()
	var headPicId int64
	err := db.QueryRow("SELECT headpic_id FROM tbl_users WHERE user_id=?", userid).Scan(&headPicId)
	switch {
	case err == sql.ErrNoRows:
		logger.Warn("DBGetHeadPicIdByUserId : why not row.")
	case err != nil:
		logger.Errorf("DBGetHeadPicIdByUserId : userid=%d, err=%s", userid, err.Error())
	default:
		logger.Infof("DBGetHeadPicIdByUserId : user_id=%d, headPicId=%d.", userid, headPicId)
	}
	return headPicId
}

func DBUpdateHeadPicIdByUserId(userId, headPicId int64) error {
	db := SP_MysqlDbPool.GetDBConn()
	res, err := db.Exec("UPDATE tbl_users SET headpic_id=? WHERE user_id=?", headPicId, userId)
	if err != nil {
		logger.Error("DBUpdateHeadPicIdByUserId : ", err.Error())
		return err
	}
	num, _ := res.RowsAffected()
	logger.Info("DBUpdateHeadPicIdByUserId : RowsAffected=", num)
	return nil
}

func DBModifyUserHeadPic(userId int64, headPic []byte) error {
	headPicId := DBGetHeadPicIdByUserId(userId)
	db := SP_MysqlDbPool.GetDBConn()
	if headPicId == 0 {
		//add new headPic
		res, err := db.Exec("INSERT INTO tbl_pictures(pic) VALUES(?)", headPic)
		if err != nil {
			logger.Error("DBModifyUserHeadPic : insert into tbl_pictures failed. ", err.Error())
			return err
		}
		num, err := res.LastInsertId()
		if err != nil {
			logger.Error("DBModifyUserHeadPic : get lastinertid failed. ", err.Error())
			return err
		}
		logger.Info("DBModifyUserHeadPic : insert LastInsertId=", num)
		return DBUpdateHeadPicIdByUserId(userId, num)
	} else {
		res, err := db.Exec("UPDATE tbl_pictures SET pic=? WHERE pic_id=?", headPic, headPicId)
		if err != nil {
			logger.Error("DBModifyUserHeadPic : update tbl_pictures failed, ", err.Error())
			return err
		}
		num, _ := res.RowsAffected()
		logger.Info("DBModifyUserHeadPic : Update RowsAffected=", num)
	}
	return nil
}

func DBModifyUserPwd(phone, pwd string) error {
	db := SP_MysqlDbPool.GetDBConn()
	res, err := db.Exec("UPDATE tbl_users SET pwd=? WHERE phone=?", util.Md5Hash(pwd), phone)
	if err != nil {
		logger.Error("DBModifyUserPwd : update tbl_pictures failed, ", err.Error())
		return err
	}
	num, _ := res.RowsAffected()
	logger.Info("DBModifyUserPwd : RowsAffected=", num)
	if num == 1 {
		return nil
	} else {
		return errors.New("not exist phone")
	}
}

func DBGetUserHeadPic(picId int64) []byte {
	db := SP_MysqlDbPool.GetDBConn()
	var pic []byte
	err := db.QueryRow("SELECT pic FROM tbl_pictures WHERE pic_id=?", picId).Scan(&pic)
	switch {
	case err == sql.ErrNoRows:
		logger.Warn("DBGetUserHeadPic : why not row, picId=", picId)
	case err != nil:
		logger.Error("DBGetUserHeadPic : ", err.Error())
	default:
		logger.Info("DBGetUserHeadPic : success get pic_id=", picId)
	}
	return pic
}
