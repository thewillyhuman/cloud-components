package main

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// MemoryStore implements DNSRecordStore in memory
type MemoryStore struct {
	records map[string]*dns.Msg
}

// NewMemoryStore initializes and returns a new MemoryStore
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{records: make(map[string]*dns.Msg)}
}

// Get retrieves a DNS record for a given domain and query type
func (s *MemoryStore) Get(domain string, qType uint16) (*dns.Msg, bool) {
	msg, ok := s.records[key(domain, qType)]
	return msg, ok
}

// Set stores a DNS record for a given domain and query type
func (s *MemoryStore) Set(domain string, qType uint16, msg *dns.Msg) {
	s.records[key(domain, qType)] = msg
}

// GetAll retrieves all stored DNS records
func (s *MemoryStore) GetAll() map[string]*dns.Msg {
	return s.records
}

func key(domain string, qType uint16) string {
	return fmt.Sprintf("%s:%d", strings.ToLower(domain), qType)
}
