package problem

import "time"

func WithTypeURI(t string) RegisteredMembersOption {
	return func(rm *RegisteredMembers) {
		rm.TypeURI = t
	}
}

func WithTypeTitle(t string) RegisteredMembersOption {
	return func(rm *RegisteredMembers) {
		rm.TypeTitle = t
	}
}

func WithStatusCode(s int) RegisteredMembersOption {
	return func(rm *RegisteredMembers) {
		rm.StatusCode = s
	}
}

func WithDetailMessage(d string) RegisteredMembersOption {
	return func(rm *RegisteredMembers) {
		rm.DetailMessage = d
	}
}

func WithInstanceURI(i string) RegisteredMembersOption {
	return func(rm *RegisteredMembers) {
		rm.InstanceURI = i
	}
}

type (
	RegisteredMembersOption func(*RegisteredMembers)
	Option                  func(any)
)

func WithRetryAfter[T time.Duration | time.Time](t T) Option {
	return func(a any) {
		if p, ok := a.(hasRetryAfter); ok {
			ta := any(t)
			switch t := ta.(type) {
			case time.Time:
				p.SetRetryAfterTime(t)
			case time.Duration:
				p.SetRetryAfterDuration(t)
			}
		}
	}
}

type hasRetryAfter interface {
	SetRetryAfterTime(t time.Time)
	SetRetryAfterDuration(d time.Duration)
}
