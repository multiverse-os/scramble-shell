package scrambleshell

// NOTE: The design policy for scrambleshell will be to NEVER build out
// components we are not actively utilizing. This could become massive as it
// mirror every aspect of a machine, and it may eventually or an alternative
// version may, but this is exclusively to model out components we are
// controlling for the purpose of managing the unified cluster shell desktop
// environment.

type Coordinate struct {
	X int
	Y int
}

type MachineType int

const (
	Physical MachineType = iota
	Virtual
)

type Machine struct {
	Type    MachineType
	Windows []string
}

type Process struct {
	Machine *Machine
	ID      int
}

type Desktop struct {
}

// TODO: Interested in experimenting with each window having its own rest
// server, its own IP address, possibly actor model or entity/component model
type Window struct {
	Process    *Process
	Desktop    *Desktop
	Title      string
	Width      int
	Height     int
	Focus      bool
	Position   Coordinate
	Parent     *Window
	Children   *Window
	Collisions []*Window
}
