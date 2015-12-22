# Social Machines
Social Machines is a server-side programming language with a syntax greatly inspired by Smalltalk. Every object in Social Machines is a concurrent unit of computation, and libraries can be shared like peers in a BitTorrent network.

 * [Goals](#goals)
 * [Semantics](#semantics)
 * [Syntax](#syntax)
 * [Getting Started](#getting-started)
 * [License](#license)

## Goals
Social Machines is an exercise (read 'experiment') in language design and the evaluation of assumptions. The two assumptions that affected the design are as follows:

1. Every object is an isolated, concurrent unit, easing the burden on the programmer by removing the need to choose whether to thread code or not.
2. All modern computer languages are designed in the context that a language is designed for writing a single application running on a single machine. The network is an afterthought and therefore relegated to APIs. This is invalid due to the ubiquity of the internet.

## Semantics
The core of Social Machines is Carl Hewitt's [Actor Model](https://en.wikipedia.org/wiki/Actor_model). All Social Machines objects exhibit three core behaviors:

1. Create Objects
2. Send Messages
3. Receive Messages

All message passing in the Actor Model was done asynchronously. To make a program slightly easier to reason about, Social Machines adds [Promises](https://en.wikipedia.org/wiki/Futures_and_promises). All Social Machines Promises are first class.  Promises allow the source to behave sequentially at the potential expense of dead locks/live locks.

## Syntax
The syntax is greatly inspired by Smalltalk.  An example of the ```True``` object is listed below. The ```+``` indicates the defining of an External Behavior.
```smalltalk
    + True ifFalse: fBlock -> Nil.

    + True ifTrue: tBlock -> tBlock value.

    + True ifTrue: tBlock ifFalse: fBlock -> tBlock value.

    + True not -> False.

    + True & aBool ->
      aBool ifTrue: { self } ifFalse: { False }.

    + True | aBool -> self.

    + True ^ aBool ->
      aBool ifTrue: { False } ifFalse: { self }.
```
## Getting Started
```bash
    $ go get github.com/socialmachines/soma
    $ cd $GOPATH/src/github.com/socialmachines/soma
    $ go install

    # ensure $GOROOT/bin is on your path, if not
    $ export PATH=$PATH:$GOROOT/bin
    $ soma
```

#### Testing
```bash
    $ cd $GOPATH/src/github.com/socialmachines/soma/test
    $ go test
```

## License
Unless otherwise specified all source code is licensed under the [MIT license](https://github.com/socialmachines/soma/tree/master/MIT).

`parse/interface.go` and `scan/errors.go` are licensed under [Go's BSD license](https://golang.org/LICENSE).
