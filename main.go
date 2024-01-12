package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

// func (b *Block) CalculateBlockHash() {

// 	if b.PrevBlock == nil {
// 		b.Blockhash = "0"
// 	} else {
// 		combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce) + b.Transactions[0] + b.Transactions[1] + b.Transactions[2]
// 		hash := sha256.Sum256([]byte(combinedData))
// 		b.Blockhash = fmt.Sprintf("%x", hash)
// 	}
// }

func (b *Block) CalculateBlockHash() {
	if b.PrevBlock == nil {
		b.Blockhash = "0"
	} else {
		combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce)

		// Combine transactions based on the actual number of transactions in the current block
		for i := 0; i < len(b.Transactions); i++ {
			combinedData += b.Transactions[i]
		}

		hash := sha256.Sum256([]byte(combinedData))
		b.Blockhash = fmt.Sprintf("%x", hash)
	}
}

// func (b *Block) CalculateBlockHashforverification() string {

// 	combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce) + b.Transactions[0] + b.Transactions[1] + b.Transactions[2]
// 	hash := sha256.Sum256([]byte(combinedData))
// 	var newstring = fmt.Sprintf("%x", hash)
// 	return newstring
// }

func (b *Block) CalculateBlockHashforverification() string {
	combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce)

	// Combine transactions based on the actual number of transactions in the current block
	for i := 0; i < len(b.Transactions); i++ {
		combinedData += b.Transactions[i]
	}

	hash := sha256.Sum256([]byte(combinedData))
	return fmt.Sprintf("%x", hash)
}

// func (b *Block) VerifyBlock() bool {
// 	var current_hash = b.Blockhash
// 	var hashString string
// 	if b.PrevBlock == nil {
// 		b.Blockhash = "0"
// 		fmt.Println("The Block is verified")
// 		return true
// 	} else {
// 		combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce) + b.Transactions[0] + b.Transactions[1] + b.Transactions[2]
// 		hash := sha256.Sum256([]byte(combinedData))
// 		hashString := fmt.Sprintf("%x", hash)
// 		if current_hash == hashString {
// 			fmt.Println("Current hash:   " + b.Blockhash)
// 			fmt.Printf("Calculated hash: %s\n", hashString)
// 			fmt.Println("The Block is verified")
// 			return true
// 		}

// 	}
// 	fmt.Println("Current hash:   " + b.Blockhash)
// 	fmt.Printf(hashString)
// 	return false
// }

func (b *Block) VerifyBlock() bool {
	var currentHash = b.Blockhash
	var hashString string

	if b.PrevBlock == nil {
		b.Blockhash = "0"
		fmt.Println("The Block is verified")
		return true
	} else {
		combinedData := b.PrevBlock.Blockhash + strconv.Itoa(b.Nonce)

		// Combine transactions based on the actual number of transactions in the current block
		for i := 0; i < len(b.Transactions); i++ {
			combinedData += b.Transactions[i]
		}

		hash := sha256.Sum256([]byte(combinedData))
		hashString = fmt.Sprintf("%x", hash)

		if currentHash == hashString {
			fmt.Println("Current hash:   " + b.Blockhash)
			fmt.Printf("Calculated hash: %s\n", hashString)
			fmt.Println("The Block is verified")
			return true
		}
	}

	fmt.Println("Current hash:   " + b.Blockhash)
	fmt.Printf(hashString)
	return false
}

func (blockchain *Blockchain) VerifyChain() {
	currentBlock := blockchain.Head
	index := 0

	for currentBlock != nil {
		fmt.Printf("Verifying Block #%d:\n", index)
		var verification = currentBlock.VerifyBlock()
		if verification == false {
			fmt.Println("\nThe block number " + strconv.Itoa(index) + " seems to be tampered and can't be verified")
			return
		}
		fmt.Println("------------------------")
		index++
		currentBlock = currentBlock.NextBlock
	}
	return
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
		fmt.Println()
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

type Node struct {
	ID           int
	IP           string
	Port         int
	Transactions []string
	Peers        map[int]*Node
	Output       chan string
	Next         *Node
	mutex        sync.Mutex
	localblock   *Block
}

func NewNode(id, port int, blockchain *Blockchain) *Node {
	return &Node{
		ID:         id,
		IP:         "127.0.0.1",
		Port:       port,
		Peers:      make(map[int]*Node),
		Output:     make(chan string),
		Next:       nil,
		localblock: nil,
	}
}

func (n *Node) AddTransaction(transaction string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.Transactions = append(n.Transactions, transaction)
}

func (n *Node) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := strings.TrimSpace(scanner.Text())

		// Store the received message in the Transactions array
		n.AddTransaction(message)

		// Broadcast the received transaction to all peers
		n.BroadcastTransaction(message)

		n.Output <- message
	}
}

func (n *Node) listenForConnections() {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(n.Port))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go n.handleConnection(conn)
	}
}

func (n *Node) StartNode() {
	go n.listenForConnections()
	for {
		select {
		case message := <-n.Output:
			if !strings.HasPrefix(message, "REGISTER") {
				n.AddTransaction(message)
			}
			// Process incoming messages
			// fmt.Printf("Node %d received message: %s\n", n.ID, message)
		}
	}
}

