package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg = sync.WaitGroup{}

type Node struct {
	nb       int
	distance []int
}

func affichage_matrice(tab [][]int) {
	for i := 0; i < len(tab); i++ {
		for j := 0; j < len(tab[0]); j++ {
			fmt.Printf("%d ", tab[i][j])
		}
		fmt.Print("\n")
	}
}

func Dijkstra(graph [][]int, start int, wg *sync.WaitGroup, ch chan Node) {
	defer wg.Done()

	n := len(graph)
	dist := make([]int, n)
	for i := 0; i < n; i++ {
		dist[i] = 9999
	}
	dist[start] = 0

	unvisited := make(map[int]bool)
	for i := 0; i < n; i++ {
		unvisited[i] = true
	}

	current := start
	for current != -1 {
		unvisited[current] = false

		for i, val := range graph[current] {
			if val == 1 {
				newDistance := dist[current] + 1
				if newDistance < dist[i] {
					dist[i] = newDistance
				}
			}
		}
		current = -1
		for i, visited := range unvisited {
			if visited && (current == -1 || dist[i] < dist[current]) {
				current = i
			}
		}
	}
	ch <- Node{start, dist}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer ln.Close()

	fmt.Println("Listening on port 8080...")

	for {
		conn, err := ln.Accept() //wait for connections
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// Create a new file to save the uploaded file
	file, err := os.Create("uploaded.txt")
	if err != nil {
		fmt.Println("Error uploading file:", err.Error())
		return
	}
	defer file.Close()

	// Copy the received file to the new file
	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Println("Error copying file:", err.Error())
		return
	}
	fmt.Print("File received.")

	file, err = os.Open("uploaded.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file) //type of *bufio.Scanner
	var matrix [][]int

	for scanner.Scan() { //Scan used in a loop to read lines of text
		line := scanner.Text() //Text is used to print each line
		var l []int
		for _, valeur := range strings.Fields(line) {

			val, err := strconv.Atoi(valeur)
			if err != nil {
				log.Fatal(err)
			}
			l = append(l, val)
		}
		matrix = append(matrix, l)
	}

	affichage_matrice(matrix)
	ch := make(chan Node, 10)

	result := make([][]int, len(matrix), len(matrix[0]))

	for i := 0; i < len(matrix); i++ {
		wg.Add(1)
		go Dijkstra(matrix, i, &wg, ch)
	}

	for i := 0; i < len(matrix); i++ {
		node := <-ch
		result[node.nb] = node.distance
	}

	wg.Wait()

	affichage_matrice(result)

}
