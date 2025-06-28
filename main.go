package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)


func resolveDomain(domain string) (string, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return "", fmt.Errorf("no IPv4 address found")
}

func isRateLimited(output string, err error) bool {
	if err != nil {
		return false
	}
	
	output = strings.ToLower(output)
	rateLimitIndicators := []string{
		"rate limit",
		"too many",
		"query limit",
		"exceeded",
		"throttled",
		"temporarily unavailable",
		"service unavailable",
	}
	
	for _, indicator := range rateLimitIndicators {
		if strings.Contains(output, indicator) {
			return true
		}
	}
	
	return false
}

func getOrganizationWithRetry(ip string, maxRetries int) (string, error) {
	for attempt := 0; attempt < maxRetries; attempt++ {
		cmd := exec.Command("whois", ip)
		cmd.Env = os.Environ()
		
		// Set timeout for whois command
		if attempt > 0 {
			// Add exponential backoff delay before retry
			backoffDelay := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			time.Sleep(backoffDelay)
		}
		
		output, err := cmd.Output()
		outputStr := string(output)
		
		// Check for rate limiting
		if isRateLimited(outputStr, err) {
			if attempt < maxRetries-1 {
				continue // Retry with backoff
			}
			return "", fmt.Errorf("rate limited after %d attempts", maxRetries)
		}
		
		// If command failed for other reasons, return empty (don't retry)
		if err != nil {
			return "", nil
		}

		// Parse successful output - handle different registry formats
		// Priority order: most specific to least specific
		var candidates []string
		scanner := bufio.NewScanner(strings.NewReader(outputStr))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "OrgName:") ||        // ARIN format (high priority)
			   strings.HasPrefix(line, "org-name:") ||       // RIPE format (high priority)
			   strings.HasPrefix(line, "Organization:") {    // Other formats (medium priority)
				cleaned := regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")
				parts := strings.Split(cleaned, " ")
				if len(parts) > 1 {
					candidates = append(candidates, strings.Join(parts[1:], " "))
				}
			}
		}
		
		// Return the first (highest priority) candidate
		if len(candidates) > 0 {
			return candidates[0], nil
		}
		
		// No OrgName found, but no error - return empty
		return "", nil
	}
	
	return "", fmt.Errorf("max retries exceeded")
}

func getOrganization(ip string) (string, error) {
	return getOrganizationWithRetry(ip, 3) // Max 3 attempts with exponential backoff
}

func readDomainsFromStdin() ([]string, error) {
	var domains []string
	scanner := bufio.NewScanner(os.Stdin)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.Contains(line, "*") {
			domains = append(domains, line)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	
	return domains, nil
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

type Result struct {
	subdomain string
	ip        string
	org       string
}


func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("dnsrecon %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	if flag.NArg() != 0 {
		fmt.Println("This tool reads domains from stdin/pipe input.")
		fmt.Printf("Usage: cat domains.txt | %s [-version]\n", os.Args[0])
		fmt.Printf("   or: chaos -d example.com | %s\n", os.Args[0])
		os.Exit(1)
	}

	// Read domains from stdin
	domains, err := readDomainsFromStdin()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading domains from stdin: %v\n", err)
		os.Exit(1)
	}

	if len(domains) == 0 {
		fmt.Fprintf(os.Stderr, "No domains provided in stdin\n")
		os.Exit(1)
	}

	// Remove duplicates
	domains = removeDuplicates(domains)

	// Process domains one at a time to avoid overwhelming whois servers
	for i, domain := range domains {
		ip, err := resolveDomain(domain)
		if err != nil {
			continue
		}

		org, _ := getOrganization(ip)
		fmt.Printf("%s;%s;%s\n", domain, ip, org)
		
		// Add small delay between requests to be respectful to whois servers
		if i < len(domains)-1 { // Don't delay after the last domain
			time.Sleep(100 * time.Millisecond)
		}
	}
}