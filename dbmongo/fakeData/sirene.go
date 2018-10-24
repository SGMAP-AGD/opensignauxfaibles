package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
)

func shuffleFile(fileName string) [][2]int64 {

	var p int64
	var op int64
	p = 0
	op = 0
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	var positions [][2]int64

	// avoid line title
	a, _, err := reader.ReadLine()
	p += int64(len(a)) + 2
	first := [][2]int64{{0, p}}
	op = p
	for {
		a, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		p += int64(len(a)) + 2
		positions = append(positions, [2]int64{op, p})
		op = p
	}

	for i := range positions {
		j := rand.Intn(i + 1)
		positions[i], positions[j] = positions[j], positions[i]
	}
	positions = append(first, positions...)

	return positions
}

func printShuffleFile(fileName string, positions *[][2]int64) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	ioreader := io.Seeker(file)
	for _, pos := range *positions {
		ioreader.Seek(pos[0], 0)
		buff := make([]byte, pos[1]-pos[0]-1)
		file.Read(buff)
		fmt.Println(string(buff))
	}
}
