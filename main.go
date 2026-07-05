package main

import (
	"encoding/binary"
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

// Raw DNS query bypassing /etc/hosts
func dnsQuery(domain string) (string, error) {
	// Build DNS query
	txid := []byte{0x01, 0x02}
	flags := []byte{0x01, 0x00}
	qdcount := []byte{0x00, 0x01}
	header := append(txid, flags...)
	header = append(header, qdcount...)
	header = append(header, 0, 0, 0, 0, 0, 0) // ancount, nscount, arcount

	// Encode domain
	qname := []byte{}
	for _, part := range splitDomain(domain) {
		qname = append(qname, byte(len(part)))
		qname = append(qname, []byte(part)...)
	}
	qname = append(qname, 0)

	qtype := []byte{0x00, 0x01} // A
	qclass := []byte{0x00, 0x01} // IN

	query := append(header, qname...)
	query = append(query, qtype...)
	query = append(query, qclass...)

	// Send UDP query to 8.8.8.8
	conn, err := net.DialTimeout("udp", "8.8.8.8:53", 5*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	if _, err := conn.Write(query); err != nil {
		return "", err
	}

	resp := make([]byte, 512)
	n, err := conn.Read(resp)
	if err != nil {
		return "", err
	}
	resp = resp[:n]

	// Parse response: skip header (12) + question
	pos := 12
	for resp[pos] != 0 {
		pos += int(resp[pos]) + 1
	}
	pos += 5 // null + qtype + qclass

	// Parse answer
	if pos+12 > len(resp) {
		return "", fmt.Errorf("no answer")
	}
	rdlen := binary.BigEndian.Uint16(resp[pos+10 : pos+12])
	if rdlen == 4 {
		ip := net.IPv4(resp[pos+12], resp[pos+13], resp[pos+14], resp[pos+15])
		return ip.String(), nil
	}
	return "", fmt.Errorf("no A record")
}

func splitDomain(domain string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(domain); i++ {
		if domain[i] == '.' {
			parts = append(parts, domain[start:i])
			start = i + 1
		}
	}
	parts = append(parts, domain[start:])
	return parts
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

	// Raw DNS query bypassing hosts file
	ipStr, err := dnsQuery(sni)
	if err != nil || ipStr == "" {
		return
	}

	target, err := net.DialTimeout("tcp", ipStr+":443", timeout)
	if err != nil {
		return
	}
	defer target.Close()
	target.SetDeadline(time.Now().Add(timeout))

	if _, err := target.Write(peek); err != nil {
		return
	}

	fmt.Printf("[加速成功] %s -> %s\n", sni, ipStr)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); pipe(conn, target) }()
	go func() { defer wg.Done(); pipe(target, conn) }()
	wg.Wait()
}

func main() {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[错误] 启动失败: %v (需要管理员权限)\n", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("========================================")
	fmt.Println("  GitHub 加速器已启动")
	fmt.Println("  监听端口: 443")
	fmt.Printf("  加速域名: %d 个\n", len(domains))
	fmt.Println("  按 Ctrl+C 停止")
	fmt.Println("========================================")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}
