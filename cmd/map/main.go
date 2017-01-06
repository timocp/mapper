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

import "github.com/timocp/mapper"
import "github.com/timocp/nbt"

// a 16x16 image fragment for a single chunk
type chunkImage struct {
	x   int
	z   int
	img image.Image
}

func main() {
	//fmt.Println("mapper")
	optType := flag.String("type", "terrain", "type of map to generate(biomes, height, terrain)")
	flag.Parse()
	var images []chunkImage
	var minx, maxx, minz, maxz int
	for _, fn := range flag.Args() {
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
					switch *optType {
					case "biomes":
						images = append(images, genBiomesImage(chunk))
					case "terrain":
						images = append(images, genTerrainImage(chunk))
					case "height":
						images = append(images, genHeightImage(chunk))
					default:
						panic(fmt.Sprintf("%s: invalid type", *optType))
					}
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
	output, err := os.Create(*optType + ".png")
	defer output.Close()
	must(err)
	png.Encode(output, img)
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
