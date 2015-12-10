package sip

import (
	"fmt"
	"io"
)

type Tracer interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

type tracer struct {
	out io.Writer
}

func TraceOn(w io.Writer) Tracer {
	return &tracer{out: w}
}
func (this *tracer) Println(a ...interface{}) {
	this.out.Write([]byte(fmt.Sprint(a...)))
	this.out.Write([]byte("\n"))
}
func (this *tracer) Printf(format string, a ...interface{}) {
	this.out.Write([]byte(fmt.Sprintf(format, a...)))
}

type nilTracer struct {
}

func TraceOff() Tracer {
	return &nilTracer{}
}
func (this *nilTracer) Println(a ...interface{}) {
}
func (this *nilTracer) Printf(format string, a ...interface{}) {
}
