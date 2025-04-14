package main

import (
	"flag"
	"log"

	"github.com/miekg/dns"
)

// StartDNSUDPServer starts the DNS UDP server
func StartDNSUDPServer(localDomain string, localStore DNSRecordStore, cacheStore DNSRecordStore) {
	server := &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", DNSHandler(localDomain, localStore, cacheStore))

	log.Println("Starting DNS UDP server on :53")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start UDP server: %v", err)
	}
}

// StartDNSTCPServer starts the DNS TCP server
func StartDNSTCPServer(localDomain string, localStore DNSRecordStore, cacheStore DNSRecordStore) {
	server := &dns.Server{Addr: ":53", Net: "tcp"}
	dns.HandleFunc(".", DNSHandler(localDomain, localStore, cacheStore))

	log.Println("Starting DNS TCP server on :53")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
}

// StartFrontendServer starts the frontend (e.g., Prometheus metrics and UI)
func StartFrontendServer(localDomain string, localStore DNSRecordStore, cacheStore DNSRecordStore) {
	log.Println("Starting frontend server")
	StartFrontend(localDomain, localStore, cacheStore)
}

func main() {
	localDomain := flag.String("local-domain", "local.", "The local domain to use (e.g., 'mydomain')")
	flag.Parse()

	localStore := NewMemoryStore()
	cacheStore := NewMemoryStore()

	go StartFrontendServer(*localDomain, localStore, cacheStore)
	go StartDNSUDPServer(*localDomain, localStore, cacheStore)
	StartDNSTCPServer(*localDomain, localStore, cacheStore)
}
