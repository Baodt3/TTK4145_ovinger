package main

import (
	"elevator/elevio"
	"fmt"
	"time"
)

func findDirection(floor int, desired_floor int) {
	if floor == desired_floor {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		elevio.SetButtonLamp(elevio.BT_Cab, desired_floor, false)
		elevio.SetButtonLamp(elevio.BT_HallDown, desired_floor, false)
		elevio.SetButtonLamp(elevio.BT_HallUp, desired_floor, false)
		time.Sleep(1 * time.Second)
		elevio.SetDoorOpenLamp(false)

	} else if floor < desired_floor {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	var floor int
	var desired_floor int
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	elevio.SetMotorDirection(elevio.MD_Down)
	time.Sleep(10000)
	floor = 0

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			desired_floor = a.Floor
			findDirection(floor, desired_floor)

		case a := <-drv_floors:
			floor = a
			elevio.SetFloorIndicator(floor)
			findDirection(floor, desired_floor)
			/*
			   fmt.Printf("%+v\n", a)
			   if a == numFloors-1 {
			       d = elevio.MD_Down
			   } else if a == 0 {
			       d = elevio.MD_Up
			   }
			   elevio.SetMotorDirection(d)
			*/
		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}