func (n *Node) SendMessage(peerID int, message string) {
	if peer, ok := n.Peers[peerID]; ok {
		peer.Output <- message
		// Update local Transactions array with the sent message
		n.AddTransaction(message)

		// Broadcast the transaction to all peers
		n.BroadcastTransaction(message)
	} else {
		fmt.Printf("Node %d: Peer %d not found\n", n.ID, peerID)
	}
}

type P2PNetwork struct {
	Head *Node
}

func NewP2PNetwork() *P2PNetwork {
	return &P2PNetwork{
		Head: nil,
	}
}

func (p *P2PNetwork) AddNode(node *Node) {
	if p.Head == nil {
		p.Head = node
	} else {
		current := p.Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = node
	}

	// Connect the new node to existing nodes
	current := p.Head
	for current != nil {
		if current != node {
			node.Peers[current.ID] = current
			current.Peers[node.ID] = node
		}
		current = current.Next
	}
}

func (p *P2PNetwork) DisplayNetworkInfo() {
	current := p.Head
	for current != nil {
		fmt.Printf("Node %d - IP %s, Port %d\n", current.ID, current.IP, current.Port)
		current = current.Next
	}
}

var (
	bootstrapIP   = "127.0.0.1" // Replace with your bootstrap node's IP
	bootstrapPort = 5000        // Replace with your bootstrap node's port
	nodeCount     = 0           // Tracks the number of nodes in the network
	mutex         sync.Mutex    // Mutex for synchronizing access to the network map
)

func (n *Node) RegisterWithBootstrap() {
	conn, err := net.Dial("tcp", bootstrapIP+":"+strconv.Itoa(bootstrapPort))
	if err != nil {
		fmt.Println("Failed to connect to bootstrap node:", err)
		return
	}
	defer conn.Close()

	// msg := "REGISTER:" + n.IP + ":" + strconv.Itoa(n.Port)
	// _, err = fmt.Fprintf(conn, msg)
	// if err != nil {
	// 	fmt.Println("Failed to register with bootstrap node:", err)
	// 	return
	// }
	// fmt.Println("Registered with bootstrap node")
}

func (n *Node) PrintPeers() {
	fmt.Printf("Node %d's Peers: ", n.ID)
	for peerID := range n.Peers {
		fmt.Printf("%d ", peerID)
	}
	fmt.Println()
}

func (n *Node) BroadcastTransaction(transaction string) {
	for _, peer := range n.Peers {
		peer.Output <- fmt.Sprintf(transaction)
	}
}

func (p *P2PNetwork) PrintAllNodeTransactions() {
	current := p.Head
	for current != nil {
		fmt.Printf("Node %d Transactions: %v\n", current.ID, current.Transactions)
		current = current.Next
	}
}

func (p *P2PNetwork) WaitForBroadcastCompletion() {
	// Wait for some time to ensure broadcasting is completed
	time.Sleep(5 * time.Second)

	// Print all node transactions after broadcasting
	p.PrintAllNodeTransactions()
}

func processTransactions(node *Node) {
	node.mutex.Lock()
	defer node.mutex.Unlock()

	uniqueTransactions := make(map[string]bool)
	var newTransactions []string

	for _, transaction := range node.Transactions {
		// Skip transactions starting with "REGISTER"
		if strings.HasPrefix(transaction, "REGISTER") {
			continue
		}

		// Skip repeating transactions
		if _, exists := uniqueTransactions[transaction]; !exists {
			uniqueTransactions[transaction] = true
			newTransactions = append(newTransactions, transaction)
		}
	}

	// Update node's transactions list
	node.Transactions = newTransactions
}

func p2pnetwork(blockchain *Blockchain) {
	p2pNetwork := NewP2PNetwork()

	bootstrapNode := NewNode(1, 5000, blockchain)
	p2pNetwork.AddNode(bootstrapNode)
	go bootstrapNode.StartNode()

	// Start nodes in separate goroutines
	node1 := NewNode(2, 8080, blockchain)
	node2 := NewNode(3, 8081, blockchain)
	node3 := NewNode(4, 8082, blockchain)

	// Add nodes to the P2P network before registering with the bootstrap node
	p2pNetwork.AddNode(node1)
	p2pNetwork.AddNode(node2)
	p2pNetwork.AddNode(node3)

	// Register nodes with the bootstrap node
	node1.RegisterWithBootstrap()
	node2.RegisterWithBootstrap()
	node3.RegisterWithBootstrap()

	// Start nodes after registering
	go node1.StartNode()
	go node2.StartNode()
	go node3.StartNode()

	p2pNetwork.DisplayNetworkInfo()

	bootstrapNode.PrintPeers()
	node1.PrintPeers()
	node2.PrintPeers()
	node3.PrintPeers()

	// Ask the user to send a message between node1 and node2 three times
	for i := 0; i < 3; i++ {
		fmt.Println("Enter a message to send from Node 1 to Node 2:")
		message, _ := AskInput("Message: ", bufio.NewReader(os.Stdin))
		go node1.SendMessage(node2.ID, message)
		currentNode := p2pNetwork.Head
		for currentNode != nil {
			go processTransactions(currentNode)
			currentNode = currentNode.Next
		}
		fmt.Println("BROADCASTING TO THE OTHER NODES")
		p2pNetwork.PrintAllNodeTransactions()
	}

	fmt.Println("MINING A NEW BLOCK")
	blockchain.AddTransactionToBlock(node1.Transactions[0])
	blockchain.AddTransactionToBlock(node1.Transactions[1])
	blockchain.AddTransactionToBlock(node1.Transactions[2])

	current := p2pNetwork.Head
	for current != nil {
		current.Transactions = nil
		current = current.Next
	}
	fmt.Println("Node transactions after block is mined")
	p2pNetwork.PrintAllNodeTransactions()

	// Process transactions for each node in the P2P network
	// Wait for user input to exit the program
	fmt.Println("Press Enter to exit.")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return
}

