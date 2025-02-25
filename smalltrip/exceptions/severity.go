package exceptions

import (
	"encoding"
	"fmt"
)

type Severity int

var (
	_ fmt.Stringer             = (Severity)(0)
	_ encoding.TextMarshaler   = (Severity)(0)
	_ encoding.TextUnmarshaler = (*Severity)(nil)
)

const (
	DEBUG Severity = -4
	INFO  Severity = 0
	WARN  Severity = 4
	ERROR Severity = 8
	FATAL Severity = 16

	stringDEBUG = "DEBUG"
	stringINFO  = "INFO"
	stringWARN  = "WARN"
	stringERROR = "ERROR"
	stringFATAL = "FATAL"

	stringUNDEFINED = "UNDEFINED"
)

func (s Severity) String() string {
	switch s {
	case DEBUG:
		return stringDEBUG
	case INFO:
		return stringINFO
	case WARN:
		return stringWARN
	case ERROR:
		return stringERROR
	case FATAL:
		return stringFATAL
	default:
		return fmt.Sprintf(stringUNDEFINED)
	}
}

func (s Severity) MarshalText() ([]byte, error) {
	str := s.String()
	if str == stringUNDEFINED {
		return nil, fmt.Errorf("severity of value %q does not exists", s)
	}
	return []byte(str), nil
}

func (s *Severity) UnmarshalText(text []byte) error {
	switch string(text) {
	case stringDEBUG:
		*s = DEBUG
	case stringINFO:
		*s = INFO
	case stringWARN:
		*s = WARN
	case stringERROR:
		*s = ERROR
	case stringFATAL:
		*s = FATAL
	default:
		return fmt.Errorf("severity level %d does not exists", s)
	}
	return nil
}
