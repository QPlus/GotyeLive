package service

import (
	"database/sql"
	"fmt"
	"gotye_protocol"
	"strings"

	"github.com/futurez/litego/logger"
)

func DBCreateLiveroom(userId, liveroomId int64, name, desc, topic, anchorPwd, assistPwd, userPwd,
	playRtmpUrls, playHlsUrls, playFlvUrls string) error {

	db := SP_MysqlDbPool.GetDBConn()
	res, err := db.Exec(`INSERT INTO tbl_liverooms(user_id,liveroom_id,liveroom_name,liveroom_desc,liveroom_topic,anchor_pwd,
        assist_pwd,user_pwd,play_rtmp_urls,play_hls_urls,play_flv_urls) VALUES(?,?,?,?,?,?,?,?,?,?,?)`,
		userId, liveroomId, name, desc, topic, anchorPwd, assistPwd, userPwd, playRtmpUrls, playHlsUrls, playFlvUrls)
	if err != nil {
		logger.Error("DBCreateLiveroom : ", err.Error())
		return err
	}
	num, _ := res.RowsAffected()
	logger.Info("DBCreateLiveroom : RowsAffected=", num)
	return nil
}

func DBGetLiveRoomByUserId(user_id int64) (liveroomInfo gotye_protocol.LiveRoomInfo, ok bool) {

	db := SP_MysqlDbPool.GetDBConn()
	err := db.QueryRow(`SELECT liveroom_id, liveroom_name, liveroom_desc, liveroom_topic, anchor_pwd, user_pwd,
      play_rtmp_urls, play_hls_urls, play_flv_urls FROM tbl_liverooms WHERE user_id=?`,
		user_id).Scan(&liveroomInfo.LiveRoomId, &liveroomInfo.LiveRoomName, &liveroomInfo.LiveRoomDesc,
		&liveroomInfo.LiveRoomTopic, &liveroomInfo.LiveAnchorPwd, &liveroomInfo.LiveUserPwd,
		&liveroomInfo.PlayRtmpUrl, &liveroomInfo.PlayHlsUrl, &liveroomInfo.PlayFlvUrl)

	switch {
	case err == sql.ErrNoRows:
		logger.Warnf("DBGetLiveRoomByUserId : user_id=%d not have liveroom_id.", user_id)
		ok = false
	case err != nil:
		logger.Error("DBGetLiveRoomByUserId : ", err.Error())
		ok = false
	default:
		logger.Info("DBGetLiveRoomByUserId : liveroomInfo", liveroomInfo)
		ok = true
	}
	return
}

func DBGetLiveroomIdByUserId(user_id int64) (liveroom_id int64) {
	db := SP_MysqlDbPool.GetDBConn()
	err := db.QueryRow("select liveroom_id from tbl_liverooms where user_id=?", user_id).Scan(&liveroom_id)
	switch {
	case err == sql.ErrNoRows:
		logger.Infof("GetLiveroomIdByUserId : user_id=%d not have liveroom_id.", user_id)
	case err != nil:
		logger.Error("GetLiveroomIdByUserId : ", err.Error())
	default:
		logger.Infof("GetLiveroomIdByUserId : user_id=%d have liveroom_id = %d.", user_id, liveroom_id)
	}
	return
}

