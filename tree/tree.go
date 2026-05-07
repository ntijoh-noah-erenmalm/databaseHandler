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
	currentNode := findInsertNode(value, tree.Root)

	//2. insert at right index
	index := findNodeIndex(value, currentNode.Keys)
	currentNode.Keys = append(currentNode.Keys, Key{}) //append empty key to current keys
	copy(currentNode.Keys[index+1:], currentNode.Keys[index:])
	currentNode.Keys[index] = Key{Value: value}

	//3. balance children
	if len(currentNode.Keys) >=tree.Degree {
			// split the node, then insert
			balance(tree, currentNode)
	}
}

func findNodeIndex(value int, keys []Key) int {
	for i, k := range keys {
		if value < k.Value {
			return i
		}
	}
	return len(keys)
}

func findInsertNode(value int, node *Node) *Node {

	if node.IsLeaf {
		return node
	}

	for i := 0 ; i < len(node.Keys); i++ {
		
		if value < node.Keys[i].Value {
			return findInsertNode(value, node.Children[i])
		}
	}				
	return findInsertNode(value, node.Children[len(node.Children)-1])
}

func balance(tree *Tree, node *Node) {
	median := len(node.Keys)/2
	medianKey := node.Keys[median]
	leftKeys := node.Keys[:median]
	rightKeys := node.Keys[median+1:]

	left := &Node{
		Keys: leftKeys,
		IsLeaf: node.IsLeaf,
		Parent: node.Parent,
	}
	right := &Node{
			Keys: rightKeys,
			IsLeaf: node.IsLeaf,
			Parent: node.Parent,
	}
	
	if !node.IsLeaf {
		left.Children = node.Children[:median+1]
		right.Children = node.Children[median+1:]
	}

	if node == tree.Root {
		newRoot := &Node{
				Keys: []Key{medianKey},
				IsLeaf: false,
				Parent: nil,
				Children: []*Node{left, right},
		}

		left.Parent = newRoot
		right.Parent = newRoot
		tree.Root = newRoot

		return
	}

	parent := node.Parent
	index := findNodeIndex(medianKey.Value, parent.Keys)
	parent.Keys = append(parent.Keys, Key{})
	copy(parent.Keys[index+1:], parent.Keys[index:])
	parent.Keys[index] = medianKey

	for i, child := range parent.Children {
		if child == node {
			parent.Children = append(parent.Children[:i], append([]*Node{left, right}, parent.Children[i+1:]...)...)
			break
		}
	}

	if len(parent.Keys) >= tree.Degree {
		balance(tree, parent)
	}
}
