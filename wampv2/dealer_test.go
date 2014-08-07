package wampv2

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type TestCallee struct {
	received Message
}

func (c *TestCallee) SendError(msg *Error)               { c.received = msg }
func (c *TestCallee) SendRegistered(msg *Registered)     { c.received = msg }
func (c *TestCallee) SendUnregistered(msg *Unregistered) { c.received = msg }
func (c *TestCallee) SendInvocation(msg *Invocation)     { c.received = msg }

func TestRegister(t *testing.T) {
	Convey("Registering a procedure", t, func() {
		dealer := NewDefaultDealer()
		callee := &TestCallee{}
		testProcedure := URI("turnpike.test.endpoint")
		msg := &Register{Request: 123, Procedure: testProcedure}
		dealer.Register(callee, msg)

		Convey("The callee should have received a REGISTERED message", func() {
			reg := callee.received.(*Registered).Registration
			So(reg, ShouldNotEqual, 0)
		})

		Convey("The dealer should have the endpoint registered", func() {
			reg := callee.received.(*Registered).Registration
			reg2, ok := dealer.registrations[testProcedure]
			So(ok, ShouldBeTrue)
			So(reg, ShouldEqual, reg2)
			proc, ok := dealer.procedures[reg]
			So(ok, ShouldBeTrue)
			So(proc.Procedure, ShouldEqual, testProcedure)
		})

		Convey("The same procedure cannot be registered more than once", func() {
			msg := &Register{Request: 321, Procedure: testProcedure}
			dealer.Register(callee, msg)
			So(callee.received, ShouldHaveSameTypeAs, &Error{})
		})
	})
}
