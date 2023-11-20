package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-blockchain/utils"
)

type MerkleNode struct {
	Data  string
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
}

func NewMerkleNode(data string) *MerkleNode {
	return &MerkleNode{
		Data: data,
		Hash: calculateHash(data),
	}
}

func calculateHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func buildMerkleTree(data []string) *MerkleNode {
	if len(data) == 0 {
		return nil
	}
	if len(data) == 1 {
		return NewMerkleNode(data[0])
	}

	// Create a new node with the combined hash of two child nodes
	mid := len(data) / 2
	left := buildMerkleTree(data[:mid])
	right := buildMerkleTree(data[mid:])
	return NewMerkleNode(left.Hash + right.Hash)
}

type Block struct {
	transaction string
	nonce       int
	merkleRoot  *MerkleNode
	prevBlock   *Block
}

// 10 marks
func createNewBlock() {

}

// 40 marks
// All the transactions of each block are arranged in a Merkel Tree
func createMerkleTree() {

}

// 20 marks
// A method to find the nonce value for the block. The target shall be adjustable as the number
// of trailing zeros in the 256-bit output string.
func mineBlock() {

}

// no marks
func displayBlocks() {

}

// no marks
func displayMerkleTree(node *MerkleNode, level int) {
	if node == nil {
		return
	}

	// Print the current node's hash and data at the given level
	fmt.Printf("%s Level %d: %s\n Reset Hash=%s,\n Data=%s\n", utils.Blue, level, utils.Reset, node.Hash, node.Data)

	// Recursively display the left and right subtrees
	displayMerkleTree(node.Left, level+1)
	displayMerkleTree(node.Right, level+1)
}

// 20 marks
// Functionality to verify a block and the chain. The verification will consider the changes to the transactions stored in the
// Merkel tree
func verifyBlock() {

}

func verifyChain() {

}

// no marks
// Function to change one or multiple transactions of the given block ref.
func changeBlock() {

}

// 10 marks
// func calculateHash() {

// }

func main() {
	data := []string{"Transaction1", "Transaction2", "Transaction3", "Transaction4"}

	// Build the Merkle tree
	merkleTree := buildMerkleTree(data)

	// Print the root hash of the Merkle tree
	fmt.Println("Merkle Root Hash:", merkleTree.Hash)

	displayMerkleTree(merkleTree, 0)
}
