package engine

type Model interface {
	// ParseEventData(e Event) interface{}
	Add(a, b interface{}) interface{}
}
