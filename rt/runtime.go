package rt

import (
	crypto "crypto/rand"
	"math/rand"
	"net"
	"runtime"
)

var NIL Value
var FALSE Value
var TRUE Value

var RT *Runtime

const (
	PROMISE = 0
	OBJECT  = 1
)

type Expr interface {
	Eval(*Scope) Value
	Visit(*Scope) Value
}

type Value interface {
	Address() Mailbox
	LookupBehavior(string) Value
	OID() uint64
	Return(am *AsyncMsg)
	String() string
}

type Mailbox chan Message

type Runtime struct {
	Globals *Scope
	Heap    *Heap
	ID      uint64
	IPAddr  net.IP
	Peers   map[uint64]*Peer
	Port    int
}

func InitRuntime() *Runtime {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ipAddr, _ := LocalIP()
	rtid := uint64(randomID()&0xFFFFFFF0) << 32
	return &Runtime{Globals: NewScope(nil), Heap: NewHeap(), IPAddr: ipAddr, ID: rtid, Peers: map[uint64]*Peer{}}
}

// |----- 28bits -----| ----- 28 bits ----- | ----- 7bits -----|----- 1bit -----|
// |    Runtime ID    |      Object ID      |   Primitive Type |   Object Type  |
// |------------------|---------------------|------------------|----------------|
//
func NewID(t uint64) (oid uint64) {
	for {
		obj := uint64(randomID() & 0xFFFFFFF0)
		oid = (RT.ID | (obj << 4)) | t

		if exists := RT.Heap.Lookup(oid); exists == nil {
			// reserve a spot for the new object
			RT.Heap.Insert(oid, nil)
			break
		}
	}

	return oid
}

func randomID() uint32 {
	i, _ := crypto.Prime(crypto.Reader, 64)
	seed := i.Int64()
	r := rand.New(rand.NewSource(seed))
	return r.Uint32() & 0xFFFFFFF0
}

func init() {
	RT = InitRuntime()
}
