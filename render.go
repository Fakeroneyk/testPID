package main

import (
	"fmt"
	"time"
)

//模拟下游渲染
func Render(msgChan chan *int) (err error) {
	for true{
		num, ok := <- msgChan
		if !ok {
			fmt.Println("get msgChan error")
		}
		//模拟渲染耗时区间为 [RenderMinTime-1,RenderMaxTime-1]
		SleepNum := *num % RenderMinNum + RenderMaxNum - RenderMinNum
		time.Sleep(time.Duration(SleepNum) * time.Second)

		GetRedisNum().CompleteOne()
	}
	return
}
