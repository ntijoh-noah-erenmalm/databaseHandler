package tests

import (
    "testing"
		"databaseHandler/tree"  // replace with your actual module name
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
