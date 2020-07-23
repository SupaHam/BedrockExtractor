package gen_structs

import (
	"bytes"
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"os"
)

var blocksFile *os.File

func setupBlocks(generator *Generator) {
	blocksFile = generator.AddFile("blockstates.json")
	generator.AddListener(blocks)
}

func blocks(header packet.Header, payload []byte, src, dst net.Addr) {
	if header.PacketID == packet.IDStartGame {
		ppacket := packet.StartGame{}
		if err := ppacket.Unmarshal(bytes.NewBuffer(payload)); err != nil {
			panic(err.Error())
		}
		encoder := json.NewEncoder(blocksFile)
		states := ppacket.Blocks
		if err := encoder.Encode(states); err != nil {
			panic(err.Error())
		}
		if err := blocksFile.Sync(); err != nil {
			panic(err.Error())
		}
	}
}
