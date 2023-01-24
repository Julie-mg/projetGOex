package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func send(conn net.Conn, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		return
	}
	defer file.Close()

	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("Error seeking file:", err.Error())
		return
	}
	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file:", err.Error())
		return
	}

	fmt.Println("File sent.")
}

func main() {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error dialing:", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	file, err := os.Open("matriceA.txt")
	if err != nil {
		log.Fatal(err)
	}

	send(conn, "matriceA.txt")

}
