package util

import (
	"fmt"
	"testing"
)

func TestTim(t *testing.T) {
	InitDatabase()
	time_s := "2019-02-21 14:58"
	st,_ := StrToTimeExpand(time_s)
	another_t := GetLastTime()
	lst,err := StrToTimeExpand(another_t)
	if err != nil {
		println("err!")
	}
	result := st.After(lst)
	println(result)
}

func TestStrToTime(t *testing.T) {
	time,e := StrToTimeExpand("2023-03-12 20:38")
	if e != nil {
		fmt.Println("Something went wrong:" + e.Error())
	} else {
		fmt.Println(time)
	}
}