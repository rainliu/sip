package sip

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request interface {
	Message

	GetMethod() string
	SetMethod(method string) error
	GetRequestURI() string
	SetRequestURI(uri string) error
}

const (
	ACK       = "ACK"
	BYE       = "BYE"
	CANCEL    = "CANCEL"
	INVITE    = "INVITE"
	OPTIONS   = "OPTIONS"
	REGISTER  = "REGISTER"
	NOTIFY    = "NOTIFY"
	SUBSCRIBE = "SUBSCRIBE"
	MESSAGE   = "MESSAGE"
	REFER     = "REFER"
	INFO      = "INFO"
	PRACK     = "PRACK"
	UPDATE    = "UPDATE"
)

////////////////////////////////////////////////////////////////////////////////
type request struct {
	message

	method     string
	requestURI string
}

func NewRequest(method, requestURI string, body io.Reader) *request {
	this := &request{
		message: message{
			sipVersion: "SIP/2.0",
			header:     make(Header),
			body:       body,
		},
		method:     method,
		requestURI: requestURI,
	}
	this.StartLineWriter = this
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			this.SetContentLength(int64(v.Len()))
		case *bytes.Reader:
			this.SetContentLength(int64(v.Len()))
		case *strings.Reader:
			this.SetContentLength(int64(v.Len()))
		}
	}

	return this
}

func (this *request) GetMethod() string {
	return this.method
}

func (this *request) SetMethod(method string) error {
	this.method = method
	return nil
}

func (this *request) GetRequestURI() string {
	return this.requestURI
}

func (this *request) SetRequestURI(requestURI string) error {
	this.requestURI = requestURI
	return nil
}

//Method RequestURI SIP/2.0
func (this *request) StartLineWrite(w io.Writer) (err error) {
	if _, err = fmt.Fprintf(w, "%s %s SIP/2.0\r\n", this.GetMethod(), this.GetRequestURI()); err != nil {
		return err
	}
	return nil
}
