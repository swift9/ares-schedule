package schedule

import (
	"fmt"
	"testing"
	"time"
)

func TestDelaySchedule(t *testing.T) {
	sched := NewDelaySchedule(1*time.Second, func() {
		fmt.Println(11)
	})
	sched.Start()
	time.Sleep(4 * time.Second)
	sched.Stop()
	time.Sleep(3 * time.Second)
}

func TestOnceSchedule(t *testing.T) {
	sched := NewOnceSchedule(2*time.Second, func() {
		fmt.Println(11)
	})
	sched.Start()
	time.Sleep(4 * time.Second)
	sched.Stop()
	time.Sleep(3 * time.Second)
}

func TestCronSchedule(t *testing.T) {
	sched := NewCronSchedule("*/1 * * * * *", func() {
		fmt.Println(11)
	})
	sched.Start()
	time.Sleep(4 * time.Second)
	sched.Stop()
	time.Sleep(3 * time.Second)
}
