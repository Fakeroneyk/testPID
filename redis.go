package main

import "sync/atomic"

type Num struct {
	submitNum    int64
	completerNum int64
}

var redisNum *Num

func init(){
	redisNum = &Num{
		submitNum: 0,
		completerNum: 0,
	}
}

func (t *Num) SubmitOne() {
	atomic.AddInt64(&t.submitNum,1)
}

func (t *Num) CompleteOne() {
	atomic.AddInt64(&t.completerNum,1)
}

func (t *Num) GetSubmitNum() int64{
	return atomic.LoadInt64(&t.submitNum)
}

func (t *Num) GetCompleteNum() int64{
	return atomic.LoadInt64(&t.completerNum)
}

func GetRedisNum() *Num {
	return redisNum
}
