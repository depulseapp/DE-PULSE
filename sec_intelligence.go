package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

func secFormMeaning(form string) (string, string) {
	f := strings.ToUpper(strings.TrimSpace(form))
	base := strings.TrimSuffix(f, "/A")
	switch {
	case base == "4":
		return "Insider Transaction", "insider"
	case base == "8-K":
		return "Material Company Update", "material"
	case base == "10-Q":
		return "Quarterly Report", "report"
	case base == "10-K":
		return "Annual Report", "report"
	case base == "6-K":
		return "Foreign Issuer Update", "material"
	case base == "20-F":
		return "Foreign Annual Report", "report"
	case base == "13F-HR" || base == "13F-NT":
		return "Institutional Holdings Report", "institutional"
	case strings.Contains(base, "13D") || strings.Contains(base, "13G"):
		return "Major Shareholder Disclosure", "ownership"
	case base == "S-1":
		return "Registration / Share Offering", "offering"
	case base == "S-3":
		return "Shelf Registration / Offering", "offering"
	case strings.HasPrefix(base, "424B"):
		return "Offering Terms", "offering"
	case base == "144":
		return "Planned Insider Sale", "insider"
	case base == "DEF 14A":
		return "Proxy Statement", "governance"
	default:
		return "SEC filing", "other"
	}
}

func isSECRelevantForm(form string) bool {
	_, category := secFormMeaning(form)
	return category != "other"
}

type secXMLValue struct {
	Value string `xml:"value"`
}

type secForm4Transaction struct {
	TransactionDate secXMLValue `xml:"transactionDate"`
	Coding          struct {
		Code string `xml:"transactionCode"`
	} `xml:"transactionCoding"`
	Amounts struct {
		Shares           secXMLValue `xml:"transactionShares"`
		Price            secXMLValue `xml:"transactionPricePerShare"`
		AcquiredDisposed secXMLValue `xml:"transactionAcquiredDisposedCode"`
	} `xml:"transactionAmounts"`
	PostTransactionAmounts struct {
		SharesOwned secXMLValue `xml:"sharesOwnedFollowingTransaction"`
		ValueOwned  secXMLValue `xml:"valueOwnedFollowingTransaction"`
	} `xml:"postTransactionAmounts"`
}

type secForm4Document struct {
	ReportingOwners []struct {
		OwnerID struct {
			Name string `xml:"rptOwnerName"`
		} `xml:"reportingOwnerId"`
		Relationship struct {
			IsDirector        string `xml:"isDirector"`
			IsOfficer         string `xml:"isOfficer"`
			IsTenPercentOwner string `xml:"isTenPercentOwner"`
			OfficerTitle      string `xml:"officerTitle"`
		} `xml:"reportingOwnerRelationship"`
	} `xml:"reportingOwner"`
	NonDerivativeTable struct {
		Transactions []secForm4Transaction `xml:"nonDerivativeTransaction"`
	} `xml:"nonDerivativeTable"`
	DerivativeTable struct {
		Transactions []secForm4Transaction `xml:"derivativeTransaction"`
	} `xml:"derivativeTable"`
}

func form4TransactionMeaning(code string) (string, string) {
	c := strings.ToUpper(strings.TrimSpace(code))

	switch c {
	case "P":
		return "BUY", "Open-market/private purchase"
	case "S":
		return "SELL", "Open-market/private sale"
	case "A":
		return "OTHER", "Award / grant"
	case "M":
		return "OTHER", "Option exercise / conversion"
	case "F":
		return "OTHER", "Tax withholding / payment"
	case "G":
		return "OTHER", "Gift"
	case "D":
		return "OTHER", "Disposition to issuer"
	case "C":
		return "OTHER", "Conversion"
	case "X":
		return "OTHER", "Exercise of in-the-money derivative"
	case "L":
		return "OTHER", "Small acquisition"
	case "W":
		return "OTHER", "Acquisition/disposition by will or inheritance"
	case "Z":
		return "OTHER", "Deposit/withdrawal from voting trust"
	case "J":
		return "OTHER", "Other transaction"
	case "K":
		return "OTHER", "Equity swap / similar instrument"
	case "U":
		return "OTHER", "Tender of shares in change-of-control transaction"
	case "I":
		return "OTHER", "Discretionary transaction"
	case "E", "H", "O":
		return "OTHER", "Derivative expiration / exercise-related transaction"
	default:
		if c == "" {
			return "OTHER", "Form 4 transaction"
		}
		return "OTHER", "Form 4 transaction (" + c + ")"
	}
}

