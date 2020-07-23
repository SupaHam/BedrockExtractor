package gen_structs

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"log"
	"net"
	"os"
)

var DefaultGenerator = &Generator{}
var DefaultPool = packet.NewPool()

type Listener func(header packet.Header, payload []byte, src, dst net.Addr)

type Generator struct {
	Files []*os.File
	Funcs []Listener
}

func GeneratorListenerSetup() func() {
	setupOutput(DefaultGenerator)
	setupBlocks(DefaultGenerator)
	setupItems(DefaultGenerator)
	setupEntities(DefaultGenerator)
	setupBiomes(DefaultGenerator)
	setupCreative(DefaultGenerator)
	setupRecipes(DefaultGenerator)
	return func() {
		for _, file := range DefaultGenerator.Files {
			if err := file.Close(); err != nil {
				log.Println(err.Error())
			}
		}
	}
}

func GeneratorListenerPacketFunc(header packet.Header, payload []byte, src, dst net.Addr) {
	if header.PacketID <= 0 {
		return
	}
	for _, f := range DefaultGenerator.Funcs {
		f(header, payload, src, dst)
	}
}

func (gen *Generator) AddFile(fileName string) *os.File {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err.Error())
	}
	return file
}

func (gen *Generator) AddListener(listener Listener)  {
	gen.Funcs = append(gen.Funcs, listener)
}
