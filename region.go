package mapper

import "bytes"
import "compress/zlib"
import "encoding/binary"
import "fmt"
import "io"
import "os"
import "path"
import "strconv"
import "strings"
import "time"

type Region struct {
	X        int64
	Z        int64
	filename string
	file     *os.File
	header   [8192]byte
}

func (r *Region) Open(fn string) (err error) {
	r.filename = fn
	components := strings.Split(path.Base(fn), ".")
	r.X, err = strconv.ParseInt(components[1], 10, 64)
	r.Z, err = strconv.ParseInt(components[2], 10, 64)
	//fmt.Printf("Opening %s (x=%d, z=%d)\n", fn, r.X, r.Z)
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	r.file = file
	// read the header
	n, err := r.file.Read(r.header[:])
	if err != nil {
		return err
	}
	if n != 8192 {
		return fmt.Errorf("open: header size %d != 8192", n)
	}
	return nil
}

func header_offset(x int, z int) int {
	return 4 * ((x % 32) + (z%32)*32)

}

// first three bytes are a (big-endian) offset in 4KiB sectors from the start of the file
// remaining byte which gives the length of the chunk (also in 4KiB sectors, rounded up)
func (r *Region) chunk_location(x int, z int) (location int, length int) {
	offset := header_offset(x, z)
	// this is what "(bigEndian) Uint24(b []byte) uint32" would do if it existed
	location = int(uint32(r.header[offset+2]) | uint32(r.header[offset+1])<<8 | uint32(r.header[offset])<<16)
	length = int(uint32(r.header[offset+3]))
	return
}

// returns a Buffer which is the uncompressed chunk data
func (r *Region) ChunkData(x int, z int) (data bytes.Buffer, err error) {
	var chunk_header [5]byte
	location, _ := r.chunk_location(x, z)
	r.file.Seek(int64(location*4096), io.SeekStart)
	r.file.Read(chunk_header[:])
	//fmt.Printf("%q\n", chunk_header)
	length := binary.BigEndian.Uint32(chunk_header[0:4])
	compression_type := int(chunk_header[4])
	if compression_type != 2 {
		err = fmt.Errorf("chunk_data: unsupported compressed type %d", compression_type)
		return
	}
	//fmt.Printf("length of compressed data = %d\n", length)
	compressed := make([]byte, length)
	_, err = r.file.Read(compressed[:])
	if err != nil {
		return
	}
	//fmt.Printf("length compressed: %d\n", n)
	reader, _ := zlib.NewReader(bytes.NewBuffer(compressed[:]))
	io.Copy(&data, reader)
	//fmt.Printf("length uncompressed: %d\n", data.Len())
	return
}

func (r *Region) chunk_timestamp(x int, z int) time.Time {
	offset := header_offset(x, z) + 4096
	bytes := r.header[offset : offset+4]
	uts := binary.BigEndian.Uint32(bytes)
	return time.Unix(int64(uts), 0)
}
