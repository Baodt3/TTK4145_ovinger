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