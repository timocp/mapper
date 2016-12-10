package main

import "bytes"
import "fmt"
import "image"
import "image/color"
import "image/draw"
import "image/png"
import "log"
import "os"

import "github.com/timocp/mapper"
import "github.com/timocp/mapper/nbt"

// a 16x16 image fragment for a single chunk
type chunkImage struct {
	x   int
	z   int
	img image.Image
}

func main() {
	//fmt.Println("mapper")
	var images []chunkImage
	var minx, maxx, minz, maxz int
	for _, fn := range os.Args[1:] {
		r := new(mapper.Region)
		fmt.Printf("Reading %s\n", fn)
		must(r.Open(fn))
		for x := 0; x < 32; x++ {
			for z := 0; z < 32; z++ {
				chunk_data, err := r.ChunkData(x, z)
				must(err)
				if chunk_data.Len() > 0 {
					chunk := mapper.NewChunk(nbt.Parse(bytes.NewReader(chunk_data.Bytes())), r, x, z)
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
					images = append(images, genHeightImage(chunk))
				}
			}
		}
	}
	fmt.Printf("imaged %d chunks (x: %d..%d, z: %d..%d)\n", len(images), minx, maxx, minz, maxz)
	// create an image large enough to place all the chunks onto
	img := image.NewRGBA(image.Rect(0, 0, (maxx-minx+1)*16, (maxz-minz+1)*16))
	// now compose all the mini images into a big one
	for _, ci := range images {
		xoffset := (ci.x - minx) * 16
		zoffset := (ci.z - minz) * 16
		r := image.Rect(xoffset, zoffset, xoffset+16, zoffset+16)
		draw.Draw(img, r, ci.img, image.Point{0, 0}, draw.Src)
	}
	output, err := os.Create("output.png")
	defer output.Close()
	must(err)
	png.Encode(output, img)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// return a 16x16 chunkImage where brightness is relative to height
func genHeightImage(c *mapper.Chunk) chunkImage {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for i, v := range c.HeightMap() {
		img.Set(i%16, i/16, color.RGBA{uint8(v), uint8(v), uint8(v), 255})
	}
	return chunkImage{c.X(), c.Z(), img}
}
