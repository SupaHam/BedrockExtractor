package gen_structs

import (
	"bytes"
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"os"
)

var itemsFile *os.File

func setupItems(generator *Generator) {
	itemsFile = generator.AddFile("items.json")
	generator.AddListener(items)
}

func items(header packet.Header, payload []byte, src, dst net.Addr) {
	if header.PacketID == packet.IDStartGame {
		ppacket := packet.StartGame{}
		if err := ppacket.Unmarshal(bytes.NewBuffer(payload)); err != nil {
			panic(err.Error())
		}
		encoder := json.NewEncoder(itemsFile)
		states := ppacket.Items
		if err := encoder.Encode(states); err != nil {
			panic(err.Error())
		}
		if err := itemsFile.Sync(); err != nil {
			panic(err.Error())
		}
	}
}
