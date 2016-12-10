package nbt

import "encoding/binary"
import "fmt"
import "io"
import "math"

const (
	tag_end        = 0
	tag_byte       = 1
	tag_short      = 2
	tag_int        = 3
	tag_long       = 4
	tag_float      = 5
	tag_double     = 6
	tag_byte_array = 7
	tag_string     = 8
	tag_list       = 9
	tag_compound   = 10
	tag_int_array  = 11
)

// simple tags are represented by structs which implement tag_value.
// you should check the type and coerce to the right thing before reading the
// value property
type simple_tag interface {
	isTag()
}

type end_t struct{}

func (t end_t) isTag() {}

type byte_t struct {
	name  string
	value int8
}

func (t byte_t) isTag() {}

func (t byte_t) Get() int8 {
	return t.value
}

type short_t struct {
	name  string
	value int16
}

func (t short_t) isTag() {}

type int_t struct {
	name  string
	value int32
}

func (t int_t) isTag() {}

type long_t struct {
	name  string
	value int64
}

func (t long_t) isTag() {}

type float_t struct {
	name  string
	value float32
}

func (t float_t) isTag() {}

type double_t struct {
	name  string
	value float64
}

func (t double_t) isTag() {}

type string_t struct {
	name  string
	value string
}

func (t string_t) isTag() {}

// complex types are compound and the various array types.  they implement
// simple_tag but also have a Size() and At(i) functions for accessing their
// contents (which could themselves be simple or complex types)
type complex_tag interface {
	Size() int
	At(int) simple_tag
}

type byte_array_t struct {
	name  string
	value []byte
}

func (t byte_array_t) isTag() {}

type list_t struct {
	name     string
	tag_type int
	value    []simple_tag
}

func (t list_t) isTag() {}

type compound_t struct {
	name  string
	value []simple_tag
}

func (t compound_t) isTag() {}

type int_array_t struct {
	name  string
	value []simple_tag
}

func (t int_array_t) isTag() {}

// the NBT parser
// expects a value of a type implementing io.Reader.  returns the root tag of
// the data, which contains all other tags
func Parse(data io.Reader) simple_tag {
	name := ""
	tag := read_tag(data)
	if tag != tag_end {
		name = read_string(data)
	}
	switch tag {
	case tag_end:
		return end_t{}
	case tag_byte:
		return byte_t{name, read_int8(data)}
	case tag_short:
		return short_t{name, read_int16(data)}
	case tag_int:
		return int_t{name, read_int32(data)}
	case tag_long:
		return long_t{name, read_int64(data)}
	case tag_float:
		return float_t{name, read_float32(data)}
	case tag_double:
		return double_t{name, read_float64(data)}
	case tag_byte_array:
		size := int(read_int32(data))
		return byte_array_t{name, read_bytes(data, size)}
	case tag_string:
		return string_t{name, read_string(data)}
	case tag_list:
		list_type := read_tag(data)
		size := int(read_int32(data))
		return list_t{name, list_type, read_list_values(data, list_type, size)}
	case tag_compound:
		return compound_t{name, read_compound_values(data)}
	case tag_int_array:
		size := int(read_int32(data))
		return int_array_t{name, read_list_values(data, tag_int, size)}
	default:
		panic(fmt.Sprintf("Parse: unknown tag %d", tag))
	}
}

func read_compound_values(data io.Reader) []simple_tag {
	var values []simple_tag
	var tag simple_tag
	for {
		tag = Parse(data)
		if _, end := tag.(end_t); end {
			break
		}
		values = append(values, tag)
	}
	return values
}

func read_list_values(data io.Reader, list_type int, size int) []simple_tag {
	var values []simple_tag
	for i := 0; i < size; i++ {
		switch list_type {
		case tag_int:
			values = append(values, int_t{"", read_int32(data)})
		case tag_float:
			values = append(values, float_t{"", read_float32(data)})
		case tag_double:
			values = append(values, double_t{"", read_float64(data)})
		case tag_compound:
			values = append(values, compound_t{"", read_compound_values(data)})
		default:
			panic(fmt.Sprintf("read_list_values: unhandled list type %d", list_type))
		}
	}
	return values
}

func read_tag(data io.Reader) int {
	tag := int(read_int8(data))
	if tag < 0 || tag > tag_int_array {
		panic(fmt.Sprintf("read_tag: unknown tag %d", tag))
	}
	return tag
}

func read_string(data io.Reader) string {
	size := read_int16(data)
	if size > 0 {
		buf := make([]byte, size)
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

func read_bytes(data io.Reader, size int) []byte {
	buf := make([]byte, size)
	data.Read(buf)
	return buf
}
