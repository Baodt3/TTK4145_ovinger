package fsm

import (
	. "elevator/elevio"
	"fmt"
)

type ElevatorBehaviour int

const (
	EB_Idle     = 0
	EB_DoorOpen = 1
	EB_Moving   = 2
)

type Elevator struct {
	floor     int
	dirn      MotorDirection
	requests  [N_FLOORS][N_BUTTONS]int
	behaviour ElevatorBehaviour
}

var elevator Elevator

func fsm_onInitBetweenFloors(elevator *Elevator) {
	SetMotorDirection(MD_Down)
	for elevator.floor != 0 {
	}
	elevator.dirn = MD_Down
	elevator.behaviour = EB_Idle
	elevator.floor = 0
}

func setAllLights(es Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			SetButtonLamp(btn, floor, es.requests[floor][btn])
		}
	}
}

func fsm_onRequestButtonPress(btn_floor int, btn_type ButtonType) {
	//fmt.Println("\n\n%s(%d, %s)\n", __FUNCTION__, btn_floor, elevio_button_toString(btn_type))
	//elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		if requests.requests_shouldClearImmediately(elevator, btn_floor, btn_type) {
			//timer_start(elevator.config.doorOpenDuration_s)
		} else {
			elevator.requests[btn_floor][btn_type] = 1
		}

	case EB_Moving:
		elevator.requests[btn_floor][btn_type] = 1

	case EB_Idle:
		elevator.requests[btn_floor][btn_type] = 1
		pair := requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour
		switch pair.behaviour {
		case EB_DoorOpen:
			outputDevice.doorLight(1)
			//timer_start(elevator.config.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)

		case EB_Moving:
			outputDevice.motorDirection(elevator.dirn)

		case EB_Idle:

		}

	}

	setAllLights(elevator)

	fmt.Println("\nNew state:\n")
	elevator_print(elevator)
}

func fsm_onFloorArrival(newFloor int) {
	fmt.Println("\n\n%s(%d)\n", __FUNCTION__, newFloor)
	elevator_print(elevator)

	elevator.floor = newFloor

	outputDevice.floorIndicator(elevator.floor)

	switch elevator.behaviour {
	case EB_Moving:
		if requests_shouldStop(elevator) {
			outputDevice.motorDirection(D_Stop)
			outputDevice.doorLight(1)
			elevator = requests_clearAtCurrentFloor(elevator)
			//timer_start(elevator.config.doorOpenDuration_s)
			setAllLights(elevator)
			elevator.behaviour = EB_DoorOpen
		}

	default:

	}

	fmt.Println("\nNew state:\n")
	elevator_print(elevator)
}

func fsm_onDoorTimeout(void) {
	fmt.Println("\n\n%s()\n", __FUNCTION__)
	elevator_print(elevator)

	switch elevator.behaviour {
	case EB_DoorOpen:
		pair := requests_chooseDirection(elevator)
		elevator.dirn = pair.dirn
		elevator.behaviour = pair.behaviour

		switch elevator.behaviour {
		case EB_DoorOpen:
			//timer_start(elevator.config.doorOpenDuration_s)
			elevator = requests_clearAtCurrentFloor(elevator)
			setAllLights(elevator)

		case EB_Moving:
		case EB_Idle:
			outputDevice.doorLight(0)
			outputDevice.motorDirection(elevator.dirn)

		}

	default:

	}

	fmt.Println("\nNew state:\n")
	elevator_print(elevator)
}
