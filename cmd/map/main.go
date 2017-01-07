package main

import "bytes"
import "flag"
import "fmt"
import "image"
import "image/color"
import "image/draw"
import "image/png"
import "log"
import "os"
import "sync"

import "github.com/timocp/mapper"
import "github.com/timocp/nbt"

// a 16x16 image fragment for a single chunk, which knows its world x/z offset
type chunkImage struct {
	x   int
	z   int
	img image.Image
}

func main() {
	optType := flag.String("type", "terrain", "type of map to generate(biomes, height, terrain)")
	flag.Parse()
	var images []chunkImage
	var wg sync.WaitGroup
	chImages := make(chan chunkImage)
	var minx, maxx, minz, maxz int
	for _, fn := range flag.Args() {
		wg.Add(1)
		go func(fn string) {
			defer wg.Done()
			imageRegion(fn, optType, chImages)
		}(fn)
	}
	go func() {
		// separate goroutine to close the channel when all files have
		// been read
		wg.Wait()
		close(chImages)
	}()
	// main routine reads from the channel until it is closed
	for ci := range chImages {
		images = append(images, ci)
	}
	// work out the world min/max offsets so we know where the image 0,0 is
	for _, ci := range images {
		if ci.x < minx {
			minx = ci.x
		}
		if ci.x > maxx {
			maxx = ci.x
		}
		if ci.z < minz {
			minz = ci.z
		}
		if ci.z > maxz {
			maxz = ci.z
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
	output, err := os.Create(*optType + ".png")
	defer output.Close()
	must(err)
	png.Encode(output, img)
}

// parses a region file, generating an image for each chunk and sending them to c
func imageRegion(fn string, optType *string, c chan chunkImage) {
	r := new(mapper.Region)
	fmt.Printf("Reading %s\n", fn)
	must(r.Open(fn))
	for x := 0; x < 32; x++ {
		for z := 0; z < 32; z++ {
			chunkData, err := r.ChunkData(x, z)
			must(err)
			if chunkData.Len() > 0 {
				chunk := mapper.NewChunk(nbt.Parse(bytes.NewReader(chunkData.Bytes())), r, x, z)
				switch *optType {
				case "biomes":
					c <- genBiomesImage(chunk)
				case "terrain":
					c <- genTerrainImage(chunk)
				case "height":
					c <- genHeightImage(chunk)
				default:
					panic(fmt.Sprintf("%s: invalid type", *optType))
				}
			}
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func genBiomesImage(c *mapper.Chunk) chunkImage {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for i, v := range c.Biomes() {
		img.Set(i%16, i/16, biomeColour(v))
	}
	return chunkImage{c.X(), c.Z(), img}
}

// return a 16x16 chunkImage where brightness is relative to height
func genHeightImage(c *mapper.Chunk) chunkImage {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for i, v := range c.HeightMap() {
		img.Set(i%16, i/16, color.RGBA{uint8(v), uint8(v), uint8(v), 255})
	}
	return chunkImage{c.X(), c.Z(), img}
}

func genTerrainImage(c *mapper.Chunk) chunkImage {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	//fmt.Printf("x=%d z=%d\n", c.X(), c.Z())
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			// start looking at top of highest in-chunk section
			// don't use heightmap, because that is about max light
			y := c.MaxSection()*16 + 15
			block := c.BlockAt(x, y, z)
			// go down until we find something other than air
			for block.Id == mapper.Air && y > 0 {
				y--
				block = c.BlockAt(x, y, z)
			}
			//fmt.Printf("%d.%d.%d is %s\n", c.X()*16+x, y, c.Z()*16+z, block.Name())
			img.Set(x, z, block.Colour())
		}
	}
	return chunkImage{c.X(), c.Z(), img}
}

func biomeColour(id byte) color.RGBA {
	switch id {
	case 1:
		return color.RGBA{0, 255, 127, 255} // Plains (Spring Green)
	case 3:
		return color.RGBA{152, 251, 152, 255} // Extreme Hills (Pale Green)
	case 4, 18:
		return color.RGBA{34, 139, 34, 255} // Forest [Hills] (Forest Green)
	case 5:
		return color.RGBA{0, 128, 0, 255} // Taiga (Green)
	case 6:
		return color.RGBA{107, 142, 35, 255} // Swampland (OliveDrab)
	case 7:
		return color.RGBA{0, 0, 205, 255} // River (MediumBlue)
	case 29:
		return color.RGBA{34, 139, 34, 255} // Roofed Forest (ForestGreen)
	case 34:
		return color.RGBA{255, 250, 250, 255} // Extreme Hills+ (Snow)
	case 131:
		return color.RGBA{220, 220, 220, 255} // Extreme Hills M (Gainsboro)
	case 133:
		return color.RGBA{34, 139, 34, 255} // Taiga M (Sea Green)
	case 157:
		return color.RGBA{85, 107, 47, 255} // Roofed Forest M (DarkOliveGreen)
	default:
		panic(fmt.Sprintf("unhandled biome: %d", id))
	}
}
