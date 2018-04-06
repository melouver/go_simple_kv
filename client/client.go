package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)
type Client struct {}

func (cli *Client)find(start string) Iter {
	it := Iter{"", "", strings.TrimSuffix(start, "\n"), nil}
	return it
}

type Iter struct {
	k string
	v string
	k2 string
	conn *net.TCPConn
}

func (it *Iter)key() string{
	return it.k
}

func (it *Iter)value() string {
	return it.v
}

func (it *Iter)next() bool{
	it.k = it.k2
	if it.k == "" {
		return false //到了最后一个，无next，返回false
	}
	b := []byte(it.k+"\n")
	_, err := it.conn.Write(b) //写一个key
	if err != nil {
		return false
	}
	fmt.Printf("write %s\n", b[:])
	input := bufio.NewReader(it.conn)

	text, _ := input.ReadString('\n')
	text = strings.TrimSuffix(text, "\n") //读value+下一个key
	fmt.Printf("receive %s\n", text)

	kvs := strings.Split(text, ",")
	fmt.Printf("kv1: |%s| kv2: |%s|\n", kvs[0], kvs[1])
	if kvs[0] == "" { //无对应value
		fmt.Printf("no value, so return\n")
		return false
	}

	if kvs[1] == "" {//到了最后一个kv,再次尝试next会因为key为空而返回false
		fmt.Printf("value exists, no next return\n")
		it.v = kvs[0]
		it.k2 = ""
		return true
	}
	fmt.Printf("has value, has next return\n")
	it.v = kvs[0]
	it.k2 = kvs[1]
	return true
}



func main() {
	CalledByMain()
}

func CalledByMain() string {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8888")


	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("dial error")
		return ""
	}
	fmt.Println("connected!")

	s := onMessageRecived(conn)

	return s
}

func onMessageRecived(conn *net.TCPConn) string {
	cli := Client{}
	inputFile, inputError := os.Open("input.dat")
	if inputError != nil {
		fmt.Printf("An error occurred on opening the inputfile\n")
		return "" // exit the function on error
	}
	defer inputFile.Close()
	stinput := bufio.NewReader(inputFile)

	text, err := stinput.ReadString('\n')
	if err != nil {
		return ""
	}

	text = strings.TrimSuffix(text, "\n")
	it := cli.find(text)
	it.conn = conn
	filename := conn.LocalAddr().String()

	outputFile, outputError := os.OpenFile(filename,
		os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return ""
	}
	defer outputFile.Close()
	for i := 0; i < 100 && it.next(); i++ {
		fmt.Printf("--------------key:%s value:%s\n", it.key(), it.value())
		outputFile.WriteString(it.key() + " " + it.value() + " ")
	}
	return filename
}