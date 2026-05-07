package tests

import (
    "testing"
		"databaseHandler/tree"  
		"databaseHandler/invariant"
)

func TestNewTree(t *testing.T) {
    tree := tree.NewTree(3)

    if tree.Root != nil {
        t.Errorf("expected root to be nil, got %v", tree.Root)
    }
    if tree.Degree != 3 {
        t.Errorf("expected degree 3, got %d", tree.Degree)
    }
}

func TestInsert(t *testing.T) {

	tr := tree.NewTree(7)
	tree.AddKey(5, &tr)
	tree.AddKey(3, &tr)
	tree.AddKey(7, &tr)
	tree.AddKey(8, &tr)
	invariant.CheckInvariants(t, &tr)	
}

func TestInvariantCatchesBadLargeTree(t *testing.T) {
    tr := tree.NewTree(3)

    // manually build a broken tree
    leaf1 := &tree.Node{
        Keys:   []tree.Key{{Value: 1}, {Value: 2}},
        IsLeaf: true,
    }
    leaf2 := &tree.Node{
        Keys:   []tree.Key{{Value: 8}, {Value: 4}}, // unsorted, violates invariant
        IsLeaf: true,
    }
    leaf3 := &tree.Node{
        Keys:   []tree.Key{{Value: 5}, {Value: 6}, {Value: 7}}, // overfull for degree 3
        IsLeaf: true,
    }
    internal := &tree.Node{
        Keys:     []tree.Key{{Value: 3}, {Value: 5}},
        IsLeaf:   false,
        Children: []*tree.Node{leaf1, leaf2, leaf3},
    }
    tr.Root = &tree.Node{
        Keys:     []tree.Key{{Value: 10}},
        IsLeaf:   false,
        Children: []*tree.Node{internal},  // only 1 child for 1 key, should be 2
    }

    leaf1.Parent = internal
    leaf2.Parent = internal
    leaf3.Parent = internal
    internal.Parent = tr.Root

    invariant.CheckInvariants(t, &tr)
}
