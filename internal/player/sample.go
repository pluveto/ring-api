package player

import (
	"encoding/binary"
	"io"
	"math"
)

func SineSample(f io.Writer) (*int, error) {

	sampleRate := 44100
	duration := 1    // seconds
	frequency := 440 // Hz, A4 note

	numSamples := sampleRate * duration

	// WAV file header (same as before, but adjusted for data size)
	header := []byte{
		0x52, 0x49, 0x46, 0x46, // ChunkID: "RIFF"
		36, 0, 0, 0, // ChunkSize: Will be updated later
		0x57, 0x41, 0x56, 0x45, // Format: "WAVE"
		0x66, 0x6d, 0x74, 0x20, // Subchunk1ID: "fmt "
		16, 0, 0, 0, // Subchunk1Size: 16 for PCM
		1, 0, // AudioFormat: 1 (PCM)
		1, 0, // NumChannels: 1 (Mono)
		0x44, 0xac, 0, 0, // SampleRate: 44100
		0x88, 0x58, 0x01, 0, // ByteRate: SampleRate * NumChannels * BitsPerSample/8
		2, 0, // BlockAlign: NumChannels * BitsPerSample/8
		16, 0, // BitsPerSample: 16
		0x64, 0x61, 0x74, 0x61, // Subchunk2ID: "data"
		0, 0, 0, 0, // Subchunk2Size: Will be updated later
	}

	// Calculate data size and update header
	dataSize := numSamples * 2                                      // 2 bytes per sample (16-bit)
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+dataSize)) // Update ChunkSize
	binary.LittleEndian.PutUint32(header[40:44], uint32(dataSize))  // Update Subchunk2Size

	if _, err := f.Write(header); err != nil {
		return nil, err
	}

	// Generate sine wave data
	for i := 0; i < numSamples; i++ {
		// Generate a sine wave
		sample := math.Sin(2 * math.Pi * float64(frequency) * float64(i) / float64(sampleRate))

		// Convert to 16-bit signed integer (range: -32768 to 32767) and reduce volume to 1/4
		intValue := int16(sample * 32767 * 0.25)

		// Write the sample to the file
		err := binary.Write(f, binary.LittleEndian, &intValue)
		if err != nil {
			return nil, err
		}
	}

	return &numSamples, nil
}

// SawtoothSample function generates a sawtooth wave audio sample and writes it to the provided io.Writer.
func SawtoothSample(f io.Writer) (*int, error) {
	sampleRate := 44100
	duration := 1    // seconds
	frequency := 440 // Hz, A4 note

	numSamples := sampleRate * duration

	// WAV file header (same as before, but adjusted for data size)
	header := []byte{
		0x52, 0x49, 0x46, 0x46, // ChunkID: "RIFF"
		36, 0, 0, 0, // ChunkSize: Will be updated later
		0x57, 0x41, 0x56, 0x45, // Format: "WAVE"
		0x66, 0x6d, 0x74, 0x20, // Subchunk1ID: "fmt "
		16, 0, 0, 0, // Subchunk1Size: 16 for PCM
		1, 0, // AudioFormat: 1 (PCM)
		1, 0, // NumChannels: 1 (Mono)
		0x44, 0xac, 0, 0, // SampleRate: 44100
		0x88, 0x58, 0x01, 0, // ByteRate: SampleRate * NumChannels * BitsPerSample/8
		2, 0, // BlockAlign: NumChannels * BitsPerSample/8
		16, 0, // BitsPerSample: 16
		0x64, 0x61, 0x74, 0x61, // Subchunk2ID: "data"
		0, 0, 0, 0, // Subchunk2Size: Will be updated later
	}

	// Calculate data size and update header
	dataSize := numSamples * 2                                      // 2 bytes per sample (16-bit)
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+dataSize)) // Update ChunkSize
	binary.LittleEndian.PutUint32(header[40:44], uint32(dataSize))  // Update Subchunk2Size

	if _, err := f.Write(header); err != nil {
		return nil, err
	}

	// Generate sawtooth wave data
	for i := 0; i < numSamples; i++ {
		// Generate a sawtooth wave
		sample := math.Mod(float64(i)*float64(frequency)/float64(sampleRate), 1.0)*2 - 1

		// Convert to 16-bit signed integer (range: -32768 to 32767) and reduce volume to 1/4
		intValue := int16(sample * 32767 * 0.25)

		// Write the sample to the file
		err := binary.Write(f, binary.LittleEndian, &intValue)
		if err != nil {
			return nil, err
		}
	}

	return &numSamples, nil
}
