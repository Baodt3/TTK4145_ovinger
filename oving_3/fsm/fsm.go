package fsm

import (
	. "elevator/elevio"
	"elevator/requests"
	"elevator/timer"
	"fmt"
)

func Fsm_onInitBetweenFloors(elevator *Elevator, floor chan int) {
	SetMotorDirection(MD_Down)

	f := <-floor

	for f != 0 {
		f = <-floor
	}

	SetMotorDirection(MD_Stop)
	elevator.Dirn = MD_Stop
	elevator.Behaviour = EB_Idle
	elevator.Floor = 0
}

func SetAllLights(es *Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			SetButtonLamp(ButtonType(btn), floor, es.Requests[floor][btn] == 1)
		}
	}
}

func Fsm_onRequestButtonPress(btn_floor int, btn_type ButtonType, elevator *Elevator) {

	switch elevator.Behaviour {
	case EB_DoorOpen:
		if requests.Requests_shouldClearImmediately(*elevator, btn_floor, btn_type) {
			timer.Timer_start()
		} else {
			elevator.Requests[btn_floor][btn_type] = 1
		}

	case EB_Moving:
		elevator.Requests[btn_floor][btn_type] = 1

	case EB_Idle:
		elevator.Requests[btn_floor][btn_type] = 1
		pair := requests.Requests_chooseDirection(*elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour
		switch pair.Behaviour {
		case EB_DoorOpen:
			SetDoorOpenLamp(true)
			timer.Timer_start()
			*elevator = requests.Requests_clearAtCurrentFloor(*elevator)

		case EB_Moving:
			SetMotorDirection(elevator.Dirn)

		case EB_Idle:
		}

	}

	SetAllLights(elevator)
}

func Fsm_onFloorArrival(newFloor int, elevator *Elevator) {

	elevator.Floor = newFloor

	SetFloorIndicator(elevator.Floor)

	switch elevator.Behaviour {
	case EB_Moving:
		if requests.Requests_shouldStop(*elevator) {
			SetMotorDirection(MD_Stop)
			SetDoorOpenLamp(true)
			*elevator = requests.Requests_clearAtCurrentFloor(*elevator)
			timer.Timer_start()
			SetAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen

		}

	default:

	}
}

func Fsm_onDoorTimeout(elevator *Elevator) {
	fmt.Printf("Elevator behaviour: %d", elevator.Behaviour)
	switch elevator.Behaviour {
	case EB_DoorOpen:
		pair := requests.Requests_chooseDirection(*elevator)
		fmt.Printf("Pair direction: %d \n", pair.Dirn)
		fmt.Printf("Pair behaviour: %d \n", pair.Behaviour)

		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			timer.Timer_start()
			*elevator = requests.Requests_clearAtCurrentFloor(*elevator)
			SetAllLights(elevator)

		case EB_Moving:
			SetDoorOpenLamp(false)
			SetMotorDirection(elevator.Dirn)
		case EB_Idle:
			SetDoorOpenLamp(false)
			SetMotorDirection(elevator.Dirn)

		}
	default:

	}
}
