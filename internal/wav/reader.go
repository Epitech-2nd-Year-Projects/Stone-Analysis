package wav

import (
	"encoding/binary"
	"fmt"
	"os"
)

func NewWavReader() *WavReader {
	return &WavReader{
		Current:    0,
		File:       nil,
		endianness: binary.LittleEndian,
	}
}

func (w *WavReader) ReadHeader() (WavHeader, error) {
	_, err := w.File.Seek(0, 0)
	if err != nil {
		return WavHeader{}, fmt.Errorf("w.File.Seek(0, 0): %w", err)
	}

	header := WavHeader{}

	err = binary.Read(w.File, w.endianness, &header.ChunkID)
	if err != nil {
		return WavHeader{}, fmt.Errorf("binary.Read(w.File, w.endianness, &header.ChunkID): %w", err)
	}

	if !header.ChunkID.Equals("RIFF") {
		return WavHeader{}, fmt.Errorf("header.ChunkID != 'RIFF'")
	}

	err = binary.Read(w.File, w.endianness, &header.ChunkSize)
	if err != nil {
		return WavHeader{}, fmt.Errorf("binary.Read(w.File, w.endianness, &header.ChunkSize): %w", err)
	}

	err = binary.Read(w.File, w.endianness, &header.Format)
	if err != nil {
		return WavHeader{}, fmt.Errorf("binary.Read(w.File, w.endianness, &header.Format): %w", err)
	}
	if !header.Format.Equals("WAVE") {
		return WavHeader{}, fmt.Errorf("header.Format != 'WAVE'")
	}

	w.Current += 12

	return header, nil
}

func (w *WavReader) ReadFmtChunk() (FmtSubChunk, error) {
	err := seekToChunk(w.File, "fmt ")
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("w.File.Seek(w.Current, 0): %w", err)
	}

	fmtChunk := FmtSubChunk{}

	err = binary.Read(w.File, w.endianness, &fmtChunk.SubChunkID)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.SubChunkID): %w", err)
	}

	if !fmtChunk.SubChunkID.Equals("fmt ") {
		return FmtSubChunk{}, fmt.Errorf("fmtChunk.SubChunkID != 'fmt '")
	}

	err = binary.Read(w.File, w.endianness, &fmtChunk.SubChunkSize)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.SubChunkSize): %w", err)
	}

	err = binary.Read(w.File, w.endianness, &fmtChunk.AudioFormat)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.AudioFormat): %w", err)
	}

	err = binary.Read(w.File, w.endianness, &fmtChunk.NumChannels)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.NumChannels): %w", err)
	}

	err = binary.Read(w.File, w.endianness, &fmtChunk.SampleRate)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.SampleRate): %w", err)
	}

	err = binary.Read(w.File, w.endianness, &fmtChunk.ByteRate)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.ByteRate): %w", err)
	}

	err = binary.Read(w.File, w.endianness, &fmtChunk.BlockAlign)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.BlockAlign): %w", err)
	}

	err = binary.Read(w.File, w.endianness, &fmtChunk.BitsPerSample)
	if err != nil {
		return FmtSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &fmtChunk.BitsPerSample): %w", err)
	}

	w.Current += 16

	return fmtChunk, nil
}

func (w *WavReader) ReadDataChunk() (DataSubChunk, error) {
	err := seekToChunk(w.File, "data")
	if err != nil {
		return DataSubChunk{}, fmt.Errorf("w.File.Seek(w.Current, 0): %w", err)
	}

	dataChunk := DataSubChunk{}

	err = binary.Read(w.File, w.endianness, &dataChunk.SubChunkID)
	if err != nil {
		return DataSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &dataChunk.SubChunkID): %w", err)
	}

	if !dataChunk.SubChunkID.Equals("data") {
		return DataSubChunk{}, fmt.Errorf("dataChunk.SubChunkID != 'data'")
	}

	err = binary.Read(w.File, w.endianness, &dataChunk.SubChunkSize)
	if err != nil {
		return DataSubChunk{}, fmt.Errorf("binary.Read(w.File, w.endianness, &dataChunk.SubChunkSize): %w", err)
	}

	dataChunk.Data = make([]byte, dataChunk.SubChunkSize)
	_, err = w.File.Read(dataChunk.Data)
	if err != nil {
		return DataSubChunk{}, fmt.Errorf("w.File.Read(dataChunk.Data): %w", err)
	}

	w.Current += 8 + int64(dataChunk.SubChunkSize)

	return dataChunk, nil
}

func (w *WavReader) ConvertToSamples(data []byte) []float64 {
	var samples = make([]float64, len(data)/2)

	for i := 0; i < len(samples); i++ {
		var sample = make([]byte, 2)
		copy(sample, data[i*2:(i+1)*2])
		samples[i] = float64(binary.LittleEndian.Uint16(sample)) / float64(1<<15)
	}

	return samples
}

func ReadWavFile(filePath string) (*WavFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("os.Open(%s): %w", filePath, err)
	}
	defer file.Close()

	reader := NewWavReader()
	reader.File = file

	wavFile := &WavFile{}

	wavFile.Header, err = reader.ReadHeader()
	if err != nil {
		return nil, fmt.Errorf("reader.ReadHeader(): %w", err)
	}

	wavFile.FmtChunk, err = reader.ReadFmtChunk()
	if err != nil {
		return nil, fmt.Errorf("reader.ReadFmtChunk(): %w", err)
	}

	wavFile.DataChunk, err = reader.ReadDataChunk()
	if err != nil {
		return nil, fmt.Errorf("reader.ReadDataChunk(): %w", err)
	}

	if err := ValidateWavFormat(wavFile); err != nil {
		return nil, fmt.Errorf("ValidateWavFormat(): %w", err)
	}

	wavFile.Samples = reader.ConvertToSamples(wavFile.DataChunk.Data)

	return wavFile, nil
}
