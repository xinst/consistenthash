package consistenthash

import (
    "fmt"
)

// INode ...
type INode interface {
    Name() string
}

// VNode is the virtual node
type VNode struct {
    name string
    id   int
}

// NewVNode create a new virtual node with name and the id
func NewVNode(name string, id int) *VNode {
    return &VNode{
        name: name,
        id:   id,
    }
}

// Name implenment the Name interface
// the Name will be used to hash for the node
func (vn *VNode) Name() string {
    return fmt.Sprintf("%s-%d", vn.name, vn.id)
}
