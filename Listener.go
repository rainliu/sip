package sip

type Listener interface {
	ProcessRequest(requestEvent RequestEvent)
	ProcessResponse(responseEvent ResponseEvent)
	ProcessTimeout(timeoutEvent TimeoutEvent)
}
