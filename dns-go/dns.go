package main

import (
	"log"
	"strings"

	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
)

// DNSRecordStore defines storage interface for DNS records
type DNSRecordStore interface {
	Get(domain string, qType uint16) (*dns.Msg, bool)
	Set(domain string, qType uint16, msg *dns.Msg)
	GetAll() map[string]*dns.Msg // To fetch all records for UI
}

// DNSHandler processes incoming DNS queries
func DNSHandler(localDomain string, localStore, cacheStore DNSRecordStore) dns.HandlerFunc {
	return func(w dns.ResponseWriter, r *dns.Msg) {
		// Log the DNS request
		log.Printf("Received DNS request: %s", r.Question[0].Name)

		// Count DNS request in Prometheus
		dnsRequests.WithLabelValues("query").Inc()

		// Create a response message
		response := new(dns.Msg)
		response.SetReply(r)

		// Process each question in the request
		handled := false
		for _, q := range r.Question {
			domain := q.Name

			var store DNSRecordStore
			if strings.HasSuffix(domain, localDomain) {
				store = localStore
				log.Printf("Local domain found in local store: %s", domain)
				log.Printf("Cache: %s", store.GetAll())
			} else {
				store = cacheStore
				log.Printf("Local domain found in local cache: %s", domain)
			}

			if msg, ok := store.Get(domain, q.Qtype); ok {
				// Cache hit
				log.Printf("Cache hit for %s", domain)
				log.Printf("Cache hit: %s", msg)
				response.Answer = append(response.Answer, msg.Answer...)
				handled = true
				continue
			}

			// Forward request to upstream DNS (8.8.8.8 as an example)
			client := new(dns.Client)
			msg, _, err := client.Exchange(r, "8.8.8.8:53")
			if err != nil {
				log.Printf("Failed to resolve %s: %v", domain, err)
				dns.HandleFailed(w, r)
				return
			}

			// Store the result in the appropriate store
			store.Set(domain, q.Qtype, msg)
			response.Answer = append(response.Answer, msg.Answer...)
			handled = true
		}

		// If handled, send the response
		if handled {
			w.WriteMsg(response)
			// Log the DNS response
			log.Printf("Responded with %d answers", len(response.Answer))
			dnsRequests.WithLabelValues("response").Inc()
		} else {
			dns.HandleFailed(w, r)
		}
	}
}

var dnsRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "dns_requests_total",
		Help: "Total number of DNS requests and responses",
	},
	[]string{"type"},
)

func init() {
	prometheus.MustRegister(dnsRequests)
}
