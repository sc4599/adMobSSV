package risk_control

import (
	"fmt"
	"testing"
	"awesomeProject/extension/global"
)

func TestNewRiskControl(t *testing.T) {
	global.Conf = global.NewConf()
	global.RedisPool = global.InitRedisPool()
	rc := NewRiskControl(101)  // 初始化 风险控制器
	// 一共两个方法给游戏使用
	rc.ResetRiskControlByRoundEnd(2000,-1000)
	for i:=0;i<100;i++{
		r,_ := rc.CheckRiskTrigger(2000)
		fmt.Println("RsikControl r=",r)
	}

}
