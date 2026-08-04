package types

type ByteType int64

const (
	Byte ByteType = 1
	KB   ByteType = 1024 * Byte
	MB   ByteType = 1024 * KB
	GB   ByteType = 1024 * MB
	TB   ByteType = 1024 * GB
)

func (by ByteType) ToInt64() int64 {
	return int64(by)
}
