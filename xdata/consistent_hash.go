package xdata

import (
	"errors"
	"hash/fnv"
	"sort"
	"strconv"
	"sync"
)

const (
	defaultReplicas = 10
	salt            = "712&%BF^*(@"
)

type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// ConsistentHash 创建结构体保存一致性hash信息
type consistentHash struct {
	//hash环，key为哈希值，值存放节点的信息
	circle map[uint32]string
	//已经排序的节点hash切片
	sortedHashes units
	//虚拟节点个数，用来增加hash的平衡性
	virtualNode int
	//map 读写锁
	sync.RWMutex
	//hash 方法
	hash HashFunc
	//salt 盐
	salt string
}

type HashFunc func(data []byte) uint32

func defaultHash(data []byte) uint32 {
	f := fnv.New32()
	_, _ = f.Write(data)
	return f.Sum32()
}

//NewConsistent 创建一致性hash算法结构体，设置默认节点数量
func NewConsistent(replicas int, salt string, hf HashFunc) *consistentHash {
	h := consistentHash{
		//初始化变量
		circle: make(map[uint32]string),
		//设置虚拟节点个数
		virtualNode: replicas,
		hash:        hf,
		salt:        salt,
	}
	//设置虚拟节点个数
	if h.virtualNode <= 0 {
		h.virtualNode = defaultReplicas
	}
	if h.hash == nil {
		h.hash = defaultHash
	}
	if h.hash == nil {
		h.hash = defaultHash
	}
	return &h
}

//自动生成key值
func (c *consistentHash) generateKey(element string, index int) string {
	//副本key生成逻辑
	return c.salt + element + strconv.Itoa(index)
}

//获取hash位置 计算key 在hash环中对应的位置
func (c *consistentHash) hashKey(key string) uint32 {
	return c.hash([]byte(key))
}

//updateSortedHashes 更新排序，方便查找 因为后面我们使用的是sort.Search进行查找 sort.Search使用的是二分法进行查找，所以这里需要排序
func (c *consistentHash) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	//判断切片容量，是否过大，如果过大则重置
	if cap(c.sortedHashes)/(c.virtualNode*4) > len(c.circle) {
		hashes = nil
	}

	//添加hashes
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	//对所有节点hash值进行排序，
	//方便之后进行二分查找
	sort.Sort(hashes)
	//重新赋值
	c.sortedHashes = hashes
}

// Add 向hash环中添加节点
func (c *consistentHash) Add(element string) {
	//加锁
	c.Lock()
	//解锁
	defer c.Unlock()
	c.add(element)
}

//添加节点
func (c *consistentHash) add(element string) {
	//循环虚拟节点，设置副本
	for i := 0; i < c.virtualNode; i++ {
		//根据生成的节点添加到hash环中
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	//更新排序
	c.updateSortedHashes()
}

func (c *consistentHash) remove(element string) {
	for i := 0; i < c.virtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

func (c *consistentHash) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

func (c *consistentHash) search(key uint32) int {
	//查找算法
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//使用"二分查找"算法来搜索指定切片满足条件的最小值
	i := sort.Search(len(c.sortedHashes), f)
	//如果超出范围则设置i=0
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *consistentHash) Get(name string) (string, error) {
	//添加锁
	c.RLock()
	//解锁
	defer c.RUnlock()
	//如果为零则返回错误
	if len(c.circle) == 0 {
		return "", errors.New("hash环没有数据")
	}
	//计算hash值
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
