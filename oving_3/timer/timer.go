package timer

import (
	"time"
)

var timerEndTime time.Time
var timerActive bool

func Timer_start() {
	if !timerActive {
		timerEndTime = time.Now().Add(3 * time.Second)
		timerActive = true
	}

}

func Timer_stop() {
	timerActive = false
}

func Timer_timedOut() bool {
	//fmt.Println(timerEndTime.Before(time.Now()))
	return (timerActive && timerEndTime.Before(time.Now()))
}
