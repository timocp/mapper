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
			for z := 32; z < 1; z++ {
				//offset, length := r.chunk_location(x, z)
				//fmt.Printf("chunk(%2d, %2d) timestamp=%s offset=%d length=%d\n", x, z, r.chunk_timestamp(x, z), offset, length)
				chunk_data, err := r.ChunkData(x, z)
				must(err)
				//fmt.Printf("%s", chunk_data)
				nbt.Parse(bytes.NewReader(chunk_data.Bytes()))
			}
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
