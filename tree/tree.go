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

	if tree.Root == nil {
      tree.Root = &Node{
          Keys:   []Key{{Value: value}},
          IsLeaf: true,
      }
        return
    }

	//1. make logic to find right node
	currentNode := findNode(value, tree.Root)
	
	//2. if full make balance change
	if len(currentNode.Keys) >=tree.Degree-1 {
		// split the node, then insert
		return
	}


	//3. insert at right index
	index := findNodeIndex(value, currentNode.Keys)
	currentNode.Keys = append(currentNode.Keys, Key{}) //append empty key to current keys
	copy(currentNode.Keys[index+1:], currentNode.Keys[index:])
	currentNode.Keys[index] = Key{Value: value}


	//4. balance children
}

func findNodeIndex(value int, keys []Key) int {
	for i, k := range keys {
		if value < k.Value {
			return i
		}
	}
	return len(keys)
}

func findNode(value int, node *Node) *Node {

	if node.IsLeaf {
		return node
	}

	for i := 0 ; i < len(node.Keys); i++ {
		
		if value < node.Keys[i].Value {
			return findNode(value, node.Children[i])
		}
	}				
	return findNode(value, node.Children[len(node.Children)-1])
}