func DBModifyLiveRoomInfo(roomId int64, roomName, anchorPwd, assistPwd, userPwd, anchorDesc, contentDesc string) error {
	db := SP_MysqlDbPool.GetDBConn()

	var setValue []string
	if len(roomName) > 0 {
		setValue = append(setValue, fmt.Sprintf("liveroom_name='%s'", roomName))
	}
	if len(anchorPwd) > 0 {
		setValue = append(setValue, fmt.Sprintf("anchor_pwd='%s'", anchorPwd))
	}
	if len(assistPwd) > 0 {
		setValue = append(setValue, fmt.Sprintf("assist_pwd='%s'", assistPwd))
	}
	if len(userPwd) > 0 {
		setValue = append(setValue, fmt.Sprintf("user_pwd='%s'", userPwd))
	}
	if len(anchorDesc) > 0 {
		setValue = append(setValue, fmt.Sprintf("liveroom_desc='%s'", anchorDesc))
	}
	if len(contentDesc) > 0 {
		setValue = append(setValue, fmt.Sprintf("liveroom_topic='%s'", contentDesc))
	}
	setData := strings.Join(setValue, ",")
	sql := fmt.Sprintf("UPDATE tbl_liverooms SET %s WHERE liveroom_id=%d", setData, roomId)

	logger.Info("DBModifyLiveRoomInfo SQL=", sql)
	result, err := db.Exec(sql)
	if err != nil {
		logger.Error("DBModifyLiveRoomInfo : ", err.Error())
		return err
	}
	num, _ := result.RowsAffected()
	logger.Info("DBModifyLiveRoomInfo : RowsAffected=", num)
	return nil
}

func DBAddFollowLiveRoom(userId int64, liveRoomId int64) error {
	db := SP_MysqlDbPool.GetDBConn()
	result, err := db.Exec("INSERT INTO tbl_follow_liverooms(user_id,liveroom_id) VALUES(?,?)", userId, liveRoomId)
	if err != nil {
		logger.Error("DBAddFollowLiveRoom : ", err.Error())
		return err
	}
	num, _ := result.LastInsertId()
	logger.Info("DBAddFollowLiveRoom : LastInsertId=", num)
	return nil
}

func DBDelFollowLiveRoom(userId int64, liveRoomId int64) error {
	db := SP_MysqlDbPool.GetDBConn()
	result, err := db.Exec("DELETE FROM tbl_follow_liverooms WHERE user_id=? AND liveroom_id=?", userId, liveRoomId)
	if err != nil {
		logger.Error("DBDelFollowLiveRoom : ", err.Error())
		return err
	}
	num, _ := result.RowsAffected()
	logger.Info("DBDelFollowLiveRoom : RowsAffected=", num)
	return nil
}

func DBGetFollowCount(liveRoomId int64) (count int) {
	db := SP_MysqlDbPool.GetDBConn()
	err := db.QueryRow("SELECT COUNT(*) as count FROM tbl_follow_liverooms where liveroom_id=?", liveRoomId).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		logger.Errorf("DBGetFollowCount : get liveroom_id=%d follow error.", liveRoomId)
	case err != nil:
		logger.Error("DBGetFollowCount : ", err.Error())
	default:
		logger.Infof("DBGetFollowCount : liveroom_id=%d,count=%d", liveRoomId, count)
	}
	return
}

func DBReloadOnlineLiveroom() *[]onlineLiveroom {
	db := SP_MysqlDbPool.GetDBConn()
	rows, err := db.Query("SELECT liveroom_id, player_num, need_up FROM tbl_online_liverooms")
	if err != nil {
		logger.Error("DBReloadOnlineLiveroom : query err=", err.Error())
		return nil
	}
	defer rows.Close()

	var (
		rooms []onlineLiveroom
		room  onlineLiveroom
	)

	for rows.Next() {
		err = rows.Scan(&room.liveroomId, &room.playerNum, &room.needup)
		if err != nil {
			logger.Error("DBReloadOnlineLiveroom : scan err=", err.Error())
			return nil
		}
		logger.Debug("DBReloadOnlineLiveroom : liveroomId=", room.liveroomId, ",playnum=", room.playerNum)
		rooms = append(rooms, room)
	}
	return &rooms
}

func DBAddOnlineLiveRoom(liveRoomId int64) error {
	db := SP_MysqlDbPool.GetDBConn()
	result, err := db.Exec("INSERT INTO tbl_online_liverooms(liveroom_id) VALUES(?)", liveRoomId)
	if err != nil {
		logger.Error("DBAddOnlineLiveRoom : ", err.Error())
		return err
	}
	num, _ := result.LastInsertId()
	logger.Info("DBAddOnlineLiveRoom : LastInsertId=", num)
	return nil
}

