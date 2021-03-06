package main

import (
	"fakerone/rate_limit"
	"fmt"
	"golang.org/x/time/rate"
	"strconv"
	"sync/atomic"
	"time"
)

var tokenBucketLimiter *rate_limit.TokenBucketLimiter

var bucketNum int64

var subNums = make([]int64, 0, PIDDataNum)

var cnt = 0

func GetNowBucketNum() (num int64) {
	return atomic.LoadInt64(&bucketNum)
}

func init() {
	bucketNum = BucketDefaultNum
	tokenBucketLimiter = rate_limit.GetBucketLimit(float64(bucketNum))
}

func sum() float64 {
	sumNum := int64(0)
	for _, num := range subNums {
		sumNum += num
	}
	return float64(sumNum)
}

func sub() float64 {
	return float64(subNums[PIDDataNum-1] - subNums[PIDDataNum-2])
}

func getNum() {
	if len(subNums) < PIDDataNum {
		return
	}
	up := Kp * float64(subNums[PIDDataNum-1])
	ui := Ki * sum()
	ud := Kd * sub()
	uk := int64(up + ui + ud)
	fmt.Printf("up:%v ui:%v ud:%v uk:%v\n", up, ui, ud, uk)
	if uk < 1 {
		uk = 1
	}
	atomic.StoreInt64(&bucketNum, uk)
	cnt++
	if cnt%3 == 0 {
		cnt = 0
		tokenBucketLimiter.Limiter.SetLimit(rate.Limit(uk))
		tokenBucketLimiter.Limiter.SetBurst(int(uk))
		fmt.Println(uk)
	}
	//fmt.Println(tokenBucketLimiter.Limiter.Burst(),tokenBucketLimiter.Limiter.Limit())
}

//动态更新令牌桶大小
func UpdateBucketNum() {
	//定时器，每两秒执行一次
	timeTicker := time.Tick(UpdateTick * time.Second)
	for true {
		submitNum := GetRedisNum().GetSubmitNum()
		completeNum := GetRedisNum().GetCompleteNum()
		redisSubNum := submitNum - completeNum
		subNum := MaxQueueNum - redisSubNum
		if len(subNums) < PIDDataNum {
			subNums = append(subNums, subNum)
		} else {
			for i := 0; i < PIDDataNum-1; i++ {
				subNums[i] = subNums[i+1]
			}
			subNums[PIDDataNum-1] = subNum
		}
		getNum()
		fmt.Println("message: ", redisSubNum, "  v:", GetNowBucketNum(), " sub: ", submitNum, "complete: ", completeNum)
		write.WriteString(strconv.FormatInt(redisSubNum, 10) + "\n")
		write.Flush()
		<-timeTicker
	}
}
