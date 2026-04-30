package tree

type Node struct {
	Keys []Key
	Children []*Node
	Parent *Node
	IsLeaf bool
}

type Key struct {
	Value int
}

type Tree struct {
	Root *Node
	Degree int
}


func NewTree(degree int) Tree {
	return Tree{Root: nil, Degree: degree}
}

func AddKey(value int, tree *Tree) {

	currentNode := tree.Root
	//make logic to find right node
	
	//if full make balance change


	index := findNodeIndex(value, currentNode.Keys)
	//insert at right index
	currentNode.Keys = append(currentNode.Keys, Key{}) //append empty key to current keys
	copy(currentNode.Keys[index+1], currentNode.Keys[index:])
	currentNode.Keys[index] = Key{Value: value}


	//balance children
}

func findNodeIndex(value int, keys []Key) int {
	for i, k := range keys {
		if value < k.Value {
			return i
		}
	}
	return len(keys)
}
