package consistenthash

import (
    "crypto/rand"
    "hash/crc32"
    "sort"
    "sync"
)

// DefaultVirtualNodeCount is the virtual node numbers
const DefaultVirtualNodeCount = 200

// KeyType is the hash value of the node
type KeyType uint32

// KeyTypeSlice for sort interface
type KeyTypeSlice []KeyType

func (k KeyTypeSlice) Len() int           { return len(k) }
func (k KeyTypeSlice) Less(i, j int) bool { return k[i] < k[j] }
func (k KeyTypeSlice) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }

type sortedHashMap struct {
    keys KeyTypeSlice
    m    map[KeyType]INode
}

// Set to add the the node to list
func (s *sortedHashMap) Set(k KeyType, value INode) {
    if _, ok := s.m[k]; !ok {
        s.keys = append(s.keys, k)
        sort.Sort(s.keys)
    }
    s.m[k] = value
}

// LowerBound Returns an INode
// in the container whose key is not considered to go before k
// (i.e., either it is equivalent or goes after).
func (s *sortedHashMap) LowerBound(k KeyType) INode {
    if len(s.keys) == 0 || len(s.m) == 0 {
        return nil
    }
    for _, v := range s.keys {
        if v >= k {
            return s.m[v]
        }
    }
    first := s.keys[0]
    return s.m[first]
}

// Exist return true if k is exist
func (s *sortedHashMap) Exist(k KeyType) bool {
    if _, ok := s.m[k]; ok {
        return true
    }
    return false
}

// Delete the node if the k is exist
func (s *sortedHashMap) Delete(k KeyType) {
    if _, ok := s.m[k]; ok {
        delete(s.m, k)
        index := -1
        for i, v := range s.keys {
            if v == k {
                index = i
                break
            }
        }
        if index >= 0 {
            s.keys = append(s.keys[:index], s.keys[index+1:]...)
        }
    }
}

// HashFunc to hash the array byte to KeyType
type HashFunc func([]byte) KeyType

// HashCRC32 use the crc32 package
func HashCRC32(data []byte) KeyType {
    v := crc32.ChecksumIEEE(data)
    return KeyType(v)
}

// ConsistentHash is the nodes manager
type ConsistentHash struct {
    mtx        sync.Mutex
    nodes      *sortedHashMap
    HashF      HashFunc
    VNodeCount int
}

// NewConsistentHash return a consistenthash manager with default value if
// options is empty
func NewConsistentHash(defaultVNodeCount int,
    options ...func(*ConsistentHash)) *ConsistentHash {
    c := &ConsistentHash{
        nodes: &sortedHashMap{
            m: make(map[KeyType]INode),
        },
        VNodeCount: DefaultVirtualNodeCount,
        HashF:      HashCRC32,
    }

    for _, option := range options {
        option(c)
    }

    return c
}

// AddNode to the node list
// it will add numbers of virtual node
func (c *ConsistentHash) AddNode(nodeName string) {
    newNodes := make([]INode, c.VNodeCount)
    for i := 0; i < c.VNodeCount; i++ {
        newNodes[i] = NewVNode(nodeName, i)
    }

    c.mtx.Lock()
    defer c.mtx.Unlock()

    for _, node := range newNodes {
        hash := c.HashF([]byte(node.Name()))
        bRet := c.nodes.Exist(hash)
        if !bRet {
            c.nodes.Set(hash, node)
        }
    }
}

// RemoveNode from the node list
func (c *ConsistentHash) RemoveNode(nodeName string) {
    nodes := make([]INode, c.VNodeCount)
    for i := 0; i < c.VNodeCount; i++ {
        nodes[i] = NewVNode(nodeName, i)
    }

    c.mtx.Lock()
    defer c.mtx.Unlock()

    for _, node := range nodes {
        hash := c.HashF([]byte(node.Name()))
        c.nodes.Delete(hash)
    }
}

func randBytes(n int) []byte {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        for i := 0; i < n; i++ {
            b[i] = 'A'
        }
    }
    return b
}

// GetRandomSuitNode return a node or nil if the node list is empty
func (c *ConsistentHash) GetRandomSuitNode() INode {
    randData := randBytes(32)
    h := c.HashF(randData)

    return c.GetSuitNode(h)
}

// GetSuitNode return a node or nil if the node list is empty
func (c *ConsistentHash) GetSuitNode(hashValue KeyType) INode {
    c.mtx.Lock()
    defer c.mtx.Unlock()

    return c.nodes.LowerBound(hashValue)
}
