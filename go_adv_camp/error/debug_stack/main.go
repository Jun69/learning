package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
)

/**
定义需要结构体
*/
type Teacher struct {
	Name    string
	ID      int
	Age     int
	Address string
}

func dump() (string,error) {

	s1 := Teacher{
		Name:    "Jason Yin",
		ID:      001,
		Age:     18,
		Address: "北京",

	}

	data, err := json.Marshal(&s1)
	if err != nil {
		debug.Stack()
		return "",err
	}
	return string(data),nil
}

func openFile()  {
	_,err:=os.Open("/int_vs_float/data/info")
	if err != nil {
		fmt.Printf("%s \n %s",err.Error(),debug.Stack())
	}
}

func main() {
	openFile()
}