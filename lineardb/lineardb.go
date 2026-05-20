package lineardb

import (
	// "databaseHandler/catalog"
	"databaseHandler/storage"
)

type List struct {
	Root *Node
	Last *Node
}

type Node struct {
	Value int
	Offset int64
	Size int
	Next *Node
}

func NewList() List {
	return List{Root: nil}
}

func ListAdd(value int, offset int64, size int, list *List){
	if list.Root == nil {
		list.Root = &Node{
			Value: value,
			Offset: offset,
			Size: size,
			Next: nil,
		}
		return
	}

	newNode := &Node{			
		Value: value,
		Offset: offset,
		Size: size,
		Next: nil,
	}

	list.Last.Next = newNode
	list.Last = newNode
}

func ListRead(value int, list *List) *Node {
	currentNode := list.Root
	for currentNode != nil {
		if currentNode.Value == value {
			return currentNode
		}
		currentNode = currentNode.Next
	}	
	return nil
}


