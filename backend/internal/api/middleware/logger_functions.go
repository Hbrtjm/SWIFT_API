package middleware

import "strings"

// shouldLogIP checks if the IP should be logged based on the filter configuration
func shouldLogIP(ip string, filterIPs []string) bool {
	// If no filter is set, log all IPs
	if len(filterIPs) == 0 {
		return true
	}

	// Check if the IP is in the filter list
	for _, filteredIP := range filterIPs {
		if filteredIP == ip || filteredIP == "*" {
			return true
		}
		// Support for CIDR-like patterns (simplified, just prefix matching)
		if strings.HasSuffix(filteredIP, "*") {
			prefix := strings.TrimSuffix(filteredIP, "*")
			if strings.HasPrefix(ip, prefix) {
				return true
			}
		}
	}
	return false
}
