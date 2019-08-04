package main

import (
	"errors"
	"fmt"
	"time"
)

func innerFunc() {
	fmt.Println("Enter innerFunc")
	panic(errors.New("Occur a panic!"))
	fmt.Println("Quit innerFunc")
}

func outerFunc() {
	fmt.Println("Enter outerFunc")
	innerFunc()
	fmt.Println("Quit outerFunc")
}

func main() {
	defer func(){
		if p:= recover(); p!=nil{
			fmt.Printf("Fatal error: %v", p)
		}
	}()
	fmt.Println("Enter main")
	outerFunc()
	fmt.Println("Quit main")
	go func() {
		fmt.Println("1")
	}()
	go func() {
		fmt.Println("2")
	}()
	go func() {
		fmt.Println("3")
	}()
	time.Sleep(100)
}