func secTruthy(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true" || v == "yes"
}

func secGetBytes(ctx context.Context, client *http.Client, rawURL string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}
	return io.ReadAll(io.LimitReader(resp.Body, 2<<20))
}

func enrichForm4(ctx context.Context, client *http.Client, headers map[string]string, rawURL string, item *FilingItem) {
	if item == nil || rawURL == "" {
		return
	}
	b, err := secGetBytes(ctx, client, rawURL, headers)
	if err != nil {
		return
	}
	var doc secForm4Document
	if xml.Unmarshal(b, &doc) != nil {
		return
	}
	if len(doc.ReportingOwners) > 0 {
		names := []string{}
		roles := []string{}
		for _, o := range doc.ReportingOwners {
			if name := strings.TrimSpace(o.OwnerID.Name); name != "" {
				names = append(names, name)
			}
			role := strings.TrimSpace(o.Relationship.OfficerTitle)
			if role == "" && secTruthy(o.Relationship.IsOfficer) {
				role = "Officer"
			}
			if role == "" && secTruthy(o.Relationship.IsDirector) {
				role = "Director"
			}
			if role == "" && secTruthy(o.Relationship.IsTenPercentOwner) {
				role = "10% Owner"
			}
			if role != "" {
				roles = append(roles, role)
			}
		}
		item.Actor = strings.Join(uniqueStrings(names), " / ")
		item.Role = strings.Join(uniqueStrings(roles), " / ")
	}

	actor, role := item.Actor, item.Role
	item.Transactions = nil

	buys, sells := 0, 0
	buyShares, sellShares := 0.0, 0.0
	buyValue, sellValue := 0.0, 0.0
	totalShares, totalValue := 0.0, 0.0
	latestDate := ""
	latestOwnership := 0.0
	codes := []string{}
	types := []string{}
	seenCodes := map[string]bool{}
	seenTypes := map[string]bool{}
	otherCount := 0

	type form4TransactionRow struct {
		tx   secForm4Transaction
		kind string
	}
	allTransactions := make([]form4TransactionRow, 0, len(doc.NonDerivativeTable.Transactions)+len(doc.DerivativeTable.Transactions))
	for _, tx := range doc.NonDerivativeTable.Transactions {
		allTransactions = append(allTransactions, form4TransactionRow{tx: tx, kind: "non-derivative"})
	}
	for _, tx := range doc.DerivativeTable.Transactions {
		allTransactions = append(allTransactions, form4TransactionRow{tx: tx, kind: "derivative"})
	}

	seenTransactionKind := map[string]string{}
	for _, row := range allTransactions {
		tx := row.tx
		code := strings.ToUpper(strings.TrimSpace(tx.Coding.Code))
		class, meaning := form4TransactionMeaning(code)
		sh := toFloat(tx.Amounts.Shares.Value)
		pr := toFloat(tx.Amounts.Price.Value)
		value := sh * pr
		date := strings.TrimSpace(tx.TransactionDate.Value)
		ownership := toFloat(tx.PostTransactionAmounts.SharesOwned.Value)
		if ownership == 0 {
			ownership = toFloat(tx.PostTransactionAmounts.ValueOwned.Value)
		}

		txKey := fmt.Sprintf("%s|%s|%.8f|%.8f|%.8f", date, code, sh, pr, ownership)
		if priorKind, seen := seenTransactionKind[txKey]; seen && priorKind != row.kind {
			continue
		}
		if _, seen := seenTransactionKind[txKey]; !seen {
			seenTransactionKind[txKey] = row.kind
		}
		item.Transactions = append(item.Transactions, SECInsiderTransaction{
			Actor: actor, Role: role, Classification: class, Code: code, Meaning: meaning,
			TransactionDate: date, Shares: sh, Price: pr, Value: value, OwnershipAfter: ownership,
			FiledAt: item.FiledAt, URL: item.URL,
		})
		if !seenCodes[code] && code != "" {
			seenCodes[code] = true
			codes = append(codes, code)
		}
		if !seenTypes[meaning] {
			seenTypes[meaning] = true
			types = append(types, meaning)
		}
		if date >= latestDate {
			latestDate = date
			if ownership > 0 {
				latestOwnership = ownership
			}
		}
		switch class {
		case "BUY":
			buys++
			buyShares += sh
			buyValue += value
		case "SELL":
			sells++
			sellShares += sh
			sellValue += value
		default:
			otherCount++
		}
		if sh > 0 {
			totalShares += sh
		}
		if value > 0 {
			totalValue += value
		}
	}

	item.TransactionDate = latestDate
	item.TransactionCode = strings.Join(codes, "/")
	item.TransactionType = strings.Join(types, " · ")
	item.OwnershipAfter = latestOwnership

	switch {
	case buys > 0 && sells == 0:
		item.Signal = "Buy"
		item.Shares = buyShares
		item.Value = buyValue
		if buyShares > 0 && buyValue > 0 {
			item.Price = buyValue / buyShares
		}
		if otherCount == 0 {
			item.Meaning = "Insider Buy"
		} else {
			item.Meaning = "Insider Buy + Other"
		}
	case sells > 0 && buys == 0:
		item.Signal = "Sell"
		item.Shares = sellShares
		item.Value = sellValue
		if sellShares > 0 && sellValue > 0 {
			item.Price = sellValue / sellShares
		}
		if otherCount == 0 {
			item.Meaning = "Insider Sell"
		} else {
			item.Meaning = "Insider Sell + Other"
		}
	case buys > 0 && sells > 0:
		item.Signal = "Mixed"
		item.Shares = buyShares + sellShares
		item.Value = buyValue + sellValue
		if item.Shares > 0 && item.Value > 0 {
			item.Price = item.Value / item.Shares
		}
		item.Meaning = "Insider Mixed"
	case otherCount > 0:
		item.Signal = "Other"
		item.Shares = totalShares
		item.Value = totalValue
		if totalShares > 0 && totalValue > 0 {
			item.Price = totalValue / totalShares
		}
		item.Meaning = "Insider Other"
	}

	parts := []string{}
	if item.Signal != "" {
		parts = append(parts, strings.ToUpper(item.Signal))
	}
	if item.TransactionType != "" {
		parts = append(parts, item.TransactionType)
	}
	if item.Role != "" {
		parts = append(parts, item.Role)
	}
	if item.Shares > 0 {
		parts = append(parts, fmt.Sprintf("%.0f shares", item.Shares))
	}
	if item.Price > 0 {
		parts = append(parts, fmt.Sprintf("$%.2f avg", item.Price))
	}
	if item.Value > 0 {
		parts = append(parts, fmt.Sprintf("$%.0fK", item.Value/1000))
	}
	if len(parts) > 0 {
		item.Description = strings.Join(parts, " · ")
	}
}

