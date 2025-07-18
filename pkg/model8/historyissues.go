package model8

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Status string

const (
	Unreviewed    Status = "U"  // All vulnerabilities found by a scan that need to be reviewed. This is your starting point when reviewing scan results.
	Verified      Status = "V"  // Vulnerabilities that users investigated and determined to have legitimate risk. Users change the status to Verified to show that it needs to be remediated.
	Ignored       Status = "I"  // Vulnerabilities that were identified as potentially harmful, but users reviewed and marked Ignored after. You may want to replay the attack or otherwise verify that the Ignored vulnerabilities do not pose a threat.
	Falsepositive Status = "FP" // Vulnerabilities that users flagged as having been incorrectly found by the InsightAppSec. Users can change the status of False Positive vulnerabilities during the investigation process. This status does not change in subsequent scans.
	Remediated    Status = "R"  // Vulnerabilities that were identified, investigated, and fixed. Users and validation scans can change the status to Remediated. If an issue is rediscovered in a subsequent scan, the status reverts back to Unreviewed.
)

type Historyissue8 struct {
	Id             uuid.UUID `json:"id"`
	Endpointid     uuid.UUID `json:"Endpointid"`
	Issue          Issue8    `json:"issue"`
	Url            string    `json:"url"`
	Signature      string    `json:"signature"`
	Status         Status    `json:"status"`
	FoundFirsttime time.Time `json:"foundFirstTime"`
}

func DifferenceHistoryissues8(slice1, slice2 []Historyissue8) []Historyissue8 {
	// Create a map to hold the elements of slice2 for easy lookup
	lookupMap := make(map[string]bool)
	for _, item := range slice2 {
		lookupMap[string(item.Signature)] = true
	}

	// Iterate through slice1 and add elements that are not in slice2
	var result []Historyissue8
	for _, item := range slice1 {
		if !lookupMap[string(item.Signature)] {
			result = append(result, item)
		}
	}

	return result
}

func RemoveDuplicatesHistoryissues8(slice []Historyissue8) []Historyissue8 {
	seen := make(map[string]bool)
	result := []Historyissue8{}
	for _, h8 := range slice {
		val := string(h8.Signature)
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, h8)
		}
	}
	return result
}

// Return highest risk severity found: critical, high or normal
func ExistCriticalOrHighRiskSeverityHistoryissue8(slice []Historyissue8) string {
	var severity = "normal"
	for _, h8 := range slice {
		if h8.Issue.Severity == "high" {
			severity = "high"
			continue
		}
		if h8.Issue.Severity == "critical" {
			severity = "critical"
			return severity
		}
	}
	return severity
}
