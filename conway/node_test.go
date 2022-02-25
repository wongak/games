package main

import (
	"math/rand"
	"testing"
)

type neighbourTestcase struct {
	Serialized uint16
	Neighbours uint8
}

func provideTestNeighbours() []neighbourTestcase {
	return []neighbourTestcase{
		{Serialized: 0b000101000 << 7, Neighbours: 2},
		{Serialized: 0b001101000 << 7, Neighbours: 3},
		{Serialized: 0b111101000 << 7, Neighbours: 5},
	}
}

func TestCalculateNeighbours(t *testing.T) {

	for _, c := range provideTestNeighbours() {
		if serializedNeighbours(c.Serialized) != c.Neighbours {
			t.Errorf("expect %b to have %d neighbours, got %d", c.Serialized, c.Neighbours, serializedNeighbours(c.Serialized))
		}
	}
}

type quadTestcase struct {
	Serialized uint16
	Expect     uint16
}

func provideNWQuadCases() []quadTestcase {
	return []quadTestcase{
		{Serialized: 0b001_0_000_0_111_00010,
			Expect: 0b001_000_111 << 7,
		},
		{Serialized: 0b101_0_010_0_110_00010,
			Expect: 0b101_010_110 << 7,
		},
	}
}

func TestNWQuadCases(t *testing.T) {
	for _, c := range provideNWQuadCases() {
		if nwQuad(c.Serialized) != c.Expect {
			t.Error("expect")
			visualizeSerialized(t, c.Serialized)
			t.Error("to be")
			visualizeSerialized(t, c.Expect)
			t.Error("got")
			visualizeSerialized(t, nwQuad(c.Serialized))
		}
	}
}

func visualizeSerialized(t *testing.T, ser uint16) {
	vis := func(a uint16) string {
		if a > 0 {
			return "x"
		}
		return "#"
	}
	t.Logf("%s %s %s %s\n",
		vis(ser&0b1000000000000000),
		vis(ser&0b0100000000000000),
		vis(ser&0b0010000000000000),
		vis(ser&0b0001000000000000),
	)
	t.Logf("%s %s %s %s\n",
		vis(ser&0b0000100000000000),
		vis(ser&0b0000010000000000),
		vis(ser&0b0000001000000000),
		vis(ser&0b0000000100000000),
	)
	t.Logf("%s %s %s %s\n",
		vis(ser&0b0000000010000000),
		vis(ser&0b0000000001000000),
		vis(ser&0b0000000000100000),
		vis(ser&0b0000000000010000),
	)
	t.Logf("%s %s %s %s\n",
		vis(ser&0b0000000000001000),
		vis(ser&0b0000000000000100),
		vis(ser&0b0000000000000010),
		vis(ser&0b0000000000000001),
	)
}
func visualizeLeaf(t *testing.T, ser uint8) {
	vis := func(a uint8) string {
		if a > 0 {
			return "x"
		}
		return "#"
	}
	t.Logf("%s %s\n",
		vis(ser&(0b1000<<4)),
		vis(ser&(0b0100<<4)),
	)
	t.Logf("%s %s\n",
		vis(ser&(0b0010<<4)),
		vis(ser&(0b0001<<4)),
	)
}

func BenchmarkQuadGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		calcQuadNextGen(uint16(i))
	}
}

func TestQuads(t *testing.T) {
	l := calcQuadNextGen(0)
	for i := 0; i < 10; i++ {
		index := uint16(rand.Intn(0xffff))
		l = calcQuadNextGen(index)
		visualizeSerialized(t, index)
		t.Logf("to %b", *l)
		visualizeLeaf(t, uint8(*l))
		t.Log("\n\n")
	}
}

func TestNodeLeafEvolution(t *testing.T) {
}
