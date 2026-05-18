package tests

import (
    "testing"
    "os"
    "databaseHandler/tree"  
    "databaseHandler/invariant"
)

// ... your existing tests ...

func TestSaveAndLoad(t *testing.T) {
    // 1. build a known tree
    tr := tree.NewTree(3)
    keys := []int{10, 20, 5, 6, 12, 30, 7, 17}
    for _, k := range keys {
        tree.AddKey(k, &tr)
    }

    // 2. save it
    path := "test_tree.bin"
    defer os.Remove(path)
    if err := tree.SaveTree(&tr, path); err != nil {
        t.Fatalf("SaveTree failed: %v", err)
    }

    // 3. load it back
    loaded, err := tree.LoadTree(path)
    if err != nil {
        t.Fatalf("LoadTree failed: %v", err)
    }

    // 4. all keys must still be searchable
    for _, k := range keys {
        if tree.Search(k, loaded.Root) == nil {
            t.Errorf("key %d not found after load", k)
        }
    }

    // 5. degree must be preserved
    if loaded.Degree != tr.Degree {
        t.Errorf("degree mismatch: got %d, want %d", loaded.Degree, tr.Degree)
    }

    // 6. tree must still satisfy all b-tree invariants
    invariant.CheckInvariants(t, loaded)

    // 7. a key that was never inserted must not appear
    if tree.Search(999, loaded.Root) != nil {
        t.Error("found key 999 which was never inserted")
    }
}

func TestSaveAndLoadLarge(t *testing.T) {
    tr := tree.NewTree(4)
    for i := 1; i <= 200; i++ {
        tree.AddKey(i, &tr)
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