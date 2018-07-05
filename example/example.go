package main

import (
    "fmt"
    "github.com/xinst/consistenthash"
    "strings"
)

func main() {
    constHash := consistenthash.NewConsistentHash(100)

    objects := make(map[string]int)
    objects["node1"] = 0
    objects["node2"] = 0
    objects["node3"] = 0
    objects["node4"] = 0

    for k := range objects {
        constHash.AddNode(k)
    }

    for index := 0; index < 10000; index++ {
        node := constHash.GetRandomSuitNode()
        if index == 984 {
            objects["node5"] = 0
            constHash.AddNode("node5")
        }
        if index == 10 {
            constHash.RemoveNode("node3")
        }
        if node != nil {
            names := strings.Split(node.Name(), "-")
            objects[names[0]]++
        }
    }
    fmt.Println("nodes visit result:")
    for k, v := range objects {
        fmt.Printf("%s:%d\n", k, v)
    }
}
