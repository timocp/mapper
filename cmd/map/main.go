package main

import "bytes"
import "fmt"
import "log"
import "os"

import "github.com/timocp/mapper"
import "github.com/timocp/mapper/nbt"

func main() {
	//fmt.Println("mapper")
	for _, fn := range os.Args[1:] {
		r := new(mapper.Region)
		fmt.Printf("Reading %s\n", fn)
		must(r.Open(fn))
		for x := 0; x < 32; x++ {
			for z := 0; z < 32; z++ {
				//offset, length := r.chunk_location(x, z)
				//fmt.Printf("chunk(%2d, %2d) timestamp=%s offset=%d length=%d\n", x, z, r.chunk_timestamp(x, z), offset, length)
				fmt.Printf("Loading chunk %d.%d ... ", x, z)
				chunk_data, err := r.ChunkData(x, z)
				must(err)
				if chunk_data.Len() == 0 {
					fmt.Printf("not in file!\n")
				} else {
					fmt.Printf("parsing...\n")
					chunk := mapper.Chunk{nbt.Parse(bytes.NewReader(chunk_data.Bytes()))}
					chunk.Debug()
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
