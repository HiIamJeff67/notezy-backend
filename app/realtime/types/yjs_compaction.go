package realtimetypes

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

type YjsCompactionInput struct {
	Snapshot                   []byte
	StateVector                []byte
	BaseCompactedUntilSequence int64
	CutoffSequence             int64
	Updates                    []YjsDocumentUpdate
}

func (i YjsCompactionInput) MarshalBytes() ([]byte, error) {
	if len(i.Snapshot) > math.MaxUint32 || len(i.StateVector) > math.MaxUint32 || len(i.Updates) > math.MaxUint32 || i.BaseCompactedUntilSequence < 0 || i.CutoffSequence <= i.BaseCompactedUntilSequence {
		return nil, errors.New("invalid yjs compaction input")
	}

	payload := bytes.NewBuffer(make([]byte, 0, 28+len(i.Snapshot)+len(i.StateVector)))
	for _, value := range []int64{i.BaseCompactedUntilSequence, i.CutoffSequence} {
		if err := binary.Write(payload, binary.BigEndian, value); err != nil {
			return nil, err
		}
	}
	for _, value := range []uint32{uint32(len(i.Snapshot)), uint32(len(i.StateVector)), uint32(len(i.Updates))} {
		if err := binary.Write(payload, binary.BigEndian, value); err != nil {
			return nil, err
		}
	}
	payload.Write(i.Snapshot)
	payload.Write(i.StateVector)

	expectedSequence := i.BaseCompactedUntilSequence + 1
	for _, update := range i.Updates {
		if update.UpdateSequence != expectedSequence || len(update.Payload) > math.MaxUint32 {
			return nil, errors.New("invalid yjs compaction update")
		}
		if err := binary.Write(payload, binary.BigEndian, update.UpdateSequence); err != nil {
			return nil, err
		}
		if err := binary.Write(payload, binary.BigEndian, uint32(len(update.Payload))); err != nil {
			return nil, err
		}
		payload.Write(update.Payload)
		expectedSequence++
	}
	if expectedSequence-1 != i.CutoffSequence {
		return nil, errors.New("yjs compaction input has an incomplete update range")
	}

	return payload.Bytes(), nil
}

func (i *YjsCompactionInput) UnmarshalBytes(payload []byte) error {
	*i = YjsCompactionInput{}
	if len(payload) < 28 {
		return errors.New("invalid yjs compaction input payload")
	}

	i.BaseCompactedUntilSequence = int64(binary.BigEndian.Uint64(payload[0:8]))
	i.CutoffSequence = int64(binary.BigEndian.Uint64(payload[8:16]))
	snapshotLength := binary.BigEndian.Uint32(payload[16:20])
	stateVectorLength := binary.BigEndian.Uint32(payload[20:24])
	updateCount := binary.BigEndian.Uint32(payload[24:28])
	if i.BaseCompactedUntilSequence < 0 || i.CutoffSequence <= i.BaseCompactedUntilSequence {
		return errors.New("invalid yjs compaction input")
	}

	offset := 28
	if uint64(snapshotLength) > uint64(len(payload)-offset) {
		return errors.New("invalid yjs compaction input")
	}
	i.Snapshot = append(i.Snapshot, payload[offset:offset+int(snapshotLength)]...)
	offset += int(snapshotLength)
	if uint64(stateVectorLength) > uint64(len(payload)-offset) {
		return errors.New("invalid yjs compaction input")
	}
	i.StateVector = append(i.StateVector, payload[offset:offset+int(stateVectorLength)]...)
	offset += int(stateVectorLength)

	i.Updates = make([]YjsDocumentUpdate, 0, updateCount)
	expectedSequence := i.BaseCompactedUntilSequence + 1
	for index := uint32(0); index < updateCount; index++ {
		if len(payload)-offset < 12 {
			return errors.New("invalid yjs compaction input")
		}

		updateSequence := int64(binary.BigEndian.Uint64(payload[offset : offset+8]))
		updateLength := binary.BigEndian.Uint32(payload[offset+8 : offset+12])
		offset += 12
		if updateSequence != expectedSequence || uint64(updateLength) > uint64(len(payload)-offset) {
			return errors.New("invalid yjs compaction input")
		}

		i.Updates = append(i.Updates, YjsDocumentUpdate{
			UpdateSequence: updateSequence,
			Payload:        append([]byte{}, payload[offset:offset+int(updateLength)]...),
		})
		offset += int(updateLength)
		expectedSequence++
	}
	if offset != len(payload) || expectedSequence-1 != i.CutoffSequence {
		return errors.New("invalid yjs compaction input")
	}

	return nil
}

type YjsCompactionResult struct {
	Snapshot                   []byte
	StateVector                []byte
	BaseCompactedUntilSequence int64
	CutoffSequence             int64
}

func (r YjsCompactionResult) MarshalBytes() ([]byte, error) {
	if len(r.Snapshot) > math.MaxUint32 || len(r.StateVector) > math.MaxUint32 || r.BaseCompactedUntilSequence < 0 || r.CutoffSequence <= r.BaseCompactedUntilSequence {
		return nil, errors.New("invalid yjs compaction result")
	}

	payload := bytes.NewBuffer(make([]byte, 0, 24+len(r.Snapshot)+len(r.StateVector)))
	for _, value := range []int64{r.BaseCompactedUntilSequence, r.CutoffSequence} {
		if err := binary.Write(payload, binary.BigEndian, value); err != nil {
			return nil, err
		}
	}
	for _, value := range []uint32{uint32(len(r.Snapshot)), uint32(len(r.StateVector))} {
		if err := binary.Write(payload, binary.BigEndian, value); err != nil {
			return nil, err
		}
	}
	payload.Write(r.Snapshot)
	payload.Write(r.StateVector)

	return payload.Bytes(), nil
}

func (r *YjsCompactionResult) UnmarshalBytes(payload []byte) error {
	if len(payload) < 24 {
		return errors.New("invalid yjs compaction result payload")
	}

	r.BaseCompactedUntilSequence = int64(binary.BigEndian.Uint64(payload[0:8]))
	r.CutoffSequence = int64(binary.BigEndian.Uint64(payload[8:16]))
	snapshotLength := binary.BigEndian.Uint32(payload[16:20])
	stateVectorLength := binary.BigEndian.Uint32(payload[20:24])
	if r.BaseCompactedUntilSequence < 0 || r.CutoffSequence <= r.BaseCompactedUntilSequence || uint64(snapshotLength)+uint64(stateVectorLength) != uint64(len(payload)-24) {
		return errors.New("invalid yjs compaction result")
	}

	r.Snapshot = append(r.Snapshot[:0], payload[24:24+snapshotLength]...)
	r.StateVector = append(r.StateVector[:0], payload[24+snapshotLength:]...)

	return nil
}
