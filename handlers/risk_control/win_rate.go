package risk_control

import (
	"awesomeProject/extension/global/redis"
	"github.com/sc4599/fireflygo/logger"
	"math/rand"
	"sort"
	"strconv"
)

type WinRate struct {
	BaseControl
	Rate map[int64]int64
	betLevel []int64		// 真是玩家下注数量级别数组[1000 2000 10000 50000 200000]
}

/**
检查此风控是否出发必杀
@params userTotalBet int64
@return int64		0: 不出发必杀， 1：出发必杀
 */
func (t *WinRate)Check(userTotalBet int64)int64{
	if t.IsOn == 0{
		return 0
	}
	t.getKillRate(userTotalBet)
	currPoint := int64(rand.Intn(100))// 通过随机检查此环节是否出发必杀
	if currPoint < t.KillRate{
		return 1
	}
	return 0
}

func (t *WinRate)getKillRate(userTotalBet int64){
	for _,v := range t.betLevel{
		if userTotalBet >= v{
			t.KillRate = t.Rate[v]
			return
		}
	}
	t.KillRate=0
}

func (t *WinRate)setSelfConf(appId int64)int64{
	r, err:= redis.GetWinRateConf(appId)
	if err!=nil{
		logger.Error("RiskControl WinRate setSelfConf error err=", err)
		return 2
	}
	m := r["Rate"].(map[string]interface{})
	list := make([]int64, len(m))
	i:=0
	ms := make(map[int64]int64)
	for k,v := range m {
		if betCount,err :=strconv.Atoi(k);err==nil{
			list[i] = int64(betCount)
			ms[int64(betCount)] = int64(v.(float64))
		}
		i++
	}
	// 将查处来的数组排序
	sort.Slice(list, func(i, j int) bool {
		if list[i]<list[j]{
			return false
		}
		return true
	})
	t.betLevel = list
	t.Rate = ms
	return 0
}