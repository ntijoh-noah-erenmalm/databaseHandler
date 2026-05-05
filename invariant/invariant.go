package invariant

import (
	"testing"
	"databaseHandler/tree"

)

func CheckInvariants(t *testing.T, tree *tree.Tree) {
	if tree.Root == nil {
		return 
	}
	checkNode(t, tree.Root, tree.Degree)
}

func checkNode(t *testing.T, node *tree.Node, degree int) {
	
	for i := 1 ; i < len(node.Keys); i++ {
		if node.Keys[i].Value < node.Keys[i-1].Value {
			t.Errorf("keys not sorted: %v", node.Keys)
		}
	}	

	if len(node.Keys) >= degree {
		t.Errorf("node has too many keys: %d, nax %d", len(node.Keys), degree-1)
	}

	for _, child := range node.Children {
		checkNode(t, child, degree)
	}
}
