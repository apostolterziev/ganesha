package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

type Resolver struct {
	records map[string]string
}

func (r *Resolver) parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			ip := r.records[q.Name]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		default:
			fmt.Println("Unknown query type -> " + q.Name)
		}
	}
}

func (r *Resolver) handleDnsRequest(w dns.ResponseWriter, resp *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(resp)
	m.Compress = false

	switch resp.Opcode {
	case dns.OpcodeQuery:
		r.parseQuery(m)
	default:
		fmt.Println("Unknown opcode -> " + string(resp.Opcode))
	}

	w.WriteMsg(m)
}

func (r *Resolver) UpdateDatabase() {
	resolverRecords := GlobalStorage.GetAllResolverRecords()
	r.records = make(map[string]string)
	for _, record := range resolverRecords {
		r.records[record.FQDN + "."] = record.IP
	}
}

func (r *Resolver) run(pattern string) {
	// attach request handler func
	dns.HandleFunc(pattern, r.handleDnsRequest)
	// start server
	port := GlobalConfiguration["resolver.port"]
	listenOn := GlobalConfiguration["resolver.listen_on"]
	server := &dns.Server{Addr: listenOn + ":" + port, Net: "udp"}
	log.Printf("Starting at %s\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}