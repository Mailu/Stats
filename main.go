package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Info struct {
	Timestamp    time.Time
	CustomValues []string
}

var (
	saveChan = make(chan Info)
	logPath  = getEnvOrFallback("LOGPATH", "/output.log")

	host          = getEnvOrFallback("HOST", "0.0.0.0")
	port          = getEnvOrFallback("PORT", "53")
	listenAddress = net.JoinHostPort(host, port)

	domain        = getEnvOrFallback("DOMAIN", "stats.mailu.io.")
	matchingParts = dns.CountLabel(domain)

	valueCount int
)

func fatalErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getEnvOrFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func init() {
	n, err := strconv.Atoi(getEnvOrFallback("VALUECOUNT", "2"))
	fatalErr(err)
	valueCount = n
}

func main() {
	// open log file and channel to write to it
	go func() {
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		fatalErr(err)
		defer f.Close()

		for v := range saveChan {
			line := fmt.Sprintf("%s,%d\n", strings.Join(v.CustomValues, ","), v.Timestamp.Unix())
			fmt.Print(line)
			if _, err := f.WriteString(line); err != nil {
				log.Printf("Could not write to file: %v\n", err)
			}
		}
	}()

	// start DNS-server
	go func() {
		fatalErr(dns.ListenAndServe(listenAddress, "udp", dns.HandlerFunc(handleQuery)))
	}()
	log.Printf("serving queries for subdomains of %s on %s", domain, listenAddress)

	// block until interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}

func handleQuery(w dns.ResponseWriter, r *dns.Msg) {
	// respond to queries with SERVFAIL, then close the connection
	// -> clients don't cache the result and won't block until they timeout
	defer w.Close()
	defer w.WriteMsg(&dns.Msg{
		MsgHdr: dns.MsgHdr{
			Response:           true,
			Opcode:             dns.OpcodeQuery,
			RecursionAvailable: true,
			RecursionDesired:   true,
			Rcode:              dns.RcodeServerFailure,
		},
		Question: r.Question,
	})

	for _, m := range r.Question {
		// not a subdomain of domain
		if dns.CompareDomainName(domain, m.Name) < matchingParts {
			return
		}

		// more values than we asked for
		if dns.CountLabel(m.Name) > valueCount+matchingParts {
			return
		}

		// getting only the "value"-domainparts
		domainParts := dns.SplitDomainName(m.Name)
		values := domainParts[0 : len(domainParts)-matchingParts]

		// push to channel for saving
		saveChan <- Info{
			Timestamp:    time.Now(),
			CustomValues: values,
		}
	}
}
