package elasticsearch

type Document interface {
	ElkName() string
	ElkBody() string
}
