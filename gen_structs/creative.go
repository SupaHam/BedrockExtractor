package gen_structs

import (
	"bytes"
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"os"
)

var creativeFile *os.File

func setupCreative(generator *Generator) {
	creativeFile = generator.AddFile("creative.json")
	generator.AddListener(creative)
}

func creative(header packet.Header, payload []byte, src, dst net.Addr) {
	if header.PacketID == packet.IDCreativeContent {
		ppacket := packet.CreativeContent{}
		if err := ppacket.Unmarshal(bytes.NewBuffer(payload)); err != nil {
			panic(err.Error())
		}
		encoder := json.NewEncoder(creativeFile)
		if err := encoder.Encode(ppacket.Items); err != nil {
			panic(err.Error())
		}
		if err := entitiesFile.Sync(); err != nil {
			panic(err.Error())
		}
	}
}
