package gen_structs

import (
	"bytes"
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"os"
)

var biomesFile *os.File

func setupBiomes(generator *Generator) {
	biomesFile = generator.AddFile("biomes.json")
	generator.AddListener(biomes)
}

func biomes(header packet.Header, payload []byte, src, dst net.Addr) {
	if header.PacketID == packet.IDBiomeDefinitionList {
		ppacket := packet.BiomeDefinitionList{}
		if err := ppacket.Unmarshal(bytes.NewBuffer(payload)); err != nil {
			panic(err.Error())
		}
		encoder := json.NewEncoder(biomesFile)
		mmap := make(map[string]interface{})
		if err := nbt.Unmarshal(ppacket.SerialisedBiomeDefinitions, &mmap); err != nil {
			panic(err.Error())
		}
		if err := encoder.Encode(mmap); err != nil {
			panic(err.Error())
		}
		if err := entitiesFile.Sync(); err != nil {
			panic(err.Error())
		}
	}
}
