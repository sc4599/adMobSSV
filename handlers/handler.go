package handlers

import (
	"awesomeProject/extension/model"
	"fmt"
	"github.com/kataras/iris"
	"strconv"
	"strings"
)


func AdMobSuccessGetHandler(ctx iris.Context){
	fmt.Println(ctx.Request().Body)
	r:=ctx.GetReferrer()
	fmt.Println("GetReferrer=",r)
	fmt.Println("GetReferrer type=",r.Type)
	queryStr := ctx.Request().URL.RawQuery
	fmt.Println("queryStr=",queryStr)
	if !strings.Contains(queryStr, "&signature="){
		_, _ = ctx.JSON(ApiResource(2, nil, "signature pas error"))
		return
	}
	sinParams := queryStr[0:strings.Index(queryStr, "&signature=")]
	fmt.Println("sinParams=",sinParams)

	fmt.Println("Real-IP=",ctx.Request().Header.Get("X-Real-IP"))
	fmt.Println("HOST=",ctx.Request().Header.Get("Host"))
	fmt.Println("X-Forwarded-For=",ctx.Request().Header.Get("X-Forwarded-For"))
	//fmt.Println(ctx.Request().FormValue("ad_network"))
	//fmt.Println(ctx.Request().FormValue("ad_unit"))
	//fmt.Println(ctx.Request().FormValue("reward_amount"))
	//fmt.Println(ctx.Request().FormValue("reward_item"))
	//fmt.Println(ctx.Request().FormValue("timestamp"))
	//fmt.Println(ctx.Request().FormValue("transaction_id"))
	//fmt.Println(ctx.Request().FormValue("user_id"))
	//fmt.Println(ctx.Request().FormValue("signature"))
	signature:= ctx.Request().FormValue("signature")
	fmt.Println(ctx.Request().FormValue("key_id"))
	if Verification(sinParams, []byte(signature)){
		// TODO 此处可以调用游戏服务器的 增加广告激励对应奖励接口
		adUnitId := Conf.Get("admob.unitId") // 从配置文件获取广告ID
		fmt.Println("adUnit ID =",adUnitId)
		// TODO 判断广告单元ID
		if adUnitId == ctx.Request().FormValue("ad_unit"){
			// 处理此类广告的激励奖励逻辑
		}
		_, _ = ctx.JSON(ApiResource(0, nil, "success"))
	}else{
		_, _ = ctx.JSON(ApiResource(1, nil, "fail"))
	}

}

/*的摄像头
获取测试
*/
func GetJson(ctx iris.Context){
	data:=`{"keys":[{"keyId":3335741209,"pem":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+nzvoGqvDeB9+SzE6igTl7TyK4JB\nbglwir9oTcQta8NuG26ZpZFxt+F2NDk7asTE6/2Yc8i1ATcGIqtuS5hv0Q==\n-----END PUBLIC KEY-----","base64":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+nzvoGqvDeB9+SzE6igTl7TyK4JBbglwir9oTcQta8NuG26ZpZFxt+F2NDk7asTE6/2Yc8i1ATcGIqtuS5hv0Q=="}]}`
	_, _ = ctx.Text(data)
}


func GetRiskControl(ctx iris.Context){
	bet :=ctx.Request().FormValue("bet")
	if betInt,err:=strconv.Atoi(bet);err!=nil{
		_, _ = ctx.Text("bet格式错误,必须是一个正整数")
	}else{
		r, trigger := model.Rc.CheckRiskTrigger(int64(betInt))
		r1 :=strconv.Itoa(int(r))

		r2:=model.Rc.DebugGameStockInfo()
		r3:=model.Rc.DebugBetWinRateInfo()
		r4:=model.Rc.DebugWinRateInfo()
		_, _ = ctx.Text(r1+"【"+trigger+"】\n"+r2 +"\n" + r3+"\n" + r4)
	}
}

func NextRound(ctx iris.Context){
	currBet :=ctx.Request().FormValue("currBet")
	currSysWin :=ctx.Request().FormValue("currSysWin")
	bet,err:=strconv.Atoi(currBet)
	if err!=nil{
		_, _ = ctx.Text("bet格式错误,必须是一个正整数")
		return
	}
	sysWin,err:=strconv.Atoi(currSysWin)
	if err!=nil{
		_, _ = ctx.Text("bet格式错误,必须是一个正整数")
		return
	}
	r:=model.Rc.ResetRiskControlByRoundEnd(int64(bet),int64(sysWin))
	switch r {
	case 1:
		_, _ = ctx.Text("Switch get error")
	case 2:
		_, _ = ctx.Text("redis winRate get error")
	case 3:
		_, _ = ctx.Text("RiskControl BetWinRate setSelfConf error")
	case 4:
		_, _ = ctx.Text("RiskControl BetWinRate setSelfConf error not has")
	default :
		_, _ = ctx.Text("OK")
	}

}