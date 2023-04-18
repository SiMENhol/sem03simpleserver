package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/SiMENhol/is105sem03/mycrypt"
)

func main() {
	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:8080")
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
					log.Println("Kryptert melding: ", string(buf[:n]))
					//Krypter melding
					kryptertMelding := mycrypt.Krypter([]rune(os.Args[1]), mycrypt.ALF_SEM03, 4)
					log.Println("Kryptert melding: ", string(kryptertMelding))
					_, err = conn.Write([]byte(string(kryptertMelding)))
					log.Println("os.Args[1] = ", os.Args[1])

					switch msg := string(dekryptertMelding); msg {
					case "ping":
						_, err = c.Write([]byte("pong"))
						if err != nil {
							log.Println(err)
							return // fra for-løkke
						}
					}
				}
			}(conn)
		}
	}()

	wg.Wait()
}
