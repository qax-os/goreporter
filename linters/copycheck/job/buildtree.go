package job

import (
	"360.cn/apollo/apollo/gocode/copycode/suffixtree"
	"360.cn/apollo/apollo/gocode/copycode/syntax"
)

func BuildTree(schan chan []*syntax.Node) (t *suffixtree.STree, d *[]*syntax.Node, done chan bool) {
	t = suffixtree.New()
	data := make([]*syntax.Node, 0, 100)
	done = make(chan bool)
	go func() {
		for seq := range schan {
			data = append(data, seq...)
			for _, node := range seq {
				t.Update(node)
			}
		}
		done <- true
	}()
	return t, &data, done
}
