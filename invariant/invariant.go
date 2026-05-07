package invariant

import (
	"testing"
	"databaseHandler/tree"

)

func CheckInvariants(t *testing.T, tr *tree.Tree) {
    if tr.Root == nil {
        return
    }
    checkNode(t, tr.Root, tr.Degree, true)
    checkLeafDepth(t, tr.Root, 0, -1)
}

func checkNode(t *testing.T, node *tree.Node, degree int, isRoot bool) {
    // keys are sorted
    for i := 1; i < len(node.Keys); i++ {
        if node.Keys[i].Value < node.Keys[i-1].Value {
            t.Errorf("keys not sorted: %v", node.Keys)
        }
    }

    // max keys
    if len(node.Keys) >= degree {
        t.Errorf("node overfull: has %d keys, max is %d", len(node.Keys), degree-1)
    }

    // min keys (root is exempt)
    minKeys := (degree / 2) - 1
    if !isRoot && len(node.Keys) < minKeys {
        t.Errorf("node underfull: has %d keys, min is %d", len(node.Keys), minKeys)
    }

    // children count must be keys+1 for internal nodes
    if !node.IsLeaf {
        if len(node.Children) != len(node.Keys)+1 {
            t.Errorf("internal node has %d keys but %d children, expected %d", len(node.Keys), len(node.Children), len(node.Keys)+1)
        }
        // check keys are consistent with children
        for i, child := range node.Children {
            if i < len(node.Keys) && len(child.Keys) > 0 {
                if child.Keys[len(child.Keys)-1].Value >= node.Keys[i].Value {
                    t.Errorf("child max key %d violates parent key %d", child.Keys[len(child.Keys)-1].Value, node.Keys[i].Value)
                }
            }
            checkNode(t, child, degree, false)
        }
    }
}

// all leaves must be at the same depth
func checkLeafDepth(t *testing.T, node *tree.Node, depth int, leafDepth int) int {
    if node.IsLeaf {
        if leafDepth == -1 {
            return depth
        }
        if depth != leafDepth {
            t.Errorf("leaf at depth %d, expected %d", depth, leafDepth)
        }
        return leafDepth
    }
    for _, child := range node.Children {
        leafDepth = checkLeafDepth(t, child, depth+1, leafDepth)
    }
    return leafDepth
}
