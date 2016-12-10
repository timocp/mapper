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
type Tag interface {
	GetName() string
}

type EndTag struct{}

func (t EndTag) GetName() string { return "" }

type ByteTag struct {
	name  string
	Value int8
}

func (t ByteTag) GetName() string { return t.name }

type ShortTag struct {
	name  string
	Value int16
}

func (t ShortTag) GetName() string { return t.name }

type IntTag struct {
	name  string
	Value int32
}

func (t IntTag) GetName() string { return t.name }

type LongTag struct {
	name  string
	Value int64
}

func (t LongTag) GetName() string { return t.name }

type FloatTag struct {
	name  string
	Value float32
}

func (t FloatTag) GetName() string { return t.name }

type DoubleTag struct {
	name  string
	Value float64
}

func (t DoubleTag) GetName() string { return t.name }

type StringTag struct {
	name  string
	Value string
}

func (t StringTag) GetName() string { return t.name }

type ByteArrayTag struct {
	name   string
	Values []byte
}

func (t ByteArrayTag) GetName() string { return t.name }

type ListTag struct {
	name    string
	TagType int
	Values  []Tag
}

func (t ListTag) GetName() string { return t.name }

type CompoundTag struct {
	name   string
	Values []Tag
}

func (t CompoundTag) GetName() string { return t.name }

func (t CompoundTag) ChildByName(name string) Tag {
	for _, v := range t.Values {
		if v.GetName() == name {
			return v
		}
	}
	return nil
}

type IntArrayTag struct {
	name   string
	Values []int32
}

func (t IntArrayTag) GetName() string { return t.name }

// the NBT parser
// expects a value of a type implementing io.Reader.  returns the root tag of
// the data, which contains all other tags
func Parse(data io.Reader) Tag {
	name := ""
	tag := read_tag(data)
	if tag != tag_end {
		name = read_string(data)
	}
	switch tag {
	case tag_end:
		return EndTag{}
	case tag_byte:
		return ByteTag{name, read_int8(data)}
	case tag_short:
		return ShortTag{name, read_int16(data)}
	case tag_int:
		return IntTag{name, read_int32(data)}
	case tag_long:
		return LongTag{name, read_int64(data)}
	case tag_float:
		return FloatTag{name, read_float32(data)}
	case tag_double:
		return DoubleTag{name, read_float64(data)}
	case tag_byte_array:
		size := int(read_int32(data))
		return ByteArrayTag{name, read_bytes(data, size)}
	case tag_string:
		return StringTag{name, read_string(data)}
	case tag_list:
		list_type := read_tag(data)
		size := int(read_int32(data))
		return ListTag{name, list_type, read_list_values(data, list_type, size)}
	case tag_compound:
		return CompoundTag{name, read_compound_values(data)}
	case tag_int_array:
		size := int(read_int32(data))
		return IntArrayTag{name, read_int_array_values(data, size)}
	default:
		panic(fmt.Sprintf("Parse: unknown tag %d", tag))
	}
}

func read_compound_values(data io.Reader) []Tag {
	var values []Tag
	var tag Tag
	for {
		tag = Parse(data)
		if _, end := tag.(EndTag); end {
			break
		}
		values = append(values, tag)
	}
	return values
}

func read_list_values(data io.Reader, list_type int, size int) []Tag {
	var values []Tag
	for i := 0; i < size; i++ {
		switch list_type {
		//case tag_int:
		//	values = append(values, IntTag{"", read_int32(data)})
		case tag_float:
			values = append(values, FloatTag{"", read_float32(data)})
		case tag_double:
			values = append(values, DoubleTag{"", read_float64(data)})
		case tag_compound:
			values = append(values, CompoundTag{"", read_compound_values(data)})
		default:
			panic(fmt.Sprintf("read_list_values: unhandled list type %d", list_type))
		}
	}
	return values
}

func read_int_array_values(data io.Reader, size int) []int32 {
	var values []int32
	for i := 0; i < size; i++ {
		values = append(values, read_int32(data))
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

func Debug(tag Tag, depth int) {
	indent(depth)
	switch t := tag.(type) {
	case ByteTag:
		fmt.Printf("Byte(%s): %d\n", t.name, t.Value)
	case ShortTag:
		fmt.Printf("Short(%s): %d\n", t.name, t.Value)
	case IntTag:
		fmt.Printf("Int(%s): %d\n", t.name, t.Value)
	case LongTag:
		fmt.Printf("Long(%s): %d\n", t.name, t.Value)
	case FloatTag:
		fmt.Printf("Float(%s): %f\n", t.name, t.Value)
	case DoubleTag:
		fmt.Printf("Double(%s): %f\n", t.name, t.Value)
	case ByteArrayTag:
		fmt.Printf("ByteArray(%s): [", t.name)
		for i, b := range t.Values {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%d", b)
		}
		fmt.Printf("]\n")
	case StringTag:
		fmt.Printf("String(%s): %s\n", t.name, t.Value)
	case ListTag:
		fmt.Printf("List(%s) (type %d): [\n", t.name, t.TagType)
		for _, v := range t.Values {
			Debug(v, depth+1)
		}
		indent(depth)
		fmt.Printf("]\n")
	case CompoundTag:
		fmt.Printf("Compound(%s):\n", t.name)
		for _, v := range t.Values {
			Debug(v, depth+1)
		}
	case IntArrayTag:
		fmt.Printf("IntArray(%s): [", t.name)
		for i, v := range t.Values {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%d", v)
		}
		fmt.Printf("]\n")
	default:
		panic(fmt.Sprintf("printNbt: unhandled type %T", tag))
	}
}

func indent(i int) {
	for ; i > 0; i-- {
		fmt.Printf(" ")
	}
}