func DBDelOnlineLiveRoom(liveRoomId int64) error {
	db := SP_MysqlDbPool.GetDBConn()
	result, err := db.Exec("DELETE FROM tbl_online_liverooms WHERE liveroom_id=?", liveRoomId)
	if err != nil {
		logger.Error("DBDelOnlineLiveRoom : ", err.Error())
		return err
	}
	num, _ := result.RowsAffected()
	logger.Info("DBDelOnlineLiveRoom : RowsAffected=", num)
	return nil
}

func DBIsOnlineLiveRoom(liveRoomId int64) int8 {
	db := SP_MysqlDbPool.GetDBConn()
	var count int8
	err := db.QueryRow("SELECT COUNT(*) AS count FROM tbl_online_liverooms WHERE liveroom_id=?", liveRoomId).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		logger.Errorf("DBIsOnlineLiveRoom : get liveroom_id=%d follow error.", liveRoomId)
	case err != nil:
		logger.Error("DBIsOnlineLiveRoom : ", err.Error())
	default:
		logger.Infof("DBIsOnlineLiveRoom : liveroom_id=%d,count=%d", liveRoomId, count)
	}
	return count
}

func DBUpdateOnlineLiveRoom(liveroomId int64, num int) error {
	db := SP_MysqlDbPool.GetDBConn()
	result, err := db.Exec("UPDATE tbl_online_liverooms SET player_num=? WHERE liveroom_id=?", num, liveroomId)
	if err != nil {
		logger.Warn("DBUpdateOnlineLiveRoom : ", err.Error())
		return err
	}
	line, _ := result.RowsAffected()
	logger.Infof("DBUpdateOnlineLiveRoom : liveroomid=%d, num=%d, RowsAffected=%d", liveroomId, num, line)
	return nil
}

func DBIsFollowLiveRoom(userId, liveroomId int64) int8 {
	db := SP_MysqlDbPool.GetDBConn()
	var count int8
	err := db.QueryRow("SELECT COUNT(*) as count FROM tbl_follow_liverooms WHERE user_id=? AND liveroom_id=?", userId, liveroomId).Scan(&count)
	switch {
	case err == sql.ErrNoRows:
		logger.Errorf("DBGetFollowCount : get liveroom_id=%d follow error.", liveroomId)
	case err != nil:
		logger.Error("DBGetFollowCount : ", err.Error())
	default:
		logger.Infof("DBGetFollowCount : userId=%d, liveroom_id=%d,count=%d", userId, liveroomId, count)
	}
	return count
}

func DBGetLiveRoomByLiveroomId(resp *gotye_protocol.SearchLiveStreamResponse, liveroomId int64) error {
	db := SP_MysqlDbPool.GetDBConn()
	err := db.QueryRow(`SELECT a.liveroom_name, a.liveroom_desc, a.liveroom_topic, a.anchor_pwd, a.user_pwd, 
    a.play_rtmp_urls, a.play_hls_urls, a.play_flv_urls, b.nickname, b.headpic_id FROM tbl_liverooms a INNER
    JOIN tbl_users b ON a.user_id=b.user_id WHERE a.liveroom_id=?`, liveroomId).Scan(&resp.LiveRoomName,
		&resp.LiveRoomDesc, &resp.LiveRoomTopic, &resp.LiveAnchorPwd, &resp.LiveUserPwd, &resp.PlayRtmpUrl,
		&resp.PlayHlsUrl, &resp.PlayFlvUrl, &resp.AnchorName, &resp.HeadPicId)

	switch {
	case err == sql.ErrNoRows:
		logger.Warn("DBGetLiveRoomByLiveroomId : not have liveroom_id=", liveroomId)
		return err
	case err != nil:
		logger.Error("DBGetLiveRoomByLiveroomId : ", err.Error())
		return err
	default:
		logger.Info("DBGetLiveRoomByLiveroomId : search success liveroomId=", liveroomId)
		resp.LiveRoomId = liveroomId
		return nil
	}
}

