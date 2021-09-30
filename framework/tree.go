package framework

import (
	"errors"
	"strings"
)

// 代表树结构
type Tree struct {
	root *node // 根节点
}

// 代表节点
type node struct {
	isLast   bool               // 该节点是否能成为一个独立的 url,是否自身就是一个终极节点
	segment  string             // url 中的字符串
	handlers []ControllerHandle // 节点中的控制器
	childs   []*node
}

// 初始化一个节点
func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

// 初始化一个根节点
func NewTree() *Tree {
	root := newNode()
	return &Tree{
		root: root,
	}
}

// 判断一个segment 是否是通用segment,即以:开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	// 如果segment是通配符，则所以下一层子节点都满足要求
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment {
			nodes = append(nodes, cnode)
		}
	}
	return nodes
}

// 判断路由是否已经在节点的所有子节点树中存在了
func (n *node) matchNode(url string) *node {
	// 使用分隔符将 url 分割为两个部分
	segments := strings.SplitN(url, "/", 2)
	// 第一个部分用于匹配下一层子节点
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	cnodes := n.filterChildNodes(segment)
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}
	// 如果只有一个segment,则是最后一个标记
	if len(segments) == 1 {
		// 如果 segment 是最后一个节点，判断这些 cnodes 是否有isLast标志
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
		// 都不是最后一个节点
		return nil
	}

	// 如果有 2个 segment,则递归每个子节点继续进行查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}
	return nil
}

// 增加路由节点, 路由节点有先后顺序
/*
/book/list
/book/:id (冲突)
/book/:id/name
/book/:student/age
/:user/name
/:user/name/:age (冲突)
*/
func (tree *Tree) AddRouter(url string, handler ...ControllerHandle) error {
	n := tree.root
	if n.matchNode(url) != nil {
		return errors.New("route exist: " + url)
	}

	segments := strings.Split(url, "/")
	// 对每个segment
	for index, segment := range segments {

		// 最终进入 node segment 的字段
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1

		var objNode *node // 标记是否有合适的子节点
		childNodes := n.filterChildNodes(segment)

		// 如果有匹配的子节点
		for _, cnode := range childNodes {
			if cnode.segment == segment {
				objNode = cnode
				break
			}
		}

		if objNode == nil {
			// 创建一个当前node的节点
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = handler
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode
	}

	return nil
}

// 根据 url 获取handler
func (tree *Tree) FindHandler(url string) []ControllerHandle {
	matchNode := tree.root.matchNode(url)
	if matchNode == nil {
		return nil
	}
	return matchNode.handlers
}
