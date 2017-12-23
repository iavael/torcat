package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/yawning/bulb"
	"github.com/yawning/bulb/utils/pkcs1"
)

var (
	control string
	listen  uint
)

func init() {
	flag.StringVar(&control, "control", "localhost:9051", "Control port")
	flag.UintVar(&listen, "l", 0, "Listen")

	flag.Parse()
}

func main() {
	c, err := bulb.Dial("tcp", control)
	if err != nil {
		log.Fatalf("Failed to connect to control port: %s", err)
	}
	defer c.Close()

	if err := c.Authenticate(os.Getenv("TORCAT_COOKIE")); err != nil {
		log.Fatalf("Authentication failed: %s", err)
	}

	if listen != 0 {
		if pk, err := rsa.GenerateKey(rand.Reader, 1024); err != nil {
			log.Fatalf("Failed to generate RSA key")
		} else if id, err := pkcs1.OnionAddr(&pk.PublicKey); err != nil {
			log.Fatalf("Failed to derive onion ID: %v", err)
		} else {
			os.Stderr.WriteString(id)
			os.Stderr.WriteString(".onion\n")
			cfg := &bulb.NewOnionConfig{
				DiscardPK:  true,
				PrivateKey: pk,
			}
			if l, err := c.NewListener(cfg, uint16(listen)); err != nil {
				log.Fatalf("Failed to listen port: %s", err)
			} else {
				defer l.Close()
				if conn, err := l.Accept(); err != nil {
					log.Fatalf("Failed to accept connection: %s", err)
				} else {
					defer conn.Close()

					select {
					case err := <-runIO(conn):
						if err != nil {
							log.Fatalf("Failed conversation: %s", err)
						}
						break
					}
				}
			}
		}
	} else {
		var dest string
		if len(os.Args) != 3 {
			log.Fatalf("Invalid arguments. Must be `%s host port'", os.Args[0])
		} else if addr := os.Args[1]; len(addr) == 0 {
			log.Fatalf("Empty destination address")
		} else if port, err := strconv.Atoi(os.Args[2]); err != nil {
			log.Fatalf("Invalid port number: %s", err)
		} else {
			dest = fmt.Sprintf("%s:%d", addr, port)
		}

		if dialer, err := c.Dialer(nil); err != nil {
			log.Fatalf("Failed to get Dialer: %v", err)
		} else if conn, err := dialer.Dial("tcp", dest); err != nil {
			log.Fatalf("Connection to %s failed", err)
		} else {
			defer conn.Close()

			select {
			case err := <-runIO(conn):
				if err != nil {
					log.Fatalf("Failed conversation: %s", err)
				}
				break
			}
		}
	}
}

func runIO(conn io.ReadWriter) <-chan error {
	errchan := make(chan error)
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			errchan <- err
		} else {
			errchan <- nil
		}
	}()
	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			errchan <- err
		} else {
			errchan <- nil
		}
	}()
	return errchan
}
