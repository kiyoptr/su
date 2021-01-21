package adapter

type Adapter interface {
	Write(message string)
	Close() error
}
