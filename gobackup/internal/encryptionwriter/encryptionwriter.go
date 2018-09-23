package encryptionwriter

import (
	"bytes"
	"io"
	"log"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

const compressionLevel = packet.BestCompression

//EncryptionWriter implements a writer that can be used to encrypt data
type EncryptionWriter struct {
	output io.WriteCloser
}

func New(destination io.Writer, secret string) *EncryptionWriter {
	hints := &openpgp.FileHints{IsBinary: true}
	compressionConfig := &packet.CompressionConfig{compressionLevel}
	config := &packet.Config{DefaultCompressionAlgo: packet.CompressionZIP, CompressionConfig: compressionConfig}
	key := []byte(secret)
	writerCloser, err := openpgp.SymmetricallyEncrypt(destination, key, hints, config)
	if err != nil {
		log.Fatalf("Error creating encwriter: %v", err)
	}
	return &EncryptionWriter{writerCloser}
}

func (encWriter *EncryptionWriter) Write(p []byte) (n int, err error) {
	reader := bytes.NewReader(p)
	bytesWritten, err := io.Copy(encWriter.output, reader)
	return int(bytesWritten), err

}

func (encWriter *EncryptionWriter) Close() error {
	return encWriter.output.Close()
}

func (encWriter *EncryptionWriter) Flush() {
	//nop
}
