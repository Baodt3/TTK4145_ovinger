package fsm

type ElevatorBehaviour int

const (
	EB_Idle   ElevatorBehaviour
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	Floor  int
	Dirn dirn
	requests[N_FLOORS][N_BUTTONS] int
	behaviour ElevatorBehaviour
}

var elevator Elevator

func fsm_onInitBetweenFloors(elevator *Elevator) {
	elevio.SetMotorDirection(elevio.MD_Down)
	for floor != 0 {}
	elevator.dirn = MD_Down
	elevator.behaviour = EB_Idle
	floor = 0
}