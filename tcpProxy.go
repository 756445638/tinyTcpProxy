package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	local  = flag.String("local", "", "local port")
	remote = flag.String("remote", "", "remote ip addr")
)

func main() {
	flag.Parse()
	if *local == "" || *remote == "" {
		fmt.Printf("Usage:%v -local LOCAL -remote REMOTE", os.Args[0])
		os.Exit(1)
	}

	ln, err := net.Listen("tcp4", *local)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accep conntion failed,err:%v\n", err)
			continue
		}
		go func() {
			defer conn.Close()
			fmt.Printf("accept conn,remote addr:%v\n", conn.RemoteAddr())
			err = serverConn(conn)
			if err != nil {
				fmt.Printf("server conn failed,err:%v\n", err)
			}
		}()

	}

}

func serverConn(conn net.Conn) error {
	//trying to connect to backend server
	remoteconn, err := net.DialTimeout("tcp4", *remote, time.Second*5)
	if err != nil {
		return err
	}
	defer remoteconn.Close()
	ok := true
	go func() {
		defer func() {
			ok = false
		}()
		for ok  {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("read for local conn failed,err:%v\n", err)
				return
			}

			n, err = remoteconn.Write(buf[0:n])
			if err != nil {
				fmt.Printf("write to remote conn failed,err:%v\n", err)
			}
		}

	}()
	defer func(){
		ok = false
	}()
	buf := make([]byte, 1024)
	for ok {
		n, err := remoteconn.Read(buf)
		if err != nil {
			fmt.Printf("read for remote conn failed,err:%v\n", err)
			return err
		}
		n, err = conn.Write(buf[0:n])
		if err != nil {
			fmt.Printf("write to local conn failed,err:%v\n", err)
			return err
		}
	}
	return nil
}
