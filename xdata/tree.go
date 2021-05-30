package xdata

import "errors"

const (
	BPTreePageType      = 1 // 叶子页
	BPTreePageTypeIndex = 2 // 索引页
)

var (
	NotFoundError = errors.New("NotFound")
)

type BPTree struct {
	Order    int
	RootPage *BPTreePage
	LeafPage *BPTreePage
}

type BPTreePage struct {
	Type       int
	Pre, Next  *BPTreePage
	ParentNode *BPTreeNode
	HeadNode   *BPTreeNode
	Len        int
	Cap        int
}

type BPTreeNode struct {
	Index     int
	Data      interface{}
	Next      *BPTreeNode
	Page      *BPTreePage
	ChildPage *BPTreePage
}

type Response struct {
	Code int
	Data interface{}
}

//
// NewBPTree
// @Description: B+树实现，B+树的索引节点在下层节点中，每一个索引节点对应下层的一个子分页，这与B树不同，B树是左右子分页的结构
// @param order int
// @return BPTree
//
func NewBPTree(order int) BPTree {
	return BPTree{
		Order:    order,
		LeafPage: &BPTreePage{Cap: order - 1, Type: BPTreePageType},
		RootPage: &BPTreePage{Cap: order - 1, Type: BPTreePageTypeIndex},
	}
}

func (bp *BPTree) Insert(index int, data interface{}) {
	page := bp.LeafPage
	if bp.RootPage.Len == 0 {
		parentNode := &BPTreeNode{Index: index, Page: bp.RootPage, ChildPage: page}
		page.ParentNode = parentNode
		bp.RootPage.Insert(parentNode)
	} else {
		page = bp.searchLeafPage(index, bp.RootPage)
	}
	node := &BPTreeNode{Index: index, Data: data}
	bp.addNodeToPage(node, page)
}

func (bp *BPTree) searchLeafPage(index int, page *BPTreePage) *BPTreePage {
	if page.Type == BPTreePageType {
		return page
	}
	cur := page.HeadNode
	for cur != nil && cur.Next != nil && cur.Next.Index <= index {
		cur = cur.Next
	}
	childPage := cur.ChildPage
	return bp.searchLeafPage(index, childPage)
}

func (bp *BPTree) searchNode(index int) *BPTreeNode {
	page := bp.LeafPage
	if bp.RootPage.Len > 0 {
		page = bp.searchLeafPage(index, bp.RootPage)
	}
	if page == nil {
		return nil
	}
	cur := page.HeadNode
	for cur != nil {
		if cur.Index == index {
			return cur
		}
		cur = cur.Next
	}
	return nil
}

func (page *BPTreePage) Insert(node *BPTreeNode) {
	cur := page.HeadNode
	if cur == nil || cur.Index > node.Index {
		node.Next = cur
		page.HeadNode = node
	} else {
		for cur.Next != nil && cur.Next.Index < node.Index {
			cur = cur.Next
		}
		node.Next = cur.Next
		cur.Next = node
	}
	node.Page = page
	page.Len++
}

func (page *BPTreePage) Split() (*BPTreePage, *BPTreePage) {
	// 找中心
	cur := page.HeadNode
	n := 1
	for n < (page.Len+1)/2 {
		cur = cur.Next
		n++
	}

	newHeadNode := cur.Next
	cur.Next = nil // 切断原链表

	// 新page
	newPage := &BPTreePage{
		Type:     page.Type,
		Pre:      page,
		Next:     page.Next,
		HeadNode: newHeadNode,
		Cap:      page.Cap,
		Len:      page.Len - n,
	}

	// 调整原page
	page.Next = newPage
	page.Len = n

	// 调整new page节点归属
	for newHeadNode.Next != nil {
		newHeadNode.Page = newPage
		newHeadNode = newHeadNode.Next
	}

	return page, newPage
}

func (page *BPTreePage) Up() {
	if page.ParentNode != nil && page.ParentNode.Index != page.HeadNode.Index {
		page.ParentNode.Index = page.HeadNode.Index
		page.ParentNode.Page.Up()
	}
}

func (page *BPTreePage) DeleteByIndex(index int) *BPTreeNode {
	var pre *BPTreeNode
	cur := page.HeadNode
	for cur != nil {
		if cur.Index == index {
			break
		}
		pre = cur
		cur = cur.Next
	}
	if cur == nil || cur.Index != index {
		return nil
	}
	if pre == nil { // 头部
		page.HeadNode = cur.Next
	} else {
		pre.Next = cur.Next
	}
	cur.Next = nil
	cur.Page = nil
	page.Len--
	if page.Len == 0 { // 删光了，去掉空页
		n := page.Next
		p := page.Pre
		if p != nil {
			p.Next = n
		}
		if n != nil {
			n.Pre = p
		}
		page.Next = nil
		page.Pre = nil
		if page.ParentNode != nil { // 向上递归删除
			page.ParentNode.Page.DeleteByIndex(index)
		}
		return cur
	}
	page.Up() // 向上递归检查
	return cur
}

func (bp *BPTree) addNodeToPage(node *BPTreeNode, page *BPTreePage) {
	// 加入page
	page.Insert(node)
	page.Up()

	// 页面未满
	if page.Len <= page.Cap {
		return
	}

	// 页满，需要拆分
	leftPage, rightPage := page.Split()
	// 新页上层node
	rightIndexNode := &BPTreeNode{
		Index:     rightPage.HeadNode.Index,
		ChildPage: rightPage,
	}
	rightPage.ParentNode = rightIndexNode

	if leftPage.ParentNode == nil { // 顶层
		bp.RootPage = &BPTreePage{Cap: bp.Order - 1, Type: BPTreePageTypeIndex} // 当前可能是root，需要重置
		leftIndexNode := &BPTreeNode{
			Index:     leftPage.HeadNode.Index,
			ChildPage: leftPage,
		}
		leftPage.ParentNode = leftIndexNode
		bp.addNodeToPage(leftIndexNode, bp.RootPage)
	}
	// 递归向上
	bp.addNodeToPage(rightIndexNode, leftPage.ParentNode.Page)
}

func (bp *BPTree) Delete(index int) (data interface{}, err error) {
	page := bp.LeafPage
	if bp.RootPage.Len > 0 {
		page = bp.searchLeafPage(index, bp.RootPage)
	}
	if page == nil {
		return nil, NotFoundError
	}
	node := page.DeleteByIndex(index)
	if node == nil {
		return nil, NotFoundError
	}
	return node.Data, nil
}

func (bp *BPTree) Search(index int) (data interface{}, err error) {
	node := bp.searchNode(index)
	if node == nil {
		return nil, NotFoundError
	}
	return node.Data, nil
}
