package rate_limit

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
)

const (
	EPS  = 0.000001
)

func IsEqual(f1, f2 float64) bool {
	if f1 > f2 {
		return f1-f2 < EPS
	} else {
		return f2-f1 < EPS
	}
}

type TokenBucketLimiter struct {
	Limiter *rate.Limiter
	isLimit bool
}

var rateLimiter *TokenBucketLimiter

func GetBucketLimit(RateLimitVal float64) *TokenBucketLimiter {
	rateLimiter = newBucketLimit(RateLimitVal)
	return rateLimiter
}

func newBucketLimit(rateLimitVal float64) (limiter *TokenBucketLimiter) {
	limiter = &TokenBucketLimiter{
		Limiter: rate.NewLimiter(rate.Limit(rateLimitVal), int(rateLimitVal)),
	}
	if !IsEqual(rateLimitVal, 0.0) {
		limiter.isLimit = true
	}
	return
}

func (l *TokenBucketLimiter) Wait(ctx context.Context, doSubmit func(ctx context.Context, params ...interface{}) error, params ...interface{}) (err error) {
	// 尝试获取 token，没有获取到则阻塞
	if l.isLimit {
		err = l.Limiter.Wait(ctx)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	err = doSubmit(ctx, params...)
	return
}
