package types

import (
	"fmt"
	"os"
	"testing"
)

func TestLoad2(t *testing.T) {
	data, err := os.ReadFile("cve/sample_report.json")
	if err != nil {
		t.Fatal(fmt.Errorf("failed to read report: %s", err.Error()))
	}
	r, err := NewCveReport(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("medium %d\n", r.Medium())
	fmt.Printf("low %d\n", r.Low())
	for _, cve := range r.Cves {
		fmt.Printf("CVE: %s (fixed: %t)\n", cve.Id, cve.Fixed())
	}
}
