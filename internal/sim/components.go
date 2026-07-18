package sim

// Transform2D is world-space pose on the ground plane.
type Transform2D struct {
	X, Y     float32
	Rotation float32 // radians
	Scale    float32 // 1 is default; 0 treated as 1 when drawing
}

// AgentBrain is an FSM instance: which definition to run, current state, and blackboard.
type AgentBrain struct {
	Machine string // registry key (e.g. "debug")
	State   StateID
	BB      Blackboard
}

// Role identifies high-level agent type for systems and spawning.
type Role int

const (
	RoleNone Role = iota
	RoleCivilian
	RolePolice
	RoleDebug
)

// RoleTag tags an entity with a gameplay role.
type RoleTag struct {
	Role Role
}
