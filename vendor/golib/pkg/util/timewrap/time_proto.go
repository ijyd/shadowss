package timewrap

import "time"

// Timestamp returns the Time as a new Timestamp value.
func (m *Time) ProtoTime() *Timestamp {
	if m == nil {
		return &Timestamp{}
	}
	return &Timestamp{
		Seconds: m.Time.Unix(),
		Nanos:   int32(m.Time.Nanosecond()),
	}
}

// Size implements the protobuf marshalling interface.
func (m *Time) Size() (n int) {
	if m == nil || m.Time.IsZero() {
		return 0
	}
	return m.ProtoTime().Size()
}

// Reset implements the protobuf marshalling interface.
func (m *Time) Unmarshal(data []byte) error {
	if len(data) == 0 {
		m.Time = time.Time{}
		return nil
	}
	p := Timestamp{}
	if err := p.Unmarshal(data); err != nil {
		return err
	}
	m.Time = time.Unix(p.Seconds, int64(p.Nanos)).Local()
	return nil
}

// Marshal implements the protobuf marshalling interface.
func (m *Time) Marshal() (data []byte, err error) {
	if m == nil || m.Time.IsZero() {
		return nil, nil
	}
	return m.ProtoTime().Marshal()
}

// MarshalTo implements the protobuf marshalling interface.
func (m *Time) MarshalTo(data []byte) (int, error) {
	if m == nil || m.Time.IsZero() {
		return 0, nil
	}
	return m.ProtoTime().MarshalTo(data)
}
