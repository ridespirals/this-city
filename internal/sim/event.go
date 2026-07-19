package sim

// EventKind classifies placeable events / POIs.
type EventKind int

const (
	EventCrime EventKind = iota
	EventDistress
	EventAttraction
	EventBench
)

// EventKindCount is the number of event kinds (for cycling).
const EventKindCount = 4

// EventSource marks an entity as a perceptible event or POI.
type EventSource struct {
	Kind     EventKind
	Priority int
	Lifetime float32 // seconds; <=0 means persistent
	Active   bool
	Age      float32 // accumulated live time
}

// EventKindName returns a short label for UI/debug.
func EventKindName(k EventKind) string {
	switch k {
	case EventCrime:
		return "crime"
	case EventDistress:
		return "distress"
	case EventAttraction:
		return "attraction"
	case EventBench:
		return "bench"
	default:
		return "event"
	}
}
