package main

import (
	"fmt"
	"sort"
	"crypto/sha1"
)


type hashRing struct {
	replicas int
	ring map[int]string
	sortedKeys []int
}

func NewHashRing(replicas int, nodes ...string) *hashRing {
	h := hashRing{
		replicas: replicas,
	}
	h.addNodes(nodes)
	return &h
}

func (h *hashRing) addNode(node string) {
	for i := 0; i < h.replicas; i ++ {
		key := h.genKey(fmt.Sprintf("%s:%s", node, i))
		h.ring[key] = node
		h.sortedKeys = append(h.sortedKeys, key)
	}
	sort.Ints(h.sortedKeys)
}

func (h *hashRing) addNodes(nodes []string) {
	for _, node := range nodes {
		h.addNode(node)
	}
}

func (h *hashRing) removeNode(node string) {
	for i := 0; i < h.replicas; i ++ {
		key := h.genKey(fmt.Sprintf("%s:%s", node, i))
		delete(h.ring, key)
		pos := sort.SearchInts(h.sortedKeys, key)
		h.sortedKeys = append(h.sortedKeys[:pos], h.sortedKeys[pos+1:]...)
	}
}

func (h *hashRing) removeNodes(nodes []string) {
	for _, node := range nodes {
		h.removeNode(node)
	}
}

func (h *hashRing) genKey(stringKey string) int{
	hash := sha1.New()
	hash.Write([]byte(stringKey))
	bs := hash.Sum(nil)
	return (int(bs[3]) << 24) | (int(bs[2]) << 16) | (int(bs[1]) << 8) | (int(bs[0]))
}

func (h *hashRing) getNode(stringKey string) *string {
	pos, key := h.getNodePos(stringKey)
	if pos == nil {
		return nil
	}
	node := h.ring[*key]
	return &node
}

func (h *hashRing) getNodePos(stringKey string) (*int, *int) {
	if len(h.ring) == 0 {
		return nil, nil
	}
	key := h.genKey(stringKey)
	pos := 0
	for i, nodeKey := range h.sortedKeys {
		if key > nodeKey {
			return &pos, &nodeKey
		}
		pos = i
		continue
	}
	return &pos, &h.sortedKeys[pos]
}