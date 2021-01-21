package log

type Mode string

const (
	Info    Mode = "info"
	Error   Mode = "error"
	Warning Mode = "warning"
	Debug   Mode = "debug"
	Trace   Mode = "trace"
)
