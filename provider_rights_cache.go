package main

import (
	"encoding/json"
	"time"
)

// hostedRightsFilteredMarketCache is deliberately narrower than the desktop
// cache. HOST-001..003 cannot prove provider-specific legal provenance for the
// historical mixed cache fields, so hosted mode retains only provider-approved
// canonical quotes plus operational capability state. Desktop serialization is
// byte-for-byte field-equivalent to the legacy MarketCache shape.
func hostedRightsFilteredMarketCache(c MarketCache) MarketCache {
	if !isHostedRuntime() {
		return c
	}
	out := MarketCache{
		Quotes:                   map[string]Quote{},
		ProviderCapabilityStates: c.ProviderCapabilityStates,
		SavedAt:                  c.SavedAt,
	}
	now := time.Now()
	for symbol, q := range c.Quotes {
		if providerQuoteHostedRightsAllowed(q, providerHostedUseProductionServing, now) {
			out.Quotes[symbol] = q
		}
	}
	return out
}

// MarshalJSON makes the existing canonical market-cache owner rights-aware
// without adding a second cache. In hosted mode, data whose retention rights
// cannot be proven is absent from the persisted artifact.
func (c MarketCache) MarshalJSON() ([]byte, error) {
	type marketCacheAlias MarketCache
	filtered := hostedRightsFilteredMarketCache(c)
	return json.Marshal(marketCacheAlias(filtered))
}

// UnmarshalJSON applies the same decision on restart, so an expired/downgraded
// provider record cannot resurrect previously retained data. The decision is
// re-evaluated at every load; a working credential is irrelevant.
func (c *MarketCache) UnmarshalJSON(data []byte) error {
	type marketCacheAlias MarketCache
	var decoded marketCacheAlias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	filtered := hostedRightsFilteredMarketCache(MarketCache(decoded))
	*c = filtered
	return nil
}
