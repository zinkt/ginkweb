package gink

import (
	"strings"
)

type node struct {
	pattern  string  //到本届点处，的待匹配路由，如/p/:lang
	part     string  //即本节点对应的part，路由中的一部分，如:lang
	children []*node //子节点，如[doc, tutorial, intro]
	isWild   bool    //是否参数匹配（:或*为true）
}

// 辅助函数
// 返回n.children中第一个成功匹配的节点，用于插入新part
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 辅助函数
// 返回n.children中所有与part匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 添加路由时调用，递归插入节点，用于构建Trie
// parts是已被分割的pattern
// height是递归插入的“进度”
// 尝试递归地将parts[height:len(parts)]的part插入到节点n的子孙节点中
// 并将叶子节点的pattern赋值为最终的待匹配路由
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { // 所有part都插入完成
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil { // 如果没有儿子节点匹配成功，则新建节点并将其添加到子孙节点组中
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 尝试从节点n开始，与parts[height]匹配，返回成功匹配的节点
// 递归地沿子孙节点匹配，直到
// 1.匹配到*，2.匹配到len(parts)层，即匹配完成，
// 3.遍历完子孙节点，匹配失败，或当前节点的值为""，匹配失败
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
