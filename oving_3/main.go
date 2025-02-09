package main

import (
	. "elevator/elevio"
	"elevator/fsm"
	"elevator/timer"
	"fmt"
	"time"
)

func findDirection(floor int, desired_floor int) {
	if floor == desired_floor {
		SetMotorDirection(MD_Stop)
		SetDoorOpenLamp(true)
		SetButtonLamp(BT_Cab, desired_floor, false)
		SetButtonLamp(BT_HallDown, desired_floor, false)
		SetButtonLamp(BT_HallUp, desired_floor, false)
		time.Sleep(1 * time.Second)
		SetDoorOpenLamp(false)

	} else if floor < desired_floor {
		SetMotorDirection(MD_Up)
	} else {
		SetMotorDirection(MD_Down)
	}
}

func PollDoorTimer(e *Elevator, obstruction *bool) {

	for {
		if *obstruction {
			timer.Timer_start()
		} else {
			if timer.Timer_timedOut() {
				fmt.Println("Stopping timer")
				timer.Timer_stop()
				fsm.Fsm_onDoorTimeout(e)

			}
		}

	}

}

func updateObstruction(obstruction *bool, drv_obstr <-chan bool) {
	for {
		*obstruction = <-drv_obstr
	}
}

func main() {

	numFloors := 4
	var elevator = Elevator{Floor: -1, Dirn: MD_Stop, Behaviour: EB_Idle}
	obstruction := false
	Init("localhost:15657", numFloors)
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go PollButtons(drv_buttons)
	go PollFloorSensor(drv_floors)
	go PollObstructionSwitch(drv_obstr)
	go PollStopButton(drv_stop)
	go PollDoorTimer(&elevator, &obstruction)
	go updateObstruction(&obstruction, drv_obstr)

	fsm.Fsm_onInitBetweenFloors(&elevator, drv_floors)
	prev := -1

	for {
		select {
		case b := <-drv_buttons:
			fsm.Fsm_onRequestButtonPress(b.Floor, b.Button, &elevator)
		case floor_sensed := <-drv_floors:
			if floor_sensed != -1 && floor_sensed != prev {
				fsm.Fsm_onFloorArrival(floor_sensed, &elevator)
			}

			prev = floor_sensed

			//case v := <-drv_obstr:
			//obstruction = v
		}

	}
}
