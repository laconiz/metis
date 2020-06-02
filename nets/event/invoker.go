package event

type Invoker interface {
	Invoke(*Event)
}
