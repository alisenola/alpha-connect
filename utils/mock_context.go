package utils

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/ctxext"
	"time"
)

type MockContext struct {
	*actor.RootContext
	ParentF                  func() *actor.PID
	SelfF                    func() *actor.PID
	ActorF                   func() actor.Actor
	ActorSystemF             func() *actor.ActorSystem
	ReceiveTimeoutF          func() time.Duration
	ChildrenF                func() []*actor.PID
	RespondF                 func(response interface{})
	StashF                   func()
	WatchF                   func(pid *actor.PID)
	UnwatchF                 func(pid *actor.PID)
	SetReceiveTimeoutF       func(d time.Duration)
	CancelReceiveTimeoutF    func()
	ForwardF                 func(pid *actor.PID)
	ReenterAfterF            func(f *actor.Future, continuation func(res interface{}, err error))
	MessageF                 func() interface{}
	MessageHeaderF           func() actor.ReadonlyMessageHeader
	SenderF                  func() *actor.PID
	SendF                    func(pid *actor.PID, message interface{})
	RequestF                 func(pid *actor.PID, message interface{})
	RequestWithCustomSenderF func(pid *actor.PID, message interface{}, sender *actor.PID)
	RequestFutureF           func(pid *actor.PID, message interface{}, timeout time.Duration) *actor.Future
	ReceiveF                 func(envelope *actor.MessageEnvelope)
	SpawnF                   func(props *actor.Props) *actor.PID
	SpawnPrefixF             func(props *actor.Props, prefix string) *actor.PID
	SpawnNamedF              func(props *actor.Props, id string) (*actor.PID, error)
	StopF                    func(pid *actor.PID)
	StopFutureF              func(pid *actor.PID) *actor.Future
	PoisonF                  func(pid *actor.PID)
	PoisonFutureF            func(pid *actor.PID) *actor.Future
	GetF                     func(id ctxext.ContextExtensionID) ctxext.ContextExtension
	SetF                     func(ext ctxext.ContextExtension)
}

func NewMockContext(root *actor.RootContext) *MockContext {
	return &MockContext{RootContext: root}
}

func (m MockContext) Parent() *actor.PID {
	return m.ParentF()
}

func (m MockContext) Self() *actor.PID {
	if m.SelfF != nil {
		return m.SelfF()
	} else {
		return m.RootContext.Self()
	}
}

func (m MockContext) Actor() actor.Actor {
	return m.ActorF()
}

func (m MockContext) ActorSystem() *actor.ActorSystem {
	return m.ActorSystemF()
}

func (m MockContext) ReceiveTimeout() time.Duration {
	return m.ReceiveTimeoutF()
}

func (m MockContext) Children() []*actor.PID {
	return m.ChildrenF()
}

func (m MockContext) Respond(response interface{}) {
	m.RespondF(response)
}

func (m MockContext) Stash() {
	m.StashF()
}

func (m MockContext) Watch(pid *actor.PID) {
	m.WatchF(pid)
}

func (m MockContext) Unwatch(pid *actor.PID) {
	m.UnwatchF(pid)
}

func (m MockContext) SetReceiveTimeout(d time.Duration) {
	m.SetReceiveTimeoutF(d)
}

func (m MockContext) CancelReceiveTimeout() {
	m.CancelReceiveTimeoutF()
}

func (m MockContext) Forward(pid *actor.PID) {
	m.ForwardF(pid)
}

func (m MockContext) ReenterAfter(f *actor.Future, continuation func(res interface{}, err error)) {
	m.ReenterAfterF(f, continuation)
}

func (m MockContext) Message() interface{} {
	return m.MessageF()
}

func (m MockContext) MessageHeader() actor.ReadonlyMessageHeader {
	return m.MessageHeaderF()
}

func (m MockContext) Sender() *actor.PID {
	return m.SenderF()
}

func (m MockContext) Send(pid *actor.PID, message interface{}) {
	if m.SendF != nil {
		m.SendF(pid, message)
	} else {
		m.RootContext.Send(pid, message)
	}
}

func (m MockContext) Request(pid *actor.PID, message interface{}) {
	m.RequestF(pid, message)
}

func (m MockContext) RequestWithCustomSender(pid *actor.PID, message interface{}, sender *actor.PID) {
	m.RequestWithCustomSenderF(pid, message, sender)
}

func (m MockContext) RequestFuture(pid *actor.PID, message interface{}, timeout time.Duration) *actor.Future {
	return m.RequestFutureF(pid, message, timeout)
}

func (m MockContext) Receive(envelope *actor.MessageEnvelope) {
	m.ReceiveF(envelope)
}

func (m MockContext) Spawn(props *actor.Props) *actor.PID {
	return m.SpawnF(props)
}

func (m MockContext) SpawnPrefix(props *actor.Props, prefix string) *actor.PID {
	return m.SpawnPrefixF(props, prefix)
}

func (m MockContext) SpawnNamed(props *actor.Props, id string) (*actor.PID, error) {
	return m.SpawnNamedF(props, id)
}

func (m MockContext) Stop(pid *actor.PID) {
	m.StopF(pid)
}

func (m MockContext) StopFuture(pid *actor.PID) *actor.Future {
	return m.StopFutureF(pid)
}

func (m MockContext) Poison(pid *actor.PID) {
	m.PoisonF(pid)
}

func (m MockContext) PoisonFuture(pid *actor.PID) *actor.Future {
	return m.PoisonFutureF(pid)
}

func (m MockContext) Get(id ctxext.ContextExtensionID) ctxext.ContextExtension {
	return m.GetF(id)
}

func (m MockContext) Set(ext ctxext.ContextExtension) {
	m.SetF(ext)
}

func init() {
	mc := &MockContext{}
	var ac actor.Context = mc
	_ = ac
}