func filingAgeDays(filed string, now time.Time) int {
	t, err := time.Parse("2006-01-02", filed)
	if err != nil {
		return 9999
	}
	d := int(now.Sub(t).Hours() / 24)
	if d < 0 {
		return 0
	}
	return d
}

func buildSECIntelligence(symbol string, items []FilingItem) SECIntelligenceSummary {
	result := SECIntelligenceSummary{Symbol: symbol, OfferingRisk: "Low", Institutional: "13F · Quarterly", Signals: []SECSignal{}, UpdatedAt: time.Now().UnixMilli()}
	now := time.Now()
	for _, f := range items {
		age := filingAgeDays(f.FiledAt, now)
		if result.LatestFiledAt == "" || f.FiledAt > result.LatestFiledAt {
			result.LatestFiledAt, result.LatestForm = f.FiledAt, f.Form
		}
		switch f.Category {
		case "insider":
			if len(f.Transactions) > 0 {
				for _, tx := range f.Transactions {
					switch strings.ToUpper(strings.TrimSpace(tx.Classification)) {
					case "BUY":
						result.InsiderBuys++
						result.InsiderBuyValue += tx.Value
					case "SELL":
						result.InsiderSells++
						result.InsiderSellValue += tx.Value
					default:
						result.InsiderOthers++
					}
					if len(result.RecentInsiderTransactions) < 100 {
						result.RecentInsiderTransactions = append(result.RecentInsiderTransactions, tx)
					}
				}
			} else {
				if strings.EqualFold(f.Signal, "buy") {
					result.InsiderBuys++
					result.InsiderBuyValue += f.Value
				} else if strings.EqualFold(f.Signal, "sell") {
					result.InsiderSells++
					result.InsiderSellValue += f.Value
				} else {
					result.InsiderOthers++
				}
			}
			if f.Actor != "" || f.TransactionCode != "" || len(f.Transactions) > 0 {
				result.RecentTransactions = append(result.RecentTransactions, f)
				if len(result.RecentTransactions) > 6 {
					result.RecentTransactions = result.RecentTransactions[:6]
				}
			}
		case "ownership":
			if age <= 120 {
				result.OwnershipChanges++
			}
		case "material":
			if age <= 90 {
				result.MaterialEvents++
			}
		case "offering":
			if age <= 30 {
				result.OfferingRisk = "High"
			} else if age <= 90 && result.OfferingRisk != "High" {
				result.OfferingRisk = "Watch"
			}
		}
		if len(result.Signals) < 3 {
			label, detail, tone := f.Form, f.Meaning, "neutral"
			if f.Signal != "" {
				label = strings.ToUpper(f.Signal)
				detail = strings.TrimSpace(strings.Join([]string{f.Role, func() string {
					if f.Value > 0 {
						return fmt.Sprintf("$%.0fK", f.Value/1000)
					}
					return ""
				}()}, " · "))
				detail = strings.Trim(detail, " ·")
			}
			if strings.EqualFold(f.Signal, "buy") {
				tone = "positive"
			}
			if strings.EqualFold(f.Signal, "sell") || f.Category == "offering" {
				tone = "negative"
			}
			if f.Category == "material" || f.Category == "ownership" {
				tone = "watch"
			}
			result.Signals = append(result.Signals, SECSignal{Label: label, Detail: detail, Tone: tone, Date: f.FiledAt, URL: f.URL})
		}
	}
	return result
}

