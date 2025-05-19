package wav

import (
	"encoding/binary"
	"fmt"
	"os"
)

func NewWavWriter() *WavWriter {
	return &WavWriter{
		File:       nil,
		endianness: binary.LittleEndian,
		current:    0,
		headerSize: 0,
	}
}

func (w *WavWriter) WriteHeader(file *os.File, header WavHeader) error {
	_, err := file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("file.Seek(0, 0): %w", err)
	}

	w.headerSize = 44

	err = binary.Write(file, w.endianness, header.ChunkID)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, header.ChunkID): %w", err)
	}

	err = binary.Write(file, w.endianness, uint32(w.headerSize))
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, uint32(w.headerSize)): %w", err)
	}

	err = binary.Write(file, w.endianness, header.Format)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, header.Format): %w", err)
	}

	w.current += 8

	return nil
}

func (w *WavWriter) WriteFmtChunk(file *os.File, fmtChunk FmtSubChunk) error {
	_, err := file.Seek(w.current, 0)
	if err != nil {
		return fmt.Errorf("file.Seek(w.current, 0): %w", err)
	}

	err = binary.Write(file, w.endianness, FourCC{'f', 'm', 't', ' '})
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, FourCC{'f', 'm', 't', ' '}): %w", err)
	}

	err = binary.Write(file, w.endianness, uint32(16))
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, uint32(16)): %w", err)
	}

	err = binary.Write(file, w.endianness, fmtChunk.AudioFormat)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, fmtChunk.AudioFormat): %w", err)
	}

	err = binary.Write(file, w.endianness, fmtChunk.NumChannels)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, fmtChunk.NumChannels): %w", err)
	}

	err = binary.Write(file, w.endianness, fmtChunk.SampleRate)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, fmtChunk.SampleRate): %w", err)
	}

	err = binary.Write(file, w.endianness, fmtChunk.ByteRate)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, fmtChunk.ByteRate): %w", err)
	}

	err = binary.Write(file, w.endianness, fmtChunk.BlockAlign)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, fmtChunk.BlockAlign): %w", err)
	}

	err = binary.Write(file, w.endianness, fmtChunk.BitsPerSample)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, fmtChunk.BitsPerSample): %w", err)
	}

	w.current += 24

	return nil
}

func (w *WavWriter) WriteDataChunk(file *os.File, dataChunk DataSubChunk) error {
	_, err := file.Seek(w.current, 0)
	if err != nil {
		return fmt.Errorf("file.Seek(w.current, 0): %w", err)
	}

	err = binary.Write(file, w.endianness, FourCC{'d', 'a', 't', 'a'})
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, FourCC{'d', 'a', 't', 'a'}): %w", err)
	}

	dataSize := uint32(len(dataChunk.Data))
	err = binary.Write(file, w.endianness, dataSize)
	if err != nil {
		return fmt.Errorf("binary.Write(file, w.endianness, uint32(len(dataChunk.Data))): %w", err)
	}

	_, err = file.Write(dataChunk.Data)
	if err != nil {
		return fmt.Errorf("file.Write(dataChunk.Data): %w", err)
	}

	if len(dataChunk.Data)%2 != 0 {
		_, err = file.Write([]byte{0x00})
		if err != nil {
			return fmt.Errorf("file.Write([]byte{0x00}): %w", err)
		}
	}

	w.current += int64(dataSize) + 8
	if dataSize%2 != 0 {
		w.current++
	}

	return nil
}

func (w *WavWriter) ConvertFromSamples(samples []float64) []byte {
	var data = make([]byte, len(samples)*2)

	for i, sample := range samples {
		binary.LittleEndian.PutUint16(data[i*2:], uint16(sample*32767))
	}

	return data
}

func WriteWavFile(filePath string, wavFile *WavFile) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("os.Create(%s): %w", filePath, err)
	}
	defer file.Close()

	writer := NewWavWriter()
	writer.File = file

	if len(wavFile.Samples) > 0 {
		wavFile.DataChunk.Data = writer.ConvertFromSamples(wavFile.Samples)
		wavFile.DataChunk.SubChunkSize = uint32(len(wavFile.DataChunk.Data))
	}

	if err := writer.WriteHeader(file, wavFile.Header); err != nil {
		return fmt.Errorf("writer.WriteHeader(): %w", err)
	}

	if err := writer.WriteFmtChunk(file, wavFile.FmtChunk); err != nil {
		return fmt.Errorf("writer.WriteFmtChunk(): %w", err)
	}

	if err := writer.WriteDataChunk(file, wavFile.DataChunk); err != nil {
		return fmt.Errorf("writer.WriteDataChunk(): %w", err)
	}

	return nil
}
