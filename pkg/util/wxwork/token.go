package wxwork

import (
	"context"
	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// wxwork API 相关代码来自 https://github.com/xen0n/go-workwx

// tokenInfo 用来从API返回的JSON中直接获取token信息
type tokenInfo struct {
	token     string
	expiresIn time.Duration
}

// token 主要用来自动刷新token
type token struct {
	mutex *sync.RWMutex
	tokenInfo
	lastRefresh  time.Time
	getTokenFunc func() (tokenInfo, error) // 获取Token的方法
}

func (t *token) setGetTokenFunc(f func() (tokenInfo, error)) {
	t.getTokenFunc = f
}

// 获取token值
func (t *token) getToken() string {
	// intensive mutex juggling action
	t.mutex.RLock()
	if t.token == "" {
		t.mutex.RUnlock() // RWMutex doesn't like recursive locking
		// TODO: what to do with the possible error?
		_ = t.syncToken()
		t.mutex.RLock()
	}
	tokenToUse := t.token
	t.mutex.RUnlock()
	return tokenToUse
}

// 同步Token信息
func (t *token) syncToken() error {
	get, err := t.getTokenFunc()
	if err != nil {
		return err
	}
	// 加锁，禁止Client同步时获取Token
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.token = get.token
	t.expiresIn = get.expiresIn * time.Second
	t.lastRefresh = time.Now()
	return nil
}

// 自动刷新Token
// refreshTimeWindow：在token超时前半小时刷新
func (t *token) tokenRefresher(ctx context.Context) {
	const refreshTimeWindow = 30 * time.Minute
	const minRefreshDuration = 5 * time.Second

	var waitDuration time.Duration = 0
	for {
		select {
		case <-time.After(waitDuration): // 等待一段时间后更新token
			// backoff 是一个退避重试的工具
			retryer := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
			if err := backoff.Retry(t.syncToken, retryer); err != nil {
				logrus.Error(errors.Wrapf(err, "刷新token失败!"))
			}

			waitUntilTime := t.lastRefresh.Add(t.expiresIn).Add(-refreshTimeWindow)
			waitDuration = time.Until(waitUntilTime)
			if waitDuration < minRefreshDuration {
				waitDuration = minRefreshDuration
			}
		case <-ctx.Done():
			return
		}
	}
}
