package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/SiMENhol/is105sem03/mycrypt"
	"github.com/simenhol/funtemps2/conv"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.4:8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				log.Println(err)

			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for-løkke
					}

					// Dekrypterer meldingen
					dekryptertMelding := mycrypt.Krypter([]rune(string(buf[:n])), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
					log.Println("Dekryptert melding: ", string(dekryptertMelding))

					switch msg := string(dekryptertMelding); msg {
					case "ping":
						kryptertSvar := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)+4)
						_, err = c.Write([]byte(string(kryptertSvar)))
						if err != nil {
							log.Println(err)
							return // fra for-løkke
						}
					}
					switch msgString := string(dekryptertMelding); msgString {
					case "Kjevik":
						parts := strings.Split(msgString, " ")
						celsius, _ := strconv.ParseFloat(parts[1], 64)
						fahrenheit := conv.CelsiusToFahrenheit(celsius)
						kryptertMelding := mycrypt.Krypter([]rune(fmt.Sprintf("Temperaturen i Fahrenheit er %.2f", fahrenheit)), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)-4)
						log.Println("Kryptert melding: ", string(kryptertMelding))
						_, err = c.Write([]byte(string(kryptertMelding)))
					default:
						_, err = c.Write(buf[:n])
						//kryptertSvar := mycrypt.Krypter([]rune("pong"), mycrypt.ALF_SEM03, len(mycrypt.ALF_SEM03)+4)
						//_, err = c.Write([]byte(string(kryptertSvar)))
						//if err != nil {
						//	log.Println(err)
						//	return // fra for-løkke
					}
				}

			}(conn)
		}
	}()

	wg.Wait()
}
