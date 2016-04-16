package service

import (
	"gotye_protocol/gotye_sdk"
	"sync"
	"time"

	"github.com/futurez/litego/logger"
)

const (
	LimitTimeOut  = time.Minute * 1
	LimitPushTime = time.Minute * 1
	PreRefreshNum = 20
)

type OnlineLive struct {
	liveRoomId    int64
	tick          time.Time
	limit         time.Duration
	playUserCount int
	playRtmpUrls  string
	playHlsUrls   string
	playFlvUrls   string
}

type OnlineLiveManager struct {
	sync.RWMutex
	wg      sync.WaitGroup
	liveMap map[int64]*OnlineLive
}

var SP_onlineLiveMgr = &OnlineLiveManager{liveMap: make(map[int64]*OnlineLive)}

func init() {
	time.AfterFunc(LimitTimeOut, func() { SP_onlineLiveMgr.GC() })
}

func (om *OnlineLiveManager) GC() {
	var inactiveIds []int64
	var activeIds []int64
	now := time.Now()

	om.RLock()
	for id, online := range om.liveMap {
		if now.Sub(online.tick) > online.limit {
			inactiveIds = append(inactiveIds, id)
		} else {
			activeIds = append(activeIds, id)
		}
	}
	om.RUnlock()

	go func() {
		om.Lock()
		for _, id := range inactiveIds {
			logger.Debug("GC : Overdue liveroom_id = ", id)
			delete(om.liveMap, id)
			DBDelOnlineLiveRoom(id)
		}
		om.Unlock()
	}()

	total := len(activeIds)
	if total > 0 {
		limit := (total + PreRefreshNum) / PreRefreshNum * PreRefreshNum
		logger.Debug("GC : actvieIds=", activeIds, ",total=", total, ",limit=", limit)

		for i := 0; i < limit; i += PreRefreshNum {

			var ids []int64

			if i != (limit - PreRefreshNum) {
				ids = activeIds[i : i+PreRefreshNum]
				logger.Debug("GC : last i=", i, ",ids=", ids)
			} else {
				ids = activeIds[i:total]
				logger.Debug("GC : pre i=", i, ",ids=", ids)
			}

			om.wg.Add(1)
			go func(ids []int64) {
				defer om.wg.Done()

				resp, err := GotyeGetRoomsLiveInfo(ids...)
				if err != nil {
					logger.Error("UpdateNum : GotyeGetLiveContext Failed, ", err.Error())
					return
				}

				if resp.Status == gotye_sdk.API_SUCCESS {
					om.UpStreamInfo(resp.Entities)
				} else {
					logger.Error("UpdateNum : GotyeGetLiveContext status=", resp.Status)
				}

			}(ids)
		}
		om.wg.Wait()
	}

	time.AfterFunc(LimitTimeOut, func() { om.GC() })
}

func (om *OnlineLiveManager) UpStreamInfo(entities []gotye_sdk.LiveRoomInfo) {
	om.Lock()
	defer om.Unlock()

	for _, entity := range entities {
		online, ok := om.liveMap[entity.RoomId]
		if !ok {
			continue
		}

		up := false

		if len(entity.PlayRtmpUrls) > 0 && entity.PlayRtmpUrls[0] != online.playRtmpUrls {
			online.playRtmpUrls = entity.PlayRtmpUrls[0]
			up = true
		}

		if len(entity.PlayHlsUrls) > 0 && entity.PlayHlsUrls[0] != online.playHlsUrls {
			online.playHlsUrls = entity.PlayHlsUrls[0]
			up = true
		}

		if len(entity.PlayFlvUrls) > 0 && entity.PlayFlvUrls[0] != online.playFlvUrls {
			online.playFlvUrls = entity.PlayFlvUrls[0]
			up = true
		}

		if up {
			logger.Info("UpStreamInfo : up liveroomId=", entity.RoomId)
			DBUpdateLiveroomUrls(entity.RoomId, online.playRtmpUrls, online.playHlsUrls, online.playFlvUrls)
		}

		if entity.PlayUserCount != online.playUserCount {
			logger.Info("UpStreamInfo : up liveroomId=", entity.RoomId, ", number=", online.playUserCount)
			online.playUserCount = entity.PlayUserCount
			DBUpdateOnlineLiveRoom(entity.RoomId, entity.PlayUserCount)
		}
	}
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

	limitTime := LimitPushTime
	if limit > 0 {
		limitTime = time.Second * time.Duration(limit)
	}
	online = &OnlineLive{}
	online.liveRoomId = liveroomId
	online.tick = time.Now()
	online.limit = limitTime
	online.playUserCount = 0
	online.playRtmpUrls, online.playHlsUrls, online.playFlvUrls = DBGetLiveroomUrls(liveroomId)

	logger.Infof("StartPushStream : start liveroomId=%d,limit=%ss,playRtmpUrls=%s,playHlsUrls=%s,playFlvUrls=%s",
		online.liveRoomId, online.limit, online.playRtmpUrls, online.playHlsUrls, online.playFlvUrls)
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
