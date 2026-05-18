package tree

import (
    "encoding/gob"
    "os"
)

// serialNode mirrors Node but without the Parent pointer,
// so gob has no cycles to deal with.
type serialNode struct {
    Keys     []Key
    Children []*serialNode
    IsLeaf   bool
}

type serialTree struct {
    Root   *serialNode
    Degree int
}

// --- conversion helpers ---

func toSerial(n *Node) *serialNode {
    if n == nil {
        return nil
    }
    s := &serialNode{
        Keys:   n.Keys,
        IsLeaf: n.IsLeaf,
    }
    for _, child := range n.Children {
        s.Children = append(s.Children, toSerial(child))
    }
    return s
}

func fromSerial(s *serialNode, parent *Node) *Node {
    if s == nil {
        return nil
    }
    n := &Node{
        Keys:   s.Keys,
        IsLeaf: s.IsLeaf,
        Parent: parent, // rebuild the pointer here
    }
    for _, sc := range s.Children {
        n.Children = append(n.Children, fromSerial(sc, n))
    }
    return n
}

// --- public API ---

func SaveTree(tree *Tree, path string) error {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()

    st := serialTree{
        Root:   toSerial(tree.Root),
        Degree: tree.Degree,
    }
    return gob.NewEncoder(f).Encode(st)
}

func LoadTree(path string) (*Tree, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var st serialTree
    if err := gob.NewDecoder(f).Decode(&st); err != nil {
        return nil, err
    }

    return &Tree{
        Root:   fromSerial(st.Root, nil),
        Degree: st.Degree,
    }, nil
}