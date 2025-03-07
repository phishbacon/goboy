package cart

import (
	"fmt"
	"os"
)

var NewLicCodes = map[string]string{
	"00": "None",
	"01": "Nintendo Research & Development 1",
	"08": "Capcom",
	"13": "EA (Electronic Arts)",
	"18": "Hudson Soft",
	"19": "B-AI",
	"20": "KSS",
	"22": "Planning Office WADA",
	"24": "PCM Complete",
	"25": "San-X",
	"28": "Kemco",
	"29": "SETA Corporation",
	"30": "Viacom",
	"31": "Nintendo",
	"32": "Bandai",
	"33": "Ocean Software/Acclaim Entertainment",
	"34": "Konami",
	"35": "HectorSoft",
	"37": "Taito",
	"38": "Hudson Soft",
	"39": "Banpresto",
	"41": "Ubi Soft1",
	"42": "Atlus",
	"44": "Malibu Interactive",
	"46": "Angel",
	"47": "Bullet-Proof Software2",
	"49": "Irem",
	"50": "Absolute",
	"51": "Acclaim Entertainment",
	"52": "Activision",
	"53": "Sammy USA Corporation",
	"54": "Konami",
	"55": "Hi Tech Expressions",
	"56": "LJN",
	"57": "Matchbox",
	"58": "Mattel",
	"59": "Milton Bradley Company",
	"60": "Titus Interactive",
	"61": "Virgin Games Ltd.3",
	"64": "Lucasfilm Games4",
	"67": "Ocean Software",
	"69": "EA (Electronic Arts)",
	"70": "Infogrames5",
	"71": "Interplay Entertainment",
	"72": "Broderbund",
	"73": "Sculptured Software6",
	"75": "The Sales Curve Limited7",
	"78": "THQ",
	"79": "Accolade",
	"80": "Misawa Entertainment",
	"83": "lozc",
	"86": "Tokuma Shoten",
	"87": "Tsukuda Original",
	"91": "Chunsoft Co.8",
	"92": "Video System",
	"93": "Ocean Software/Acclaim Entertainment",
	"95": "Varie",
	"96": "Yonezawa/s’pal",
	"97": "Kaneko",
	"99": "Pack-In-Video",
	"9H": "Bottom Up",
	"A4": "Konami (Yu-Gi-Oh!)",
	"BL": "MTO",
	"DK": "Kodansha",
}

var Types = [0x23]string{
	0x00: "ROM ONLY",
	0x01: "MBC1",
	0x02: "MBC1+RAM",
	0x03: "MBC1+RAM+BATTERY",
	0x05: "MBC2",
	0x06: "MBC2+BATTERY",
	0x0B: "MMM01",
	0x0C: "MMM01+RAM",
	0x0D: "MMM01+RAM+BATTERY",
	0x0F: "MBC3+TIMER+BATTERY",
	0x10: "MBC3+TIMER+RAM+BATTERY 10",
	0x11: "MBC3",
	0x12: "MBC3+RAM 10",
	0x13: "MBC3+RAM+BATTERY 10",
	0x19: "MBC5",
	0x1A: "MBC5+RAM",
	0x1B: "MBC5+RAM+BATTERY",
	0x1C: "MBC5+RUMBLE",
	0x1D: "MBC5+RUMBLE+RAM",
	0x1E: "MBC5+RUMBLE+RAM+BATTERY",
	0x20: "MBC6",
	0x22: "MBC7+SENSOR+RUMBLE+RAM+BATTERY",
}

var RAMSizes = [0x06]string{
	0x00: "0",
	0x02: "8 KB",
	0x03: "32 KB",
	0x04: "128 KB",
	0x05: "64 KB",
}

var DestCodes = [0x02]string{
	0x00: "Japan",
	0x01: "Not Japan",
}

var Logo = [0x30]byte{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
	0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
	0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

type CartHeader struct {
	Entry          [0x0004]uint8
	Logo           [0x0030]uint8
	Title          [0x0010]uint8
	NewLicCode     [0x0002]uint8
	SGBFlag        uint8
	Type           uint8
	ROMSize        uint8
	RAMSize        uint8
	DestCode       uint8
	OldLicCode     uint8
	Version        uint8
	HeaderChecksum uint8
	GlobalChecksum uint8
}

type Cart []byte

func (c CartHeader) GetCartLicName() string {
	ascii := string(c.NewLicCode[:])
	if code, ok := NewLicCodes[ascii]; ok {
		return code
	} else {
		return "UNKOWN LICENSE CODE"
	}
}

func (c CartHeader) GetCartTypeName() string {
	if c.Type < uint8(len(Types)) && Types[c.Type] != "" {
		return Types[c.Type]
	} else {
		return "UNKNOWN TYPE"
	}
}

func (c CartHeader) GetRAMSize() string {
	if c.RAMSize < uint8(len(RAMSizes)) && RAMSizes[c.RAMSize] != "" {
		return RAMSizes[c.RAMSize]
	} else {
		return "UNKNOWN RAM SIZE"
	}
}

func (c CartHeader) GetDestCode() string {
	if c.DestCode < uint8(len(DestCodes)) && DestCodes[c.DestCode] != "" {
		return DestCodes[c.DestCode]
	} else {
		return "UNKNOWN DEST CODE"
	}
}

func (c *Cart) Read(address uint16) uint8 {
	return (*c)[address]
}

func (c *Cart) Write(address uint16, value uint8) {
	(*c)[address] = value
}

func (c *Cart) DumpHex() {
	fo, err := os.Create("dump.txt")
	if err != nil {
		fmt.Print(err)
	}

	defer fo.Close()
	cart := *c
	for i := 0; i < len(cart)-16; i += 16 {
		str := fmt.Sprintf("0x%04x: ", i)
		str += fmt.Sprintf("%02x %02x %02x %02x ", cart[i], cart[i+1], cart[i+2], cart[i+3])
		str += fmt.Sprintf("%02x %02x %02x %02x ", cart[i+4], cart[i+5], cart[i+6], cart[i+7])
		str += fmt.Sprintf("%02x %02x %02x %02x ", cart[i+8], cart[i+9], cart[i+10], cart[i+11])
		str += fmt.Sprintf("%02x %02x %02x %02x\n", cart[i+12], cart[i+13], cart[i+14], cart[i+15])
		fo.WriteString(str)
	}
}

func (c *Cart) VerifyLogoDump() bool {
	for i := range 0x30 {
		if (*c)[i+0x0104] != Logo[i] {
			return false
		}
	}
	return true
}
