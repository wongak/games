package main

import "sync"

type leafMask uint8

const (
	maskNW uint8 = 0b1000
	maskNE uint8 = 0b0100
	maskSW uint8 = 0b0010
	maskSE uint8 = 0b0001
)

// serializes a grid of 4x4 (4 2x2 leafs) to their unique 16bit index
func serializeQuad(nw, ne, sw, se *Leaf) uint16 {
	var ser uint16
	ser |= uint16((uint8(*nw) & maskNW)) << 15
	ser |= uint16((uint8(*nw) & maskNE)) << 14
	ser |= uint16((uint8(*ne) & maskNW)) << 13
	ser |= uint16((uint8(*ne) & maskNE)) << 12
	ser |= uint16((uint8(*nw) & maskSW)) << 11
	ser |= uint16((uint8(*nw) & maskSE)) << 10
	ser |= uint16((uint8(*ne) & maskSW)) << 9
	ser |= uint16((uint8(*ne) & maskNE)) << 8
	ser |= uint16((uint8(*sw) & maskNW)) << 7
	ser |= uint16((uint8(*sw) & maskNE)) << 6
	ser |= uint16((uint8(*se) & maskNW)) << 5
	ser |= uint16((uint8(*se) & maskNE)) << 4
	ser |= uint16((uint8(*sw) & maskSW)) << 3
	ser |= uint16((uint8(*sw) & maskSE)) << 2
	ser |= uint16((uint8(*se) & maskSW)) << 1
	ser |= uint16((uint8(*se) & maskSE))
	return ser
}

var (
	quadCache map[uint16]*Leaf
	quadM     sync.RWMutex
)

func fillQuadCache() {
	quadM.Lock()
	if quadCache == nil {
		quadCache = make(map[uint16]*Leaf, 0xffff)
	}
	quadM.Unlock()
	var i uint16
	for i = 0; i <= 0xffff; i += 0x1111 {
		go func(n uint16) {
			var j uint16
			for j = n; j < n+0x1111; j++ {
				next := calcQuadNextGen(j)
				quadM.Lock()
				quadCache[j] = next
				quadM.Unlock()
			}
		}(i)
	}
}

func quadNextGen(nw, ne, sw, se *Leaf) *Leaf {
	if quadCache == nil {
		quadCache = make(map[uint16]*Leaf, 0xffff)
	}
	i := serializeQuad(nw, ne, sw, se)
	if l := quadCache[i]; l != nil {
		return l
	} else {
		l = calcQuadNextGen(i)
		quadCache[i] = l
		return l
	}
}

func nwQuad(quad uint16) uint16 {
	var ser uint16
	ser |= (quad & 0b1110000000000000)
	ser |= (quad & 0b0000111000000000) << 1
	ser |= (quad & 0b0000000011100000) << 2
	return ser
}

func calcQuadNextGen(quad uint16) *Leaf {
	// # #  # #
	// # #  # #
	//
	// # #  # #
	// # #  # #

	var leaf uint8
	var ser uint16
	// nw
	ser = nwQuad(quad)
	if serializedCenterNextGen(ser) {
		leaf |= 0b1000
	}
	// ne
	ser = 0
	ser |= (quad & 0b0111000000000000) << 1
	ser |= (quad & 0b0000011100000000) << 2
	ser |= (quad & 0b0000000001110000) << 3
	if serializedCenterNextGen(ser) {
		leaf |= 0b0100
	}
	// sw
	ser = 0
	ser |= (quad & 0b0000111000000000) << 4
	ser |= (quad & 0b0000000011100000) << 5
	ser |= (quad & 0b0000000000001110) << 6
	if serializedCenterNextGen(ser) {
		leaf |= 0b0010
	}
	// se
	ser = 0
	ser |= (quad & 0b0000011100000000) << 5
	ser |= (quad & 0b0000000001110000) << 6
	ser |= (quad & 0b0000000000000111) << 7
	if serializedCenterNextGen(ser) {
		leaf |= 0b0001
	}
	leaf = leaf << 4
	return Leafs[leaf]
}

// serializedCenterNextGen applies the rules of conway's game
//
// The 9 most significant bits contain the neighbourhood of the center cell
func serializedCenterNextGen(ser uint16) bool {
	neighbours := serializedNeighbours(ser)
	// if center is live
	if ser&(0b000010000<<7) > 0 {
		if neighbours == 2 || neighbours == 3 {
			return true
		}
	}
	if ser&(0b000010000<<7) == 0 {
		if neighbours == 3 {
			return true
		}
	}
	return false
}

// serializedNeighbours returns number of neighbours in serialized 9-block
func serializedNeighbours(ser uint16) uint8 {
	var count uint8
	if ser&(0b100000000<<7) > 0 {
		count++
	}
	if ser&(0b010000000<<7) > 0 {
		count++
	}
	if ser&(0b001000000<<7) > 0 {
		count++
	}
	if ser&(0b000100000<<7) > 0 {
		count++
	}
	if ser&(0b000001000<<7) > 0 {
		count++
	}
	if ser&(0b000000100<<7) > 0 {
		count++
	}
	if ser&(0b000000010<<7) > 0 {
		count++
	}
	if ser&(0b000000001<<7) > 0 {
		count++
	}
	return count
}