func (bc *Blockchain) AddBlock(block *Block) {
	currentBlock := bc.Head

	// Find the last block in the chain
	for currentBlock.NextBlock != nil {
		currentBlock = currentBlock.NextBlock
	}

	// Append the new block to the blockchain
	currentBlock.NextBlock = block
}

func (blockchain *Blockchain) AddTransactionToBlock(transaction string) {
	currentBlock := blockchain.Head

	// Find the last block in the chain
	for currentBlock.NextBlock != nil {
		currentBlock = currentBlock.NextBlock
	}

	// If the current block has fewer than 3 transactions, add the new transaction
	if len(currentBlock.Transactions) < 3 {
		currentBlock.Transactions = append(currentBlock.Transactions, transaction)
		currentBlock.Root = buildMerkleTree(currentBlock.Transactions)
		currentBlock.CalculateBlockHash()

	} else {
		// If the current block is full, create a new block
		newBlock := NewBlock([]string{transaction}, currentBlock)
		newBlock.Root = buildMerkleTree(newBlock.Transactions)
		newBlock.CalculateBlockHash()
		currentBlock.NextBlock = newBlock
	}
}

// func (blockchain *Blockchain) AddTransactionToBlock(transaction string) {
// 	currentBlock := blockchain.Head

// 	// Find the last block in the chain
// 	for currentBlock.NextBlock != nil {
// 		currentBlock = currentBlock.NextBlock
// 	}

// 	// If the current block has fewer than 3 transactions, add the new transaction
// 	if len(currentBlock.Transactions) == 3 {
// 		currentBlock.Transactions = append(currentBlock.Transactions, transaction)
// 		currentBlock.Root = buildMerkleTree(currentBlock.Transactions)
// 		currentBlock.CalculateBlockHash()
// 	} else {
// 		// If the current block is full, create a new block
// 		newBlock := NewBlock([]string{transaction}, currentBlock)
// 		currentBlock.NextBlock = newBlock
// 	}
// }

// func (n *Node) AddTransactionToBlock() {
// 	n.mutex.Lock()
// 	defer n.mutex.Unlock()

// 	// Check if the local block exists or is full
// 	if n.localblock == nil || len(n.localblock.Transactions) == 3 {
// 		// If the local block is full, create a new block and append it to the blockchain
// 		newBlock := NewBlock(n.Transactions, n.localblock)
// 		n.localblock = newBlock
// 		blockchain.AddBlock(newBlock)

// 		// Clear transactions of the node
// 		n.Transactions = nil
// 	}
// }

var (
	genesisBlock *Block
	blockchain   *Blockchain
	block2       *Block
	block3       *Block
)

func main() {

	// bootstrapNode := &Node{
	// 	ID:   nodeCount,
	// 	IP:   "127.0.0.1", // Replace with your bootstrap node's IP
	// 	Port: 5000,        // Replace with your bootstrap node's port
	// }

	var input string
	var input1 string
	genesisBlock := NewBlock([]string{"Genesis Transaction1", "Genesis Transaction2", "Genesis Transaction3"}, nil)
	blockchain := &Blockchain{Head: genesisBlock}

	block2 := NewBlock([]string{"Second Transaction1", "Second Transaction2", "Second Transaction3"}, genesisBlock)
	genesisBlock.NextBlock = block2

	block3 := NewBlock([]string{"Third Transaction1", "Third Transaction2", "Third Transaction3"}, block2)
	block2.NextBlock = block3

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Print a block + Merkle Tree")
		fmt.Println("2. Print the blockchain + Merkle trees of all blocks")
		fmt.Println("3. Verify a specific block")
		fmt.Println("4. Verify the blockchain")
		fmt.Println("5. Change a block")
		fmt.Println("6. Simulate a P2P network")
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
			fmt.Scanln(&input1)
			blocknumber, err := strconv.Atoi(input1)
			if err != nil {
				fmt.Println("Invalid input for block number. Please enter a valid integer.")
				continue
			}
			fmt.Print("Verifying Block: ")
			blockchain.printBlock(blocknumber)
			blockchain.Head.VerifyBlock()

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
		case "6":
			p2pnetwork(blockchain)
		case "0":
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}

}
