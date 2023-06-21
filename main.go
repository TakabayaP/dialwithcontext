package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	addr := net.JoinHostPort("localhost", port)

	ctx := context.Background()
	// キャンセルするためのコンテキストを作成
	ctx, cancel := context.WithCancel(ctx)

	// 実装処理
	var conn net.Conn
	conn, err := DialWithContext("tcp", addr, ctx)
	if err != nil {
		defer cancel()
		return err
	}

	go func() {
		<-time.After(1 * time.Second)
		fmt.Println("cancelled")
		cancel()
	}()

	// go func() {
	// io.Copy(os.Stdout, conn)
	// }()

	io.Copy(os.Stdout, conn)
	_, err = io.Copy(conn, os.Stdin)
	fmt.Println(err)

	return nil
}

func DialWithContext(network, address string, ctx context.Context) (net.Conn, error) {
	// TODO: implement
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	_ = context.AfterFunc(ctx, func() {
		fmt.Println("closing")
		conn.Close()
		// conn.SetDeadline(time.Now())
	})
	return conn, nil
}
