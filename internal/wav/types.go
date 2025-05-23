package wav

import (
	"encoding/binary"
	"errors"
	"os"
)

var (
	ErrUnsupportedAudioFormat = errors.New("unsupported audio format")
	ErrInvalidNumChannels     = errors.New("invalid number of channels")
	ErrInvalidSampleRate      = errors.New("invalid sample rate")
	ErrInvalidBitsPerSample   = errors.New("invalid bits per sample")
)

type FourCC [4]byte

func (f FourCC) String() string {
	return string(f[:])
}

func (f FourCC) Equals(s string) bool {
	return f.String() == s
}

type WavHeader struct {
	ChunkID   FourCC
	ChunkSize uint32
	Format    FourCC
}

type FmtSubChunk struct {
	SubChunkID    FourCC
	SubChunkSize  uint32
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

type DataSubChunk struct {
	SubChunkID   FourCC
	SubChunkSize uint32
	Data         []byte
}

type WavFile struct {
	Header    WavHeader
	FmtChunk  FmtSubChunk
	DataChunk DataSubChunk
	Samples   []float64
}

type WavReader struct {
	Current    int64
	File       *os.File
	endianness binary.ByteOrder
}

type WavWriter struct {
	File       *os.File
	endianness binary.ByteOrder
	current    int64
	headerSize int64
}
