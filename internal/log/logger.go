package log

type Logger interface {
	Info(string)
	Infof(...any)
}
