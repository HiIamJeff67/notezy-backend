package realtimetypes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

type YjsDocumentState struct {
	Snapshot               []byte
	StateVector            []byte
	LastUpdateSequence     int64
	CompactedUntilSequence int64
	ProjectedUntilSequence int64
	Updates                []YjsDocumentUpdate
}

type YjsDocumentUpdate struct {
	UpdateSequence int64
	Payload        []byte
}

func (s YjsDocumentState) MarshalBytes() ([]byte, error) {
	if len(s.Snapshot) > math.MaxUint32 || len(s.StateVector) > math.MaxUint32 || len(s.Updates) > math.MaxUint32 {
		return nil, errors.New("invalid yjs document state")
	}
	if s.LastUpdateSequence < 0 ||
		s.CompactedUntilSequence < 0 ||
		s.CompactedUntilSequence > s.LastUpdateSequence ||
		s.ProjectedUntilSequence < -1 ||
		s.ProjectedUntilSequence > s.LastUpdateSequence {
		return nil, errors.New("invalid yjs document state")
	}

	payload := bytes.NewBuffer(make([]byte, 0, 36+len(s.Snapshot)+len(s.StateVector)))
	if err := binary.Write(payload, binary.BigEndian, s.LastUpdateSequence); err != nil {
		return nil, err
	}
	if err := binary.Write(payload, binary.BigEndian, s.CompactedUntilSequence); err != nil {
		return nil, err
	}
	if err := binary.Write(payload, binary.BigEndian, s.ProjectedUntilSequence); err != nil {
		return nil, err
	}
	if err := binary.Write(payload, binary.BigEndian, uint32(len(s.Snapshot))); err != nil {
		return nil, err
	}
	if err := binary.Write(payload, binary.BigEndian, uint32(len(s.StateVector))); err != nil {
		return nil, err
	}
	if err := binary.Write(payload, binary.BigEndian, uint32(len(s.Updates))); err != nil {
		return nil, err
	}
	if _, err := payload.Write(s.Snapshot); err != nil {
		return nil, err
	}
	if _, err := payload.Write(s.StateVector); err != nil {
		return nil, err
	}

	for _, update := range s.Updates {
		if len(update.Payload) > math.MaxUint32 {
			return nil, errors.New("invalid yjs document update")
		}
		if err := binary.Write(payload, binary.BigEndian, update.UpdateSequence); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.BigEndian, uint32(len(update.Payload))); err != nil {
			return nil, err
		}
		if _, err := payload.Write(update.Payload); err != nil {
			return nil, err
		}
	}

	return payload.Bytes(), nil
}

func MarshalYjsUpdateSequence(updateSequence int64) []byte {
	payload := make([]byte, 8)
	binary.BigEndian.PutUint64(payload, uint64(updateSequence))

	return payload
}
