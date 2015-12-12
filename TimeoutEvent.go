package sip

/**
 * This class represents an Timeout event that is passed from a SipProvider to
 * its SipListener. A specific message may need retransmitted on a specific
 * transaction numerous times before it is acknowledged by the receiver. If the
 * message is not acknowledged after a specified period in the underlying
 * implementation the transaction will expire, this occurs usually
 * after seven retransmissions. The mechanism to alert an application that a
 * message for a an underlying transaction needs retransmitted (i.e. 200OK) or
 * an underlying transaction has expired is a Timeout Event.
 * <p>
 * A Timeout Event can be of two different types, namely:
 * <ul>
 * <li>{@link Timeout#RETRANSMIT}
 * <li>{@link Timeout#TRANSACTION}
 * </ul>
 * A TimeoutEvent contains the following information:
 * <ul>
 * <li>source - the SipProvider that sent the TimeoutEvent.
 * <li>transaction - the transaction that this Timeout applies to.
 * <li>isServerTransaction - boolean indicating whether the transaction refers to
 * a client or server transaction.
 * <li>timeout - indicates what type of {@link Timeout} occurred.
 * </ul>
 *
 * @see Timeout
 */
type TimeoutEvent struct {
	transaction Transaction
	timeout     Timeout
}

/**
 * Constructs a TimeoutEvent to indicate a server retransmission or transaction
 * timeout.
 *
 * @param source - the source of TimeoutEvent.
 * @param serverTransaction - the server transaction that timed out.
 * @param timeout - indicates if this is a retranmission or transaction
 * timeout event.
 */
func NewTimeoutEvent(transaction Transaction, timeout Timeout) *TimeoutEvent {
	return &TimeoutEvent{
		transaction: transaction,
		timeout:     timeout,
	}
}

/**
 * Gets the transaction associated with this TimeoutEvent.
 *
 * @return  transaction associated with this TimeoutEvent
 */
func (this *TimeoutEvent) GetTransaction() Transaction {
	return this.transaction
}

/**
 * Indicates if the transaction associated with this TimeoutEvent is a server
 * transaction.
 *
 * @return returns true if a server transaction or false if a client
 * transaction.
 */
func (this *TimeoutEvent) IsServerTransaction() bool {
	_, ok := this.transaction.(ServerTransaction)
	return ok
}

/**
 * Gets the event type of this TimeoutEvent. The event type can be used to
 * determine if this Timeout Event is one of the following types:
 * <ul>
 * <li>{@link Timeout#TRANSACTION}
 * <li>{@link Timeout#RETRANSMIT}
 * </ul>
 *
 * @return the event type of this TimeoutEvent
 */
func (this *TimeoutEvent) GetTimeout() Timeout {
	return this.timeout
}
