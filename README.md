<h1 align="center">
  <img src="https://raw.githubusercontent.com/gcapizzi/moka/master/images/logo.png" />
  <br />
  Moka
</h1>

<p align="center">
  <strong>A Go mocking framework.</strong>
  <br />
  <a href="https://godoc.org/github.com/gcapizzi/moka">
    <img alt="GoDOc" src="https://godoc.org/github.com/gcapizzi/moka?status.svg" />
  </a>
  <a href="https://travis-ci.org/gcapizzi/moka">
    <img alt="TravisCI" src="https://travis-ci.org/gcapizzi/moka.svg?branch=master" />
  </a>
</p>

<strong>Moka</strong> is a mocking framework for the [Go programming
language](https://golang.org). Moka works very well with the [Ginkgo testing
framework](http://onsi.github.io/ginkgo), but can be easily used with any other
testing framework, including the `testing` package from the standard library.

## Getting Moka

```
go get github.com/gcapizzi/moka
```

## Setting Up Moka

### Ginkgo

Moka is designed to play well with [Ginkgo](http://onsi.github.io/ginkgo). All
you'll need to do is:

* import the `moka` package;
* register Ginkgo's `Fail` function as Moka's double fail handler using
  `RegisterDoublesFailHandler`.

Here's an example:

```go
package game_test

import (
	. "github.com/gcapizzi/moka"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGame(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterDoublesFailHandler(Fail)
	RunSpecs(t, "Game Suite")
}
```

### `testing` and other frameworks

Support for the [`testing`](https://golang.org/pkg/testing) package hasn't been
added yet, but this doesn't mean you can't use Moka with it.

Here is the type for the Moka doubles fail handler:

```go
type FailHandler func(message string, callerSkip ...int)
```

This type is modeled to match Ginkgo's `Fail` function. To use Moka with the
`testing` package, just provide a doubles fail handler that makes the test
fail!

Here's an example:

```go
package game

import (
	. "github.com/gcapizzi/moka"

	"testing"
)

func TestGame(t *testing.T) {
	RegisterDoublesFailHandler(func(message string, callerSkip ...int) {
		t.Fatal(message)
	})

	// use Moka here
}
```

## Getting Started: Building Your First Double

A test double is an object that stands in for another object in your system
during tests. The first step to use Moka is to declare a double type. The type
will have to:

* implement the same interface as the replaced object: this means that only
  objects used through interfaces can be replaced by Moka doubles;
* embed the `moka.Double` type;
* delegate any method that you'll need to stub/mock to the embedded
  `moka.Double` instance, using the `Call` method.

Let's build a double type for a `Die` interface:

```go
package dice

import (
	. "github.com/gcapizzi/moka"
)

type Die interface {
	Roll(times int) []int
}

type DieDouble struct {
	Double
}

func NewDieDouble() DieDouble {
	return DieDouble{Double: NewStrictDouble()}
}

func (d DieDouble) Roll(times int) []int {
	returnValues, _ := d.Call("Roll", times)
	returnedRolls, _ := returnValues[0].([]int)
	return returnedRolls
}
```

Some notes:

* The `Double` instance we are embedding is of type `StrictDouble`: strict
  doubles will fail the test if they receive a call on a method that wasn't
  previously allowed or expected.
* If the `Call` invocation fails, it will both return an error as its second
  return value, and invoke the fail handler. Since we know the fail handler
  will immediately stop the execution of the test, we don't need to check for
  the error.
* This style of type assertions allow us to have `nil` return values.

## Allowing interactions

Now that our double type is ready, let's use it in our tests! We will test a
`Game` type, which looks like this:

```go
package game

type Game struct {
	die Die
}

func NewGame(die Die) Game {
	return Game{die: die}
}

func (g Game) Score() int {
	rolls := g.die.Roll(3)
	return rolls[0] + rolls[1] + rolls[2]
}
```

Here is the test:

```go
package game_test

import (
	. "github.com/gcapizzi/moka/examples/game"

	. "github.com/gcapizzi/moka"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Game", func() {
	var die DieDouble
	var game Game

	BeforeEach(func() {
		die = NewDieDouble()
		game = NewGame(die)
	})

	Describe("Score", func() {
		It("returns the sum of three die rolls", func() {
			AllowDouble(die).To(ReceiveCallTo("Roll").With(3).AndReturn([]int{1, 2, 3}))

			Expect(game.Score()).To(Equal(6))
		})
	})
})
```

## Typed doubles

You might be wondering: what happens if I allow a method call that would be
impossible to perform, given the type of my double? Let's say we forgot how the
`Die` interface was defined, and wrote a test like this:

```go
Describe("Score", func() {
	It("returns the sum of three die rolls", func() {
		AllowDouble(die).To(ReceiveCallTo("Cast").With(3).AndReturn([]int{1, 2, 3}))

		Expect(game.Score()).To(Equal(9))
	})
})
```

The `Die` interface has no method called `Cast`, so our configured interaction
will never happen!

To avoid this kind of problems, Moka provides _types
doubles_, which are associated with a type and will make sure that any
configured interactions actually match with the type.

To instantiate a typed double, use the `NewStrictDoubleWithTypeOf` constructor:

```go
func NewDieDouble() DieDouble {
	return DieDouble{Double: NewStrictDoubleWithTypeOf(DieDouble{})}
}
```

If run against a typed double, the previous test would fail with a message like this:

```
Invalid interaction: type 'DieDouble' has no method 'Cast'
```

## Expecting interactions

Sometimes allowing a method call is not enough. Some methods have side effects,
and we need to make sure they have been invoked in order to be confident that
our code is working.

For example, let's assume we added a `Logger` collaborator to our `Game` type:

```go
package game

import "fmt"

type Logger interface {
	Log(message string)
}

type StdoutLogger struct{}

func (l StdoutLogger) Log(message string) {
	fmt.Println(message)
}
```

We'll start, as usual, by building our double type:

```go
type LoggerDouble struct {
	Double
}

func NewLoggerDouble() LoggerDouble {
	return LoggerDouble{Double: NewStrictDoubleWithTypeOf(LoggerDouble{})}
}

func (d LoggerDouble) Log(message string) {
	d.Call("Log", message)
}
```

We can now add the `Logger` collaborator to `Game`, and test their interaction:

```go
Describe("Score", func() {
	It("returns the sum of three die rolls", func() {
		AllowDouble(die).To(ReceiveCallTo("Roll").With(3).AndReturn([]int{1, 2, 3}))
		ExpectDouble(logger).To(ReceiveCallTo("Log").With("[1, 2, 3]"))

		Expect(game.Score()).To(Equal(6))

		VerifyCalls(logger)
	})
})
```

We use `ExpectDouble()` to expect method calls on a double, and `VerifyCalls`
to verify that the calls have actually been made.
