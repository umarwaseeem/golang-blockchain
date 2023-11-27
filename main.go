package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"go-blockchain/utils"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type Block struct {
	Blocknumber  int
	Transactions []string
	PrevBlock    *Block
	NextBlock    *Block
	Root         *MerkleNode
	Nonce        int
	Blockhash    string
}

type MerkleNode struct {
	Hashes []string
	Left   *MerkleNode
	Right  *MerkleNode
}

type Node struct {
	IPAddress   string
	Port        int
	Neighbours  []Node
	transaction string // transaction data - received from other nodes or the one to be sent
	// make a function to send the transaction data to the other nodes
}

func (node *Node) FloodTransaction() {
	// send the transaction data to all the neighbours
}

func (node *Node) ReceiveTransaction() {
	// receive the transaction data from other nodes
}

func calculateHash(data []string) string {
	combinedData := ""
	for _, d := range data {
		combinedData += d
	}
	hash := sha256.Sum256([]byte(combinedData))
	return fmt.Sprintf("%x", hash)
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
		fmt.Println(indent+"Value in Merkle Node:", node.Hashes)
	}
	if node.Left != nil {
		printTree(node.Left, indent+"  ")
	}
	if node.Right != nil {
		printTree(node.Right, indent+"  ")
	}
}

type Blockchain struct {
	Head *Block
}

func NewBlock(transactions []string, prevBlock *Block) *Block {

	if prevBlock != nil {
		block := &Block{
			Blocknumber:  prevBlock.Blocknumber + 1,
			Transactions: transactions,
			PrevBlock:    prevBlock,
			NextBlock:    nil,
			Root:         nil,
			Nonce:        0,
			Blockhash:    "",
		}
		block.Root = buildMerkleTree(block.Transactions)
		block.MineBlock()
		block.CalculateBlockHash()
		return block
	}

	block := &Block{
		Blocknumber:  1,
		Transactions: transactions,
		PrevBlock:    prevBlock,
		NextBlock:    nil,
		Root:         nil,
		Nonce:        0,
		Blockhash:    "",
	}
	block.Root = buildMerkleTree(block.Transactions)
	block.MineBlock()
	block.CalculateBlockHash()
	return block

}

func (b *Block) CalculateBlockHash() {

	if b.PrevBlock == nil {
		b.Blockhash = "0"
	} else {
		combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce) + b.Transactions[0] + b.Transactions[1] + b.Transactions[2]
		hash := sha256.Sum256([]byte(combinedData))
		b.Blockhash = fmt.Sprintf("%x", hash)
	}
}

func (b *Block) CalculateBlockHashforverification() string {

	combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce) + b.Transactions[0] + b.Transactions[1] + b.Transactions[2]
	hash := sha256.Sum256([]byte(combinedData))
	var newstring = fmt.Sprintf("%x", hash)
	return newstring
}

func (b *Block) VerifyBlock() bool {
	var current_hash = b.Blockhash
	var hashString string
	if b.PrevBlock == nil {
		b.Blockhash = "0"
		fmt.Println("The Block is verified")
		return true
	} else {
		combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce) + b.Transactions[0] + b.Transactions[1] + b.Transactions[2]
		hash := sha256.Sum256([]byte(combinedData))
		hashString := fmt.Sprintf("%x", hash)
		if current_hash == hashString {
			fmt.Println("Current hash:   " + b.Blockhash)
			fmt.Printf("Calculated hash: %s\n", hashString)
			fmt.Println("The Block is verified")
			return true
		}

	}
	fmt.Println("Current hash:   " + b.Blockhash)
	fmt.Println(hashString)
	return false
}

func (blockchain *Blockchain) VerifyChain() {
	currentBlock := blockchain.Head
	index := 0

	for currentBlock != nil {
		fmt.Printf("%sVerifying Block #%d:\n%s", utils.Cyan, index, utils.Reset)
		var verification = currentBlock.VerifyBlock()
		if !verification {
			fmt.Println("The block number " + strconv.Itoa(index) + " seems to be tampered and can't be verified")
			return
		}
		fmt.Printf("%s------------------------%s", utils.Cyan, utils.Reset)
		index++
		currentBlock = currentBlock.NextBlock
	}
}

func (b *Block) MineBlock() {
	targetPrefix := "0000"
	for {
		hash := calculateHashWithNonce(b.Transactions, b.Nonce)
		if hash[:len(targetPrefix)] == targetPrefix {
			fmt.Printf("Block mined! Nonce: %d, Hash: %s\n", b.Nonce, hash)
			break
		}
		b.Nonce++
	}
}

func (bc *Blockchain) PrintBlockchain() {
	currentBlock := bc.Head
	for currentBlock != nil {
		fmt.Printf("Transactions: %v\n", currentBlock.Transactions)
		fmt.Printf("Nonce: %v\n", currentBlock.Nonce)
		fmt.Println("Block Hash Value: ", currentBlock.Blockhash)
		printTree(currentBlock.Root, "")
		currentBlock = currentBlock.NextBlock
	}
}

func (blockchain *Blockchain) printBlock(x int) {
	var currentBlock = blockchain.Head
	for currentBlock != nil {
		if currentBlock.Blocknumber != x {
			currentBlock = currentBlock.NextBlock
		}
		if currentBlock.Blocknumber == x {
			fmt.Println("Block number: ", x)
			printTree(currentBlock.Root, "")
			return
		}
	}
}

