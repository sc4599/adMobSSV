package handlers

import (
	"awesomeProject/net"
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
	//signature:= ctx.Request().FormValue("signature")
	fmt.Println(ctx.Request().FormValue("key_id"))
	userId := ctx.Request().FormValue("user_id")
	customData := ctx.Request().FormValue("custom_data")
	soul := "cptbtptp"
	url := "http://im.t3e.cn/pay/admobSucess"

	if Md5V2(userId+soul) == customData{
		// TODO 此处可以调用游戏服务器的 增加广告激励对应奖励接口
		data := map[string]interface{}{
			"user_id": userId,
			"ad_unit": ctx.Request().FormValue("ad_unit"),
		}
		fmt.Println("http data=",data)
		r:=net.Post(url,data, "text/plain")
		_, _ = ctx.JSON(ApiResource(0, r, "success"))
	}else{
		_, _ = ctx.JSON(ApiResource(1, nil, "fail"))
	}

}



func TestAdMobHandler(ctx iris.Context){
	addr:= ctx.RemoteAddr()
	if addr != "127.0.0.1"{
		_, _ = ctx.JSON(ApiResource(0, nil, "success1"))
		return
	}
	url := "http://im.t3e.cn/pay/admobSucess"
	data := map[string]interface{}{
		"user_id": 1,
		"ad_unit": "12312312",
	}
	r:=net.Post(url,data, "text/plain")
	_, _ = ctx.JSON(ApiResource(0, r, "success"))

}
/*的摄像头
获取测试
*/
func GetJson(ctx iris.Context){
	data:=`{"keys":[{"keyId":3335741209,"pem":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+nzvoGqvDeB9+SzE6igTl7TyK4JB\nbglwir9oTcQta8NuG26ZpZFxt+F2NDk7asTE6/2Yc8i1ATcGIqtuS5hv0Q==\n-----END PUBLIC KEY-----","base64":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE+nzvoGqvDeB9+SzE6igTl7TyK4JBbglwir9oTcQta8NuG26ZpZFxt+F2NDk7asTE6/2Yc8i1ATcGIqtuS5hv0Q=="}]}`
	_, _ = ctx.Text(data)
}

