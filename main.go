package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//
//Напишите код, в котором имеются два канала сообщений из целых чисел так,
//чтобы приём сообщений всегда приводил к блокировке.
//Приёмом сообщений из обоих каналов будет заниматься главная горутина.
//Сделайте так, чтобы во время такого «бесконечного ожидания» сообщений выполнялась
//фоновая работа в виде вывода текущего времени в консоль.
//
//В качестве ответа приложите архивный файл с кодом программы из Задания 17.6.2.

func intGenerator(control <-chan int) <-chan int {
	c := make(chan int)
	immer := true
	go func() {
		defer close(c)
		for immer {
			i := rand.Intn(10000)
			time.Sleep(time.Duration(i) * time.Millisecond)
			c <- i
		}
	}()
	go func() {
		<-control
		immer = false
	}()

	return c
}

func control() <-chan int {
	c := make(chan int)
	immer := true
	go func() {
		defer close(c)
		for immer {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if strings.ToLower(input) == "exit" {
				immer = false
			}
		}
	}()

	return c
}

func main() {

	ctrl := control()
	var ch1, ch2 <-chan int = intGenerator(ctrl), intGenerator(ctrl)
	var ticker *time.Ticker = time.NewTicker(time.Second * 1)
	var t time.Time
loop:
	for {
		select {
		case i1 := <-ch1:
			fmt.Println("На первом канале:", i1)
		case i2 := <-ch2:
			fmt.Println("На втором канале:", i2)
		case <-ctrl:
			fmt.Println("Кино бетте. Зрители китте")
			break loop
		default:
			t = <-ticker.C
			outputMessage := []byte("Время: ")
			// Метод AppendFormat преобразует объект time.Time
			// к заданному строковому формату (второй аргумент)
			// и добавляет полученную строку к строке, переданной в первом
			// аргументе
			outputMessage = t.AppendFormat(outputMessage, "15:04:05")
			fmt.Println(string(outputMessage))
		}
	}
}
