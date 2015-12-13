package sip

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

func NewRequest(method, requestURI string, body io.Reader) (Request, error) {
	//	u, err := uri.Parse(requestURI)
	//	if err != nil {
	//		return nil, err
	//	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	this := &request{
		message: message{
			sipVersion: "SIP/2.0",
			header:     make(Header),
			body:       rc,
		},
		method:     method,
		requestURI: requestURI,
	}
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			this.contentLength = int(v.Len())
		case *bytes.Reader:
			this.contentLength = int(v.Len())
		case *strings.Reader:
			this.contentLength = int(v.Len())
		}
	}

	return this, nil
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

// parseRequestLine parses "INVITE sip:bob@biloxi.com SIP/2.0" into its three parts.
func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}

// ReadRequest reads and parses an incoming request from b.
func ReadRequest(b *bufio.Reader) (req *request, err error) {
	tp := newTextprotoReader(b)
	req = new(request)

	// First line: INVITE sip:bob@biloxi.com SIP/2.0
	var s string
	if s, err = tp.ReadLine(); err != nil {
		return nil, err
	}
	defer func() {
		putTextprotoReader(tp)
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
	}()

	var ok bool
	req.method, req.requestURI, req.sipVersion, ok = parseRequestLine(s)
	if !ok {
		return nil, fmt.Errorf("malformed SIP request %s", s)
	}
	//rawurl := req.requestURI
	if _, _, ok = ParseSIPVersion(req.sipVersion); !ok {
		return nil, fmt.Errorf("malformed SIP version %s", req.sipVersion)
	}

	////////////////////////////////////////////////////////////////////////////
	err = ReadMessage(req, tp, b)
	if err != nil {
		return nil, err
	}

	return req, nil
}
