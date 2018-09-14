package merkletree

import (
	"encoding/hex"
	"fmt"

	"github.com/fatih/color"
)

//
// this file is only for debug purposes, to print the merkle tree in the console
//

func (mt *MerkleTree) printLevel(parent Hash, iLevel int, maxLevel int) {
	for i := mt.numLevels - 1; i > iLevel; i-- {
		fmt.Print("	")
	}
	fmt.Print("level ")
	fmt.Print(iLevel)
	fmt.Print(" - ")
	fmt.Print("'" + parent.Hex() + "' = ")
	nodeType, _, nodeBytes, err := mt.storage.Get(parent)
	if err != nil {
		color.Red(err.Error())
	}
	var node treeNode
	if nodeType == byte(normalNodeType) {
		node = parseNodeBytes(nodeBytes)
		color.Blue("'" + node.ChildL.Hex() + "' - '" + node.ChildR.Hex() + "'")
	} else if nodeType == byte(valueNodeType) {
		// claim := core.ParseClaimDefaultBytes(nodeBytes)
		color.Green("value")
	} else if nodeType == byte(finalNodeType) { //typ==FINAL_NODE
		// claim := core.ParseClaimDefaultBytes(nodeBytes)
		fmt.Print("[FinalTree]:")
		color.Cyan("final tree node: " + HashBytes(nodeBytes).Hex())
		_, _, leafNodeBytes, err := mt.storage.Get(HashBytes(nodeBytes))
		if err != nil {
			color.Red(err.Error())
		}
		for i := mt.numLevels - 1; i > iLevel; i-- {
			fmt.Print("	")
		}
		color.Cyan("					leaf value: 0x" + hex.EncodeToString(leafNodeBytes))
	} else {
		//EMPTY_NODE
		fmt.Print("[EmptyBranch]:")
		fmt.Println(EmptyNodeValue.Bytes())
	}
	iLevel--
	if len(node.ChildR) > 0 && iLevel > maxLevel && nodeType != byte(EmptyNodeType) && nodeType != byte(finalNodeType) {
		mt.printLevel(node.ChildL, iLevel, maxLevel)
		mt.printLevel(node.ChildR, iLevel, maxLevel)
	}
}

// PrintFullMT prints the tree in the terminal, all the levels with all the nodes
func (mt *MerkleTree) PrintFullMT() {
	mt.printLevel(mt.root, mt.numLevels-1, mt.numLevels-1-1000)
	fmt.Print("root: ")
	color.Yellow(mt.Root().Hex())
}

// PrintLevelsMT prints the tree in the terminal until a specified depth
func (mt *MerkleTree) PrintLevelsMT(maxLevel int) {
	mt.printLevel(mt.root, mt.numLevels-1, mt.numLevels-1-maxLevel)
	fmt.Print("root: ")
	color.Yellow(mt.Root().Hex())
}
