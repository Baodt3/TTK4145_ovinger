package requests

import (
	. "elevator/elevio"
)

type DirnBehaviourPair struct {
	dirn      MotorDirection
	behaviour ElevatorBehaviour
}

func requests_above(e Elevator) int {
	for f := e.floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func requests_below(e Elevator) int {
	for f := 0; f < e.floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.requests[f][btn] {
				return 1
			}
		}
	}
	return 0
}

func requests_here(e Elevator) int {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			return 1
		}
	}
	return 0
}

func requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.dirn {
	case D_Up:
		if requests_above(e) {
			return DirnBehaviourPair{dirn: MD_Up, behaviour: EB_Moving}
		}
		if requests_here(e) {
			return DirnBehaviourPair{dirn: MD_Down, behaviour: EB_DoorOpen}
		}
		if requests_below(e) {
			return DirnBehaviourPair{dirn: MD_Down, behaviour: EB_Moving}
		}
		return DirnBehaviourPair{dirn: MD_Stop, behaviour: EB_Idle}
	case D_Down:
		if requests_above(e) {
			return DirnBehaviourPair{dirn: MD_Down, behaviour: EB_Moving}
		}
		if requests_here(e) {
			return DirnBehaviourPair{dirn: MD_Up, behaviour: EB_DoorOpen}
		}
		if requests_below(e) {
			return DirnBehaviourPair{dirn: MD_Up, behaviour: EB_Moving}
		}
		return DirnBehaviourPair{dirn: MD_Stop, behaviour: EB_Idle}
	case D_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
		if requests_above(e) {
			return DirnBehaviourPair{dirn: MD_Stop, behaviour: EB_Moving}
		}
		if requests_here(e) {
			return DirnBehaviourPair{dirn: MD_Up, behaviour: EB_DoorOpen}
		}
		if requests_below(e) {
			return DirnBehaviourPair{dirn: MD_Down, behaviour: EB_Moving}
		}
		return DirnBehaviourPair{dirn: MD_Stop, behaviour: EB_Idle}
	default:
		return DirnBehaviourPair{dirn: MD_Stop, behaviour: EB_Idle}
	}
}

func requests_shouldStop(e Elevator) int {
	switch e.dirn {
	case MD_Down:
		return e.requests[e.floor][BT_HallDown] || e.requests[e.floor][BT_Cab] || !requests_below(e)
	case MD_Up:
		return e.requests[e.floor][BT_HallUp] || e.requests[e.floor][BT_Cab] || !requests_above(e)
	case MD_Stop:
	default:
		return 1
	}
}

func requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type Button) int {
	switch e.config.clearRequestVariant {
	case CV_All:
		return e.floor == btn_floor
	case CV_InDirn:
		return e.floor == btn_floor && ((e.dirn == D_Up && btn_type == BT_HallUp) || (e.dirn == D_Down && btn_type == BT_HallDown) || e.dirn == D_Stop || btn_type == BT_Cab)
	default:
		return 0
	}
}

func requests_clearAtCurrentFloor(e Elevator) Elevator {

	switch e.config.clearRequestVariant {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ { //kan hende btn må være Button
			e.requests[e.floor][btn] = 0
		}

	case CV_InDirn:
		e.requests[e.floor][BT_Cab] = 0
		switch e.dirn {
		case MD_Up:
			if !requests_above(e) && !e.requests[e.floor][BT_HallUp] {
				e.requests[e.floor][BT_HallDown] = 0
			}
			e.requests[e.floor][BT_HallUp] = 0

		case MD_Down:
			if !requests_below(e) && !e.requests[e.floor][BT_HallDown] {
				e.requests[e.floor][BT_HallUp] = 0
			}
			e.requests[e.floor][BT_HallDown] = 0

		case MD_Stop:
		default:
			e.requests[e.floor][BT_HallUp] = 0
			e.requests[e.floor][BT_HallDown] = 0

		}

	default:

	}

	return e
}
