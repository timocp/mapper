package mapper

import "errors"
import "fmt"

import "github.com/timocp/nbt"

var EmptySectionError = errors.New("section is empty")
var airBlock = NewBlock(0, 0)

type Chunk struct {
	root    nbt.Tag
	regionX int
	regionZ int
	chunkX  int
	chunkZ  int
}

func NewChunk(tag nbt.Tag, region *Region, x int, z int) *Chunk {
	return &Chunk{tag, region.X, region.Z, x, z}
}

func (c *Chunk) Root() nbt.CompoundTag {
	return c.root.(nbt.CompoundTag)
}

func (c *Chunk) Level() nbt.CompoundTag {
	return c.Root().ChildByName("Level").(nbt.CompoundTag)
}

func (c *Chunk) Biomes() []byte {
	return c.Level().ChildByName("Biomes").(nbt.ByteArrayTag).Values
}

// return the index of the highest section in this chunk - above this is only
// air
func (c *Chunk) MaxSection() (max int) {
	for _, s := range c.Level().ChildByName("Sections").(nbt.ListTag).Values {
		i := int(s.(nbt.CompoundTag).ChildByName("Y").(nbt.ByteTag).Value)
		if i > max {
			max = i
		}
	}
	return
}

// returns the section with index Y=y.  If it doesn't exist, returns an empty
// section
func (c *Chunk) Section(y int) (result nbt.CompoundTag, err error) {
	for _, s := range c.Level().ChildByName("Sections").(nbt.ListTag).Values {
		if int(s.(nbt.CompoundTag).ChildByName("Y").(nbt.ByteTag).Value) == y {
			result = s.(nbt.CompoundTag)
			return
		}
	}
	err = EmptySectionError
	return
}

// accessor methods to values inside chunks
func (c *Chunk) HeightMap() []int32 {
	return c.Level().ChildByName("HeightMap").(nbt.IntArrayTag).Values
}

// returns the Block found at coords x y z
func (c *Chunk) BlockAt(x int, y int, z int) Block {
	// get the section which stores this block
	section, err := c.Section(y / 16)
	if err == EmptySectionError {
		return airBlock
	} else if err != nil {
		panic(fmt.Sprintf("BlockAt: %s", err.Error()))
	}
	blockPos := (y%16)*16*16 + z*16 + x
	//fmt.Printf("x=%d y=%d z=%d blockPos=%d\n", x, y, z, blockPos)
	idA := section.ChildByName("Blocks").(nbt.ByteArrayTag).Values[blockPos]
	add := section.ChildByName("Add")
	idB := byte(0)
	if add != nil {
		idB = nibble4(add.(nbt.ByteArrayTag).Values, blockPos)
	}
	data := nibble4(section.ChildByName("Data").(nbt.ByteArrayTag).Values, blockPos)
	//fmt.Printf("newblock idA=%d idB=%d data=%d\n", idA, idB, data)
	return NewBlock(int16(idA)+int16(idB)<<8, int8(data))
}

// walk the NBT data structure and print it out
func (c *Chunk) Debug() {
	nbt.Debug(c.root, 0)
}

func (c *Chunk) X() int {
	return c.regionX*32 + c.chunkX
}

func (c *Chunk) Z() int {
	return c.regionZ*32 + c.chunkZ
}

// from examples on http://minecraft.gamepedia.com/Chunk_format#Block_format
//
// from a byte array where each byte is a 4 bit number, extract the item at
// the specified index
func nibble4(arr []byte, index int) byte {
	if index%2 == 0 {
		return arr[index/2] & 0x0F
	} else {
		return (arr[index/2] >> 4) & 0x0F
	}
}
