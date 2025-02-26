package main

import (
	"context"
	"net"
	"strings"
	"sync"
)

var dnsServers = map[string]string{
	"8.8.8.8":         "USA West Coast (Google)",
	"8.8.4.4":         "USA East Coast (Google)",
	"1.1.1.1":         "USA East Coast (Cloudflare)",
	"1.0.0.1":         "USA West Coast (Cloudflare)",
	"9.9.9.9":         "USA (Quad9)",
	"208.67.222.222":  "USA West Coast (OpenDNS)",
	"208.67.220.220":  "USA East Coast (OpenDNS)",
	"77.88.8.8":       "Russia (Yandex)",
	"185.228.168.9":   "Europe (CleanBrowsing)",
	"84.200.69.80":    "Germany (DNS Watch)",
	"180.76.76.76":    "China (Baidu)",
	"114.114.114.114": "China (114DNS)",
	"168.126.63.1":    "South Korea (KT)",
	"200.160.0.10":    "Brazil (NIC.br)",
	"196.3.81.5":      "South Africa (MTN)",
	"203.80.96.10":    "Hong Kong (PCCW)",
	"202.134.0.155":   "India (Reliance)",
	"139.99.222.72":   "Australia (OVH DNS)",
	"101.101.101.101": "Taiwan (Hinet)",
	"49.213.18.19":    "Thailand (TOT Public DNS)",
	"103.197.71.66":   "Singapore (SingNet DNS)",
	"149.112.112.112": "Canada (Quad9)",
	"45.90.28.190":    "UAE (NextDNS)",
	"91.239.100.100":  "Sweden (Freenom)",
	"195.46.39.39":    "Ukraine (SafeDNS)",
	"185.222.222.222": "Poland (DNSForge)",
	"213.80.98.2":     "Norway (PowerTech)",
	"212.81.1.1":      "Israel (Bezeq DNS)",
	"80.80.80.80":     "Switzerland (Freenom)",
	"198.51.45.22":    "Saudi Arabia (STC DNS)",
}

func lookupDNS(server, location, domain, recordType string, wg *sync.WaitGroup, results chan<- []string) {
	defer wg.Done()
	resolv := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("udp", server+":53")
		},
	}

	var records []string
	var err error
	switch strings.ToUpper(recordType) {
	case "A":
		records, err = resolv.LookupHost(context.Background(), domain)
	case "AAAA":
		ips, err := resolv.LookupIP(context.Background(), "ip6", domain)
		if err == nil {
			for _, ip := range ips {
				records = append(records, ip.String())
			}
		}
	case "CNAME":
		cname, err := resolv.LookupCNAME(context.Background(), domain)
		if err == nil {
			records = append(records, cname)
		}
	case "MX":
		mxRecords, err := resolv.LookupMX(context.Background(), domain)
		if err == nil {
			for _, mx := range mxRecords {
				records = append(records, mx.Host)
			}
		}
	case "TXT":
		records, err = resolv.LookupTXT(context.Background(), domain)
	default:
		results <- []string{server, location, recordType, "❌ Unsupported record type"}
		return
	}

	if err != nil {
		results <- []string{server, location, recordType, "❌ " + err.Error()}
	} else {
		results <- []string{server, location, recordType, "✅ " + strings.Join(records, ", ")}
	}
}
