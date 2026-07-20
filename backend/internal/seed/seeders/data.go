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
	primaryOrgSlug = "crm-lite-demo"
	demoUserEmail  = "demo@crmlite.com"
	demoUserPass   = "Password@123"
)

// workspaceDef describes a demo workspace (organization).
type workspaceDef struct {
	Name        string
	Slug        string
	Description string
	Industry    string
	CompanySize string
	Country     string
	LogoURL     string
	Timezone    string
	Currency    string
	Locale      string
	DateFormat  string
	Plan        string
}

var demoWorkspaces = []workspaceDef{
	{
		Name: "CRM Lite Demo", Slug: primaryOrgSlug,
		Description: "Primary showcase workspace with a full SaaS sales pipeline.",
		Industry: "Technology", CompanySize: "51-200", Country: "IN",
		LogoURL: "https://res.cloudinary.com/demo/image/upload/v1312461204/sample.jpg",
		Timezone: "Asia/Kolkata", Currency: "INR", Locale: "en-IN", DateFormat: "DD/MM/YYYY",
		Plan: "pro",
	},
	{
		Name: "SME Junction", Slug: "sme-junction",
		Description: "Education and SMB lead engine for college counselling partners.",
		Industry: "Education", CompanySize: "11-50", Country: "IN",
		LogoURL: "https://res.cloudinary.com/demo/image/upload/docs/models.jpg",
		Timezone: "Asia/Kolkata", Currency: "INR", Locale: "en-IN", DateFormat: "DD/MM/YYYY",
		Plan: "pro",
	},
	{
		Name: "Acme Manufacturing", Slug: "acme-manufacturing",
		Description: "Industrial B2B CRM for plant equipment and spare-parts sales.",
		Industry: "Manufacturing", CompanySize: "201-500", Country: "US",
		LogoURL: "https://res.cloudinary.com/demo/image/upload/v1312461204/sample.jpg",
		Timezone: "America/Chicago", Currency: "USD", Locale: "en-US", DateFormat: "MM/DD/YYYY",
		Plan: "pro",
	},
	{
		Name: "Bright Marketing Agency", Slug: "bright-marketing",
		Description: "Agency pipeline for retainers, campaigns, and creative retainers.",
		Industry: "Marketing", CompanySize: "11-50", Country: "GB",
		LogoURL: "https://res.cloudinary.com/demo/image/upload/docs/models.jpg",
		Timezone: "Europe/London", Currency: "GBP", Locale: "en-GB", DateFormat: "DD/MM/YYYY",
		Plan: "pro",
	},
	{
		Name: "Personal Sales CRM", Slug: "personal-sales",
		Description: "Solo seller workspace for freelance deals and warm intros.",
		Industry: "Consulting", CompanySize: "1-10", Country: "IN",
		LogoURL: "https://res.cloudinary.com/demo/image/upload/v1312461204/sample.jpg",
		Timezone: "Asia/Kolkata", Currency: "INR", Locale: "en-IN", DateFormat: "YYYY-MM-DD",
		Plan: "free",
	},
}

// Backward-compatible aliases used by older seeder helpers.
const (
	demoOrgSlug = primaryOrgSlug
	demoOrgName = "CRM Lite Demo"
)

