package requests

import (
	. "elevator/elevio"
)

type DirnBehaviourPair struct {
	Dirn      MotorDirection
	Behaviour ElevatorBehaviour
}

func Requests_above(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func Requests_below(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Requests[f][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func Requests_here(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Requests[e.Floor][btn] == 1 {
			return true
		}
	}
	return false
}

func Requests_chooseDirection(e Elevator) DirnBehaviourPair {
	switch e.Dirn {
	case MD_Up:
		if Requests_above(e) {
			return DirnBehaviourPair{Dirn: MD_Up, Behaviour: EB_Moving}
		}
		if Requests_here(e) {
			return DirnBehaviourPair{Dirn: MD_Down, Behaviour: EB_DoorOpen}
		}
		if Requests_below(e) {
			return DirnBehaviourPair{Dirn: MD_Down, Behaviour: EB_Moving}
		}
		return DirnBehaviourPair{Dirn: MD_Stop, Behaviour: EB_Idle}
	case MD_Down:
		if Requests_above(e) {
			return DirnBehaviourPair{Dirn: MD_Up, Behaviour: EB_Moving}
		}
		if Requests_here(e) {
			return DirnBehaviourPair{Dirn: MD_Up, Behaviour: EB_DoorOpen} // Er denne riktig?
		}
		if Requests_below(e) {
			return DirnBehaviourPair{Dirn: MD_Down, Behaviour: EB_Moving}
		}
		return DirnBehaviourPair{Dirn: MD_Stop, Behaviour: EB_Idle}
	case MD_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
		if Requests_above(e) {
			return DirnBehaviourPair{Dirn: MD_Up, Behaviour: EB_Moving}
		}
		if Requests_here(e) {
			return DirnBehaviourPair{Dirn: MD_Stop, Behaviour: EB_DoorOpen}
		}
		if Requests_below(e) {
			return DirnBehaviourPair{Dirn: MD_Down, Behaviour: EB_Moving}
		}
		return DirnBehaviourPair{Dirn: MD_Stop, Behaviour: EB_Idle}
	default:
		return DirnBehaviourPair{Dirn: MD_Stop, Behaviour: EB_Idle}
	}
}

func Requests_shouldStop(e Elevator) bool {
	switch e.Dirn {
	case MD_Down:
		return (e.Requests[e.Floor][BT_HallDown] == 1) || (e.Requests[e.Floor][BT_Cab] == 1) || !Requests_below(e)
	case MD_Up:
		return e.Requests[e.Floor][BT_HallUp] == 1 || e.Requests[e.Floor][BT_Cab] == 1 || !Requests_above(e)
	//case MD_Stop:
	default:
		return true
	}
}

func Requests_shouldClearImmediately(e Elevator, btn_floor int, btn_type ButtonType) bool {
	switch e.Config {
	case CV_All:
		return e.Floor == btn_floor
	case CV_InDirn:
		return e.Floor == btn_floor && ((e.Dirn == MD_Up && btn_type == BT_HallUp) || (e.Dirn == MD_Down && btn_type == BT_HallDown) || e.Dirn == MD_Stop || btn_type == BT_Cab)
	default:
		return false
	}
}

func Requests_clearAtCurrentFloor(e Elevator) Elevator {

	switch e.Config {
	case CV_All:
		for btn := 0; btn < N_BUTTONS; btn++ { //kan hende btn må være Button
			e.Requests[e.Floor][btn] = 0
		}

	case CV_InDirn:
		e.Requests[e.Floor][BT_Cab] = 0
		switch e.Dirn {
		case MD_Up:
			if !Requests_above(e) && e.Requests[e.Floor][BT_HallUp] == 0 {
				e.Requests[e.Floor][BT_HallDown] = 0
			}
			e.Requests[e.Floor][BT_HallUp] = 0

		case MD_Down:
			if !Requests_below(e) && e.Requests[e.Floor][BT_HallDown] == 0 {
				e.Requests[e.Floor][BT_HallUp] = 0
			}
			e.Requests[e.Floor][BT_HallDown] = 0

		case MD_Stop:
		default:
			e.Requests[e.Floor][BT_HallUp] = 0
			e.Requests[e.Floor][BT_HallDown] = 0

		}

	default:

	}

	return e
}
