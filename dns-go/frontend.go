package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/miekg/dns"
)

// StartFrontend starts the HTTP server for the UI and metrics
func StartFrontend(localDomain string, localStore, cacheStore DNSRecordStore) {
	http.Handle("/metrics", promhttp.Handler()) // Expose Prometheus metrics

	// Serve the static CSS file for better styling
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Status page for DNS records
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintln(w, "<link rel=\"stylesheet\" type=\"text/css\" href=\"/static/style.css\">")
		fmt.Fprintln(w, "<h1>DNS Records and Cache Status</h1>")

		// Display form to add new local records
		fmt.Fprintln(w, `<h2>Add Local DNS Record</h2>
			<form method="POST" action="/add-record">
				<label for="domain">Domain:</label><br>
				<input type="text" id="domain" name="domain" required><br><br>
				<label for="type">Record Type:</label><br>
				<select id="type" name="type" required>
					<option value="A">A</option>
					<option value="CNAME">CNAME</option>
					<option value="MX">MX</option>
					<option value="TXT">TXT</option>
				</select><br><br>
				<label for="data">Record Data:</label><br>
				<input type="text" id="data" name="data" required><br><br>
				<input type="submit" value="Add Record">
			</form>
			<hr>`)

		// Display records from local store
		fmt.Fprintln(w, "<h2>Local DNS Records</h2>")
		displayRecords(w, localStore)

		// Display records from cache store
		fmt.Fprintln(w, "<h2>Cache DNS Records</h2>")
		displayRecords(w, cacheStore)
	})

	// Handle adding new local DNS records via POST
	http.HandleFunc("/add-record", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Parse the form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}

			domain := r.FormValue("domain")
			recordType := r.FormValue("type")
			data := r.FormValue("data")

			// Add the record to the local store
			var record dns.RR
			switch recordType {
			case "A":
				record = &dns.A{Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA}, A: net.ParseIP(data)}
			case "CNAME":
				record = &dns.CNAME{Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeCNAME}, Target: data}
			case "MX":
				record = &dns.MX{Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeMX}, Mx: data, Preference: 10}
			case "TXT":
				record = &dns.TXT{Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeTXT}, Txt: []string{data}}
			}

			// Save it to the local store
			msg := new(dns.Msg)
			msg.Answer = append(msg.Answer, record)
			localStore.Set(domain, record.Header().Rrtype, msg)

			// Redirect to the status page
			http.Redirect(w, r, "/status", http.StatusFound)
		} else {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	log.Println("Starting frontend server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// displayRecords renders DNS records in a simple HTML table
func displayRecords(w http.ResponseWriter, store DNSRecordStore) {
	records := store.GetAll()
	if len(records) == 0 {
		fmt.Fprintln(w, "<p>No records found</p>")
		return
	}
	fmt.Fprintln(w, `<table border='1' cellpadding='5' cellspacing='0'>
		<tr><th>Domain</th><th>Record Type</th><th>Record Data</th></tr>`)

	for key, msg := range records {
		for _, answer := range msg.Answer {
			switch record := answer.(type) {
			case *dns.A:
				fmt.Fprintf(w, "<tr><td>%s</td><td>A</td><td>%s</td></tr>", key, record.A)
			case *dns.CNAME:
				fmt.Fprintf(w, "<tr><td>%s</td><td>CNAME</td><td>%s</td></tr>", key, record.Target)
			case *dns.MX:
				fmt.Fprintf(w, "<tr><td>%s</td><td>MX</td><td>%s (Priority: %d)</td></tr>", key, record.Mx, record.Preference)
			case *dns.TXT:
				fmt.Fprintf(w, "<tr><td>%s</td><td>TXT</td><td>%s</td></tr>", key, strings.Join(record.Txt, ", "))
			default:
				// Handle additional DNS record types here as needed
			}
		}
	}
	fmt.Fprintln(w, "</table>")
}
