package rt

import (
	"strconv"
	"time"
)

type Promise struct {
	ID   uint64
	Addr Mailbox

	Value  uint64
	Valued chan bool

	Behaviors map[string]Value

	Blocking []Message
}

func CreatePromise() *Promise {
	id := NewID(PROMISE)

	n := 128
	promise := &Promise{ID: id, Addr: make(Mailbox, n), Behaviors: map[string]Value{}, Blocking: []Message{}}
	promise.Valued = make(chan bool, 1)

	RT.Heap.Insert(id, promise)
	go promise.New()

	return promise
}

func (promise *Promise) New() {
	for {
		select {
		case <-promise.Valued:
			for _, msg := range promise.Blocking {
				forwardMessage(promise, msg)
			}
			promise.Blocking = []Message{}
		case msg := <-promise.Address():
			msg.ForwardMessage(promise)
		}
	}
}

func (p *Promise) Address() Mailbox {
	return p.Addr
}

func (p *Promise) LookupBehavior(name string) Value {
	return p.Behaviors[name]
}

func (p *Promise) OID() uint64 {
	return p.ID
}

// Any Definition body or Block will return the last expression to be
// evaluated.  If the last expression is a Message (unary, binary, or keyword)
// then the result will be a Promise. To return a promise, we don't know when
// the promised the value of the Message will be available so we must send
// it an asynchronous message (because it could be returned from a remote
// machine) request the Promise's value on behalf of the received message.
// Therefore we send the "value" asynchronous message to the Promise, but
// instead of creating a new Promise, we use the same Promise ID of the
// original message.
//
func (p *Promise) Return(am *AsyncMsg) {
	async := &AsyncMsg{[]uint64{0, 0}, "value", am.PromisedTo}
	p.Address() <- async
}

func (p *Promise) String() (repr string) {
	for p.Value == 0 {
		time.Sleep(10 * time.Millisecond)
	}

	switch p.Value & 0xF {
	case 0x7:
		integer, base := int64(p.Value)>>8, 10
		return strconv.FormatInt(integer, base)
	default:
		return RT.Heap.Lookup(p.Value).String()
	}
}
