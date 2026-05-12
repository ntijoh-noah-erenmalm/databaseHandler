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

	node, _ := findKey(value, tree.Root)
	if node != nil {
		return
	}
	//1. make logic to find right node
	currentNode := findInsertionNode(value, tree.Root)

	//2. insert at right index
	index := findInsertionNodeIndex(value, currentNode.Keys)
	currentNode.Keys = append(currentNode.Keys, Key{}) //append empty key to current keys
	copy(currentNode.Keys[index+1:], currentNode.Keys[index:])
	currentNode.Keys[index] = Key{Value: value}

	//3. balance children
	if len(currentNode.Keys) >=tree.Degree {
			// split the node, then insert
			balanceInsertion(tree, currentNode)
	}
}

func findInsertionNodeIndex(value int, keys []Key) int {
	for i, k := range keys {
		if value < k.Value {
			return i
		}
	}
	return len(keys)
}

func findInsertionNode(value int, node *Node) *Node {

	if node.IsLeaf {
		return node
	}

	for i := 0 ; i < len(node.Keys); i++ {
		
		if value < node.Keys[i].Value {
			return findInsertionNode(value, node.Children[i])
		}
	}				
	return findInsertionNode(value, node.Children[len(node.Children)-1])
}

func balanceInsertion(tree *Tree, node *Node) {
	median := len(node.Keys)/2
	medianKey := node.Keys[median]
	leftKeys := make([]Key, median)
	rightKeys := make([]Key, len(node.Keys)-median-1)
	copy(leftKeys, node.Keys[:median])
	copy(rightKeys, node.Keys[median+1:])

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

			for _, child := range left.Children {
					child.Parent = left
			}
			for _, child := range right.Children {
					child.Parent = right
			}
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
	left.Parent = parent
	right.Parent = parent

	index := findInsertionNodeIndex(medianKey.Value, parent.Keys)
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
		balanceInsertion(tree, parent)
	}
}

func Search(value int, node *Node) *Key {

	for i := 0; i < len(node.Keys); i++ {
		keyValue := node.Keys[i].Value
		if value == keyValue  {
			return &node.Keys[i]
		}

		if value < keyValue && !node.IsLeaf {
			return Search(value, node.Children[i])
		}
	}
	if !node.IsLeaf {
		return Search(value, node.Children[len(node.Children)-1])
	}

	return nil
}

func Delete(value int, tree *Tree) {
    node, index := findKey(value, tree.Root)
    if node == nil {
        return
    }

    if !node.IsLeaf {
        // find in-order predecessor (rightmost key in left subtree)
        predecessor := node.Children[index]
        for !predecessor.IsLeaf {
            predecessor = predecessor.Children[len(predecessor.Children)-1]
        }
        // replace key with predecessor
        node.Keys[index] = predecessor.Keys[len(predecessor.Keys)-1]
        // delete predecessor from leaf
        predecessor.Keys = predecessor.Keys[:len(predecessor.Keys)-1]
        if len(predecessor.Keys) < (tree.Degree+1)/2-1 {
            balanceDeletion(predecessor, tree)
        }
        return
    }

    // leaf case — same as before
    node.Keys = append(node.Keys[:index], node.Keys[index+1:]...)
    if len(node.Keys) < (tree.Degree+1)/2-1 {
        balanceDeletion(node, tree)
    }
}

func findKey(value int, node *Node) (*Node, int) {
	for i, k := range node.Keys {
        if value == k.Value {
            return node, i
        }
        if value < k.Value {
            if node.IsLeaf {
                return nil, -1
            }
            return findKey(value, node.Children[i])
        }
    }
    if node.IsLeaf {
        return nil, -1
    }
    return findKey(value, node.Children[len(node.Children)-1])
}

func balanceDeletion(node *Node, tree *Tree) {
	if node == tree.Root {
        if len(node.Keys) == 0 && len(node.Children) > 0 {
            tree.Root = node.Children[0]
            tree.Root.Parent = nil
        }
        return
    }

		parent := node.Parent

    nodeIndex := 0
    for i, child := range parent.Children {
        if child == node {
            nodeIndex = i
            break
        }
    }

		minKeys := (tree.Degree+1)/2 - 1
		// try borrow from left sibling
		if nodeIndex > 0 {
				leftSibling := parent.Children[nodeIndex-1]
				if len(leftSibling.Keys) > minKeys {
						node.Keys = append([]Key{parent.Keys[nodeIndex-1]}, node.Keys...)
						parent.Keys[nodeIndex-1] = leftSibling.Keys[len(leftSibling.Keys)-1]
						leftSibling.Keys = leftSibling.Keys[:len(leftSibling.Keys)-1]
						return
				}
		}

		// try borrow from right sibling
		if nodeIndex < len(parent.Children)-1 {
				rightSibling := parent.Children[nodeIndex+1]
				if len(rightSibling.Keys) > minKeys {
						node.Keys = append(node.Keys, parent.Keys[nodeIndex])
						parent.Keys[nodeIndex] = rightSibling.Keys[0]
						rightSibling.Keys = rightSibling.Keys[1:]
						return
				}
		}
    // no siblings can lend — merge
    if nodeIndex > 0 {
        // merge with left sibling
        leftSibling := parent.Children[nodeIndex-1]
        // pull parent separator down
        leftSibling.Keys = append(leftSibling.Keys, parent.Keys[nodeIndex-1])
        // merge current node keys into left sibling
        leftSibling.Keys = append(leftSibling.Keys, node.Keys...)
        // move children if internal node
        leftSibling.Children = append(leftSibling.Children, node.Children...)
        // remove parent separator key
        parent.Keys = append(parent.Keys[:nodeIndex-1], parent.Keys[nodeIndex:]...)
        // remove current node from parent children
        parent.Children = append(parent.Children[:nodeIndex], parent.Children[nodeIndex+1:]...)
    } else {
        // merge with right sibling
        rightSibling := parent.Children[nodeIndex+1]
        // pull parent separator down
        node.Keys = append(node.Keys, parent.Keys[nodeIndex])
        // merge right sibling keys into current node
        node.Keys = append(node.Keys, rightSibling.Keys...)
        // move children if internal node
        node.Children = append(node.Children, rightSibling.Children...)
        // remove parent separator key
        parent.Keys = append(parent.Keys[:nodeIndex], parent.Keys[nodeIndex+1:]...)
        // remove right sibling from parent children
        parent.Children = append(parent.Children[:nodeIndex+1], parent.Children[nodeIndex+2:]...)
    }

    // if parent underflowed recurse upward
    if len(parent.Keys) < (tree.Degree+1)/2-1 {
			balanceDeletion(parent, tree)
		}	
}
