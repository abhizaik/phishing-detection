package rank

func DomainRankLookup(domain string) int {
	if rank, exists := domainRankMap[domain]; exists {
		return rank
	}
	return 0
}
