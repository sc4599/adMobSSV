package redis

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sort"
	"strconv"
	"testing"
	"xggame/extension/global"
)

type Game struct {
	GameType string		`json:"gameType"`
	GameName string		`json:"gameName"`
	GameImg		string	`json:"gameImg"`
}

func TestGetGameList(t *testing.T) {
	global.Conf= global.NewConf()
	global.RedisCon = global.InitRedis()
	 r, err := redis.Strings(global.RedisCon.Do("ZRANGE", "GameList", 0, -1))
	 fmt.Println("1:",r[0])
	 if err!=nil{
		t.Fatal(err)
	}
	var games []Game
	 for _, item := range r{
	 	g := Game{}
		_ = json.Unmarshal([]byte(item), &g)
		fmt.Println(g)
		games = append(games, g)
	 }
	fmt.Println(games)
}

func TestRedis(t *testing.T) {
	global.Conf= global.NewConf()
	global.RedisCon = global.InitRedis()

}

/*
批量获取所有与风险控制相关控件
*/
func TestRedisPipelineGetRiskControl(t *testing.T) {

	global.Conf= global.NewConf()
	global.RedisCon = global.InitRedis()
	//riskSwitch ,err:=redis.StringMap(global.RedisCon.Do("hgetall", "RiskControl:Switch"))

	_ = global.RedisCon.Send("hgetall", "RiskControl:Switch")
	_ = global.RedisCon.Send("hgetall", "RiskControl:BetWinRate")
	_ = global.RedisCon.Send("hgetall", "RiskControl:WinRateOption")
	_ = global.RedisCon.Flush()
	riskSwitch, err := redis.StringMap(global.RedisCon.Receive())
	if err==nil{
		fmt.Println("this redis return1:",riskSwitch)
	}
	riskBetWinRate, err := redis.StringMap(global.RedisCon.Receive())
	if err==nil{
		fmt.Println("this redis return2:",riskBetWinRate)
	}
	riskWinRateOption, err := redis.StringMap(global.RedisCon.Receive())
	if err==nil{
		fmt.Println("this redis return3:",riskWinRateOption)
	}
}

func TestGetRiskSwitch(t *testing.T) {
	global.Conf= global.NewConf()
	global.RedisPool = global.InitRedisPool()
	r, err:= GetRiskSwitch(101)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(r)
}

func TestGetWinRateConf(t *testing.T) {
	global.Conf= global.NewConf()
	global.RedisPool = global.InitRedisPool()
	r, err:= GetWinRateConf(101)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(r)
	m := r["Rate"].(map[string]interface{})
	list := make([]uint64, len(m))
	i:=0
	ms := make(map[uint64]uint64)
	for k,v := range m {
		if betCount,err :=strconv.Atoi(k);err==nil{
			list[i] = uint64(betCount)
			ms[uint64(betCount)] = uint64(v.(float64))
		}
		i++
		fmt.Println(k,v)
	}
	fmt.Println("map keys=", list)
	sort.Slice(list, func(i, j int) bool {
		if list[i]<list[j]{
			return false
		}
		return true
	})
	fmt.Println("sort list=", list)
	fmt.Println("ms=", ms)

	for _,v:= range list{
		fmt.Println(v)
	}
}

func TestGetBetWinRateConf(t *testing.T) {
	global.Conf = global.NewConf()
	global.RedisPool = global.InitRedisPool()
	r, err := GetBetWinRateConf(101)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(r)
	updateAt := r["UpdateAt"]
	fmt.Println("updateAt=",uint64(updateAt.(float64)))
	m := r["Rate"].(map[string]interface{})
	list := make([]uint64, len(m))
	i:=0
	ms := make(map[uint64]uint64)
	for k,v := range m {
		if betCount,err :=strconv.Atoi(k);err==nil{
			list[i] = uint64(betCount)
			ms[uint64(betCount)] = uint64(v.(float64))
		}
		i++
		fmt.Println(k,v)
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i]<list[j]{
			return false
		}
		return true
	})
	fmt.Println("sort list=", list)
	fmt.Println("ms=", ms)
	var a uint64
	fmt.Println("ms=", r["MaxBigCycle"])
	a = uint64(r["MaxBigCycle"].(float64))
	fmt.Println("a", a)
}

func TestRemainder(t *testing.T) {
	for i:=0;i<100;i++{
		fmt.Println(i%4)
	}
}

func TestGetGameStock(t *testing.T) {
	global.Conf = global.NewConf()
	global.RedisPool = global.InitRedisPool()
	r, err := GetGameStock(101)
	if err!=nil{
		t.Fatal(err)
	}
	var r1 []interface{}
	r1 = r["StockLevel"].([]interface{})
	fmt.Println(r1)
	r2 := r["Rate"].(interface{}).([]interface{})
	r3:=make([]map[string]interface{},len(r2))
	for i,v := range r2{
		r3[i] = v.(map[string]interface{})
	}
	fmt.Println(r3)
}

func TestSetGameStockPool(t *testing.T) {

	err:=SetGameStockPool(101, 10000,300)
	if err!=nil{
		t.Fatal(err)
	}
}