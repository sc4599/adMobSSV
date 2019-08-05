package risk_control

import (
	"awesomeProject/extension/global/redis"
	"encoding/json"
	"fmt"
	"github.com/viphxin/xingo/logger"
)

type RiskControl struct {
	AppId int64
	betWinRate BetWinRate  // 投注盈利比控制
	gameStock GameStock		// 游戏库存控制
	winRate  WinRate		// 胜率控制(根据真实玩家下注数控制胜率)
}

func NewRiskControl (appId int64) RiskControl{
	riskControl := RiskControl{appId,BetWinRate{RoundCount:1},GameStock{RoundCount:1}, WinRate{}}
	return riskControl
}

/*
重置每一轮结束时, 调用此函数重置控制数据
@params currBet int64, 当前轮真实玩家总下注
@params currSysWin int64	当前轮系统输赢
*/
func (t *RiskControl)ResetRiskControlByRoundEnd(currBet int64, currWin, currLost int64)int64{
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Sprintf("ResetRiskControlByRoundEnd err: %v", err))
		}
	}()
	// 从redis获取开关状态
	if r:=t.setSwitchFromRedis();r!=0{
		return 1
	}
	if r:=t.setControlConf();r!=0{
		return r
	}
	t.recordMsg(currBet, currWin, currLost)
	return 0
}

/**
检查是否会触发风险控制
@params userTotalBet
@return int64   0:表示正常随机， 1：杀大开小， 2：杀小开大
*/
func (t *RiskControl)CheckRiskTrigger(userTotalBet int64)(int64, string){
	defer func() {
		if err := recover(); err != nil {
			logger.Error(fmt.Sprintf("CheckRiskTrigger err: %v", err))
		}
	}()
	r0:= t.gameStock.CheckContinuousTriggering()
	if r0!=0{
		logger.Debug("CheckRiskTrigger gameStock CheckContinuousTriggering Trigger kill")
		return r0, "游戏库存触发"
	}
	r1:=t.betWinRate.Check()
	if r1!=0{
		logger.Debug("CheckRiskTrigger betWinRate Trigger kill")
		return r1, "投注盈利比触发"
	}
	r2:=t.winRate.Check(userTotalBet)
	if r2!=0{
		logger.Debug("CheckRiskTrigger winRate Trigger kill")
		return r2, "胜率触发"
	}
	return 0, ""
}

func (t *RiskControl)DebugBetWinRateInfo()string{
	str,_:=json.Marshal(t.betWinRate)
	return string(str)
}

func (t *RiskControl)DebugGameStockInfo()string{
	str,_:=json.Marshal(t.gameStock)
	return string(str)
}

func (t *RiskControl)DebugWinRateInfo()string{
	str,_:=json.Marshal(t.winRate)
	return string(str)
}

func (t *RiskControl)setSwitchFromRedis() int64{
	r, err:= redis.GetRiskSwitch(t.AppId)
	if err!=nil{
		logger.Error("RiskControl GetRiskSwitch error err=", err)
		return 1
	}
	t.gameStock.IsOn = r["GameStock"]
	t.winRate.IsOn = r["WinRate"]
	t.betWinRate.IsOn = r["BetWinRate"]
	return 0
}

//发现开关被关闭，择重置相关控制器内的数据
func (t *RiskControl)resetSwitchOffStatus(){

}
// 根据开关状态加载风控内其他配置
func (t *RiskControl) setControlConf() int64{
	if t.winRate.IsOn == 1{
		//TODO 加载胜率控制器的数值和必杀概率
		r:=t.winRate.setSelfConf(t.AppId)
		if r!=0{
			return r
		}
	}

	if t.betWinRate.IsOn == 1{
		//TODO 加载投注盈利比控制器的数值和必杀概率
		r:=t.betWinRate.setSelfConf(t.AppId)
		if r!=0{
			return r
		}
	}else{
		t.betWinRate.IsInitConf = false
	}

	if t.gameStock.IsOn == 1{
		//TODO 加载库存控制器的数值和必杀概率
		r:=t.gameStock.setSelfConf(t.AppId)
		if r!=0{
			return r
		}
	}
	return 0
}

// 每轮上报玩家总押注 和 系统输赢 给风控
func (t *RiskControl)recordMsg(currBet int64, currWin, currLost  int64){
	t.betWinRate.recordSysWin(currBet, currWin-currLost)
	t.gameStock.record(currBet, currWin, currLost,t.AppId)
}



/**
牌局结束 上报系统输赢结果
@params totalWin int64		系统此轮游戏 盈利筹码
@params totalLost int64	系统此轮游戏 赔付筹码
 */
func (t *RiskControl)SetSystemWinLost(totalWin int64, totalLost int64){

}