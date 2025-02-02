package fsm

import (
	. "elevator/elevio"
	"elevator/requests"
	"time"
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
	//fmt.Println("\n\n%s(%d, %s)\n", __FUNCTION__, btn_floor, elevio_button_toString(btn_type))
	//elevator_print(elevator)

	switch elevator.Behaviour {
	case EB_DoorOpen:
		if requests.Requests_shouldClearImmediately(*elevator, btn_floor, btn_type) {
			//timer_start(elevator.config.doorOpenDuration_s)
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
			//outputDevice.doorLight(1)
			SetDoorOpenLamp(true)
			time.Sleep(3 * time.Second)
			SetDoorOpenLamp(false)
			//timer_start(elevator.config.doorOpenDuration_s)
			*elevator = requests.Requests_clearAtCurrentFloor(*elevator)

		case EB_Moving:
			//outputDevice.motorDirection(elevator.Dirn)
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
			SetAllLights(elevator)
			elevator.Behaviour = EB_DoorOpen

			time.Sleep(3 * time.Second)
			SetDoorOpenLamp(false)
			elevator.Behaviour = EB_Idle
		}

	default:

	}
}

func Fsm_onDoorTimeout(elevator *Elevator) {

	switch elevator.Behaviour {
	case EB_DoorOpen:
		pair := requests.Requests_chooseDirection(*elevator)
		elevator.Dirn = pair.Dirn
		elevator.Behaviour = pair.Behaviour

		switch elevator.Behaviour {
		case EB_DoorOpen:
			*elevator = requests.Requests_clearAtCurrentFloor(*elevator)
			SetAllLights(elevator)
			time.Sleep(3 * time.Second)
			SetDoorOpenLamp(false)

		case EB_Moving:
		case EB_Idle:
			SetMotorDirection(elevator.Dirn)

		}

	default:

	}
}
