package gaugin_middleware

import (
	"fmt"
	"strings"
)

type ServerTimings struct {
	durations map[string]int64
	flags     map[string]interface{}
}

func NewServerTimings() *ServerTimings {
	return &ServerTimings{
		durations: make(map[string]int64),
		flags:     make(map[string]interface{}),
	}
}

func (st *ServerTimings) Set(name string, dur int64) {
	st.durations[name] = dur
}

func (st *ServerTimings) SetFlag(name string) {
	st.flags[name] = nil
}

func (st *ServerTimings) String() string {
	entries := make([]string, 0, len(st.durations)+len(st.flags))
	for flag := range st.flags {
		entries = append(entries, flag)
	}
	for name, dur := range st.durations {
		entries = append(entries, fmt.Sprintf("%s;dur=%d", name, dur))
	}
	return strings.Join(entries, ",")
}
