package service

import (
	"sync"
	"time"

	"github.com/futurez/litego/logger"
	"github.com/futurez/litego/util"
)

const (
	HeartBeatTimeOut  = time.Hour * 1
	PreCheckSessionId = time.Minute * 30
)

type sessionManager struct {
	mux        sync.RWMutex
	sessionMap map[string]*sessionData
}

var SP_sessionMgr = &sessionManager{sessionMap: make(map[string]*sessionData)}

func init() {
	time.AfterFunc(PreCheckSessionId, func() { SP_sessionMgr.GC() })
}

func (sm *sessionManager) GC() {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	delIdList := make([]string, 0, len(sm.sessionMap))
	now := time.Now()
	for id, session := range sm.sessionMap {
		if now.Sub(session.tick) > HeartBeatTimeOut {
			delIdList = append(delIdList, id)
		}
	}

	for _, id := range delIdList {
		logger.Debug("GC : Overdue sessionid = ", id)
		delete(sm.sessionMap, id)
	}
	time.AfterFunc(PreCheckSessionId, func() { SP_sessionMgr.GC() })
}

func (sm *sessionManager) addSession(userid int64, liveroomid int64, username string, nickName string) string {
	sm.mux.Lock()
	defer sm.mux.Unlock()

	sd := &sessionData{
		session_id:   util.UUID(),
		user_id:      userid,
		liveroom_id:  liveroomid,
		account:      username,
		nickname:     nickName,
		tick:         time.Now(),
		bfcousOnline: true,
	}
	sm.sessionMap[sd.session_id] = sd
	logger.Info("addSession : session_id=", sd.session_id, ",user_id=", sd.user_id,
		",liveroom_id=", sd.liveroom_id, ",account=", sd.account)
	return sd.session_id
}

func (sm *sessionManager) deleteSession(session_id string) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	logger.Info("deleteSession : session_id=", session_id)
	delete(sm.sessionMap, session_id)
}

func (sm *sessionManager) readSession(sessionId string) (*sessionData, bool) {
	sm.mux.RLock()
	defer sm.mux.RUnlock()
	sd, ok := sm.sessionMap[sessionId]
	return sd, ok
}

type sessionData struct {
	session_id   string
	user_id      int64
	liveroom_id  int64
	account      string
	nickname     string
	tick         time.Time
	allLastId    int64
	fcousLastId  int64
	bfcousOnline bool
}

func (sd *sessionData) UpdateTick() {
	sd.tick = time.Now()
}
