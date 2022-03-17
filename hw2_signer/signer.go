package main

import (
	"fmt"
	"strings"
	"sync"
)

func main() {
	prepareWorkers()
}

func prepareWorkers() {
	wg := &sync.WaitGroup{}
	md5InChannel := make(chan int)
	crcInChannel := make(chan int)
	md5OutChannel := make(chan string)
	crcOutChannel := make(chan string)

	wg.Add(4)
	go generateData(md5InChannel, crcInChannel, wg)
	go SingleHash(md5InChannel, md5OutChannel, wg)
	go MultiHash(crcInChannel, crcOutChannel, wg)
	go concatHashes(wg, md5OutChannel, crcOutChannel)

	wg.Wait()
}

func generateData(md5Chanel chan int, crcChanel chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := 0; j <= 10; j++ {
		md5Chanel <- j
		crcChanel <- j
	}

	close(md5Chanel)
	close(crcChanel)
}

func SingleHash(in <-chan int, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range in {
		out <- DataSignerMd5(fmt.Sprint(data))
	}

	close(out)
}

func MultiHash(in <-chan int, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for data := range in {
		out <- DataSignerCrc32(fmt.Sprint(data))
	}

	close(out)
}

func concatHashes(wg *sync.WaitGroup, chanels ...chan string) {
	defer wg.Done()
	tmp := make([]string, len(chanels))

	for i, chanel := range chanels {
		for data := range chanel {
			tmp[i] = data
		}

		fmt.Println(strings.Join(tmp, "_"))
	}
}