func DBGetAllLiveRoomList(resp *gotye_protocol.GetAllLiveRoomListResponse, lastIndex int64, count int) (int64, error) {
	db := SP_MysqlDbPool.GetDBConn()

	var (
		rows *sql.Rows
		err  error
	)

	if lastIndex == 0 {
		rows, err = db.Query(`SELECT b.id, a.liveroom_id, a.liveroom_name, a.liveroom_desc, a.liveroom_topic,
         a.anchor_pwd, a.user_pwd, a.play_rtmp_urls, a.play_hls_urls, a.play_flv_urls, c.nickname, c.headpic_id
        FROM tbl_liverooms a INNER JOIN tbl_online_liverooms b INNER JOIN tbl_users c 
        ON a.liveroom_id=b.liveroom_id AND a.user_id=c.user_id 
        ORDER BY b.pushing_time DESC LIMIT ?`, count)
	} else {
		rows, err = db.Query(`SELECT b.id, a.liveroom_id, a.liveroom_name, a.liveroom_desc, a.liveroom_topic,
         a.anchor_pwd, a.user_pwd, a.play_rtmp_urls, a.play_hls_urls, a.play_flv_urls, c.nickname, c.headpic_id
        FROM tbl_liverooms a INNER JOIN tbl_online_liverooms b INNER JOIN tbl_users c 
        ON a.liveroom_id=b.liveroom_id AND a.user_id=c.user_id 
        ORDER BY b.pushing_time DESC LIMIT ?,?`, lastIndex, count)
	}
	defer rows.Close()

	if err != nil {
		logger.Error("DBGetAllLiveRoomList : ", err.Error())
		return lastIndex, err
	}
	lastId := lastIndex
	for rows.Next() {
		var info gotye_protocol.LiveRoomInfo
		if err = rows.Scan(&lastId, &info.LiveRoomId, &info.LiveRoomName, &info.LiveRoomDesc,
			&info.LiveRoomTopic, &info.LiveAnchorPwd, &info.LiveUserPwd, &info.PlayRtmpUrl, &info.PlayHlsUrl,
			&info.PlayFlvUrl, &info.AnchorName, &info.HeadPicId); err != nil {
			logger.Error("DBGetAllLiveRoomList : ", err.Error())
			resp.List = resp.List[:0]
			return lastIndex, err
		}
		info.FollowCount = DBGetFollowCount(info.LiveRoomId)
		logger.Info("DBGetAllLiveRoomList: info=", info)
		resp.List = append(resp.List, info)
	}
	if err = rows.Err(); err != nil {
		logger.Error("DBGetAllLiveRoomList : ", err.Error())
		resp.List = resp.List[:0]
		return lastIndex, err
	}
	return lastId, nil
}

