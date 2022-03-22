package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	iterations_count = 10

	md5_channel = 0
	crc_channel = 1
)

func main() {
	prepareWorkers()
}

func prepareWorkers() {
	wg := &sync.WaitGroup{}
	md5InChannel := make(chan int, 1)
	crcInChannel := make(chan int, 1)
	md5OutChannel := make(chan string, 1)
	crcOutChannel := make(chan string, 1)

	wg.Add(4)
	go generateData(md5InChannel, crcInChannel, wg)
	go SingleHash(md5InChannel, md5OutChannel, wg)
	go MultiHash(crcInChannel, crcOutChannel, wg)
	go CombineResults(md5OutChannel, crcOutChannel, wg)

	wg.Wait()
}

func generateData(md5Chanel chan int, crcChanel chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	i := iterations_count
	state := make([]bool, 2)

	for i > 0 {
		select {
		case md5Chanel <- i:
			if !state[md5_channel] {
				// md5Chanel <- i
				state[md5_channel] = true
				fmt.Printf("md5 gen %d\n", i)
			}

		case crcChanel <- i:
			if !state[crc_channel] {
				// crcChanel <- i
				state[crc_channel] = true
				fmt.Printf("crc gen %d\n", i)
			}

		default:
			if state[md5_channel] && state[crc_channel] {
				state[md5_channel] = false
				state[crc_channel] = false
				i--
			}
		}
	}

	close(md5Chanel)
	close(crcChanel)
}

func SingleHash(in <-chan int, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case val := <-in:
			fmt.Printf("md5 hash for %d\n", val)
			out <- DataSignerMd5(fmt.Sprint(val))
		default:
			time.Sleep(1 * time.Second)
		}
	}

	close(out)
}

func MultiHash(in <-chan int, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case val := <-in:
			fmt.Printf("crc hash for %d\n", val)
			out <- DataSignerMd5(fmt.Sprint(val))
		default:
			time.Sleep(1 * time.Second)
		}
	}

	close(out)
}

func CombineResults(md5Chanel chan string, crcChanel chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	i := iterations_count
	state := make([]string, 2)

	for i > 0 {
		select {
		case val := <-md5Chanel:
			state[md5_channel] = val
			fmt.Printf("md5 hash readed %s\n", val)

		case val := <-crcChanel:
			state[md5_channel] = val
			fmt.Printf("crc hash readed %s\n", val)

		default:
			if state[md5_channel] != `` && state[crc_channel] != `` {
				state[md5_channel] = ``
				state[md5_channel] = ``
				i--
				fmt.Println(state[md5_channel] + `_` + state[crc_channel])
			}
		}
	}
}