// refreshSECResearchSymbol performs a deeper, target-scoped EDGAR reconciliation for
// the active Research symbol. The normal universe refresh intentionally caps Form 4
// work for efficiency; Research needs enough canonical enriched transactions to make
// its 30/90-day BUY/SELL/OTHER summaries materially complete without multiplying that
// cost across every tracked symbol.
func (e *Engine) refreshSECResearchSymbol(ctx context.Context, symbol string) bool {
	symbol = normalizeSymbol(symbol)
	if !validUserTicker(symbol) {
		return false
	}
	e.app.mu.RLock()
	email := strings.TrimSpace(e.app.state.Settings.SECEmail)
	e.app.mu.RUnlock()
	if !strings.Contains(email, "@") {
		e.setHealth("research-sec:"+symbol, "setup required · SEC contact email")
		return false
	}
	client := &http.Client{Timeout: 20 * time.Second}
	headers := map[string]string{"User-Agent": appName + "/" + appVersion + " " + email}
	var tickerMap map[string]struct {
		CIK    int    `json:"cik_str"`
		Ticker string `json:"ticker"`
	}
	if err := getJSON(ctx, client, "https://www.sec.gov/files/company_tickers.json", headers, &tickerMap); err != nil {
		e.recordProviderFailure("SEC EDGAR", err)
		e.setHealth("research-sec:"+symbol, "degraded · ticker map unavailable")
		return false
	}
	cik := ""
	for _, row := range tickerMap {
		if strings.EqualFold(row.Ticker, symbol) {
			cik = fmt.Sprintf("%010d", row.CIK)
			break
		}
	}
	if cik == "" {
		e.setHealth("research-sec:"+symbol, "unavailable · symbol not mapped by SEC")
		return false
	}
	var data struct {
		Name    string `json:"name"`
		Filings struct {
			Recent struct {
				Accession       []string `json:"accessionNumber"`
				FilingDate      []string `json:"filingDate"`
				ReportDate      []string `json:"reportDate"`
				Form            []string `json:"form"`
				PrimaryDocument []string `json:"primaryDocument"`
				Description     []string `json:"primaryDocDescription"`
				Items           []string `json:"items"`
			} `json:"recent"`
		} `json:"filings"`
	}
	if err := getJSON(ctx, client, "https://data.sec.gov/submissions/CIK"+cik+".json", headers, &data); err != nil {
		e.recordProviderFailure("SEC EDGAR", err)
		e.setHealth("research-sec:"+symbol, "degraded · submissions unavailable")
		return false
	}
	categoryCounts := map[string]int{}
	capByCategory := map[string]int{"insider": 25, "material": 5, "report": 4, "ownership": 4, "offering": 4, "institutional": 2, "governance": 2}
	cutoff := time.Now().In(easternLocation()).AddDate(0, 0, -95)
	var symbolItems []FilingItem
	for i, form := range data.Filings.Recent.Form {
		if !isSECRelevantForm(form) || i >= len(data.Filings.Recent.Accession) {
			continue
		}
		meaning, category := secFormMeaning(form)
		limit := capByCategory[category]
		if limit <= 0 || categoryCounts[category] >= limit {
			continue
		}
		filed := ""
		if i < len(data.Filings.Recent.FilingDate) {
			filed = data.Filings.Recent.FilingDate[i]
		}
		if category == "insider" && filed != "" {
			if ft, err := time.ParseInLocation("2006-01-02", filed, easternLocation()); err == nil && ft.Before(cutoff) {
				continue
			}
		}
		acc := data.Filings.Recent.Accession[i]
		doc, report, rawDesc, items := "", "", "", ""
		if i < len(data.Filings.Recent.PrimaryDocument) {
			doc = data.Filings.Recent.PrimaryDocument[i]
		}
		if i < len(data.Filings.Recent.ReportDate) {
			report = data.Filings.Recent.ReportDate[i]
		}
		if i < len(data.Filings.Recent.Description) {
			rawDesc = strings.TrimSpace(data.Filings.Recent.Description[i])
		}
		if i < len(data.Filings.Recent.Items) {
			items = strings.TrimSpace(data.Filings.Recent.Items[i])
		}
		link := fmt.Sprintf("https://www.sec.gov/edgar/browse/?CIK=%d", mustInt(cik))
		if doc != "" {
			link = fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%d/%s/%s", mustInt(cik), strings.ReplaceAll(acc, "-", ""), doc)
		}
		desc := meaning
		if category == "material" && items != "" {
			desc = meaning + " · Items " + items
		}
		if rawDesc != "" && category == "other" {
			desc = rawDesc
		}
		item := FilingItem{ID: symbol + "-" + acc, Symbol: symbol, Company: data.Name, Form: form, FiledAt: filed, ReportDate: report, Description: desc, Meaning: meaning, Category: category, Items: items, URL: link}
		if category == "insider" && strings.HasPrefix(strings.ToUpper(strings.TrimSpace(form)), "4") && doc != "" {
			enrichForm4(ctx, client, headers, link, &item)
			if !sleepContext(ctx, 110*time.Millisecond) {
				return false
			}
		}
		symbolItems = append(symbolItems, item)
		categoryCounts[category]++
		if len(symbolItems) >= 45 {
			break
		}
	}
	sort.Slice(symbolItems, func(i, j int) bool { return symbolItems[i].FiledAt > symbolItems[j].FiledAt })
	e.mu.Lock()
	merged := make([]FilingItem, 0, len(e.filings)+len(symbolItems))
	for _, f := range e.filings {
		if !strings.EqualFold(f.Symbol, symbol) {
			merged = append(merged, f)
		}
	}
	merged = append(merged, symbolItems...)
	sort.Slice(merged, func(i, j int) bool { return merged[i].FiledAt > merged[j].FiledAt })
	if len(merged) > 220 {
		merged = merged[:220]
	}
	e.filings = merged
	if e.secIntelligence == nil {
		e.secIntelligence = map[string]SECIntelligenceSummary{}
	}
	e.secIntelligence[symbol] = buildSECIntelligence(symbol, symbolItems)
	now := time.Now().UnixMilli()
	e.lastUpdated["research-sec:"+symbol] = now
	e.health["research-sec:"+symbol] = fmt.Sprintf("healthy · SEC EDGAR · %d filings · %d Form 4 mapped", len(symbolItems), categoryCounts["insider"])
	filingsSnapshot := clone(e.filings)
	intelSnapshot := clone(e.secIntelligence)
	e.mu.Unlock()
	e.recordProviderSuccess("SEC EDGAR")
	e.app.broadcastFilings(filingsSnapshot, intelSnapshot)
	return true
}

