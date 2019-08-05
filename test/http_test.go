package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
)

func TestADMob(t *testing.T) {
	googleAdMod := "?ad_network=5450213213286189855&ad_unit=12345678&reward_amount=10&reward_item=coins&timestamp=1507770365237823&transaction_id=1234567890ABCDEF1234567890ABCDEF&user_id=1234567&signature=MEUCIQDGx44BZgQU6TU4iYEo1nyzh3NgDEvqNAUXlax-XPBQ5AIgCXSdjgKZvs_6QNYad29NJRqwGIhGb7GfuI914MDDZ1c&key_id=1268222887"
	resp, err:=http.Get("http://18.136.172.78:9080/adMobSuccess"+ googleAdMod)
	if err !=nil{
		t.Fatal(err)
	}

	defer resp.Body.Close()
	body,err:= ioutil.ReadAll(resp.Body)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(body)

}

func TestMap(t *testing.T) {
	a := map[string]interface{}{
		"name":"song",
		"age":"13",
	}
	b := []map[string]interface{}{}
	b=append(b,a)
	for k,v:= range b {
		fmt.Println(k,v)
	}
}

func TestChan(t *testing.T) {
	ch4 := make(chan int, 1)
	for i := 0; i < 4; i++ {
		select {
		case e, ok := <-ch4:
			if !ok {
				fmt.Println("End.")
				return
			}
			fmt.Println(e)
			close(ch4)

		default:
			fmt.Println("No Data!")
			ch4<-1
		}
	}
}

func TestDefer(t *testing.T) {
	for i := 1; i < 5; i++ {
		defer fmt.Println(i)
	}
}

func TestMap2(t *testing.T) {
	m:= make(map[int]bool, 2)
	m[1]=true
	m[2]=false
	fmt.Println(m)
	m[3]=true
	fmt.Println(m)
}

func TestDefer2(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%d ", fibonacci(i))
		defer func(n int){
			fmt.Printf("%d ",n)
		}(fibonacci(i))
	}
}
func fibonacci(num int) int {
	if num == 0 {
		return 0
	}
	if num < 2 {
		return 1
	}
	return fibonacci(num-1) + fibonacci(num-2)
}

func TestRandom(t *testing.T) {
	for i:=0;i<100;i++{
		d :=rand.Intn(101)
		fmt.Println(d)
	}

}

