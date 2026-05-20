package tests

import (
	"testing"
	"os"
	"databaseHandler/tree"
	"databaseHandler/invariant"
)

func TestNewTree(t *testing.T) {
	tr := tree.NewTree(3)
	if tr.Root != nil {
		t.Errorf("expected root to be nil, got %v", tr.Root)
	}
	if tr.Degree != 3 {
		t.Errorf("expected degree 3, got %d", tr.Degree)
	}
}

func TestInsert(t *testing.T) {
	tr := tree.NewTree(3)
	tree.AddKey(5, 0, 0, &tr)
	tree.AddKey(3, 0, 0, &tr)
	tree.AddKey(7, 0, 0, &tr)
	tree.AddKey(1, 0, 0, &tr)
	tree.AddKey(19, 0, 0, &tr)
	tree.AddKey(10, 0, 0, &tr)
	tree.AddKey(2, 0, 0, &tr)
	invariant.CheckInvariants(t, &tr)
}

func TestSearch(t *testing.T) {
	tr := tree.NewTree(3)
	tree.AddKey(5, 0, 0, &tr)
	tree.AddKey(3, 0, 0, &tr)
	tree.AddKey(7, 0, 0, &tr)
	tree.AddKey(8, 0, 0, &tr)
	result := tree.Search(99, tr.Root)
	if result != nil {
		t.Errorf("expected nil for missing key, got %v", result)
	}
}

func TestDeletingSingleNode(t *testing.T) {
	tr := tree.NewTree(6)
	tree.AddKey(5, 0, 0, &tr)
	tree.AddKey(2, 0, 0, &tr)
	tree.AddKey(1, 0, 0, &tr)
	tree.AddKey(9, 0, 0, &tr)
	tree.Delete(1, &tr)
	invariant.CheckInvariants(t, &tr)
}

func TestLargeTree(t *testing.T) {
	tr := tree.NewTree(4)
	values := []int{10, 20, 5, 6, 12, 30, 7, 17, 3, 25, 40, 35, 15, 22, 8, 1, 50, 45, 28, 33}
	for _, v := range values {
		tree.AddKey(v, 0, 0, &tr)
	}
	tree.Delete(17, &tr)
	invariant.CheckInvariants(t, &tr)
}

func TestDeleteInternalNode(t *testing.T) {
	tr := tree.NewTree(4)
	values := []int{10, 20, 5, 6, 12, 30, 7, 17, 3, 25}
	for _, v := range values {
		tree.AddKey(v, 0, 0, &tr)
	}
	tree.Delete(10, &tr)
	invariant.CheckInvariants(t, &tr)
	result := tree.Search(10, tr.Root)
	if result != nil {
		t.Errorf("expected 10 to be deleted but found it")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tr := tree.NewTree(3)
	keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
	for i, k := range keys {
		tree.AddKey(k, int64(i*100), 50, &tr)
	}
	path := "test_tree.bin"
	defer os.Remove(path)
	if err := tree.SaveTree(&tr, path); err != nil {
		t.Fatalf("SaveTree failed: %v", err)
	}
	loaded, err := tree.LoadTree(path)
	if err != nil {
		t.Fatalf("LoadTree failed: %v", err)
	}
	for _, k := range keys {
		if tree.Search(k, loaded.Root) == nil {
			t.Errorf("key %d not found after load", k)
		}
	}
	if loaded.Degree != tr.Degree {
		t.Errorf("degree mismatch: got %d, want %d", loaded.Degree, tr.Degree)
	}
	invariant.CheckInvariants(t, loaded)
	if tree.Search(999, loaded.Root) != nil {
		t.Error("found key 999 which was never inserted")
	}
}

func TestSaveAndLoadLarge(t *testing.T) {
	tr := tree.NewTree(4)
	for i := 1; i <= 200; i++ {
		tree.AddKey(i, int64(i*100), 50, &tr)
	}
	path := "test_tree_large.bin"
	defer os.Remove(path)
	if err := tree.SaveTree(&tr, path); err != nil {
		t.Fatalf("SaveTree failed: %v", err)
	}
	loaded, err := tree.LoadTree(path)
	if err != nil {
		t.Fatalf("LoadTree failed: %v", err)
	}
	invariant.CheckInvariants(t, loaded)
	for _, k := range []int{1, 50, 100, 150, 200} {
		if tree.Search(k, loaded.Root) == nil {
			t.Errorf("key %d not found after load", k)
		}
	}
}
