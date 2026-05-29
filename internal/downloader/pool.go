package downloader

import (
	"net/url"
	"sync"
)

type Task func() error

// Pool 按"目标域名"做并发下载限流：每个域名一个独立信号量，容量 = maxPerDomain。
// 目的：避免同一站点被并发分段打挂或触发风控，同时不阻塞其它域名的下载。
type Pool struct {
	mu           sync.Mutex
	domainLimits map[string]chan struct{}
	maxPerDomain int
}

// Group 把一组分段视作同一批次：所有 Submit 完成后通过 Wait 同步等待。
type Group struct {
	pool *Pool
	wg   sync.WaitGroup
}

func NewPool(maxPerDomain int) *Pool {
	return &Pool{
		domainLimits: make(map[string]chan struct{}),
		maxPerDomain: maxPerDomain,
	}
}

func (p *Pool) NewGroup() *Group {
	return &Group{pool: p}
}

func (p *Pool) getDomainSem(rawURL string) chan struct{} {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	domain := u.Host

	p.mu.Lock()
	defer p.mu.Unlock()

	if sem, ok := p.domainLimits[domain]; ok {
		return sem
	}

	sem := make(chan struct{}, p.maxPerDomain)
	p.domainLimits[domain] = sem
	return sem
}

// Submit 提交一个分段任务。URL 解析失败时直接同步跑（fallback，避免静默丢任务）。
func (g *Group) Submit(rawURL string, task Task) {
	sem := g.pool.getDomainSem(rawURL)
	if sem == nil {
		task()
		return
	}

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()
		task()
	}()
}

func (g *Group) Wait() {
	g.wg.Wait()
}
