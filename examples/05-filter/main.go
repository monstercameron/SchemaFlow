package main

import (
	"fmt"
	"os"

	schemaflow "github.com/monstercameron/SchemaFlow"
)

// SupportTicket represents a customer support ticket
type SupportTicket struct {
	ID          int    `json:"id"`
	Customer    string `json:"customer"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
}

func main() {
	// Initialize SchemaFlow
	if err := schemaflow.InitWithEnv(); err != nil {
		schemaflow.GetLogger().Error("Failed to initialize SchemaFlow", "error", err)
		os.Exit(1)
	}

	// Support tickets
	tickets := []SupportTicket{
		{
			ID:          1,
			Customer:    "Alice Johnson",
			Subject:     "Website is completely down",
			Description: "Our entire e-commerce website has been offline for 2 hours. We're losing sales!",
			Priority:    "unknown",
			Status:      "open",
		},
		{
			ID:          2,
			Customer:    "Bob Smith",
			Subject:     "How to change password?",
			Description: "I forgot my password and need help resetting it.",
			Priority:    "unknown",
			Status:      "open",
		},
		{
			ID:          3,
			Customer:    "Carol White",
			Subject:     "Payment processing failure",
			Description: "Customer payments are failing with error code 500. This is blocking all transactions.",
			Priority:    "unknown",
			Status:      "open",
		},
		{
			ID:          4,
			Customer:    "David Brown",
			Subject:     "Feature request: dark mode",
			Description: "It would be nice to have a dark mode option in the settings.",
			Priority:    "unknown",
			Status:      "open",
		},
		{
			ID:          5,
			Customer:    "Eve Davis",
			Subject:     "Data breach suspected",
			Description: "We detected unusual activity and potential unauthorized access to customer data.",
			Priority:    "unknown",
			Status:      "open",
		},
		{
			ID:          6,
			Customer:    "Frank Miller",
			Subject:     "Invoice copy request",
			Description: "Can you send me a copy of invoice #12345 from last month?",
			Priority:    "unknown",
			Status:      "open",
		},
	}

	fmt.Println("🎫 Filter Example - Urgent Ticket Triage")
	fmt.Println("=" + string(make([]byte, 50)))
	fmt.Printf("\n📥 Total Tickets: %d\n", len(tickets))
	fmt.Println("\nAll Tickets:")
	for _, t := range tickets {
		fmt.Printf("  #%d - %s: %s\n", t.ID, t.Customer, t.Subject)
	}

	// Filter to find URGENT tickets only
	criteria := `Identify tickets that require immediate attention:
- Affects multiple users or critical business functions
- Security-related issues
- Service outages or payment failures
Exclude routine requests and feature requests.`

	filterOpts := schemaflow.NewFilterOptions().WithCriteria(criteria)
	filterOpts.OpOptions.Intelligence = schemaflow.Fast
	filterOpts.OpOptions.Steering = "Focus on business impact and urgency"

	urgentTickets, err := schemaflow.Filter(tickets, filterOpts)

	if err != nil {
		schemaflow.GetLogger().Error("Filtering failed", "error", err)
		os.Exit(1)
	}

	// Display urgent tickets
	fmt.Println("\n🚨 URGENT Tickets (require immediate attention):")
	fmt.Println("---")
	if len(urgentTickets) == 0 {
		fmt.Println("  No urgent tickets found")
	} else {
		for _, t := range urgentTickets {
			fmt.Printf("\n⚠️  Ticket #%d - %s\n", t.ID, t.Customer)
			fmt.Printf("   Subject: %s\n", t.Subject)
			fmt.Printf("   Issue: %s\n", t.Description)
		}
	}

	fmt.Printf("\n✨ Success! Filtered %d urgent tickets from %d total\n", len(urgentTickets), len(tickets))
}
