package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type EncryptedWriter struct {
	w      io.Writer
	stream cipher.Stream
}

func NewEncryptedWriter(w io.Writer, key []byte) (*EncryptedWriter, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)

	return &EncryptedWriter{
		w:      w,
		stream: stream,
	}, nil
}

func (ew *EncryptedWriter) Write(p []byte) (n int, err error) {
	encrypted := make([]byte, len(p))

	ew.stream.XORKeyStream(encrypted, p)

	return ew.w.Write(encrypted)
}

type DecryptedReader struct {
	r      io.Reader
	stream cipher.Stream
}

func NewDecryptedReader(r io.Reader, key []byte) (*DecryptedReader, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(r, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)

	return &DecryptedReader{
		r:      r,
		stream: stream,
	}, nil
}

func (dr *DecryptedReader) Read(p []byte) (n int, err error) {
	n, err = dr.r.Read(p)

	if n > 0 {
		dr.stream.XORKeyStream(p[:n], p[:n])
	}

	return n, err
}
