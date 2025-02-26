package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
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

	results := make(chan []string, len(dnsServers))
	var wg sync.WaitGroup

	for server, location := range dnsServers {
		wg.Add(1)
		go lookupDNS(server, location, domain, recordType, &wg, results)
	}

	wg.Wait()
	close(results)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"DNS Server", "Location", "Record Type", "Result"})
	table.SetBorder(true)

	for result := range results {
		table.Append(result)
	}

	table.Render()
}