var (
	firstNames = []string{
		"Aarav", "Vivaan", "Aditya", "Vihaan", "Arjun", "Sai", "Reyansh",
		"Krishna", "Ishaan", "Rohan", "Priya", "Ananya", "Diya", "Isha",
		"Sneha", "Aisha", "Kavya", "Riya", "Meera", "Neha", "James", "Emma",
		"Oliver", "Sophia", "Liam", "Olivia", "Noah", "Ava", "Lucas", "Mia",
	}

	lastNames = []string{
		"Sharma", "Verma", "Gupta", "Reddy", "Nair", "Iyer", "Patel",
		"Singh", "Kumar", "Das", "Bose", "Mehta", "Shah", "Rao",
		"Chopra", "Malhotra", "Anderson", "Brooks", "Carter", "Dixon",
	}

	companyNames = []string{
		"Tata Digital", "Reliance Retail", "Infosys Solutions", "Wipro Systems",
		"Mahindra Logistics", "Adani Enterprises", "Bajaj Finserv", "Godrej Consumer",
		"HCL Technologies", "Zomato Foods", "Swiggy Kitchens", "Flipkart Commerce",
		"Paytm Payments", "Ola Mobility", "Freshworks CRM", "Zoho Labs",
		"Razorpay Fintech", "Dabur India", "Asian Paints", "Titan Industries",
		"Acme Components", "Brightline Media", "Northwind Traders", "Contoso Labs",
		"Fabrikam Steel", "Adventure Works", "Blue Yonder Co", "Summit Analytics",
	}

	cities = []string{
		"Mumbai", "Delhi", "Bengaluru", "Hyderabad", "Chennai",
		"Pune", "Kolkata", "Ahmedabad", "Jaipur", "Surat",
		"Chicago", "Austin", "London", "Manchester", "Singapore",
	}

	countries = []string{"IN", "US", "GB", "SG", "AE"}

	industries = []string{
		"IT Services", "Manufacturing", "Retail", "Fintech", "Healthcare",
		"Education", "Logistics", "Real Estate", "FMCG", "Automotive",
	}

	companyStatuses = []string{
		"Prospect", "Active", "Partner", "Churned",
	}

	jobTitles = []string{
		"CEO", "CTO", "Sales Manager", "Procurement Head", "Founder",
		"VP Engineering", "Marketing Lead", "Operations Manager",
	}

	leadStatuses = []string{
		"NEW", "NEW", "NEW", "CONTACTED", "CONTACTED", "CONTACTED",
		"QUALIFIED", "QUALIFIED", "WON", "LOST",
	}

	leadStatusOptions = []string{
		"NEW", "CONTACTED", "QUALIFIED", "WON", "LOST",
	}

	leadSources = []string{
		"Website", "Referral", "Cold Call", "Trade Show", "LinkedIn", "Email Campaign",
	}

	taskStatuses = []string{
		"PENDING", "PENDING", "IN_PROGRESS", "IN_PROGRESS", "COMPLETED",
	}

	taskStatusOptions = []string{
		"PENDING", "IN_PROGRESS", "COMPLETED",
	}

	taskPriorities = []string{
		"Low", "Medium", "High",
	}

	dealStages = []string{
		"Prospecting", "Qualification", "Proposal", "Negotiation", "Closed Won", "Closed Lost",
	}

	tagPool = []string{"vip", "inbound", "referral", "enterprise", "smb", "hot", "cold"}

	noteBodies = []string{
		"Discussed budget and timeline; follow up next week.",
		"Sent product brochure and pricing sheet.",
		"Decision maker joined the call — positive signal.",
		"Needs internal approval before moving to proposal.",
		"Requested a technical demo for their ops team.",
		"Competitor evaluation in progress; stay close.",
	}

	callSummaries = []string{
		"Intro call — mapped stakeholders and pain points.",
		"Discovery call — confirmed budget range and buying process.",
		"Missed call; left voicemail requesting a callback.",
		"Busy signal; will retry tomorrow morning.",
		"Closing call — reviewed contract terms and next steps.",
		"Support follow-up after onboarding kickoff.",
	}

	attachmentNames = []string{
		"proposal.pdf", "pricing-sheet.xlsx", "brand-guidelines.pdf",
		"contract-draft.docx", "demo-recording.mp4", "site-photo.jpg",
	}

	attachmentURLs = []string{
		"https://res.cloudinary.com/demo/image/upload/v1312461204/sample.jpg",
		"https://res.cloudinary.com/demo/image/upload/docs/models.jpg",
		"https://res.cloudinary.com/demo/raw/upload/sample_document.pdf",
	}
)

func pick(r *rand.Rand, items []string) string {
	return items[r.Intn(len(items))]
}

func fullName(r *rand.Rand) (string, string) {
	return pick(r, firstNames), pick(r, lastNames)
}

func phone(r *rand.Rand) string {
	return fmt.Sprintf("+91 %d%09d", 6+r.Intn(4), r.Intn(1_000_000_000))
}

func emailFrom(first, last string, r *rand.Rand) string {
	local := strings.ToLower(fmt.Sprintf("%s.%s%d", first, last, r.Intn(90)+10))
	return local + "@example.com"
}

func website(company string) string {
	slug := strings.ToLower(company)
	slug = strings.ReplaceAll(slug, " ", "")
	return "https://" + slug + ".example.com"
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
