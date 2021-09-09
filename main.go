package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	wg                  sync.WaitGroup
	r                   *rand.Rand
	file			 	*os.File
	write 				*bufio.Writer
)

func Process(ctx context.Context, taskChan chan *int) {
	for true{
		_, ok := <- taskChan
		if !ok {
			fmt.Println("get taskChan error")
			return
		}

		err := tokenBucketLimiter.Wait(ctx,doSubmitTask)
		if err != nil {
			fmt.Println(err)
			time.Sleep(20 * time.Millisecond)
		}
		wg.Done()
	}
}

var msgChan = make(chan *int, 50000)

func main() {
	var err error
	file, err = os.OpenFile(time.Now().Format("2006-01-02 15:04:05") + ".txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return
	}
	defer file.Close()
	write = bufio.NewWriter(file)
	go UpdateBucketNum()

	taskChan := make(chan *int, 2 * SubmitterGoRoutineLimit)
	for i := 0; i < RenderNode; i++ {
		go Render(msgChan)
	}
	ctx := context.Background()
	for i := 0; i < SubmitterGoRoutineLimit; i++ {
		go Process(ctx, taskChan)
	}
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	//模拟动态限流上送
	for true {
		taskNum := r.Intn(SubmitterGoRoutineLimit) + int(SubmitterGoRoutineLimit * 1.5)
		if taskNum > SubmitterGoRoutineLimit * 2 {
			taskNum = SubmitterGoRoutineLimit * 2
		}
		for i := 0; i < taskNum; i ++ {
			wg.Add(1)
			taskChan <- &i
		}
		wg.Wait()
	}
}

func doSubmitTask(ctx context.Context, params ...interface{}) (err error){
	num := r.Int()
	msgChan <- &num
	GetRedisNum().SubmitOne()
	return
}
