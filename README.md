# consistenthash

[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/xinst/consistenthash?status.svg)](https://godoc.org/github.com/xinst/consistenthash)
[![Go Report Card](https://goreportcard.com/badge/github.com/xinst/consistenthash)](https://goreportcard.com/report/github.com/xinst/consistenthash)
[![Build Status](https://travis-ci.org/xinst/consistenthash.svg?branch=master)](https://travis-ci.org/xinst/consistenthash)
![Version](https://img.shields.io/badge/version-1.0-brightgreen.svg)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/consistenthash/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link)


Implementing Consistent Hashing in Golang with Virtual Nodes

```go
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


```

#test result

nodes visit result:
node2:2045
node3:1
node4:2715
node5:2202
node1:3037
