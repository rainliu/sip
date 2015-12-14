package sip

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
)

type Message interface {
	GetSIPVersion() string
	SetSIPVersion(string) error
	GetHeader() Header
	SetHeader(Header)
	GetContentLength() int64
	SetContentLength(l int64)
	GetBody() io.Reader
	SetBody(io.Reader)
}

////////////////////////////////////////////////////////////////////////////////
type message struct {
	sipVersion    string
	header        Header
	contentLength int64
	body          io.Reader
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

func (this *message) GetHeader() Header {
	return this.header
}

func (this *message) SetHeader(header Header) {
	this.header = header
}

func (this *message) GetContentLength() int64 {
	return this.contentLength
}

func (this *message) SetContentLength(l int64) {
	this.contentLength = l
}

func (this *message) GetBody() io.Reader {
	return this.body
}

func (this *message) SetBody(body io.Reader) {
	this.body = body
}

// Headers that Request.Write handles itself and should be skipped.
var reqWriteExcludeHeader = map[string]bool{
	"Content-Length": true,
}

//	Header
//	ContentLength
//	Body
func (this *message) write(w io.Writer) (err error) {
	if err = this.header.WriteSubset(w, reqWriteExcludeHeader); err != nil {
		return err
	}

	if _, err = fmt.Fprintf(w, "%s: %d\r\n", "Content-Length", this.GetContentLength()); err != nil {
		return err
	}

	if _, err = io.WriteString(w, "\r\n"); err != nil {
		return err
	}

	// Write body
	if this.body != nil {
		if _, err = io.Copy(w, io.LimitReader(this.body, this.GetContentLength())); err != nil {
			return err
		}
	}

	return nil
}

func ReadMessage(m Message, tp *textproto.Reader, b *bufio.Reader) error {
	// Subsequent lines: Key: value.
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		return err
	}
	m.SetHeader(Header(mimeHeader))

	////////////////////////////////////////////////////////////////////////////

	contentLens := m.GetHeader()["Content-Length"]
	if len(contentLens) > 1 { // harden against SIP request smuggling. See RFC 7230.
		return errors.New("http: message cannot contain multiple Content-Length headers")
	}

	// Logic based on Content-Length
	var cl string
	if len(contentLens) == 1 {
		cl = strings.TrimSpace(contentLens[0])
	}
	if cl != "" {
		n, err := parseContentLength(cl)
		if err != nil {
			return err
		}
		m.SetContentLength(n)
	} else {
		m.GetHeader().Del("Content-Length")
		m.SetContentLength(0)
	}

	////////////////////////////////////////////////////////////////////////////

	if m.GetContentLength() > 0 {
		m.SetBody(io.LimitReader(b, int64(m.GetContentLength())))
	} else {
		m.SetBody(nil)
	}

	return nil
}

// parseContentLength trims whitespace from s and returns -1 if no value
// is set, or the value if it's >= 0.
func parseContentLength(cl string) (int64, error) {
	cl = strings.TrimSpace(cl)
	if cl == "" {
		return -1, nil
	}
	n, err := strconv.ParseInt(cl, 10, 64)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("bad Content-Length %d", cl)
	}
	return n, nil

}

var textprotoReaderPool sync.Pool

func newTextprotoReader(br *bufio.Reader) *textproto.Reader {
	if v := textprotoReaderPool.Get(); v != nil {
		tr := v.(*textproto.Reader)
		tr.R = br
		return tr
	}
	return textproto.NewReader(br)
}

func putTextprotoReader(r *textproto.Reader) {
	r.R = nil
	textprotoReaderPool.Put(r)
}

// ParseSIPVersion parses a SIP version string.
// "SIP/2.0" returns (2, 0, true).
func ParseSIPVersion(vers string) (major, minor int, ok bool) {
	const Big = 1000000 // arbitrary upper bound
	switch vers {
	case "SIP/2.0":
		return 2, 0, true
	}
	if !strings.HasPrefix(vers, "SIP/") {
		return 0, 0, false
	}
	dot := strings.Index(vers, ".")
	if dot < 0 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(vers[4:dot])
	if err != nil || major < 0 || major > Big {
		return 0, 0, false
	}
	minor, err = strconv.Atoi(vers[dot+1:])
	if err != nil || minor < 0 || minor > Big {
		return 0, 0, false
	}
	return major, minor, true
}