func (e *Engine) refreshFilings(ctx context.Context) {
	if e.highImpactModeActive() {
		e.setHealth("filings", "deferred · high-impact event mode")
		return
	}
	e.app.mu.RLock()
	email := strings.TrimSpace(e.app.state.Settings.SECEmail)
	symbols := analysisSymbolsFromState(e.app.processingStateLocked())
	e.app.mu.RUnlock()
	if !strings.Contains(email, "@") {
		e.setHealth("filings", "setup required · keeping cached data")
		return
	}
	e.setHealth("filings", "loading")
	e.mu.RLock()
	oldFinancialFingerprint := financialFilingFingerprint(e.filings)
	e.mu.RUnlock()
	client := &http.Client{Timeout: 20 * time.Second}
	headers := map[string]string{"User-Agent": appName + "/" + appVersion + " " + email}
	var tickerMap map[string]struct {
		CIK    int    `json:"cik_str"`
		Ticker string `json:"ticker"`
	}
	if err := getJSON(ctx, client, "https://www.sec.gov/files/company_tickers.json", headers, &tickerMap); err != nil {
		e.setHealth("filings", "degraded")
		e.setError("SEC", err)
		return
	}
	byTicker := map[string]string{}
	for _, row := range tickerMap {
		byTicker[strings.ToUpper(row.Ticker)] = fmt.Sprintf("%010d", row.CIK)
	}
	symbols = uniqueSymbols(symbols)
	if len(symbols) > 20 {
		symbols = symbols[:20]
	}
	var filings []FilingItem
	intel := map[string]SECIntelligenceSummary{}
	fetchedCompanies := 0
	for _, symbol := range symbols {
		cik := byTicker[symbol]
		if cik == "" {
			continue
		}
		var data struct {
			Name    string `json:"name"`
			Filings struct {
				Recent struct {
					Accession       []string `json:"accessionNumber"`
					FilingDate      []string `json:"filingDate"`
					ReportDate      []string `json:"reportDate"`
					Form            []string `json:"form"`
					PrimaryDocument []string `json:"primaryDocument"`
					Description     []string `json:"primaryDocDescription"`
					Items           []string `json:"items"`
				} `json:"recent"`
			} `json:"filings"`
		}
		if err := getJSON(ctx, client, "https://data.sec.gov/submissions/CIK"+cik+".json", headers, &data); err != nil {
			continue
		}
		fetchedCompanies++
		if !sleepContext(ctx, 130*time.Millisecond) {
			return
		}
		categoryCounts := map[string]int{}
		var symbolItems []FilingItem
		for i, form := range data.Filings.Recent.Form {
			if !isSECRelevantForm(form) || i >= len(data.Filings.Recent.Accession) {
				continue
			}
			meaning, category := secFormMeaning(form)
			capByCategory := map[string]int{"insider": 3, "material": 3, "report": 2, "ownership": 2, "offering": 2, "institutional": 1, "governance": 1}
			if categoryCounts[category] >= capByCategory[category] {
				continue
			}
			acc := data.Filings.Recent.Accession[i]
			doc, filed, report, rawDesc, items := "", "", "", "", ""
			if i < len(data.Filings.Recent.PrimaryDocument) {
				doc = data.Filings.Recent.PrimaryDocument[i]
			}
			if i < len(data.Filings.Recent.FilingDate) {
				filed = data.Filings.Recent.FilingDate[i]
			}
			if i < len(data.Filings.Recent.ReportDate) {
				report = data.Filings.Recent.ReportDate[i]
			}
			if i < len(data.Filings.Recent.Description) {
				rawDesc = strings.TrimSpace(data.Filings.Recent.Description[i])
			}
			if i < len(data.Filings.Recent.Items) {
				items = strings.TrimSpace(data.Filings.Recent.Items[i])
			}
			link := fmt.Sprintf("https://www.sec.gov/edgar/browse/?CIK=%d", mustInt(cik))
			if doc != "" {
				link = fmt.Sprintf("https://www.sec.gov/Archives/edgar/data/%d/%s/%s", mustInt(cik), strings.ReplaceAll(acc, "-", ""), doc)
			}
			desc := meaning
			if category == "material" && items != "" {
				desc = meaning + " · Items " + items
			}
			if rawDesc != "" && category == "other" {
				desc = rawDesc
			}
			item := FilingItem{ID: symbol + "-" + acc, Symbol: symbol, Company: data.Name, Form: form, FiledAt: filed, ReportDate: report, Description: desc, Meaning: meaning, Category: category, Items: items, URL: link}
			if category == "insider" && strings.HasPrefix(strings.ToUpper(strings.TrimSpace(form)), "4") && doc != "" {
				enrichForm4(ctx, client, headers, link, &item)
				if !sleepContext(ctx, 130*time.Millisecond) {
					return
				}
			}
			symbolItems = append(symbolItems, item)
			categoryCounts[category]++
			if len(symbolItems) >= 11 {
				break
			}
		}
		filings = append(filings, symbolItems...)
		intel[symbol] = buildSECIntelligence(symbol, symbolItems)
	}
	if fetchedCompanies == 0 && len(symbols) > 0 {
		e.setHealth("filings", "degraded · keeping cached data")
		return
	}
	sort.Slice(filings, func(i, j int) bool { return filings[i].FiledAt > filings[j].FiledAt })
	if len(filings) > 160 {
		filings = filings[:160]
	}
	newFinancialFingerprint := financialFilingFingerprint(filings)
	e.mu.Lock()
	e.filings = filings
	e.secIntelligence = intel
	e.health["filings"] = "healthy · SEC EDGAR"
	e.lastUpdated["filings"] = time.Now().UnixMilli()
	if oldFinancialFingerprint != "" && newFinancialFingerprint != "" && oldFinancialFingerprint != newFinancialFingerprint {
		e.lastUpdated["fundamentals"] = 0
		e.health["fundamentals"] = "stale · new financial filing"
	}
	e.mu.Unlock()
	e.app.broadcastFilings(clone(filings), clone(intel))
	e.enrichEarningsGuidanceFromEvidence()
	e.evaluateCatalystWatch(time.Now())
}
