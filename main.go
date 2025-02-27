package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hamochi/dnsprop/dnsprop"
	"github.com/hamochi/dnsprop/tui"
	"os"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [record type] domain.com")
		return
	}

	var recordType, domain string

	if len(os.Args) == 2 {
		recordType = "A"
		domain = os.Args[1]
	} else {
		recordType = os.Args[1]
		domain = os.Args[2]
	}

	fmt.Printf("🔍 Checking DNS propagation for %s (%s records)\n", domain, recordType)

	initstatuses := make(map[string]tui.Status)
	for _, server := range dnsprop.DNSServers {
		initstatuses[server.IP] = tui.Status{Name: server.Name, IP: server.IP, Location: server.Location, Flag: server.Flag, Results: "⏳ Pending"}
	}
	mut := &sync.Mutex{}
	m := tui.NewModel(initstatuses, mut)
	p := tea.NewProgram(m)

	results := make(chan map[string]string, len(dnsprop.DNSServers))
	for _, server := range dnsprop.DNSServers {
		go dnsprop.LookupDNS(server.IP, domain, recordType, results)
	}

	go func() {
		count := 0
		for result := range results {
			count++
			s, ok := m.Statuses[result["serverIP"]]
			if ok {
				mut.Lock()
				s.Results = result["status"]
				m.Statuses[result["serverIP"]] = s
				mut.Unlock()
			}
			p.Send("update")
			if count == len(dnsprop.DNSServers) {
				close(results)
			}
		}
		p.Send("quit")
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}

}
