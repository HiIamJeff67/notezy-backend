package adapterscontract

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

func (s *YjsDocumentState) UnmarshalBytes(payload []byte) error {
	*s = YjsDocumentState{}
	if len(payload) > YjsDocumentMaximumLoadPayloadBytes {
		return errors.New("Yjs document state exceeds the maximum load payload size")
	}
	if len(payload) < 36 {
		return errors.New("invalid yjs document state payload")
	}

	s.LastUpdateSequence = int64(binary.BigEndian.Uint64(payload[0:8]))
	s.CompactedUntilSequence = int64(binary.BigEndian.Uint64(payload[8:16]))
	s.ProjectedUntilSequence = int64(binary.BigEndian.Uint64(payload[16:24]))
	snapshotLength := binary.BigEndian.Uint32(payload[24:28])
	stateVectorLength := binary.BigEndian.Uint32(payload[28:32])
	updateCount := binary.BigEndian.Uint32(payload[32:36])
	if s.LastUpdateSequence < 0 ||
		s.CompactedUntilSequence < 0 ||
		s.CompactedUntilSequence > s.LastUpdateSequence ||
		s.ProjectedUntilSequence < -1 ||
		s.ProjectedUntilSequence > s.LastUpdateSequence {
		return errors.New("invalid yjs document state")
	}

	offset := 36
	if uint64(snapshotLength) > uint64(len(payload)-offset) {
		return errors.New("invalid yjs document state")
	}
	s.Snapshot = append(s.Snapshot, payload[offset:offset+int(snapshotLength)]...)
	offset += int(snapshotLength)
	if uint64(stateVectorLength) > uint64(len(payload)-offset) {
		return errors.New("invalid yjs document state")
	}
	s.StateVector = append(s.StateVector, payload[offset:offset+int(stateVectorLength)]...)
	offset += int(stateVectorLength)

	s.Updates = make([]YjsDocumentUpdate, 0, updateCount)
	for index := uint32(0); index < updateCount; index++ {
		if len(payload)-offset < 12 {
			return errors.New("invalid yjs document state")
		}

		updateSequence := int64(binary.BigEndian.Uint64(payload[offset : offset+8]))
		updateLength := binary.BigEndian.Uint32(payload[offset+8 : offset+12])
		offset += 12
		if updateSequence <= s.CompactedUntilSequence ||
			updateSequence > s.LastUpdateSequence ||
			uint64(updateLength) > uint64(len(payload)-offset) {
			return errors.New("invalid yjs document state")
		}

		s.Updates = append(s.Updates, YjsDocumentUpdate{
			UpdateSequence: updateSequence,
			Payload:        append([]byte{}, payload[offset:offset+int(updateLength)]...),
		})
		offset += int(updateLength)
	}
	if offset != len(payload) {
		return errors.New("invalid yjs document state")
	}

	return nil
}

func MarshalYjsUpdateSequence(updateSequence int64) []byte {
	payload := make([]byte, 8)
	binary.BigEndian.PutUint64(payload, uint64(updateSequence))

	return payload
}
