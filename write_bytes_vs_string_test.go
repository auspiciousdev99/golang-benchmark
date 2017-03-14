package main

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const sampleStr = "It’s the job that’s never started as takes longest to finish"

var (
	sampleBytes                   = []byte(sampleStr)
	errInvalidByteSliceConversion = errors.New("invalid byte slice passed to conversion, slice have same the same length and capacity")
)

func BenchmarkWriteBytes(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	var w io.Writer
	w = buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.Write(sampleBytes)
	}
}

func BenchmarkWriteString(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	var w io.Writer
	w = buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.Write([]byte(sampleStr))
	}
}

func BenchmarkWriteUnafeString(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	var w io.Writer
	w = buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.Write(unsafeStrToByte(sampleStr))
	}
}

func unsafeStrToByte(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	var b []byte
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	byteHeader.Data = strHeader.Data

	l := len(s)
	byteHeader.Len = l
	byteHeader.Cap = l
	return b
}

func unsafeByteToStr(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	// need to assert that the slice's length and capacity are equal to avoid a memory leak
	// when converting to a string
	if len(b) != cap(b) {
		panic(errInvalidByteSliceConversion)
	}

	var s string
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	strHeader.Data = sliceHeader.Data
	strHeader.Len = len(b)
	return s
}

func TestUnsafeStrToByte(t *testing.T) {
	s := "fizzbuzz"
	expected := []byte(s)
	assert.Equal(t, expected, unsafeStrToByte(s))
}

func TestUnsafeByteToStr(t *testing.T) {
	b := []byte{'f', 'i', 'z', 'z', 'b', 'u', 'z', 'z'}
	expected := string(b)
	assert.Equal(t, expected, unsafeByteToStr(b))
}
