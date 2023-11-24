package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

// Block represents a block in the blockchain
type Block struct {
	Transactions []string
	PrevBlock    *Block      // Pointer to the previous block
	NextBlock    *Block      // Pointer to the next block
	Root         *MerkleNode // Merkle root of the transactions in the block
	Nonce        int
}

type MerkleNode struct {
	Hashes []string
	Left   *MerkleNode
	Right  *MerkleNode
}

func calculateHash(data []string) string {
	combinedData := ""
	for _, d := range data {
		combinedData += d
	}
	hash := sha256.Sum256([]byte(combinedData))
	return fmt.Sprintf("%x", hash)
}

func buildMerkleTree(data []string) *MerkleNode {
	if len(data) == 1 {
		return &MerkleNode{Hashes: data, Left: nil, Right: nil}
	}
	mid := len(data) / 2
	left := buildMerkleTree(data[:mid])
	right := buildMerkleTree(data[mid:])
	return &MerkleNode{Hashes: []string{calculateHash(data)}, Left: left, Right: right}
}

func printTree(node *MerkleNode, indent string) {
	if node != nil {
		fmt.Println(indent+"Hash:", node.Hashes)
	}
	if node.Left != nil {
		printTree(node.Left, indent+"  ")
	}
	if node.Right != nil {
		printTree(node.Right, indent+"  ")
	}
}

// Blockchain represents a linked list of blocks
type Blockchain struct {
	Head *Block // Head of the blockchain (genesis block)
}

// NewBlock creates a new block with the given transactions and links it to the previous block
func NewBlock(transactions []string, prevBlock *Block) *Block {
	block := &Block{
		Transactions: transactions,
		PrevBlock:    prevBlock,
		NextBlock:    nil,
		Root:         nil,
		Nonce:        0,
	}

	block.Root = buildMerkleTree(block.Transactions)
	block.MineBlock()
	return block
}

func calculateHashWithNonce(data []string, nonce int) string {
	combinedData := ""
	for _, d := range data {
		combinedData += d
	}
	combinedData += strconv.Itoa(nonce)
	hash := sha256.Sum256([]byte(combinedData))
	return fmt.Sprintf("%x", hash)
}

func (b *Block) MineBlock() {
	targetPrefix := "0000" // Define the target prefix for the hash
	for {
		hash := calculateHashWithNonce(b.Transactions, b.Nonce)
		if hash[:len(targetPrefix)] == targetPrefix {
			// Valid nonce found
			fmt.Printf("Block mined! Nonce: %d, Hash: %s\n", b.Nonce, hash)
			break
		}
		// Increment nonce and try again
		b.Nonce++
	}
}

// PrintBlockchain prints the content of the blockchain
func (bc *Blockchain) PrintBlockchain() {
	currentBlock := bc.Head
	for currentBlock != nil {
		fmt.Printf("Transactions: %v\n", currentBlock.Transactions)
		fmt.Printf("Nonce: %v\n", currentBlock.Nonce)
		// fmt.Printf("Merkle Tree:\n")
		printTree(currentBlock.Root, "")
		currentBlock = currentBlock.NextBlock
	}
}

func main() {
	// Create a new blockchain with the genesis block
	genesisBlock := NewBlock([]string{"Genesis Transaction1", "Genesis Transaction2", "Genesis Transaction3"}, nil)
	blockchain := &Blockchain{Head: genesisBlock}

	// Add new blocks to the blockchain
	block2 := NewBlock([]string{"Second Transaction1", "Second Transaction2", "Second Transaction3"}, genesisBlock)
	genesisBlock.NextBlock = block2

	block3 := NewBlock([]string{"Third Transaction1", "Third Transaction2", "Third Transaction3"}, block2)
	block2.NextBlock = block3

	// Print the content of the blockchain
	fmt.Println("Blockchain:")
	blockchain.PrintBlockchain()
}
