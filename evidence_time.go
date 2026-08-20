package main

// latestProviderEvidenceMillis returns the newest genuine provider observation
// timestamp in milliseconds. A zero result means the provider did not supply
// usable evidence time; callers must not replace that unknown with retrieval
// or cache time when making a market-freshness claim.
func latestProviderEvidenceMillis(values ...string) int64 {
	var latest int64
	for _, value := range values {
		stamp := providerTimeMillis(value)
		if stamp > latest {
			latest = stamp
		}
	}
	return latest
}
