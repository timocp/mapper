package mapper

import "github.com/timocp/mapper/nbt"

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

// accessor methods to values inside chunks
func (c *Chunk) HeightMap() []int32 {
	return c.Level().ChildByName("HeightMap").(nbt.IntArrayTag).Values
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
