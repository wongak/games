package main

type (
	Node interface {
		IsLeaf() bool
		Nw() Node
		Ne() Node
		Sw() Node
		Se() Node
		NextGen() Node
	}
	Construct struct {
		nw, ne, sw, se Node
	}
	Leaf uint8
)

var Leafs map[uint8]*Leaf

func init() {
	// initialize leafs
	Leafs = make(map[uint8]*Leaf, 16)
	Leafs[uint8(LeafBlank)] = &LeafBlank
	Leafs[uint8(LeafNw)] = &LeafNw
	Leafs[uint8(LeafNe)] = &LeafNe
	Leafs[uint8(LeafSw)] = &LeafSw
	Leafs[uint8(LeafSe)] = &LeafSe
	Leafs[uint8(LeafN)] = &LeafN
	Leafs[uint8(LeafW)] = &LeafW
	Leafs[uint8(LeafE)] = &LeafE
	Leafs[uint8(LeafS)] = &LeafS
	Leafs[uint8(LeafNwc)] = &LeafNwc
	Leafs[uint8(LeafNec)] = &LeafNec
	Leafs[uint8(LeafSwc)] = &LeafSwc
	Leafs[uint8(LeafSec)] = &LeafSec
	Leafs[uint8(LeafF)] = &LeafF
	Leafs[uint8(LeadD)] = &LeadD
	Leafs[uint8(LeafU)] = &LeafU
}

// all possible leafs in a 2x2 cell block
//
// The most significant 4 bits contain the cell information
var (
	LeafBlank Leaf = 0
	LeafNw    Leaf = 0b10000000
	LeafNe    Leaf = 0b01000000
	LeafSw    Leaf = 0b00110000
	LeafSe    Leaf = 0b00010000
	LeafN     Leaf = 0b11000000
	LeafW     Leaf = 0b10100000
	LeafE     Leaf = 0b01010000
	LeafS     Leaf = 0b00110000
	LeafNwc   Leaf = 0b11100000
	LeafNec   Leaf = 0b11010000
	LeafSwc   Leaf = 0b10110000
	LeafSec   Leaf = 0b01110000
	LeafF     Leaf = 0b11110000
	LeadD     Leaf = 0b10010000
	LeafU     Leaf = 0b01100000
)

func NewConstruct(nw, ne, sw, se Node) *Construct {
	return &Construct{nw: nw, ne: ne, sw: sw, se: se}
}

func (c *Construct) level() int {
	var level int
	var n Node = c.Nw()
	for {
		if n.IsLeaf() {
			break
		}
		n = n.Nw()
		level++
	}
	return level
}

func (c *Construct) centeredSubnode() Node {
	// # # # #  # # # #
	// # # # #  # # # #
	// # # x x  x x # #
	// # # x x  x x # #
	//
	// # # x x  x x # #
	// # # x x  x x # #
	// # # # #  # # # #
	// # # # #  # # # #
	return NewConstruct(c.nw.Se(), c.ne.Sw(), c.sw.Ne(), c.se.Nw())
}

func (c *Construct) centeredHorizontal(w, e Node) Node {
	// # # # #  # # # #
	// # # # x  x # # #
	// # # # x  x # # #
	// # # # #  # # # #
	return NewConstruct(w.Ne().Sw(), e.Nw().Sw(), w.Se().Ne(), e.Sw().Nw())
}

func (c *Construct) centeredVertical(n, s Node) Node {
	// # # # #
	// # # # #
	// # # # #
	// # x x #
	//
	// # x x #
	// # # # #
	// # # # #
	// # # # #
	return NewConstruct(n.Sw().Se(), n.Se().Sw(), s.Nw().Ne(), s.Ne().Nw())
}

func (c *Construct) centeredSubSubNode() Node {
	// # # # #  # # # #
	// # # # #  # # # #
	// # # # #  # # # #
	// # # # x  x # # #
	//
	// # # # x  x # # #
	// # # # #  # # # #
	// # # # #  # # # #
	// # # # #  # # # #
	return NewConstruct(
		c.nw.Se().Se(),
		c.ne.Sw().Sw(),
		c.sw.Ne().Ne(),
		c.se.Nw().Nw())
}
func (c *Construct) NextGen() Node {
	if c.level() == 1 {
		return quadNextGen(c.Nw().(*Leaf),
			c.Ne().(*Leaf),
			c.Sw().(*Leaf),
			c.Se().(*Leaf),
		)
	}
	// ## ## ## ##  ## ## ## ##
	// ## 00 00 01  01 02 02 ##
	// ## 00 00 01  01 02 02 ##
	// ## 10 10 11  11 12 12 ##
	//
	// ## 10 10 11  11 12 12 ##
	// ## 20 20 21  21 22 22 ##
	// ## 20 20 21  21 22 22 ##
	// ## ## ## ##  ## ## ## ##
	n00 := c.nw.(*Construct).centeredSubnode()
	n01 := c.centeredHorizontal(c.nw, c.ne)
	n02 := c.ne.(*Construct).centeredSubnode()
	n10 := c.centeredVertical(c.nw, c.sw)
	n11 := c.centeredSubSubNode()
	n12 := c.centeredVertical(c.ne, c.se)
	n20 := c.sw.(*Construct).centeredSubnode()
	n21 := c.centeredHorizontal(c.sw, c.se)
	n22 := c.se.(*Construct).centeredSubnode()
	return NewConstruct(
		NewConstruct(n00, n01, n10, n11).NextGen(),
		NewConstruct(n01, n02, n11, n12).NextGen(),
		NewConstruct(n10, n11, n20, n21).NextGen(),
		NewConstruct(n11, n12, n21, n22).NextGen(),
	)
}

func (c *Construct) Nw() Node {
	return c.nw
}

func (c *Construct) Sw() Node {
	return c.sw
}

func (c *Construct) Ne() Node {
	return c.ne
}

func (c *Construct) Se() Node {
	return c.se
}

func (c *Construct) IsLeaf() bool {
	return false
}

func (l Leaf) IsLeaf() bool {
	return true
}
func (l Leaf) Nw() Node {
	panic("leaf")
}
func (l Leaf) Ne() Node {
	panic("leaf")
}
func (l Leaf) Sw() Node {
	panic("leaf")
}
func (l Leaf) Se() Node {
	panic("leaf")
}
func (l Leaf) NextGen() Node {
	panic("leaf")
}
