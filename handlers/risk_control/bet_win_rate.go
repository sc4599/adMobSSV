package risk_control

import (
	"awesomeProject/extension/global/redis"
	"github.com/sc4599/fireflygo/logger"
	"math/rand"
	"sort"
	"strconv"
)

type BetWinRate struct {
	BaseControl
	//  配置信息
	MaxBigCycle int64				// 大周期 上限
	MinBigCycle int64				// 大周期 下限
	MaxSmallCycle int64			// 小周期 上限
	MinSmallCycle int64			// 小周期 下限
	Rate map[int64]int64			// 	盈利比对应必杀率   例如    {"2": 60}  表示 盈利比2% 则 必杀率 60%
	betWinRateLevel []int64		// 真是玩家下注数量级别数组[10 8 5 2 0]
	UpdateAt	int64				// 配置更新时间
	// 游戏内通过随机后获得的数值
	BigCycle int64					// 大周期具体数值
	SmallCycle int64				// 小周期具体数值

	// 游戏内存数据
	RoundCount int64
	CumulativeBet int64			// 当前周期内累计投注
	CumulativeSysWin	int64		// 当前周期内系统累计输赢
}


func (t *BetWinRate)Check()int64{
	if t.IsOn == 0{
		return 0
	}
	// 如果此时达到大周期则，不出发任何 投注盈利比控制
	if t.BigCycle !=0 && t.RoundCount % t.BigCycle ==0 {
		return 0
	}
	if t.SmallCycle!=0 && t.RoundCount % t.SmallCycle == 0 {
		if t.CumulativeSysWin < 0 { // 当系统输赢数量为负数 则一定出发最大必杀率
			currPoint := int64(rand.Intn(101))// 通过随机检查此环节是否出发必杀
			t.KillRate = t.Rate[t.betWinRateLevel[len(t.betWinRateLevel)-1]]
			if currPoint < t.KillRate{
				return 1
			}
		}
		// 计算当前 投注盈利比
		betWinRate := t.CumulativeSysWin / t.CumulativeBet * 100
		for _, v:= range t.betWinRateLevel{
			if betWinRate > v{
				currPoint := int64(rand.Intn(101))// 通过随机检查此环节是否出发必杀
				t.KillRate = t.Rate[t.betWinRateLevel[0]]
				if currPoint < t.KillRate {
					return 1
				}
			}
		}
	}
	return 0
}
func (t *BetWinRate)setSelfConf(appId int64) int64{
	r, err := redis.GetBetWinRateConf(101)
	if err!=nil{
		logger.Error("RiskControl BetWinRate setSelfConf error err=", err)
		return 3
	}
	updateTime,has := r["UpdateAt"]
	if !has{
		logger.Error("RiskControl BetWinRate setSelfConf error not has")
		return 4
	}
	updateAt := int64(updateTime.(float64))
	// 通过配置时间，检查配置文件是否更新(更新时间是否有变动||是否初始化配置)
	if updateAt != t.UpdateAt || !t.IsInitConf{
		// 更新配置文件参数
		t.MaxBigCycle = int64(r["MaxBigCycle"].(float64))
		t.MinBigCycle = int64(r["MinBigCycle"].(float64))
		t.MaxSmallCycle = int64(r["MaxSmallCycle"].(float64))
		t.MinSmallCycle = int64(r["MinSmallCycle"].(float64))
		t.UpdateAt = updateAt
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
		sort.Slice(list, func(i, j int) bool {
			if list[i]<list[j]{
				return false
			}
			return true
		})
		t.Rate = ms
		t.betWinRateLevel = list
		t.resetBigCycle() // 重新初始化 投注盈利比控制 游戏内部数据
		return 0
	}
	// 检查是否达到大周期
	if t.RoundCount % t.BigCycle == 0{
		t.resetBigCycle()
	}
	if t.RoundCount % t.SmallCycle == 0 {
		t.resetSmallCycle()
	}
	return 0
}


// 重置大周期，  初始化 游戏轮数 重新随机获取 大周期 小周期， 初始化累计下注 和 累计系统输赢
func (t *BetWinRate)resetBigCycle(){
	t.RoundCount = 0
	if t.MaxBigCycle - t.MinBigCycle <= 0{
		t.BigCycle = t.MaxBigCycle
	}else{
		t.BigCycle = t.MinBigCycle + int64(rand.Intn(int(t.MaxBigCycle - t.MinBigCycle)))
	}
	if t.MaxSmallCycle - t.MinSmallCycle <= 0 {
		t.SmallCycle = t.MaxSmallCycle
	}else{
		t.SmallCycle = t.MinSmallCycle + int64(rand.Intn(int(t.MaxSmallCycle - t.MinSmallCycle)))
	}
	if t.BigCycle == 0 || t.SmallCycle==0{
		t.BigCycle =1
		t.SmallCycle=1
	}
	t.CumulativeBet = 0
	t.CumulativeSysWin = 0
	t.IsInitConf = true
}

// 重置小周期，  初始化累计下注 和 累计系统输赢
func (t *BetWinRate)resetSmallCycle(){
	t.CumulativeBet = 0
	t.CumulativeSysWin = 0
}

// 上报系统输赢 和 玩家总下注
func (t *BetWinRate)recordSysWin(currBet int64, currSysWin int64){
	t.CumulativeBet += currBet
	t.CumulativeSysWin += currSysWin
	t.RoundCount += 1
}