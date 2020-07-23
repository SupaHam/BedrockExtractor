package gen_structs

import (
	"bytes"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"io"
	"net"
	"os"
	"reflect"
)

var outputFile *os.File

func setupOutput(generator *Generator) {
	outputFile = generator.AddFile("output.txt")
	generator.AddListener(outputTxt)
}

func outputTxt(header packet.Header, payload []byte, src, dst net.Addr) {
	d := "<<"
	if dst.String()[len(dst.String())-5:] == "19132" {
		d = ">>"
	}
	aa := make([]byte, len(payload))
	copy(aa, payload)
	i := DefaultPool[header.PacketID]
	p := reflect.New(reflect.ValueOf(i).Elem().Type()).Interface().(packet.Packet)
	err := p.Unmarshal(bytes.NewBuffer(aa))
	if err != nil {
		panic(err.Error())
	}
	str := fmt.Sprintf("%s %v %#v\n", d, header.PacketID, p)
	_, err = io.WriteString(outputFile, str)
	if err != nil {
		panic(err.Error())
	}
	err = outputFile.Sync()
	if err != nil {
		panic(err.Error())
	}
}
