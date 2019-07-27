package handlers

import (
	"fmt"
	"github.com/kataras/iris"
	"strings"
)


func AdMobSuccessGetHandler(ctx iris.Context){
	fmt.Println(ctx.Request().Body)
	r:=ctx.GetReferrer()
	fmt.Println("GetReferrer=",r)
	fmt.Println("GetReferrer type=",r.Type)
	queryStr := ctx.Request().URL.RawQuery
	fmt.Println("queryStr=",queryStr)
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
