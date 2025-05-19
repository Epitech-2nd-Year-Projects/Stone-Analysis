package wav

import (
	"encoding/binary"
	"fmt"
	"os"
)

func ValidateWavFormat(wavFile *WavFile) error {
	if wavFile.FmtChunk.AudioFormat != 1 {
		return ErrUnsupportedAudioFormat
	}

	if wavFile.FmtChunk.NumChannels != 1 {
		return ErrInvalidNumChannels
	}

	if wavFile.FmtChunk.SampleRate != 48000 {
		return ErrInvalidSampleRate
	}

	if wavFile.FmtChunk.BitsPerSample != 16 {
		return ErrInvalidBitsPerSample
	}

	return nil
}

func readBytes(file *os.File, n int) ([]byte, error) {
	return []byte{}, nil
}

func readString(file *os.File, n int) (string, error) {
	return "", nil
}

func readUint16(file *os.File) (uint16, error) {
	return 0, nil
}

func readUint32(file *os.File) (uint32, error) {
	return 0, nil
}

func seekToChunk(file *os.File, chunkID string) error {
	_, err := file.Seek(12, 0)
	if err != nil {
		return fmt.Errorf("file.Seek(12, 0): %w", err)
	}

	var id FourCC
	var size uint32

	for {
		err = binary.Read(file, binary.LittleEndian, &id)
		if err != nil {
			return fmt.Errorf("binary.Read(file, binary.LittleEndian, &id): %w", err)
		}

		err = binary.Read(file, binary.LittleEndian, &size)
		if err != nil {
			return fmt.Errorf("binary.Read(file, binary.LittleEndian, &size): %w", err)
		}

		if id.Equals(chunkID) {
			_, err = file.Seek(-8, 1)
			if err != nil {
				return fmt.Errorf("file.Seek(-8, 1): %w", err)
			}
			return nil
		}

		skipSize := int64(size)
		if skipSize%2 != 0 {
			skipSize++
		}
		_, err = file.Seek(skipSize, 1)
		if err != nil {
			return fmt.Errorf("file.Seek(skipSize, 1): %w", err)
		}
	}
}
