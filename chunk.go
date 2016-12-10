package mapper

import "github.com/timocp/mapper/nbt"

type Chunk struct {
	Root nbt.Tag
}

// walk the NBT data structure and print it out
func (c *Chunk) Debug() {
	//nbt.Debug(c.Root, 0)
}
