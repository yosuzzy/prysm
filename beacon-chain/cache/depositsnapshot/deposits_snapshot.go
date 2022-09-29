package depositsnapshot

import (
	"github.com/prysmaticlabs/prysm/crypto/hash"
)

const (
	DepositContractDepth = 32
)

var zeroHash = [32]byte{}

type MerkleTree interface {
	GetRoot() [32]byte
	IsFull() bool
	Finalize(deposits uint, depth uint) MerkleTree
	GetFinalized(result [][32]byte) ([][32]byte, uint)
	PushLeaf(leaf [32]byte, deposits uint, depth uint) MerkleTree
}

func CreateMerkleTreeIter(leaves [][32]byte, deposits uint, depth uint) MerkleTree {
	switch {
	case deposits == 0, len(leaves) == 0:
		return &Zero{depth: depth}
	case depth == 0:
		return &Leaf{
			hash: leaves[0],
		}
	default:
		node := &Node{
			left:  nil,
			right: nil,
		}
		var split uint
		for depth > 0 {
			split = UintPow(2, depth-1)
			if deposits < split {
				next := &Node{}
				node.left = &Node{left: next, right: &Node{}}
				node.right = &Zero{depth: depth - 1}
				node = next // = node.left
				depth -= 1
			} else if deposits > split {
				node.left = &Finalized{
					deposits: deposits,
					hash:     leaves[0],
				}
				next := &Node{}
				node.right = &Node{left: &Node{}, right: next}
				leaves = leaves[1:]
				deposits -= split
				node = next // = node.right
				depth -= 1
			} else {
				node.left = &Finalized{split, leaves[0]}
				node.right = &Zero{depth: depth - 1}
				leaves = leaves[1:]
			}
		}
		return node
	}
}

type Finalized struct {
	deposits uint
	hash     [32]byte
}

func (f *Finalized) GetRoot() [32]byte {
	return f.hash
}

func (f *Finalized) IsFull() bool {
	return true
}

func (f *Finalized) Finalize(deposits uint, depth uint) MerkleTree {
	return f
}

func (f *Finalized) GetFinalized(result [][32]byte) ([][32]byte, uint) {
	return append(result, f.hash), f.deposits
}

func (f *Finalized) PushLeaf(leaf [32]byte, deposits uint, depth uint) MerkleTree {
	panic("Can't push a leaf to something finalized")
}

type Leaf struct {
	hash [32]byte
}

func (l *Leaf) GetRoot() [32]byte {
	return l.hash
}

func (l *Leaf) IsFull() bool {
	return true
}

func (l *Leaf) Finalize(deposits uint, depth uint) MerkleTree {
	return &Finalized{
		deposits: 1,
		hash:     l.hash,
	}
}

func (l *Leaf) GetFinalized(result [][32]byte) ([][32]byte, uint) {
	return result, 0
}

func (l *Leaf) PushLeaf(leaf [32]byte, deposits uint, depth uint) MerkleTree {
	panic("leaf should not be able to push another leaf")
}

type Node struct {
	left, right MerkleTree
}

func (n *Node) GetRoot() [32]byte {
	left := n.left.GetRoot()
	right := n.right.GetRoot()
	return hash.Hash(append(left[:], right[:]...))
}

func (n *Node) IsFull() bool {
	return n.right.IsFull()
}

func (n *Node) Finalize(deposits uint, depth uint) MerkleTree {
	depositsNum := UintPow(2, depth)
	if depositsNum <= deposits {
		return &Finalized{
			deposits: depositsNum,
			hash:     n.GetRoot(),
		}
	}
	n.left = n.left.Finalize(deposits, depth-1)
	if deposits > depositsNum/2 {
		remaining := deposits - depositsNum/2
		n.right = n.right.Finalize(remaining, depth-1)
	}
	return n
}

func (n *Node) GetFinalized(result [][32]byte) ([][32]byte, uint) {
	result, depositsLeft := n.left.GetFinalized(result)
	result, depositsRight := n.right.GetFinalized(result)

	return result, depositsLeft + depositsRight
}

func (n *Node) PushLeaf(leaf [32]byte, deposits uint, depth uint) MerkleTree {
	if !n.left.IsFull() {
		n.left = n.left.PushLeaf(leaf, deposits, depth-1)
	} else {
		n.right = n.right.PushLeaf(leaf, deposits, depth-1)
	}
	return n
}

type Zero struct {
	depth uint
}

func (z *Zero) GetRoot() [32]byte {
	if z.depth == DepositContractDepth {
		return hash.Hash(append(zeroHash[:], zeroHash[:]...))
	}
	return zeroHash
}

func (z *Zero) IsFull() bool {
	return false
}

func (z *Zero) Finalize(deposits uint, depth uint) MerkleTree {
	panic("finalize should not be called")
}

func (z *Zero) GetFinalized(result [][32]byte) ([][32]byte, uint) {
	return result, 0
}

func (z *Zero) PushLeaf(leaf [32]byte, deposits uint, depth uint) MerkleTree {
	return CreateMerkleTreeIter([][32]byte{leaf}, deposits, depth)
}

type DepositTree struct {
	tree                    MerkleTree
	mixInLength             uint
	finalizedExecutionblock [32]byte
}

func (d *DepositTree) FromSnapshot(finalized [][32]byte, deposits uint) MerkleTree {
	return fromSnapshotParts(finalized, deposits, DepositContractDepth)
}

func fromSnapshotParts(finalized [][32]byte, deposits uint, depth uint) MerkleTree {
	if len(finalized) < 1 || deposits == 0 {
		return &Zero{
			depth: depth,
		}
	}
	if deposits == UintPow(2, depth) {
		return &Finalized{
			deposits: deposits,
			hash:     finalized[0],
		}
	}

	node := Node{}
	if leftSubtree := UintPow(2, depth-1); deposits <= leftSubtree {
		node.left = fromSnapshotParts(finalized, deposits, depth-1)
		node.right = &Zero{depth: depth - 1}

	} else {
		node.left = &Finalized{
			deposits: leftSubtree,
			hash:     finalized[0],
		}
		node.right = fromSnapshotParts(finalized[1:], deposits-leftSubtree, depth-1)
	}
	return &node
}

func (d *DepositTree) PushLeaf(leaf [32]byte, deposits uint) {
	d.mixInLength += 1
	d.tree = d.tree.PushLeaf(leaf, deposits, DepositContractDepth)
}

func UintPow(n, m uint) uint {
	if m == 0 {
		return 1
	}
	result := n
	for i := uint(2); i <= m; i++ {
		result *= n
	}
	return result
}
