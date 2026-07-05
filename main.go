package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const port = "443"
const bufSize = 8192
const timeout = 30 * time.Second

var domains = map[string]bool{
	"github.com": true, "api.github.com": true, "raw.github.com": true,
	"raw.githubusercontent.com": true, "avatars.githubusercontent.com": true,
	"avatars0.githubusercontent.com": true, "avatars1.githubusercontent.com": true,
	"avatars2.githubusercontent.com": true, "avatars3.githubusercontent.com": true,
	"github.githubassets.com": true, "objects.githubusercontent.com": true,
	"user-images.githubusercontent.com": true, "camo.githubusercontent.com": true,
	"cloud.githubusercontent.com": true, "gist.github.com": true,
	"github.io": true, "github.dev": true, "pages.github.com": true,
	"githubapp.com": true, "www.github.io": true,
}

func extractSNI(data []byte) string {
	if len(data) < 5 || data[0] != 0x16 {
		return ""
	}
	p := 9 + 1 + int(data[9])
	if p+2 > len(data) {
		return ""
	}
	cipherLen := int(data[p])<<8 | int(data[p+1])
	p += 2 + cipherLen
	if p >= len(data) {
		return ""
	}
	p += 1 + int(data[p])
	if p+2 > len(data) {
		return ""
	}
	extEnd := p + 2 + int(data[p])<<8 | int(data[p+1])
	p += 2
	for p < extEnd-4 {
		extType := int(data[p])<<8 | int(data[p+1])
		extLen := int(data[p+2])<<8 | int(data[p+3])
		if extType == 0 {
			if p+9 > len(data) {
				return ""
			}
			sniLen := int(data[p+7])<<8 | int(data[p+8])
			if p+9+sniLen > len(data) {
				return ""
			}
			return string(data[p+9 : p+9+sniLen])
		}
		p += 4 + extLen
	}
	return ""
}

func pipe(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()
	buf := make([]byte, bufSize)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(timeout))

	peek := make([]byte, bufSize)
	n, err := conn.Read(peek)
	if err != nil || n == 0 {
		return
	}
	peek = peek[:n]

	sni := extractSNI(peek)
	if sni == "" || !domains[sni] {
		return
	}

	ips, err := net.LookupIP(sni)
	if err != nil || len(ips) == 0 {
		return
	}

	var ip net.IP
	for _, v := range ips {
		if v4 := v.To4(); v4 != nil {
			ip = v4
			break
		}
	}
	if ip == nil {
		return
	}

	target, err := net.DialTimeout("tcp", ip.String()+":443", timeout)
	if err != nil {
		return
	}
	defer target.Close()
	target.SetDeadline(time.Now().Add(timeout))

	// Send peeked data
	if _, err := target.Write(peek); err != nil {
		return
	}

	// Bidirectional forward
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); pipe(conn, target) }()
	go func() { defer wg.Done(); pipe(target, conn) }()
	wg.Wait()
}

func main() {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Listen error: %v (need root)\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Printf("GitHub accelerator :%s (%d domains)\n", port, len(domains))

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}
