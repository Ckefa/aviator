package models

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/Ckefa/aviator/internals/services"
)

type Flight struct {
	Odd   float32
	Count chan float32
	delay int
	Wait  sync.WaitGroup
	data  []float32
}

func NewFlight() *Flight {
	flight := &Flight{
		Odd:   1,
		Count: make(chan float32),
		delay: 8,
	}
	path := filepath.Join("models", "nums.txt")
	flight.LoadDataFromFile(path)
	return flight
}

func (f *Flight) LoadDataFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		println("File Not Found")
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Parse the line into a float64
		value, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			return err
		}
		// Append the value to the data slice
		f.data = append(f.data, float32(value))
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func generator(flight *Flight) {
	defer flight.Wait.Done()

	for {
		// Generate a random float64 and round to 2 decimal places
		randIndex := rand.Intn(len(flight.data))
		reach := flight.data[randIndex]

		fmt.Printf("Counting to %v\n", reach)

		for i := float32(1); i <= reach; i += 0.01 {
			flight.Count <- i
			time.Sleep(time.Millisecond * 150)
		}

		burst := map[string]interface{}{
			"name": "burst",
			"msg":  "flew away!",
		}

		go func() {
			services.BroadCastMsg(burst)
		}()

		time.Sleep(time.Second * time.Duration(flight.delay))

		start := map[string]interface{}{
			"name": "startrun",
			"msg":  "starting the run",
		}

		go func() {
			services.BroadCastMsg(start)
		}()

	}
}

func (flight *Flight) Run() {
	flight.Wait.Add(1)
	go generator(flight)

	go func() {
		for v := range flight.Count {
			data := map[string]interface{}{
				"name": "odd",
				"msg":  math.Floor(float64(v)*100) / 100,
			}

			go services.BroadCastMsg(data)
		}
	}()
}
