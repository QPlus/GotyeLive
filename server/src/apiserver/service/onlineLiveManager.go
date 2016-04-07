package service

import (
	"gotye_protocol/gotye_sdk"
	"sync"
	"time"

	"github.com/futurez/litego/logger"
)

const (
	PushingLimitTimeOut = time.Minute * 2
	PreCheckPushStatus  = time.Minute * 2
	UpdatePlayNumTimer  = time.Minute * 1
)

type OnlineLive struct {
	//userId     int64
	liveRoomId int64
	tick       time.Time
	limit      time.Duration
}

type OnlineLiveManager struct {
	sync.Mutex
	wg      sync.WaitGroup
	liveMap map[int64]*OnlineLive
}

var SP_onlineLiveMgr = &OnlineLiveManager{liveMap: make(map[int64]*OnlineLive)}

func init() {
	time.AfterFunc(PreCheckPushStatus, func() { SP_onlineLiveMgr.GC() })
	time.AfterFunc(UpdatePlayNumTimer, func() { SP_onlineLiveMgr.UpdateNum() })
}

func (om *OnlineLiveManager) UpdateNum() {
	var ids []int64
	om.Lock()
	for id, _ := range om.liveMap {
		ids = append(ids, id)
	}
	om.Unlock()

	for _, id := range ids {
		om.wg.Add(1)

		go func(id int64) {
			defer om.wg.Done()

			for i := 0; i < 2; i++ {
				apptoken, err := GotyeAccessAppToken()
				if err != nil {
					logger.Error("UpdateNum : GotyeAccessAppToken Failed, ", err.Error())
					return
				}

				num, status, err := GotyeGetLiveContext(apptoken, id)
				if err != nil {
					logger.Error("UpdateNum : GotyeGetLiveContext Failed, ", err.Error())
					return
				}

				switch status {
				case gotye_sdk.API_SUCCESS:
					DBUpdateOnlineLiveRoom(id, num)
					return

				case gotye_sdk.API_INVALID_TOKEN_ERROR:
					if i == 0 {
						GotyeClearAccessToken()
						logger.Info("UpdateNum : invalid token error, and accesstoken again.")
					} else {
						logger.Error("UpdateNum : why access new token, but return invalid")
						return
					}

				default:
					logger.Error("UpdateNum : GotyeGetLiveContext status=%d", status)
					return
				}
			}

		}(id)

		om.wg.Wait()
	}

	time.AfterFunc(PreCheckPushStatus, func() { om.UpdateNum() })
}

func (om *OnlineLiveManager) GC() {
	om.Lock()
	defer om.Unlock()

	delList := make([]int64, 0, len(om.liveMap))
	now := time.Now()
	for id, online := range om.liveMap {
		if now.Sub(online.tick) > online.limit {
			delList = append(delList, id)
		}
	}
	for _, id := range delList {
		logger.Debug("GC : Overdue liveroom_id = ", id)
		delete(om.liveMap, id)
		DBDelOnlineLiveRoom(id)
	}
	time.AfterFunc(PreCheckPushStatus, func() { om.GC() })
}

func (om *OnlineLiveManager) StartPushStream(liveroomId int64, limit int) {
	om.Lock()
	defer om.Unlock()

	online, ok := om.liveMap[liveroomId]
	if ok {
		logger.Info("StartPushStream : update liveroomid=", liveroomId)
		online.tick = time.Now()
		return
	}

	limitTime := PushingLimitTimeOut
	if limit > 0 {
		limitTime = time.Second * time.Duration(limit)
	}
	online = &OnlineLive{liveroomId, time.Now(), limitTime}
	logger.Infof("StartPushStream : new liveroomId=%d,limit=%ss", online.liveRoomId, online.limit)
	om.liveMap[liveroomId] = online
	DBAddOnlineLiveRoom(liveroomId)
}

func (om *OnlineLiveManager) StopPushStream(liveroomId int64) {
	om.Lock()
	defer om.Unlock()

	_, ok := om.liveMap[liveroomId]
	if !ok {
		logger.Warn("StopPushStream : not found liveroomid=", liveroomId)
		return
	}
	logger.Info("StopPushStream : liveroomid=", liveroomId)
	delete(om.liveMap, liveroomId)
	DBDelOnlineLiveRoom(liveroomId)
}
