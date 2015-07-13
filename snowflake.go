package snowflake

import (
	"errors"
	"fmt"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
	"sync"
	"time"
)

type state struct {
	machineID     int32
	seq           int32
	lastTimestamp int64
	m             sync.Mutex
}

// NewState returns a user defined state for snowflake ID generation.
// This function can be registered as UDSCreator.
func NewState(ctx *core.Context, params data.Map) (core.SharedState, error) {
	v, ok := params["machine_id"]
	if !ok {
		return nil, errors.New("machine_id parameter is missing")
	}

	mid, err := data.ToInt(v)
	if err != nil {
		return nil, fmt.Errorf("machine_id parameter cannot be converted to an integer: %v", err)
	}
	if mid < 0 || mid >= (1<<10) {
		return nil, fmt.Errorf("machine_id must be in [0, 1023]: %v", mid)
	}

	return &state{
		machineID: int32(mid),
	}, nil
}

func (s *state) Terminate(ctx *core.Context) error {
	return nil
}

const (
	timestampShift uint64 = 63 - 41
	machineIDShift uint64 = timestampShift - 10
)

func (s *state) gen(ctx *core.Context) (int64, error) {
	ts, seq, err := s.inc(ctx)
	if err != nil {
		return 0, err
	}
	return (ts << timestampShift) |
		(int64(s.machineID) << machineIDShift) |
		seq, nil
}

const (
	seqMax int32 = (1 << 12) - 1
)

func (s *state) inc(ctx *core.Context) (int64, int64, error) {
	// TODO: make this a CAS loop
	s.m.Lock()
	defer s.m.Unlock()
	for {
		now := time.Now().UnixNano() / int64(time.Millisecond)
		if now == s.lastTimestamp && s.seq > seqMax {
			continue // wait for at most 1ms
		}

		if now < s.lastTimestamp {
			ctx.Log().WithField("udf", "snowflake").
				Warnf("The system clock might have been changed during execution. ID generation stops for %v millseconds.", s.lastTimestamp-now)
			return 0, 0, fmt.Errorf("the systen clock may be changed during exection")

		} else if now > s.lastTimestamp {
			s.lastTimestamp = now
			s.seq = -1
		}

		s.seq++
		if s.seq <= seqMax {
			return s.lastTimestamp, int64(s.seq), nil
		}
		// sequence counter overflow
	}
}

// Snowflake generates a new ID based on snowflake ID generation algorithm.
// stateName must point to a shared state created by NewState.
func Snowflake(ctx *core.Context, stateName data.Value) (data.Value, error) {
	s, err := lookupState(ctx, stateName)
	if err != nil {
		return nil, err
	}

	id, err := s.gen(ctx)
	if err != nil {
		return nil, err
	}
	return data.Int(id), nil
}

func lookupState(ctx *core.Context, stateName data.Value) (*state, error) {
	name, err := data.AsString(stateName)
	if err != nil {
		return nil, fmt.Errorf("name of the state must be a string: %v", stateName)
	}

	st, err := ctx.SharedStates.Get(name)
	if err != nil {
		return nil, err
	}

	if s, ok := st.(*state); ok {
		return s, nil
	}
	return nil, fmt.Errorf("state '%v' cannot be converted to snowflake.state", name)
}
