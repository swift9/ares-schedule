package schedule

import (
	"github.com/robfig/cron"
	"time"
)

type Schedule struct {
	StopChan chan bool
	F        func()
	Stopped  bool
}

func (schedule *Schedule) Stop() {
	if schedule.Stopped {
		return
	}
	schedule.Stopped = true
	schedule.StopChan <- true
}

type TimeSchedule struct {
	Schedule
	TimeChan <-chan time.Time
}

func (timeSchedule *TimeSchedule) Start() {
	timeSchedule.Stopped = false
	go func() {
		for {
			select {
			case <-timeSchedule.TimeChan:
				if !timeSchedule.Stopped {
					timeSchedule.F()
				}
			case stop := <-timeSchedule.StopChan:
				if stop {
					timeSchedule.StopTime()
					return
				}
			}
		}
	}()
}

func (timeSchedule *TimeSchedule) StopTime() {
}

type DelaySchedule struct {
	TimeSchedule
	Ticker time.Ticker
}

func (schedule *DelaySchedule) StopTime() {
	schedule.Stopped = true
	schedule.Ticker.Stop()
}

func NewDelaySchedule(duration time.Duration, f func()) *DelaySchedule {
	schedule := DelaySchedule{}
	schedule.StopChan = make(chan bool, 1)
	schedule.Stopped = true
	ticker := time.NewTicker(duration)
	schedule.Ticker = *ticker
	schedule.TimeChan = ticker.C
	schedule.F = f
	return &schedule
}

type OnceSchedule struct {
	TimeSchedule
	Timer time.Timer
}

func (schedule *OnceSchedule) StopTime() {
	schedule.Stopped = true
	schedule.Timer.Stop()
}

func NewOnceSchedule(duration time.Duration, f func()) *OnceSchedule {
	schedule := OnceSchedule{}
	schedule.Stopped = true
	schedule.StopChan = make(chan bool, 1)
	timer := time.NewTimer(duration)
	schedule.Timer = *timer
	schedule.TimeChan = timer.C
	schedule.F = func() {
		if !schedule.Stopped {
			f()
			schedule.Stop()
		}
	}
	return &schedule
}

type CronSchedule struct {
	Schedule
	Cron cron.Cron
}

func (schedule *CronSchedule) Start() {
	schedule.Stopped = false
	go func() {
		schedule.Cron.Start()
		select {
		case stop := <-schedule.StopChan:
			if stop {
				schedule.Cron.Stop()
				return
			}
		}
	}()
}

func NewCronSchedule(spec string, f func()) *CronSchedule {
	schedule := CronSchedule{Cron: *cron.New()}
	schedule.Stopped = true
	schedule.StopChan = make(chan bool, 1)
	schedule.F = func() {
		if !schedule.Stopped {
			f()
		}
	}
	schedule.Cron.AddFunc(spec, f)
	return &schedule
}
