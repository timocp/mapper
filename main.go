package main

import "bytes"
import "compress/zlib"
import "encoding/binary"
import "fmt"
import "io"
import "log"
import "math"
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

func main() {
	//fmt.Println("mapper")
	r := new(Region)
	must(r.open("/home/minecraft/server/world/region/r.0.0.mca"))
	for x := 0; x < 1; x++ {
		for z := 0; z < 1; z++ {
			//offset, length := r.chunk_location(x, z)
			//fmt.Printf("chunk(%2d, %2d) timestamp=%s offset=%d length=%d\n", x, z, r.chunk_timestamp(x, z), offset, length)
			chunk_data, err := r.chunk_data(x, z)
			must(err)
			//fmt.Printf("%s", chunk_data)
			parse_nbt(bytes.NewReader(chunk_data.Bytes()), 0)
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const (
	tag_end        = iota
	tag_byte       = iota
	tag_short      = iota
	tag_int        = iota
	tag_long       = iota
	tag_float      = iota
	tag_double     = iota
	tag_byte_array = iota
	tag_string     = iota
	tag_list       = iota
	tag_compound   = iota
	tag_int_array  = iota
)

// prototype; just parse and print out what we find
func parse_nbt(data io.Reader, depth int) {
	name := ""
	for {
		indent(depth)
		tag := read_tag(data)
		if tag != tag_end {
			name = read_name(data)
		}
		switch tag {
		case tag_end:
			// this is only used to mark the end of a compound tag; it has no name
			fmt.Printf("TAG_End\n")
			return
		case tag_byte:
			fmt.Printf("TAG_Byte (%s) = %d\n", name, read_int8(data))
		case tag_short:
			fmt.Printf("TAG_Short (%s) = %d\n", name, read_int16(data))
		case tag_int:
			fmt.Printf("TAG_Int (%s) = %d\n", name, read_int32(data))
		case tag_long:
			fmt.Printf("TAG_Long (%s) = %d\n", name, read_int64(data))
		case tag_float:
			fmt.Printf("TAG_Float (%s) = %f\n", name, read_float32(data))
		case tag_double:
			fmt.Printf("TAG_Double(%s) = %f\n", name, read_float64(data))
		case tag_byte_array:
			length := int(read_int32(data))
			bytes := read_bytes(data, length)
			fmt.Printf("TAG_Byte_Array (%s) = %q\n", name, bytes)
		case tag_string:
			fmt.Printf("TAG_String (%s) = %q\n", name, read_name(data))
		case tag_list:
			list_of := read_tag(data)
			size := int(read_int32(data))
			fmt.Printf("TAG_List (%s type=%d size=%d)\n", name, list_of, size)
			parse_list(data, list_of, size, depth+1)
		case tag_compound:
			fmt.Printf("TAG_Compound (%s)\n", name)
			parse_nbt(data, depth+1)
		case tag_int_array:
			size := int(read_int32(data))
			fmt.Printf("TAG_Int_Array (%s size=%d)\n", name, size)
			parse_list(data, tag_int, size, depth+1)
		default:
			panic(fmt.Sprintf("parse_nbt: unknown tag %d", tag))
		}
	}
}

func parse_list(data io.Reader, tag int, size int, depth int) {
	for i := 0; i < size; i++ {
		indent(depth)
		switch tag {
		case tag_int:
			fmt.Printf("Int = %d\n", read_int32(data))
		case tag_float:
			fmt.Printf("Float = %f\n", read_float32(data))
		case tag_double:
			fmt.Printf("Double = %f\n", read_float64(data))
		case tag_compound:
			fmt.Printf("Compound\n")
			parse_nbt(data, depth+1)
		default:
			panic(fmt.Sprintf("parse_list: unhandled tag %d", tag))
		}
	}
}

func indent(i int) {
	for ; i > 0; i-- {
		fmt.Printf(" ")
	}
}

func read_tag(data io.Reader) int {
	var buf [1]byte
	data.Read(buf[:])
	tag := int(buf[0])
	if tag < 0 || tag > tag_int_array {
		panic(fmt.Sprintf("parse_nbt: unknown tag %d", tag))
	}
	return tag
}

func read_name(data io.Reader) string {
	length := read_int16(data)
	if length > 0 {
		buf := make([]byte, length)
		data.Read(buf[:])
		return string(buf)
	}
	return ""
}

func read_int8(data io.Reader) int8 {
	var buf [1]byte
	data.Read(buf[:])
	return int8(buf[0])
}

func read_int16(data io.Reader) int16 {
	var buf [2]byte
	data.Read(buf[:])
	return int16(binary.BigEndian.Uint16(buf[:]))
}

func read_int32(data io.Reader) int32 {
	var buf [4]byte
	data.Read(buf[:])
	return int32(binary.BigEndian.Uint32(buf[:]))
}

func read_int64(data io.Reader) int64 {
	var buf [8]byte
	data.Read(buf[:])
	return int64(binary.BigEndian.Uint64(buf[:]))
}

func read_float32(data io.Reader) float32 {
	var buf [4]byte
	data.Read(buf[:])
	return math.Float32frombits(binary.BigEndian.Uint32(buf[:]))
}

func read_float64(data io.Reader) float64 {
	var buf [8]byte
	data.Read(buf[:])
	return math.Float64frombits(binary.BigEndian.Uint64(buf[:]))
}

func read_bytes(data io.Reader, length int) []byte {
	buf := make([]byte, length)
	data.Read(buf)
	return buf
}

func (r *Region) open(fn string) (err error) {
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
func (r *Region) chunk_data(x int, z int) (data bytes.Buffer, err error) {
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
