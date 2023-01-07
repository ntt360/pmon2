package proc

import (
	"encoding/binary"
	"github.com/mdlayher/netlink/nlenc"
	"golang.org/x/net/bpf"
)

func loadBPF(filters []EventType) ([]bpf.RawInstruction, error) {
	var totals uint32
	for _, toFilter := range filters {
		totals = totals | uint32(toFilter)
	}

	inst, err := bpf.Assemble([]bpf.Instruction{
		bpf.LoadAbsolute{Off: (16 + 20), Size: 4},
		bpf.JumpIf{Cond: bpf.JumpBitsNotSet, Val: swap(totals), SkipTrue: 1},
		bpf.RetConstant{Val: 4096},
		bpf.RetConstant{Val: 0},
	})

	return inst, err
}

//if needed, change the endian-ness before we file the bitmask off to bpf
func swap(val uint32) uint32 {

	if nlenc.NativeEndian() == binary.LittleEndian {
		return (val << 24) |
			((val << 8) & 0x00ff0000) |
			((val >> 8) & 0x0000ff00) |
			((val >> 24) & 0x000000ff)
	}

	return val

}
