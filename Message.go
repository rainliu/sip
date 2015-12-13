package sip

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	GetContentLength() int
	SetContentLength(l int)
	GetBody() io.ReadCloser
	SetBody(io.ReadCloser)
}

////////////////////////////////////////////////////////////////////////////////
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

func (this *message) GetHeader() Header {
	return this.header
}

func (this *message) SetHeader(header Header) {
	this.header = header
}

func (this *message) GetContentLength() int {
	return this.contentLength
}

func (this *message) SetContentLength(l int) {
	this.contentLength = l
}

func (this *message) GetBody() io.ReadCloser {
	return this.body
}

func (this *message) SetBody(body io.ReadCloser) {
	this.body = body
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
		m.SetBody(&body{src: io.LimitReader(b, int64(m.GetContentLength()))})
	} else {
		m.SetBody(eofReader)
	}

	return nil
}

// parseContentLength trims whitespace from s and returns -1 if no value
// is set, or the value if it's >= 0.
func parseContentLength(cl string) (int, error) {
	cl = strings.TrimSpace(cl)
	if cl == "" {
		return -1, nil
	}
	n, err := strconv.ParseInt(cl, 10, 32)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("bad Content-Length %d", cl)
	}
	return int(n), nil

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

type eofReaderWithWriteTo struct{}

func (eofReaderWithWriteTo) WriteTo(io.Writer) (int64, error) { return 0, nil }
func (eofReaderWithWriteTo) Read([]byte) (int, error)         { return 0, io.EOF }

// eofReader is a non-nil io.ReadCloser that always returns EOF.
// It has a WriteTo method so io.Copy won't need a buffer.
var eofReader = &struct {
	eofReaderWithWriteTo
	io.Closer
}{
	eofReaderWithWriteTo{},
	ioutil.NopCloser(nil),
}

// body turns a Reader into a ReadCloser.
// Close ensures that the body has been fully read
// and then reads the trailer if necessary.
type body struct {
	src          io.Reader
	r            *bufio.Reader // underlying wire-format reader for the trailer
	closing      bool          // is the connection to be closed after reading body?
	doEarlyClose bool          // whether Close should stop early

	mu         sync.Mutex // guards closed, and calls to Read and Close
	sawEOF     bool
	closed     bool
	earlyClose bool // Close called and we didn't read to the end of src
}

// ErrBodyReadAfterClose is returned when reading a Request or Response
// Body after the body has been closed. This typically happens when the body is
// read after an HTTP Handler calls WriteHeader or Write on its
// ResponseWriter.
var ErrBodyReadAfterClose = errors.New("http: invalid Read on closed Body")

func (b *body) Read(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return 0, ErrBodyReadAfterClose
	}
	return b.readLocked(p)
}

// Must hold b.mu.
func (b *body) readLocked(p []byte) (n int, err error) {
	if b.sawEOF {
		return 0, io.EOF
	}
	n, err = b.src.Read(p)

	if err == io.EOF {
		b.sawEOF = true
		{
			// If the server declared the Content-Length, our body is a LimitedReader
			// and we need to check whether this EOF arrived early.
			if lr, ok := b.src.(*io.LimitedReader); ok && lr.N > 0 {
				err = io.ErrUnexpectedEOF
			}
		}
	}

	// If we can return an EOF here along with the read data, do
	// so. This is optional per the io.Reader contract, but doing
	// so helps the HTTP transport code recycle its connection
	// earlier (since it will see this EOF itself), even if the
	// client doesn't do future reads or Close.
	if err == nil && n > 0 {
		if lr, ok := b.src.(*io.LimitedReader); ok && lr.N == 0 {
			err = io.EOF
			b.sawEOF = true
		}
	}

	return n, err
}

var (
	singleCRLF = []byte("\r\n")
	doubleCRLF = []byte("\r\n\r\n")
)

func seeUpcomingDoubleCRLF(r *bufio.Reader) bool {
	for peekSize := 4; ; peekSize++ {
		// This loop stops when Peek returns an error,
		// which it does when r's buffer has been filled.
		buf, err := r.Peek(peekSize)
		if bytes.HasSuffix(buf, doubleCRLF) {
			return true
		}
		if err != nil {
			break
		}
	}
	return false
}

// unreadDataSizeLocked returns the number of bytes of unread input.
// It returns -1 if unknown.
// b.mu must be held.
func (b *body) unreadDataSizeLocked() int64 {
	if lr, ok := b.src.(*io.LimitedReader); ok {
		return lr.N
	}
	return -1
}

// maxPostHandlerReadBytes is the max number of Request.Body bytes not
// consumed by a handler that the server will read from the client
// in order to keep a connection alive.  If there are more bytes than
// this then the server to be paranoid instead sends a "Connection:
// close" response.
//
// This number is approximately what a typical machine's TCP buffer
// size is anyway.  (if we have the bytes on the machine, we might as
// well read them)
const maxPostHandlerReadBytes = 256 << 10

func (b *body) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil
	}
	var err error
	switch {
	case b.sawEOF:
		// Already saw EOF, so no need going to look for it.
	case b.doEarlyClose:
		// Read up to maxPostHandlerReadBytes bytes of the body, looking for
		// for EOF (and trailers), so we can re-use this connection.
		if lr, ok := b.src.(*io.LimitedReader); ok && lr.N > maxPostHandlerReadBytes {
			// There was a declared Content-Length, and we have more bytes remaining
			// than our maxPostHandlerReadBytes tolerance. So, give up.
			b.earlyClose = true
		} else {
			var n int64
			// Consume the body, or, which will also lead to us reading
			// the trailer headers after the body, if present.
			n, err = io.CopyN(ioutil.Discard, bodyLocked{b}, maxPostHandlerReadBytes)
			if err == io.EOF {
				err = nil
			}
			if n == maxPostHandlerReadBytes {
				b.earlyClose = true
			}
		}
	default:
		// Fully consume the body, which will also lead to us reading
		// the trailer headers after the body, if present.
		_, err = io.Copy(ioutil.Discard, bodyLocked{b})
	}
	b.closed = true
	return err
}

func (b *body) didEarlyClose() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.earlyClose
}

// bodyLocked is a io.Reader reading from a *body when its mutex is
// already held.
type bodyLocked struct {
	b *body
}

func (bl bodyLocked) Read(p []byte) (n int, err error) {
	if bl.b.closed {
		return 0, ErrBodyReadAfterClose
	}
	return bl.b.readLocked(p)
}
