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
