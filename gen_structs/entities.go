package gen_structs

import (
	"bytes"
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"os"
)

var entitiesFile *os.File

func setupEntities(generator *Generator) {
	entitiesFile = generator.AddFile("entities.json")
	generator.AddListener(entities)
}

func entities(header packet.Header, payload []byte, src, dst net.Addr) {
	if header.PacketID == packet.IDAvailableActorIdentifiers {
		ppacket := packet.AvailableActorIdentifiers{}
		if err := ppacket.Unmarshal(bytes.NewBuffer(payload)); err != nil {
			panic(err.Error())
		}
		encoder := json.NewEncoder(entitiesFile)
		mmap := make(map[string]interface{})
		if err := nbt.Unmarshal(ppacket.SerialisedEntityIdentifiers, &mmap); err != nil {
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
