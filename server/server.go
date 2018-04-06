package main

import (
	"net"
	"bufio"
	"fmt"
	"strings"
	"time"
	"math/rand"
	"os"
)

/*
编写一个 C/S Demo。服务器上存储有一个有序的 Key-Value 集合（类似于 c++ 中的 std::map<string, string>），
客户端通过网络 Scan 服务器中的数据，并对外提供迭代器风格的接口。

下面的代码展示了 client 大致的用法：
// Print up to 100 key-value pairs, starting from "abc".
auto iter = client.find(string("abc"));
for (int i = 0; i < 100 && iter.next(); i++) {
    cout << iter.key() << iter.value() << endl;
}

Next: 用来获取下一个kv pair
使用了next后还能输出当前的kv，为了少一次请求，我们存下当前kv和后一个kv，如果到了末尾，next返回false
如果数据量很大，我们似乎可以考虑一次发送接近MSS个字节的数据

提示：
1. 选择你熟悉的编程语言和网络通讯 framework/library。
2. 服务器中的数据直接用 BST 存就行了，不过设计协议时要考虑到数据量可能会很大。
3. 注意代码风格，添加必要的单元测试和文档。
4. 尝试优化性能。
 */

func main() {
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Printf("listen error: %s", err)
		return
	}

	var keys []string
	var values []string
	outputFile, _ := os.OpenFile("kvpairs.dat", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777) // 所有kv数据存储在这个文件
	defer outputFile.Close()
	outputWriter := bufio.NewWriter(outputFile)
	randomstartstr := "aIK" // 某个固定的key，用于简化测试
	root = avl_insert(root, randomstartstr, "testvalue", nil)

	for i := 0; i < 200; i++ {
		tmpkey := RandStringBytesMaskImprSrc(3) // 随机产生key，value只是附加几个字节
		tmpval := tmpkey + "'svalue"
		keys = append(keys, tmpkey)
		values = append(values, tmpval)
		root = avl_insert(root,  tmpkey, tmpval, nil)
	}

	fmt.Printf("build finished\n")
	avl_inorder(root, outputWriter)
	fmt.Println("")
	outputWriter.Flush()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accpet error: %s", err)
			return;
		}
		go HandleConn(conn) // 简单地开一个go routine用于处理请求
	}
}

/*
AVL tree implemented by someone
 */

type Node AVLTreeNode
type AVLTree *AVLTreeNode

type AVLTreeNode struct {
	key       string
	nodevalue string
	height    int
	left      *AVLTreeNode
	right     *AVLTreeNode
	parent    *AVLTreeNode
}

var root AVLTree = nil
/*
获取中序的下一个节点
 */
func getNext(avl AVLTree) AVLTree{
	if avl == nil {
		return nil
	}
	if avl.right != nil {
		var p = avl.right
		for p.left != nil {
			p = p.left
		}
		return p
	} else {
		for avl.parent != nil && avl == avl.parent.right {
			avl = avl.parent
		}
		return avl.parent
	}
}


/*

AVL 树的实现暂时使用某篇博客的代码，稍作了修改，添加了parent指针

 */
func highTree(p AVLTree) int {
	if p == nil {
		return -1
	} else {
		return p.height
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

/*Please look LL*/
func left_left_rotation(k AVLTree) AVLTree {
	var kl AVLTree
	var p AVLTree

	p = k.parent
	kl = k.left
	k.left = kl.right
	if kl.right != nil {
		kl.right.parent = k
	}
	kl.right = k
	k.parent = kl
	kl.parent = p

	k.height = max(highTree(k.left), highTree(k.right)) + 1
	kl.height = max(highTree(kl.left), k.height) + 1
	return kl
}

/*Please look RR*/
func right_right_rotation(k AVLTree) AVLTree {
	var kr AVLTree
	var p AVLTree
	kr = k.right
	p = k.parent
	k.right = kr.left
	if k.right != nil {
		k.right.parent = k
	}
	kr.left = k
	k.parent = kr
	kr.parent = p

	k.height = max(highTree(k.left), highTree(k.right)) + 1
	kr.height = max(k.height, highTree(kr.right)) + 1
	return kr
}

/*so easy*/
func left_righ_rotation(k AVLTree) AVLTree {
	k.left = right_right_rotation(k.left)
	return left_left_rotation(k)
}

func right_left_rotation(k AVLTree) AVLTree {
	k.right = left_left_rotation(k.right)
	return right_right_rotation(k)
}

func avl_insert(avl AVLTree, key string, nodeval string, par AVLTree) AVLTree {
	if avl == nil {
		avl = new(AVLTreeNode)
		if avl == nil {
			fmt.Println("avl tree create error!")
			return nil
		}else {
			avl.key = key
			avl.height = 0
			avl.left = nil
			avl.right = nil
			avl.nodevalue = nodeval
			avl.parent = par
		}
	} else if key < avl.key {
		avl.left = avl_insert(avl.left, key, nodeval, avl)
		if highTree(avl.left)-highTree(avl.right) == 2 {
			if key < avl.left.key { //LL
				avl = left_left_rotation(avl)
			} else { // LR
				avl = left_righ_rotation(avl)
			}
		}
	} else if key > avl.key {
		avl.right = avl_insert(avl.right, key, nodeval, avl)
		if (highTree(avl.right) - highTree(avl.left)) == 2 {
			if key < avl.right.key { // RL
				avl = right_left_rotation(avl)
			} else {
				fmt.Println("right right", key)
				avl = right_right_rotation(avl)
			}
		}
	} else if key == avl.key {
		fmt.Println("the key", key, "has existed!")
	}
	//notice: update height(may be this insert no rotation, so you should update height)
	avl.height = max(highTree(avl.left), highTree(avl.right)) + 1
	return avl
}

func avl_search(avl AVLTree, val string) AVLTree {
	if avl == nil {
		return nil
	}
	if val < avl.key {
		return avl_search(avl.left, val)
	} else if val > avl.key {
		return avl_search(avl.right, val)
	} else {
		return avl
	}

}

func avl_inorder(avl AVLTree, writer *bufio.Writer) {
	if avl == nil {
		return
	}
	avl_inorder(avl.left, writer)
	fmt.Printf("%s ", avl.key)
	writer.WriteString(avl.key + " ")
	avl_inorder(avl.right, writer)
}
/*
 go routine : handle conn
 */
func HandleConn(conn net.Conn) {
	nanos := time.Now().UnixNano()

	defer conn.Close()
	reader := bufio.NewReader(conn)
	filename := conn.RemoteAddr().String()
	outputFile, outputError := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return
	}
	defer outputFile.Close()
	for {
		mesg, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		//fmt.Println("receive :" + string(mesg))
		mesg = strings.TrimSuffix(mesg, "\n")

		now := avl_search(root, mesg)
		nextnode := getNext(now)

		nowret := ""
		nextret := ""

		if now != nil {
			nowret = now.nodevalue
			outputFile.WriteString(mesg+" " + nowret + " ")
		}
		if nextnode != nil {
			nextret = nextnode.key
		}
		b := []byte(nowret+","+nextret+"\n")
		conn.Write(b)
		//fmt.Printf("write %s\n", b[:])
	}
	after := time.Now().UnixNano()
	fmt.Println((after-nanos)/1000000)
	fmt.Printf("ms\n")
}

/*
helper function:
random string generator written by someone for random test
 */
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	paircnt = 100
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}
