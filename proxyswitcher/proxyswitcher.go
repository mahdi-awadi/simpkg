package proxyswitcher

import (
	"sync"

	"github.com/go-per/simpkg/format"
	"github.com/go-per/simpkg/random"
	"github.com/go-per/simpkg/str"
)

// ISwitcher interface
type ISwitcher interface {
	Load(l []string) error                           // load proxies
	Add(p string) error                              // add a proxy
	Remove(p string) error                           // remove a proxy
	All() []IProxy                                   // get all proxies
	Next(usageId ...string) IProxy                   // get next available proxy
	Random(usageId ...string) IProxy                 // get random available proxy
	Count() int                                      // proxies count
	TestAll(callback func(result string, err error)) // test all proxies
}

// Switcher struct
type Switcher struct {
	proxies []*Proxy
	locker  sync.RWMutex
}

// New returns new proxy switcher
func New() ISwitcher {
	return &Switcher{
		proxies: make([]*Proxy, 0),
		locker:  sync.RWMutex{},
	}
}

// Load load proxies list
func (ps *Switcher) Load(l []string) error {
	if l == nil || len(l) == 0 {
		return format.Error("Proxy list is empty")
	}

	for _, s := range l {
		if err := ps.Add(s); err != nil {
			return err
		}
	}

	return nil
}

// Add new proxy url to list
func (ps *Switcher) Add(p string) error {
	id := str.Checksum(p)
	exists := ps.get(id)
	if exists != nil {
		return format.Error("Proxy is duplicated [%v]", p)
	}

	proxy, err := NewProxy(id, p)
	if err != nil {
		return err
	}

	ps.locker.Lock()
	ps.proxies = append(ps.proxies, proxy)
	ps.locker.Unlock()

	return nil
}

// Remove proxy url from list
func (ps *Switcher) Remove(p string) error {
	id := str.Checksum(p)
	proxy := ps.get(id)
	if proxy == nil {
		return format.Error("Proxy not found [%v]", p)
	}

	ps.locker.Lock()
	for i, p := range ps.proxies {
		if p.id == id {
			ps.proxies = append(ps.proxies[:i], ps.proxies[i+1:]...)
			break
		}
	}
	ps.locker.Unlock()

	return nil
}

// All returns all proxies
func (ps *Switcher) All() []IProxy {
	var proxies []IProxy
	ps.locker.Lock()
	for _, p := range ps.proxies {
		proxies = append(proxies, p)
	}
	ps.locker.Unlock()

	return proxies
}

// Next get free proxy url
// usageId helps proxy not reuse in a period of time
func (ps *Switcher) Next(usageId ...string) (p IProxy) {
	ps.locker.Lock()
	for _, proxy := range ps.proxies {
		if proxy.isAvailable(usageId...) {
			if proxy.use(usageId...) {
				p = proxy
				break
			}
		}
	}
	ps.locker.Unlock()

	return
}

// Random get random free proxy url
func (ps *Switcher) Random(usageId ...string) (p IProxy) {
	ps.locker.Lock()
	max := len(ps.proxies) - 1
	for i := 0; i <= max; i++ {
		index := random.IntInRange(0, max)
		if index <= max {
			proxy := ps.proxies[index]
			if proxy.isAvailable(usageId...) {
				if proxy.use(usageId...) {
					p = proxy
					break
				}
			}
		}
	}
	ps.locker.Unlock()

	return
}

// get returns proxy by id
func (ps *Switcher) get(id string) (proxy IProxy) {
	ps.locker.Lock()
	for _, p := range ps.proxies {
		if p.id == id {
			proxy = p
			break
		}
	}
	ps.locker.Unlock()

	return
}

// Count return proxies count
func (ps *Switcher) Count() int {
	return len(ps.proxies)
}
