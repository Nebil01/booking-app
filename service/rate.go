package service

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	MaxRequestPerWindow int
	RateLimitWindow     time.Duration
	mu                  sync.Mutex
	ipRequest           map[string][]time.Time
}

func NewRateLimiter(request int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		MaxRequestPerWindow: request,
		RateLimitWindow:     window,
		ipRequest:           make(map[string][]time.Time),
	}
}

func (r *RateLimiter) LimitMiddleWare() gin.HandlerFunc {
	return r.LimitMiddleWareWithMessage("Too many requests. Please try again later.")
}

func (r *RateLimiter) LimitMiddleWareWithMessage(message string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		now := time.Now()

		r.mu.Lock()

		// grab existing timestamps (nil-slice if none)
		info := r.ipRequest[ip]

		// keep only those inside the window
		pruned := info[:0]
		for _, ts := range info {
			if now.Sub(ts) < r.RateLimitWindow {
				pruned = append(pruned, ts)
			}
		}

		// Reject once the current request would exceed the configured limit.
		if len(pruned) >= r.MaxRequestPerWindow {
			r.ipRequest[ip] = pruned
			r.mu.Unlock()
			retryAfter := int(time.Until(pruned[0].Add(r.RateLimitWindow)).Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			ctx.Header("Retry-After", strconv.Itoa(retryAfter))
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":               message,
				"retry_after_seconds": retryAfter,
			})
			return
		}

		// Append current timestamp & save
		pruned = append(pruned, now)
		r.ipRequest[ip] = pruned
		r.mu.Unlock()

		ctx.Next()
	}
}
