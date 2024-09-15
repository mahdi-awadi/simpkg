package timerange

import (
	"regexp"
	"time"

	"github.com/go-per/simpkg/format"
	"github.com/go-per/simpkg/helpers"
	"github.com/go-per/simpkg/ticker"
	"github.com/go-per/simpkg/types"
)

// Format is time format
const Format = "04:05.000"

// TimeRange struct
type TimeRange struct {
	ticker  ticker.ITicker
	ranges  []ranges
	onEnter []func()
	onExit  []func()
	entered bool
	timer   *time.Timer
}

// ranges struct
type ranges struct {
	Start time.Time
	End   time.Time
}

// New time range
func New() *TimeRange {
	tr := &TimeRange{
		ranges:  make([]ranges, 0),
		onEnter: make([]func(), 0),
		onExit:  make([]func(), 0),
		ticker:  ticker.NewTicker(time.Second),
	}

	tr.init()
	return tr
}

// init initialize ticker
func (tr *TimeRange) init() {
	tr.ticker.OnTick(func(t ticker.ITicker) {
		if !tr.entered {
			current, _ := time.Parse(Format, time.Now().Format(Format))
			for _, tt := range tr.ranges {
				if helpers.TimeInRange(tt.Start, tt.End, true, current) {
					tr.entered = true
					for _, fn := range tr.onEnter {
						go fn()
					}

					// stop timer if already exists
					if tr.timer != nil {
						tr.timer.Stop()
						tr.timer = nil
					}

					// remove entered flag
					if tr.entered {
						tr.timer = time.AfterFunc(time.Duration(tt.End.Second()+2)*time.Second, func() {
							tr.entered = false
							for _, fn := range tr.onExit {
								go fn()
							}
						})
					}

					return
				}
			}
		}
	})
}

// Load and validate time range (0x:xx)
func (tr *TimeRange) Load(start, end string, threshold ...types.Duration) error {
	re := regexp.MustCompile(`^0[0-9]:[0-5][0-9]$`)
	if !re.MatchString(start) || !re.MatchString(end) {
		return format.Error("Invalid time range format")
	}

	firstStart, err := time.Parse(Format, start+".000")
	firstEnd, err := time.Parse(Format, end+".000")
	if err != nil {
		return err
	}

	// threshold
	th := types.Duration(0)
	if len(threshold) > 0 {
		th = threshold[0]
	}

	// create ranges
	for i := 0; i < 12; i++ {
		addDuration := time.Duration(i*5) * time.Minute
		tr.ranges = append(tr.ranges, ranges{
			Start: firstStart.Add(addDuration - th.Duraion()),
			End:   firstEnd.Add(addDuration + th.Duraion()),
		})
	}

	return nil
}

// Interval set ticker ticker
func (tr *TimeRange) Interval(i time.Duration) *TimeRange {
	tr.ticker.Duration(i)
	return tr
}

// OnEnter fire if time in range
func (tr *TimeRange) OnEnter(fn func()) *TimeRange {
	tr.onEnter = append(tr.onEnter, fn)
	return tr
}

// OnExit fire if time is out of range
func (tr *TimeRange) OnExit(fn func()) *TimeRange {
	tr.onExit = append(tr.onExit, fn)
	return tr
}

// Start ticker
func (tr *TimeRange) Start() error {
	if tr.ticker.IsActive() {
		return format.Error("Ticker already started")
	}

	tr.ticker.Start()
	return nil
}

// IsStarted returns is time ticker started
func (tr *TimeRange) IsStarted() bool {
	return tr.ticker.IsActive()
}

// Stop ticker
func (tr *TimeRange) Stop() error {
	if !tr.ticker.IsActive() {
		return format.Error("Ticker not started")
	}

	tr.ticker.Stop()
	return nil
}

// Ranges returns time ranges
func (tr *TimeRange) Ranges() []ranges {
	return tr.ranges
}
