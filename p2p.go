package main

import (
	"fmt"
	"go-blockchain/utils"
	"math/rand"
)

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
	fmt.Printf("\n%sNode %v:%v added as neighbour of node %v:%v\n\n%s", utils.Green, neighbour.IPAddress, neighbour.Port, node.IPAddress, node.Port, utils.Reset)
	// randomly add the node as a neighbour to one of the neighbours of the node and vice versa
	randomIndex := rand.Intn(len(node.Neighbours))
	node.Neighbours[randomIndex].Neighbours = append(node.Neighbours[randomIndex].Neighbours, neighbour)
	neighbour.Neighbours = append(neighbour.Neighbours, node.Neighbours[randomIndex])
	fmt.Printf("\n%sNode %v:%v added as neighbour of node %v:%v\n\n%s", utils.Green, neighbour.IPAddress, neighbour.Port, node.Neighbours[randomIndex].IPAddress, node.Neighbours[randomIndex].Port, utils.Reset)
}

func displayNetwork(bootstrapNode Node) {
	fmt.Printf("\nBootstrap Node: %v:%v\n\n", bootstrapNode.IPAddress, bootstrapNode.Port)
	// loop through all neighbours
	for index, neighbour := range bootstrapNode.Neighbours {
		fmt.Printf("Neighbour %d: %v:%v\n", index+1, neighbour.IPAddress, neighbour.Port)
	}
	fmt.Println()
}
