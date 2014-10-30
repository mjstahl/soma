package lib

import (
	"github.com/socialmachines/soma/ast"
	"github.com/socialmachines/soma/parse"
	"github.com/socialmachines/soma/rt"
)

type primitiveFn func(*rt.AsyncMsg)

func init() {
	rt.NIL = loadPrimitiveObj("Nil", 0x01, nilBehaviorMap)
	rt.TRUE = loadPrimitiveObj("True", 0x03, trueBehaviorMap)
	rt.FALSE = loadPrimitiveObj("False", 0x05, falseBehaviorMap)
}

func loadPrimitiveObj(name string, id uint64, behaviorMap map[string]primitiveFn) rt.Value {
	obj := rt.CreateObject(&ast.Global{Value: name}, nil, id)
	obj.New()

	behaviorObj := rt.CreateObject(nil, nil, 0)
	startPrimitiveBehaviors(obj, behaviorObj, behaviorMap)

	rt.RT.Globals.Insert(name, obj.ID)
	return obj
}

func startPrimitiveBehaviors(recv *rt.Object, behaviorObj rt.Value, behaviorMap map[string]primitiveFn) {
	for name, _ := range behaviorMap {
		recv.Behaviors[name] = behaviorObj.OID()
	}

	go func() {
		for {
			msg := <-behaviorObj.Address()
			amsg := msg.(*rt.AsyncMsg)
			if behaviorFn, fn := behaviorMap[amsg.Behavior]; !fn {
				(rt.NIL).Return(amsg)
			} else {
				behaviorFn(amsg)
			}
		}
	}()
}

func loadBehaviors(src string) {
	exprs, _ := parse.ParseExpr(src)
	for _, expr := range exprs {
		expr.Visit(nil)
	}
}
