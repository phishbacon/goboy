package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"goboy/cart"
	"goboy/soc"
	"goboy/util"
	"os"
)

type goboy struct {
	soc     *soc.SOC
	paused  bool
	running bool
	ticks   uint64
	cart    *cart.Cart
}

func NewGoboy() *goboy {

	return &goboy{
		soc: soc.NewSOC(),
	}
}

func (g *goboy) Start() {
	g.soc.Init()

	g.running = true
	g.paused = false
	g.ticks = 0

	for g.running {
		if g.paused {
			continue
		}
		g.soc.Step()
		g.ticks++
	}
}

func (g *goboy) LoadCart(fileName string) {
	dump, dumpErr := os.ReadFile(fileName)
	if dumpErr != nil {
		fmt.Print(dumpErr)
		os.Exit(-1)
	}
	// give bus reference to cart so soc components can read from it
	g.cart = (*cart.Cart)(&dump)
	g.soc.Bus.ConnectCart(&dump)

	var cartHeader cart.CartHeader
	headerErr := binary.Read(bytes.NewReader((*g.cart)[0x0100:0x0150]), binary.LittleEndian, &cartHeader)
	if headerErr != nil {
		fmt.Print(headerErr)
		os.Exit(-1)
	}

	fmt.Printf("\nTITLE      %s\n", string(cartHeader.Title[:]))
	fmt.Printf("LIC        %s\n", cartHeader.GetCartLicName())
	fmt.Printf("SGB        %x\n", cartHeader.SGBFlag)
	fmt.Printf("TYPE       %s\n", cartHeader.GetCartTypeName())
	fmt.Printf("ROM SIZE   %d KB\n", 32<<cartHeader.ROMSize)
	fmt.Printf("RAM SIZE   %s\n", cartHeader.GetRAMSize())
	fmt.Printf("DEST CODE  %s\n", cartHeader.GetDestCode())
	fmt.Printf("VERSION    %d\n", cartHeader.Version)

	var checksum uint8 = 0
	for address := 0x0134; address <= 0x014C; address++ {
		checksum = checksum - (*g.cart)[address] - 1
	}

	checksumPassed := util.If(cartHeader.HeaderChecksum == (checksum&0x00FF), "PASSED", "FAILED")
	fmt.Printf("CHECKSUM   %s\n", checksumPassed)
	g.cart.DumpHex()
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Rom required")
		os.Exit(-1)
	}

	goboy := NewGoboy()
	goboy.LoadCart(args[1])

	fmt.Printf("Loading %s\n", args[1])
	if goboy.cart.VerifyLogoDump() {
		// jump to address 0x0100
	} else {
		fmt.Println("Failed to verify logo")
		os.Exit(-1)
	}

	goboy.Start()

}
