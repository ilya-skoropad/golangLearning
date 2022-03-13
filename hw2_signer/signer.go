package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
)

func main() {
	prepareWorkers()
}

func prepareWorkers() {
	wg := &sync.WaitGroup{}
	dataChanel := make(chan int)

	wg.Add(2)
	go generateData(dataChanel, wg)
	go ProcessData(dataChanel, wg)

	wg.Wait()
}

func generateData(dataChanel chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i <= 10; i++ {
		dataChanel <- i
	}
	close(dataChanel)
}

func ProcessData(dataChanel <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range dataChanel {
		fmt.Printf("%x\n", GetMD5Hash(fmt.Sprint(data)))
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
