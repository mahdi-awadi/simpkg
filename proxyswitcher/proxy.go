package proxyswitcher

import (
	urlPkg "net/url"
	"sync"
	"time"

	"github.com/go-per/simpkg/helpers"
)

// timers locker
var timersLocker sync.RWMutex

// IProxy interface
type IProxy interface {
	Id() string                                       // returns proxy id
	Url() urlPkg.URL                                  // returns proxy url
	RawUrl() string                                   // returns proxy raw url
	IsInUse() bool                                    // determine if proxy is in use
	IsUsedFor(usageId ...string) bool                 // determine if proxy is used for usage id
	Release()                                         // release proxy and make it available
	RemoveUsage(usageId string)                       // remove proxy usage id
	RemoveUsageAfter(usageId string, d time.Duration) // remove proxy usage id after a while
}

// Proxy struct
type Proxy struct {
	id        string
	url       string
	isInUse   bool
	usageIds  []string
	timers    map[string]*time.Timer
	parsedUrl *urlPkg.URL
}

// NewProxy returns new proxy instance
func NewProxy(id, u string) (*Proxy, error) {
	up, err := urlPkg.Parse(u)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		id:        id,
		url:       u,
		isInUse:   false,
		parsedUrl: up,
		usageIds:  make([]string, 0),
		timers:    make(map[string]*time.Timer),
	}, nil
}

// Id returns proxy id
func (proxy *Proxy) Id() string {
	return proxy.id
}

// RawUrl returns proxy url
func (proxy *Proxy) RawUrl() string {
	return proxy.url
}

// Url returns proxy url
func (proxy *Proxy) Url() urlPkg.URL {
	return *proxy.parsedUrl
}

// IsInUse determine if proxy is in use
func (proxy *Proxy) IsInUse() bool {
	return proxy.isInUse
}

// IsUsedFor returns
func (proxy *Proxy) IsUsedFor(usageId ...string) bool {
	if usageId == nil || len(usageId) == 0 {
		return false
	}

	return helpers.Includes(proxy.usageIds, usageId[0])
}

// isAvailable returns true if proxy is available
func (proxy *Proxy) isAvailable(usageId ...string) bool {
	return !proxy.isInUse && !proxy.IsUsedFor(usageId...)
}

// use uses a proxy and return true, if false means proxy is not available
func (proxy *Proxy) use(usageId ...string) bool {
	if proxy.isInUse {
		return false
	}
	if usageId != nil && len(usageId) > 0 && proxy.IsUsedFor(usageId[0]) {
		return false
	}

	proxy.isInUse = true
	if usageId != nil && len(usageId) > 0 {
		proxy.usageIds = append(proxy.usageIds, usageId[0])
	}

	return true
}

// Release releases proxy
func (proxy *Proxy) Release() {
	proxy.isInUse = false
}

// RemoveUsage remove proxy usage id
func (proxy *Proxy) RemoveUsage(usageId string) {
	timersLocker.Lock()
	timer, ok := proxy.timers[usageId]
	if ok && timer != nil {
		timer.Stop()
	}
	proxy.usageIds = helpers.RemoveItem(proxy.usageIds, usageId)
	timersLocker.Unlock()
}

// RemoveUsageAfter remove proxy usage id after a while
func (proxy *Proxy) RemoveUsageAfter(usageId string, d time.Duration) {
	timersLocker.Lock()
	if proxy.timers == nil {
		proxy.timers = make(map[string]*time.Timer)
	}
	proxy.timers[usageId] = time.AfterFunc(d, func() {
		proxy.RemoveUsage(usageId)
	})
	timersLocker.Unlock()
}
