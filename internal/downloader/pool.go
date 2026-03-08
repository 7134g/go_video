package downloader

import (
	"net/url"
	"sync"
)

type Task func() error

type Pool struct {
	mu           sync.Mutex
	domainLimits map[string]chan struct{}
	maxPerDomain int
}

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

func (g *Group) Submit(rawURL string, task Task) {
	sem := g.pool.getDomainSem(rawURL)
	if sem == nil {
		task() // URL解析失败也要执行任务
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
