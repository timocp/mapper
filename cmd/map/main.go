package main

import "bytes"
import "fmt"
import "image"
import "image/color"
import "image/png"
import "log"
import "os"

import "github.com/timocp/mapper"
import "github.com/timocp/mapper/nbt"

type hm struct {
	x      int
	z      int
	values []int32
}

func main() {
	//fmt.Println("mapper")
	var heightMaps []hm
	var minx, maxx, minz, maxz int
	for _, fn := range os.Args[1:] {
		r := new(mapper.Region)
		fmt.Printf("Reading %s\n", fn)
		must(r.Open(fn))
		for x := 0; x < 32; x++ {
			for z := 0; z < 32; z++ {
				//offset, length := r.chunk_location(x, z)
				//fmt.Printf("chunk(%2d, %2d) timestamp=%s offset=%d length=%d\n", x, z, r.chunk_timestamp(x, z), offset, length)
				//fmt.Printf("Loading chunk %d.%d ... ", x, z)
				chunk_data, err := r.ChunkData(x, z)
				must(err)
				if chunk_data.Len() == 0 {
					//fmt.Printf("not in file!\n")
				} else {
					//fmt.Printf("parsing...\n")
					chunk := mapper.NewChunk(nbt.Parse(bytes.NewReader(chunk_data.Bytes())), r, x, z)
					//chunk := mapper.Chunk{nbt.Parse(bytes.NewReader(chunk_data.Bytes())), r.X, r.Z, x, z}
					if chunk.X() < minx {
						minx = chunk.X()
					}
					if chunk.X() > maxx {
						maxx = chunk.X()
					}
					if chunk.Z() < minz {
						minz = chunk.Z()
					}
					if chunk.Z() > maxz {
						maxz = chunk.Z()
					}
					//chunk.Debug()
					// maybe better than this would be an
					// array of images, then compose them
					// at the end
					heightMaps = append(heightMaps, hm{chunk.X(), chunk.Z(), chunk.HeightMap()})
					//os.Exit(0)
					// making an array of all chunks (~76k)
					// takes too much memory.
				}
			}
		}
	}
	fmt.Printf("found %d heightMaps...\n", len(heightMaps))
	fmt.Printf("x: %d..%d\n", minx, maxx)
	fmt.Printf("z: %d..%d\n", minz, maxz)
	chunkWidth := maxx - minx + 1
	chunkHeight := maxz - minz + 1
	img := image.NewRGBA(image.Rect(0, 0, chunkWidth*16, chunkHeight*16))
	for _, hm := range heightMaps {
		// hm.x and hm.z is the xz coord of the chunk, which can be
		// used to determine the offset into the image to use
		for i, v := range hm.values {
			//if i == 10 && hm.x == 10 && hm.z == 10 {
			//	fmt.Printf("hm.x=%d, hm.z=%d, i=%d, x=%d, z=%d at=(%d,%d) v=%d\n", hm.x, hm.z, i, i%16, i/16, (hm.x-minx)*16+(i%16), (hm.z-minz)*16+(i/16), v)
			//}
			img.Set((hm.x-minx)*16+i%16, (hm.z-minz)*16+i/16, color.RGBA{uint8(v), uint8(v), uint8(v), 255})
		}
	}
	output, err := os.Create("output.png")
	defer output.Close()
	must(err)
	png.Encode(output, img)

	// 76156 chunks seen in 102 regions
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
