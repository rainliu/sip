package sip

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestReadRequest(t *testing.T) {
	var tvi = []string{
		"REGISTER sip:nist.gov SIP/2.0\r\n" +
			"Via: SIP/2.0/UDP 129.6.55.182:14826\r\n" +
			"Max-Forwards: 70\r\n" +
			"From: <sip:mranga@nist.gov>;tag=6fcd5c7ace8b4a45acf0f0cd539b168b;epid=0d4c418ddf\r\n" +
			"To: <sip:mranga@nist.gov>\r\n" +
			"Call-ID: c5679907eb954a8da9f9dceb282d7230@129.6.55.182\r\n" +
			"CSeq: 1 REGISTER\r\n" +
			"Contact: <sip:129.6.55.182:14826>;methods=\"INVITE, MESSAGE, INFO, SUBSCRIBE, OPTIONS, BYE, CANCEL, NOTIFY, ACK, REFER\"\r\n" +
			"User-Agent: RTC/(Microsoft RTC)\r\n" +
			"Event:  registration\r\n" +
			"Allow-Events: presence\r\n" +
			"Content-Length: 0\r\n\r\n",

		"INVITE sip:littleguy@there.com:5060 SIP/2.0\r\n" +
			"Via: SIP/2.0/UDP 65.243.118.100:5050\r\n" +
			"From: M. Ranganathan  <sip:M.Ranganathan@sipbakeoff.com>;tag=1234\r\n" +
			"To: \"littleguy@there.com\" <sip:littleguy@there.com:5060> \r\n" +
			"Call-ID: Q2AboBsaGn9!?x6@sipbakeoff.com \r\n" +
			"CSeq: 1 INVITE \r\n" +
			"Content-Length: 247\r\n\r\n" +
			"v=0\r\n" +
			"o=4855 13760799956958020 13760799956958020 IN IP4  129.6.55.78\r\n" +
			"s=mysession session\r\n" +
			"p=+46 8 52018010\r\n" +
			"c=IN IP4  129.6.55.78\r\n" +
			"t=0 0\r\n" +
			"m=audio 6022 RTP/AVP 0 4 18\r\n" +
			"a=rtpmap:0 PCMU/8000\r\n" +
			"a=rtpmap:4 G723/8000\r\n" +
			"a=rtpmap:18 G729A/8000\r\n" +
			"a=ptime:20\r\n",
	}
	//	var tvo = []string{
	//		"REGISTER sip:nist.gov SIP/2.0\r\n" +
	//			"Via: SIP/2.0/UDP 129.6.55.182:14826\r\n" +
	//			"Max-Forwards: 70\r\n" +
	//			"From: <sip:mranga@nist.gov>;tag=6fcd5c7ace8b4a45acf0f0cd539b168b;epid=0d4c418ddf\r\n" +
	//			"To: <sip:mranga@nist.gov>\r\n" +
	//			"Call-ID: c5679907eb954a8da9f9dceb282d7230@129.6.55.182\r\n" +
	//			"CSeq: 1 REGISTER\r\n" +
	//			"Contact: <sip:129.6.55.182:14826>;methods=\"INVITE, MESSAGE, INFO, SUBSCRIBE, OPTIONS, BYE, CANCEL, NOTIFY, ACK, REFER\"\r\n" +
	//			"User-Agent: RTC/(Microsoft RTC)\r\n" +
	//			"Event: registration\r\n" +
	//			"Allow-Events: presence\r\n" +
	//			"Content-Length: 0\r\n\r\n",

	//		"INVITE sip:littleguy@there.com:5060 SIP/2.0\r\n" +
	//			"Via: SIP/2.0/UDP 65.243.118.100:5050\r\n" +
	//			"From: \"M. Ranganathan\" <sip:M.Ranganathan@sipbakeoff.com>;tag=1234\r\n" +
	//			"To: \"littleguy@there.com\" <sip:littleguy@there.com:5060>\r\n" +
	//			"Call-ID: Q2AboBsaGn9!?x6@sipbakeoff.com\r\n" +
	//			"CSeq: 1 INVITE\r\n" +
	//			"Content-Length: 247\r\n\r\n" +
	//			"v=0\r\n" +
	//			"o=4855 13760799956958020 13760799956958020 IN IP4  129.6.55.78\r\n" +
	//			"s=mysession session\r\n" +
	//			"p=+46 8 52018010\r\n" +
	//			"c=IN IP4  129.6.55.78\r\n" +
	//			"t=0 0\r\n" +
	//			"m=audio 6022 RTP/AVP 0 4 18\r\n" +
	//			"a=rtpmap:0 PCMU/8000\r\n" +
	//			"a=rtpmap:4 G723/8000\r\n" +
	//			"a=rtpmap:18 G729A/8000\r\n" +
	//			"a=ptime:20\r\n",
	//	}

	for i := 0; i < len(tvi); i++ {
		b := bufio.NewReader(strings.NewReader(tvi[i]))
		req, err := ReadMessage(b)
		if err != nil {
			t.Log(req)
		} else {
			var buffer bytes.Buffer
			err = req.Write(&buffer)
			if err != nil {
				t.Log(err)
			} else {
				t.Log(buffer.String())
			}
		}
	}
}