func (blockchain *Blockchain) changeBlock(x int) {
	var currentBlock = blockchain.Head
	var changestringnum int
	var newtransaction string
	for currentBlock != nil {
		if currentBlock.Blocknumber != x {
			currentBlock = currentBlock.NextBlock
		}
		if currentBlock.Blocknumber == x {
			fmt.Println("Here are the transactions in this block:")
			fmt.Printf("Transactions: %v\n", currentBlock.Transactions)
			fmt.Println("Which transaction to change?")
			fmt.Scanln(&changestringnum)
			fmt.Println("Enter the new transaction")
			fmt.Scanln(&newtransaction)
			currentBlock.Transactions[changestringnum] = newtransaction
			// var currentBlockHash = currentBlock.CalculateBlockHashforverification()
			// if currentBlock.Blockhash != currentBlockHash {
			// 	fmt.Printf("The block number", currentBlock.Blocknumber)
			// 	fmt.Println(" has been found to be tampered with. Blockchain verification failed")
			// 	return
			// }
			return

		}
	}
}

func AskInput(prompt string, r *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := r.ReadString('\n')

	return strings.TrimSpace(input), err
}

func newPeerNode() Node {
	currentPort++
	newNode := Node{
		IPAddress:   localIP.String(),
		Port:        currentPort,
		transaction: "",
		Neighbours:  []Node{},
	}
	return newNode
}

func addNeighbourToNode(node *Node, neighbour Node) {
	node.Neighbours = append(node.Neighbours, neighbour)
	fmt.Printf("\n%sA New Node %v:%v added as neighbour of Node %v:%v\n\n%s", utils.Green, neighbour.IPAddress, neighbour.Port, node.IPAddress, node.Port, utils.Reset)
	// randomly add the node as a neighbour to one of the neighbours of the node and vice versa
	randomIndex := rand.Intn(len(node.Neighbours) - 1)
	node.Neighbours[randomIndex].Neighbours = append(node.Neighbours[randomIndex].Neighbours, neighbour)
	neighbour.Neighbours = append(neighbour.Neighbours, node.Neighbours[randomIndex])
	fmt.Printf("\n%sOld Node %v:%v added as neighbour of New Node %v:%v\n\n%s", utils.Green, neighbour.IPAddress, neighbour.Port, node.Neighbours[randomIndex].IPAddress, node.Neighbours[randomIndex].Port, utils.Reset)
}

func displayNetwork(bootstrapNode Node) {
	fmt.Printf("\nBootstrap Node: %v:%v\n\n", bootstrapNode.IPAddress, bootstrapNode.Port)
	// loop through all neighbours
	for index, neighbour := range bootstrapNode.Neighbours {
		fmt.Printf("Neighbour %d: %v:%v\n", index+1, neighbour.IPAddress, neighbour.Port)
		// loop through all neighbours of the neighbour
		for index2, neighbour2 := range neighbour.Neighbours {
			fmt.Printf("  Neighbour %d: %v:%v\n", index2+1, neighbour2.IPAddress, neighbour2.Port)
		}
	}
	fmt.Println()
}

func GetLocalIP() net.IP {
	// udp connection to google dns
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	// after connection done, we get the local ip to which the socket is bound
	return localAddr.IP
}

var bootstrapNodePort = 8080
var currentPort = 8080
var localIP = GetLocalIP()

func main() {

	node1, node2, node3, node4 := newPeerNode(), newPeerNode(), newPeerNode(), newPeerNode()
	node5 := newPeerNode()

	bootstrapNode := Node{
		IPAddress:   localIP.String(),
		Port:        bootstrapNodePort,
		transaction: "",
		Neighbours:  []Node{node1, node2, node3, node4},
	}

	displayNetwork(bootstrapNode)

	addNeighbourToNode(&bootstrapNode, node5)
	fmt.Printf("%sAfter adding a new node: %s \n", utils.Cyan, utils.Reset)
	displayNetwork(bootstrapNode)

	var input string
	var input1 string
	genesisBlock := NewBlock([]string{"Genesis Transaction1", "Genesis Transaction2", "Genesis Transaction3"}, nil)
	blockchain := &Blockchain{Head: genesisBlock}

	block2 := NewBlock([]string{"Second Transaction1", "Second Transaction2", "Second Transaction3"}, genesisBlock)
	genesisBlock.NextBlock = block2

	block3 := NewBlock([]string{"Third Transaction1", "Third Transaction2", "Third Transaction3"}, block2)
	block2.NextBlock = block3

	// 3ca38309fea3736cb0c91b4cf53048ea770bdebf35b6cbd037d532e4515bb228
	// c1f62d9b30f2f9a6121f8e060143321869296467ebde2bde7e8645593bc70beb
	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Print a block + Merkle Tree")
		fmt.Println("2. Print the blockchain + Merkle trees of all blocks")
		fmt.Println("3. Verify a specific block")
		fmt.Println("4. Verify the blockchain")
		fmt.Println("5. Change a block")
		fmt.Println("0. Exit the program")

		fmt.Scanln(&input)

		switch input {
		case "1":
			fmt.Println("Choose the block number you would like to print")
			fmt.Scanln(&input1)
			blocknumber, err := strconv.Atoi(input1)
			if err != nil {
				fmt.Println("Invalid input for block number. Please enter a valid integer.")
				continue
			}
			fmt.Print("Printing Block: ")
			blockchain.printBlock(blocknumber)

		case "2":
			blockchain.PrintBlockchain()
		case "3":
			fmt.Print("Enter block index to verify: ")
		case "4":
			blockchain.VerifyChain()
		case "5":
			fmt.Println("Choose the block number you would like to change")
			fmt.Scanln(&input1)
			blocknumber, err := strconv.Atoi(input1)
			blocknumber++
			if err != nil {
				fmt.Println("Invalid input for block number. Please enter a valid integer.")
				continue
			}
			fmt.Print("Printing Block: ")
			blockchain.changeBlock(blocknumber)

		case "0":
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}

}
