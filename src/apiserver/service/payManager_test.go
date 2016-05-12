package service

import (
	"sync"
	//	"gotye_protocol"
	"testing"
)

func TestDBChargeRMB(t *testing.T) {
	InitMysqlDbPool("192.168.1.10", "gotye_open_live", "appuser", "gotye2013")

	var wait sync.WaitGroup

	for i := 0; i < 20; i++ {
		wait.Add(1)
		go func(n int) {
			defer wait.Done()
			qin_coin, err := dbChargeRMB(5, 1000)
			if err != nil {
				t.Error("TestdbChargeRMB = ", err.Error())
				return
			}
			t.Log("n=", n, ": userid=5 pay 1000rmb, total_qin_coin=", qin_coin)
		}(i)
	}
	wait.Wait()
}

func TestDBUpdateJiaCoin(t *testing.T) {
	InitMysqlDbPool("192.168.1.10", "gotye_open_live", "appuser", "gotye2013")

	var wait sync.WaitGroup

	for i := 0; i < 20; i++ {
		wait.Add(1)
		go func(n int) {
			defer wait.Done()
			errno := dbUpdateJiaCoin(90, 5, 1000)
			if errno == -1 {
				t.Error("TestDBUpdateJiaCoin = ", errno)
				return
			}
			t.Log("n=", n, ": userid=5 pay 1000qincoin to userid=90, errno=", errno)
		}(i)
	}
	wait.Wait()
}
