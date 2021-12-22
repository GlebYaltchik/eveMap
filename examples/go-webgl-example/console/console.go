package console

import (
	"honnef.co/go/js/dom"
)

type Console struct {
	console *dom.Console
}

func New(w dom.Window) Console {
	return Console{
		console: w.Console(),
	}
}

func (c Console) Info(args ...interface{}) {
	c.console.Call("info", args...)
}

func (c Console) Log(args ...interface{}) {
	c.console.Call("log", args...)
}

func (c Console) Debug(args ...interface{}) {
	c.console.Call("debug", args...)
}

func (c Console) Error(args ...interface{}) {
	c.console.Call("error", args...)
}

func (c Console) Warn(args ...interface{}) {
	c.console.Call("warn", args...)
}

func (c Console) Assert(cond bool, args ...interface{}) {
	c.console.Call("assert", cond, append(append([]interface{}(nil), cond), args...))
}
