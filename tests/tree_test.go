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

	tr := tree.NewTree(3)
	tree.AddKey(5, &tr)
	tree.AddKey(3, &tr)
	tree.AddKey(7, &tr)
	tree.AddKey(1, &tr)
	tree.AddKey(19, &tr)
	tree.AddKey(10, &tr)
	tree.AddKey(2, &tr)
	invariant.CheckInvariants(t, &tr)	
}

func TestSearch(t *testing.T) {

	tr := tree.NewTree(3)
	tree.AddKey(5, &tr)
	tree.AddKey(3, &tr)
	tree.AddKey(7, &tr)
	tree.AddKey(8, &tr)

	result := tree.Search(99, tr.Root)

	if result != nil {
		t.Errorf("expected nil for missing key, got %v", result)
	}
}

func TestDeletingSingleNode(t *testing.T) {
	tr := tree.NewTree(6)
	tree.AddKey(5, &tr)
	tree.AddKey(2, &tr)
	tree.AddKey(1, &tr)
	tree.AddKey(9, &tr)
  tree.Delete(1, &tr)
	invariant.CheckInvariants(t, &tr)
}

func TestLargeTree(t *testing.T) {
    tr := tree.NewTree(4)

    values := []int{10, 20, 5, 6, 12, 30, 7, 17, 3, 25, 40, 35, 15, 22, 8, 1, 50, 45, 28, 33}
    for _, v := range values {
        tree.AddKey(v, &tr)
    }

    tree.Delete(17, &tr)

    invariant.CheckInvariants(t, &tr)
}

func TestDeleteInternalNode(t *testing.T) {
    tr := tree.NewTree(4)
    values := []int{10, 20, 5, 6, 12, 30, 7, 17, 3, 25}
    for _, v := range values {
        tree.AddKey(v, &tr)
    }

    // delete a key that is likely in an internal node
    tree.Delete(10, &tr)

    invariant.CheckInvariants(t, &tr)

    result := tree.Search(10, tr.Root)
    if result != nil {
        t.Errorf("expected 10 to be deleted but found it")
    }
}
