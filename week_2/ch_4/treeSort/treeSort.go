package main

import "fmt"

type node struct {
	value int
	left  *node
	right *node
}

type bst struct {
	root *node
}

func (t *bst) insert(v int) {
	if t.root == nil {
		t.root = &node{v, nil, nil}
		return
	}
	current := t.root
	for {
		if v < current.value {
			if current.left == nil {
				current.left = &node{v, nil, nil}
				return
			}
			current = current.left
		} else {
			if current.right == nil {
				current.right = &node{v, nil, nil}
				return
			}
			current = current.right
		}
	}
}

func (t *bst) inorder(visit func(int)) {
	var traverse func(*node)
	traverse = func(current *node) {
		if current == nil {
			return
		}
		traverse(current.left)
		visit(current.value)
		traverse(current.right)
	}
	traverse(t.root)
}

func (t *bst) slice() []int {
	sliced := []int{}
	t.inorder(func(v int) {
		sliced = append(sliced, v)
	})
	return sliced
}

func treesort(values []int) []int {
	tree := bst{}
	for _, v := range values {
		tree.insert(v)
	}
	return tree.slice()
}

func main() {
	fmt.Println(treesort([]int{2, 4, 3, 1, 9, 7, 8}))
}
