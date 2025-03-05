// Copyright 2025-present Gustavo "Guz" L. de Mello
// Copyright 2025-present The Lored.dev Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exception

import (
	"encoding"
	"fmt"
	"log/slog"
)

type Severity int

var (
	_ fmt.Stringer             = (Severity)(0)
	_ encoding.TextMarshaler   = (Severity)(0)
	_ encoding.TextUnmarshaler = (*Severity)(nil)
)

const (
	DEBUG Severity = Severity(slog.LevelDebug)
	INFO  Severity = Severity(slog.LevelInfo)
	WARN  Severity = Severity(slog.LevelWarn)
	ERROR Severity = Severity(slog.LevelError)
	FATAL Severity = Severity(slog.LevelError * 2)

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