func DBGetOnlineFocusLiveRoomList(
	resp *gotye_protocol.GetFcousLiveRoomListResponse, userId, lastIndex int64, count int) (int64, error) {

	db := SP_MysqlDbPool.GetDBConn()
	var (
		rows *sql.Rows
		err  error
	)

	if lastIndex == 0 {
		rows, err = db.Query(`SELECT b.id, a.liveroom_id, a.liveroom_name, a.liveroom_desc, a.liveroom_topic,
         a.anchor_pwd, a.user_pwd, a.play_rtmp_urls, a.play_hls_urls, a.play_flv_urls, c.nickname, c.headpic_id
        FROM tbl_liverooms a INNER JOIN tbl_online_liverooms b INNER JOIN tbl_users c 
        ON a.liveroom_id=b.liveroom_id AND a.user_id=c.user_id 
        WHERE a.liveroom_id IN (SELECT liveroom_id FROM tbl_follow_liverooms WHERE user_id=?)
        ORDER BY b.pushing_time DESC LIMIT ?`, userId, count)
	} else {
		rows, err = db.Query(`SELECT b.id, a.liveroom_id, a.liveroom_name, a.liveroom_desc, a.liveroom_topic,
         a.anchor_pwd, a.user_pwd, a.play_rtmp_urls, a.play_hls_urls, a.play_flv_urls, c.nickname, c.headpic_id 
        FROM tbl_liverooms a INNER JOIN tbl_online_liverooms b INNER JOIN tbl_users c 
        ON a.liveroom_id=b.liveroom_id AND a.user_id=c.user_id 
        WHERE a.liveroom_id IN (SELECT liveroom_id FROM tbl_follow_liverooms WHERE user_id=?)
        ORDER BY b.pushing_time DESC LIMIT ?,?`, userId, lastIndex, count)
	}

	if err != nil {
		logger.Error("DBGetAllLiveRoomList : ", err.Error())
		return lastIndex, err
	}

	lastId := lastIndex
	for rows.Next() {
		var info gotye_protocol.LiveRoomInfo
		if err = rows.Scan(&lastId, &info.LiveRoomId, &info.LiveRoomName, &info.LiveRoomDesc,
			&info.LiveRoomTopic, &info.LiveAnchorPwd, &info.LiveUserPwd, &info.PlayRtmpUrl, &info.PlayHlsUrl,
			&info.PlayFlvUrl, &info.AnchorName, &info.HeadPicId); err != nil {
			logger.Error("DBGetAllLiveRoomList : ", err.Error())
			resp.OnlineList = resp.OnlineList[:0]
			return lastIndex, err
		}
		info.FollowCount = DBGetFollowCount(info.LiveRoomId)
		info.IsFollow = 1
		logger.Info("DBGetAllLiveRoomList: info=", info)
		resp.OnlineList = append(resp.OnlineList, info)
	}
	if err = rows.Err(); err != nil {
		logger.Error("DBGetAllLiveRoomList : ", err.Error())
		resp.OnlineList = resp.OnlineList[:0]
		return lastIndex, err
	}
	return lastId, nil
}

func DBGetOfflineFocusLiveRoomList(
	resp *gotye_protocol.GetFcousLiveRoomListResponse, userId, lastIndex int64, count int) (int64, error) {

	db := SP_MysqlDbPool.GetDBConn()
	var (
		rows *sql.Rows
		err  error
	)

	if lastIndex == 0 {
		rows, err = db.Query(`SELECT a.user_id, a.nickname, a.headpic_id, b.liveroom_id, b.liveroom_name, b.liveroom_desc,
         b.liveroom_topic, b.anchor_pwd, b.user_pwd, b.play_rtmp_urls, b.play_hls_urls, b.play_flv_urls
        FROM tbl_users a INNER JOIN tbl_liverooms b ON a.user_id=b.user_id 
        WHERE b.liveroom_id IN (SELECT liveroom_id FROM tbl_follow_liverooms 
        WHERE user_id=? AND (liveroom_id NOT IN (SELECT liveroom_id FROM tbl_online_liverooms)))
        ORDER BY a.user_id DESC LIMIT ?`, userId, count)
	} else {
		rows, err = db.Query(`SELECT a.user_id, a.nickname, a.headpic_id, b.liveroom_id, b.liveroom_name, b.liveroom_desc,
         b.liveroom_topic, b.anchor_pwd, b.user_pwd, b.play_rtmp_urls, b.play_hls_urls, b.play_flv_urls
        FROM tbl_users a INNER JOIN tbl_liverooms b ON a.user_id=b.user_id 
        WHERE b.liveroom_id IN (SELECT liveroom_id FROM tbl_follow_liverooms 
        WHERE user_id=? AND (liveroom_id NOT IN (SELECT liveroom_id FROM tbl_online_liverooms)))
        ORDER BY a.user_id DESC LIMIT ?,?`, userId, lastIndex, count)
	}

	if err != nil {
		logger.Error("DBGetAllLiveRoomList : ", err.Error())
		return lastIndex, err
	}

	lastuserId := lastIndex
	for rows.Next() {
		var info gotye_protocol.LiveRoomInfo
		if err = rows.Scan(&lastuserId, &info.AnchorName, &info.HeadPicId, &info.LiveRoomId,
			&info.LiveRoomName, &info.LiveRoomDesc, &info.LiveRoomTopic, &info.LiveAnchorPwd,
			&info.LiveUserPwd, &info.PlayRtmpUrl, &info.PlayHlsUrl, &info.PlayFlvUrl); err != nil {
			logger.Error("DBGetAllLiveRoomList : ", err.Error())
			resp.OfflineList = resp.OfflineList[:0]
			return lastIndex, err
		}
		info.FollowCount = DBGetFollowCount(info.LiveRoomId)
		info.IsFollow = 1
		logger.Info("DBGetAllLiveRoomList: info=", info)
		resp.OfflineList = append(resp.OfflineList, info)
	}
	if err = rows.Err(); err != nil {
		logger.Error("DBGetAllLiveRoomList : ", err.Error())
		resp.OfflineList = resp.OfflineList[:0]
		return lastIndex, err
	}
	return lastuserId, nil
}

