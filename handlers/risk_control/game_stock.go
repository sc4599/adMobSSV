package risk_control

import (
	"awesomeProject/extension/global/redis"
	"github.com/sc4599/fireflygo/logger"
	"math"
	"math/rand"
)

type GameStock struct {
	BaseControl
	RoundCount    int64
	StandardStock int64 // 标准库存
	MinStock      int64 // 最低存库		库存保护
	KeepRound     int64 // 达到 最低存库	库存保护 持续触发轮数
	PumpPercent   int64 // 每轮系统输赢绝对值扣除百分比   例如  3  = 3%
	UpdateAt      int64 // 配置更新时间戳
	/*
		"Rate":[{"MinKeepRound":1,"MaxKeepRound":2,"KillRate":40},
		{"MinKeepRound":1,"MaxKeepRound":3,"KillRate":35},
		{"MinKeepRound":1,"MaxKeepRound":3,"KillRate":30},
		{"MinKeepRound":1,"MaxKeepRound":4,"KillRate":25},
		{"MinKeepRound":1,"MaxKeepRound":4,"KillRate":20},
		{"MinKeepRound":2,"MaxKeepRound":5,"KillRate":15},
		{"MinKeepRound":1,"MaxKeepRound":4,"KillRate":10},
		{"MinKeepRound":1,"MaxKeepRound":4,"KillRate":0},
		{"MinKeepRound":1,"MaxKeepRound":4,"KillRate":-10},
		{"MinKeepRound":1,"MaxKeepRound":3,"KillRate":-20},
		{"MinKeepRound":1,"MaxKeepRound":2,"KillRate":-30}]
	*/
	Rate       []map[string]interface{} // 不同级别触发配置表
	StockLevel []interface{}                 // 库存级别  例如：[6000 7500 8500 9000 9500 10500 11000 11500 12500 14000]

	//游戏内存数据============================================
	CurrKeepRound int64 // 剩余持续轮数 当次数值存在时直接触发对应必杀率，游戏每循环一轮此数值则-1
	// 每局上报redis 数据
	CurrStock int64 // 开关开放后   到目前的库存（系统每轮输赢值的累计-每轮扣除）
	TotalPump int64 // 每轮按照系统输赢绝对值总和的百分之PumpPercent 的累加
}

func (t *GameStock) Check() int64 {
	return 0
}

func (t *GameStock) CheckContinuousTriggering() int64 {
	if t.IsOn == 0{
		return 0
	}
	if t.CurrKeepRound>0{
		currPoint := int64(rand.Intn(100))// 通过随机检查此环节是否出发必杀
		if t.KillRate > 0{
			if currPoint < t.KillRate{
				return 1
			}
		}else if t.KillRate < 0{
			if currPoint < int64(math.Abs(float64(t.KillRate))){
				return 2
			}
		}
	}
	return 0
}

/**
从redis 获取 风控配置
*/
func (t *GameStock) setSelfConf(appId int64) int64 {

	r, err := redis.GetGameStock(101)
	if err != nil {
		logger.Error("RiskControl GameStock setSelfConf error err=", err)
		return 4
	}
	if updateAt, has := r["UpdateAt"]; has {
		if t.UpdateAt != updateAt || !t.IsInitConf { // 如果配置更新， 或者 没有初始化则更新
			if v, has := r["MinStock"]; has {
				t.MinStock = int64(v.(float64))
			}
			if v, has := r["KeepRound"]; has {
				t.KeepRound = int64(v.(float64))
			}
			if v, has := r["PumpPercent"]; has {
				t.PumpPercent = int64(v.(float64))
			}
			if v, has := r["StandardStock"]; has {
				t.StandardStock = int64(v.(float64))
			}
			if v, has := r["StockLevel"]; has {
				t.StockLevel = v.([]interface{})
			}
			if v, has := r["Rate"]; has {
				redisValue:=v.(interface{}).([]interface{})
				t.Rate = make([]map[string]interface{},len(redisValue))
				for i, v := range redisValue{
					t.Rate[i] = v.(map[string]interface{})
				}
			}
		}
	} else {
		logger.Error("RiskControl GameStock UpdateAt not has")
		return 5
	}

	return 0
}

func (t *GameStock) record(currBet, currWin, currLost int64, appId int64) {
	r,err := redis.GetGameStockPool(appId)
	if err != nil {
		logger.Error("RiskControl GameStock GetGameStockPool error err=", err)
		return
	}
	abs := math.Abs(float64(currWin)) + math.Abs(float64(currLost)) // 系统输赢绝对值总和
	currPump := int64(abs * float64(t.PumpPercent) / 100)           // 当前轮系统扣除数量
	currSysWin := currWin - currLost - currPump                     // 计算出本轮系统输赢数量
	t.CurrStock = int64(r["CurrStock"].(float64))					// 从redis获取
	t.TotalPump = int64(r["TotalPumpAmount"].(float64))
	t.CurrStock += currSysWin                                       // 累加当前库存
	t.TotalPump += currPump                                         // 累加当前轮系统扣除数量
	t.RoundCount += 1                                               // 开关打开到现在储蓄多少轮
	if t.CurrKeepRound > 0 {
		t.CurrKeepRound -= 1
	}else{
		if t.CurrStock >= int64(t.StockLevel[len(t.StockLevel)-1].(float64)){  // >14000
			t.KillRate = int64(t.Rate[0]["KillRate"].(float64))
			minKeepRound := int64(t.Rate[0]["MinKeepRound"].(float64))
			maxKeepRound := int64(t.Rate[0]["MaxKeepRound"].(float64))
			t.CurrKeepRound = minKeepRound + int64(rand.Intn(int(maxKeepRound-minKeepRound)))
		}else{
			for i,v := range t.StockLevel{
				if t.CurrStock < int64(v.(float64)){
					t.KillRate = int64(t.Rate[i]["KillRate"].(float64))
					minKeepRound := int64(t.Rate[i]["MinKeepRound"].(float64))
					maxKeepRound := int64(t.Rate[i]["MaxKeepRound"].(float64))
					t.CurrKeepRound = minKeepRound + int64(rand.Intn(int(maxKeepRound-minKeepRound)))
				}
			}
		}
	}


	if err:=redis.SetGameStockPool(appId, t.CurrStock, t.TotalPump);err!=nil{
		logger.Error("RiskControl GameStock SetGameStockPool error err=", err)
	}

}
