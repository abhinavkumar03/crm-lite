package seeders

import (
	"fmt"
	"math/rand"
	"strings"
)

// Shared, deterministic datasets used by the demo seeders. A fixed RNG seed
// keeps generated data stable across runs so counts and screenshots are
// reproducible.

const (
	demoOrgSlug = "crm-lite-demo"
	demoOrgName = "CRM Lite Demo Co"
)

var (
	firstNames = []string{
		"Aarav", "Vivaan", "Aditya", "Vihaan", "Arjun", "Sai", "Reyansh",
		"Krishna", "Ishaan", "Rohan", "Priya", "Ananya", "Diya", "Isha",
		"Sneha", "Aisha", "Kavya", "Riya", "Meera", "Neha",
	}

	lastNames = []string{
		"Sharma", "Verma", "Gupta", "Reddy", "Nair", "Iyer", "Patel",
		"Singh", "Kumar", "Das", "Bose", "Mehta", "Shah", "Rao",
		"Chopra", "Malhotra",
	}

	companyNames = []string{
		"Tata Digital", "Reliance Retail", "Infosys Solutions", "Wipro Systems",
		"Mahindra Logistics", "Adani Enterprises", "Bajaj Finserv", "Godrej Consumer",
		"HCL Technologies", "Zomato Foods", "Swiggy Kitchens", "Flipkart Commerce",
		"Paytm Payments", "Ola Mobility", "Freshworks CRM", "Zoho Labs",
		"Razorpay Fintech", "Dabur India", "Asian Paints", "Titan Industries",
	}

	cities = []string{
		"Mumbai", "Delhi", "Bengaluru", "Hyderabad", "Chennai",
		"Pune", "Kolkata", "Ahmedabad", "Jaipur", "Surat",
	}

	industries = []string{
		"IT Services", "Manufacturing", "Retail", "Fintech", "Healthcare",
		"Education", "Logistics", "Real Estate", "FMCG", "Automotive",
	}

	jobTitles = []string{
		"CEO", "CTO", "Sales Manager", "Procurement Head", "Founder",
		"VP Engineering", "Marketing Lead", "Operations Manager",
	}

	// Weighted so the pipeline looks realistic (more early-stage than closed).
	leadStatuses = []string{
		"NEW", "NEW", "NEW", "CONTACTED", "CONTACTED", "CONTACTED",
		"QUALIFIED", "QUALIFIED", "WON", "LOST",
	}

	taskStatuses = []string{
		"PENDING", "PENDING", "IN_PROGRESS", "IN_PROGRESS", "COMPLETED",
	}

	dealStages = []string{
		"Prospecting", "Qualification", "Proposal", "Negotiation", "Closed Won", "Closed Lost",
	}

	tagPool = []string{"vip", "inbound", "referral", "enterprise", "smb", "hot", "cold"}
)

func pick(r *rand.Rand, items []string) string {
	return items[r.Intn(len(items))]
}

func fullName(r *rand.Rand) (string, string) {
	return pick(r, firstNames), pick(r, lastNames)
}

func phone(r *rand.Rand) string {
	// Indian mobile format: +91 followed by a 10-digit number starting 6-9.
	return fmt.Sprintf("+91 %d%09d", 6+r.Intn(4), r.Intn(1_000_000_000))
}

func emailFrom(first, last string, r *rand.Rand) string {
	local := strings.ToLower(fmt.Sprintf("%s.%s%d", first, last, r.Intn(90)+10))
	return local + "@example.in"
}

func website(company string) string {
	slug := strings.ToLower(company)
	slug = strings.ReplaceAll(slug, " ", "")
	return "https://" + slug + ".example.in"
}

func pickTags(r *rand.Rand) []string {
	n := 1 + r.Intn(2)
	seen := map[string]bool{}
	out := make([]string, 0, n)
	for len(out) < n {
		t := pick(r, tagPool)
		if !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}