func DBUpdateLiveroomUrls(liveroomId int64, playRtmpUrls, playHlsUrls, playFlvUrls string) error {
	db := SP_MysqlDbPool.GetDBConn()
	_, err := db.Exec("UPDATE tbl_liverooms SET play_rtmp_urls=?, play_hls_urls=?, play_flv_urls=? WHERE liveroom_id=?",
		playRtmpUrls, playHlsUrls, playFlvUrls, liveroomId)

	if err != nil {
		logger.Error("DBUpdateLiveroomUrls : ", err.Error())
		return err
	}
	logger.Info("DBUpdateLiveroomUrls : success. liveroomId=", liveroomId, ", playRtmpUrls=", playRtmpUrls,
		", playHlsUrls=", playHlsUrls, ", playFlvUrls=", playFlvUrls)
	return nil
}

func DBGetLiveroomUrls(liveroomId int64) (playRtmpUrls, playHlsUrls, playFlvUrls string) {
	db := SP_MysqlDbPool.GetDBConn()

	err := db.QueryRow("SELECT play_rtmp_urls, play_hls_urls, play_flv_urls FROM tbl_liverooms WHERE liveroom_id=?",
		liveroomId).Scan(&playRtmpUrls, &playHlsUrls, &playFlvUrls)

	if err != nil {
		logger.Error("DBGetLiveroomUrls : ", err.Error())
		return
	}

	logger.Info("DBGetLiveroomUrls : success. liveroomId=", liveroomId, ", playRtmpUrls=", playRtmpUrls,
		", playHlsUrls=", playHlsUrls, ", playFlvUrls=", playFlvUrls)
	return
}

func DBGetLiveroomIds() []int64 {
	db := SP_MysqlDbPool.GetDBConn()
	rows, err := db.Query("SELECT liveroom_id FROM tbl_liverooms")
	if err != nil {
		logger.Error("DBGetLiveroomIds : query err=", err.Error())
		return nil
	}
	defer rows.Close()

	var (
		liveroomId  int64
		liveroomIds []int64
	)

	for rows.Next() {
		err = rows.Scan(&liveroomId)
		if err != nil {
			logger.Error("DBGetLiveroomIds : scan err=", err.Error())
			return nil
		}
		liveroomIds = append(liveroomIds, liveroomId)
	}
	return liveroomIds
}
