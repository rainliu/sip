package sip

import (
	"errors"
	"io"
)

// A Request represents an SIP request received by a server
// or to be sent by a client.
//
// The field semantics differ slightly between client and server
// usage. In addition to the notes on the fields below, see the
// documentation for Request.Write and RoundTripper.

type Message interface {
	// The protocol version for incoming requests.
	// Client requests always use SIP/2.0.
	GetSIPVersion() string // "SIP/2.0"
	SetSIPVersion(string) error

	// A header maps request lines to their values.
	// If the header says
	//
	//	accept-encoding: gzip, deflate
	//	Accept-Language: en-us
	//	Connection: keep-alive
	//
	// then
	//
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Connection": {"keep-alive"},
	//	}
	//
	// HTTP defines that header names are case-insensitive.
	// The request parser implements this by canonicalizing the
	// name, making the first character and any characters
	// following a hyphen uppercase and the rest lowercase.
	//
	// For client requests certain headers are automatically
	// added and may override values in Header.
	//
	// See the documentation for the Request.Write method.
	Header() Header

	// ContentLength records the length of the associated content.
	// The value -1 indicates that the length is unknown.
	// Values >= 0 indicate that the given number of bytes may
	// be read from Body.
	// For client requests, a value of 0 means unknown if Body is not nil.
	GetContentLength() int
	SetContentLength(l int)

	// Body is the request's body.
	//
	// For client requests a nil body means the request has no
	// body, such as a GET request. The HTTP Client's Transport
	// is responsible for calling the Close method.
	//
	// For server requests the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.
	Body() io.ReadCloser
}

////////////////////////////////////////
type message struct {
	sipVersion    string
	header        Header
	contentLength int
	body          io.ReadCloser
}

func (this *message) GetSIPVersion() string {
	return this.sipVersion
}

func (this *message) SetSIPVersion(s string) error {
	if s != "SIP/2.0" {
		return errors.New("Wrong SIP Version")
	} else {
		this.sipVersion = s
		return nil
	}
}

func (this *message) Header() Header {
	return this.header
}

func (this *message) GetContentLength() int {
	return this.contentLength
}

func (this *message) SetContentLength(l int) {
	this.contentLength = l
}

func (this *message) Body() io.ReadCloser {
	return this.body
}
