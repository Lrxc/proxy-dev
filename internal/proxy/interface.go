package proxy

type Logger interface {
	Printf(format string, a ...any)
	Println(a ...any)
}

type DefaultLog struct{}

func (s DefaultLog) Printf(format string, a ...any) {
	//fmt.Printf("[proxy] "+format, a...)
}

func (s DefaultLog) Println(a ...any) {
	//fmt.Println("[proxy]", a)
}
