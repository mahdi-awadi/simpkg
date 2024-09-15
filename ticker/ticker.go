package ticker

import (
	"time"
)

// ITicker interface
type ITicker interface {
	Start()
	OnStart(fn startFunc) ITicker
	OnTick(fn callbackFunc) ITicker
	OnComplete(fn callbackFunc) ITicker
	Duration(time.Duration)
	Stop()
	IsActive() bool
	Ticks() int
}

type startFunc func()
type callbackFunc func(t ITicker)
type Ticker struct {
	onStart          startFunc
	onTick           callbackFunc
	onComplete       callbackFunc
	duration         time.Duration
	timeTicker       *time.Ticker
	startImmediately bool
	doneCh           chan bool
	started          bool
	ticks            int
}

// NewTicker returns Ticker instance
func NewTicker(d time.Duration, im ...bool) ITicker {
	imd := true
	if len(im) > 0 && !im[0] {
		imd = false
	}
	return &Ticker{duration: d, startImmediately: imd, onTick: func(t ITicker) {}}
}

// Duration sets the duration
func (ticker *Ticker) Duration(d time.Duration) {
	ticker.duration = d
	if ticker.timeTicker != nil {
		ticker.timeTicker.Reset(d)
	}
}

// Start starts time ticker
func (ticker *Ticker) Start() {
	if ticker.started {
		return
	}

	ticker.started = true
	if ticker.onStart != nil {
		go ticker.onStart()
	}
	if ticker.startImmediately {
		ticker.ticks++
		go ticker.onTick(ticker)
	}

	ticker.doneCh = make(chan bool, 1)
	ticker.timeTicker = time.NewTicker(ticker.duration)
	go func() {
		for {
			select {
			case <-ticker.doneCh:
				if ticker.onComplete != nil {
					go ticker.onComplete(ticker)
				}
				return
			case <-ticker.timeTicker.C:
				if ticker.started {
					ticker.ticks++
					go ticker.onTick(ticker)
				}
			}
		}
	}()
}

// Ticks returns ticks count
func (ticker *Ticker) Ticks() int {
	return ticker.ticks
}

// OnTick fires on every tick
func (ticker *Ticker) OnTick(fn callbackFunc) ITicker {
	ticker.onTick = fn
	return ticker
}

// OnStart fires on start
func (ticker *Ticker) OnStart(fn startFunc) ITicker {
	ticker.onStart = fn
	return ticker
}

// OnComplete fires on complete
func (ticker *Ticker) OnComplete(fn callbackFunc) ITicker {
	ticker.onComplete = fn
	return ticker
}

// Stop stops the time ticker
func (ticker *Ticker) Stop() {
	ticker.started = false
	if ticker.doneCh != nil {
		ticker.doneCh <- true
	}
	if ticker.timeTicker != nil {
		ticker.timeTicker.Stop()
	}
	ticker.ticks = 0
}

// IsActive returns true if ticker started
func (ticker *Ticker) IsActive() bool {
	return ticker.started
}
