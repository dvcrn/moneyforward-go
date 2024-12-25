package main

import (
	"fmt"
	"log"

	"github.com/dvcrn/moneyforward-go"
)

func main() {
	// Initialize client
	cookieString := "" // Replace with your MF cookie string
	client := moneyforward.NewClient(cookieString)

	// Get account summaries
	accounts, err := client.GetAccountSummaries()
	if err != nil {
		log.Fatal("Failed to get accounts:", err)
	}

	// Print account balances
	for _, account := range accounts.Accounts {
		fmt.Printf("Account: %s, Balance: ¥%.2f\n", account.Name, account.Amount)
	}

	// Get recent transactions
	transactions, err := client.GetUserAssetActs(moneyforward.UserAssetActsParams{
		Size:         16,
		IsNew:        true,
		IsContinuous: true,
	})
	if err != nil {
		log.Fatal("Failed to get transactions:", err)
	}

	// Print recent transactions
	fmt.Println("\nRecent Transactions:")
	for _, act := range transactions.UserAssetActs {
		fmt.Printf("%s: ¥%.2f - %s\n",
			act.RecognizedAt.Format("2006-01-02"),
			act.Amount,
			act.Content,
		)
	}

	// Get home timeline
	timeline, err := client.GetHomeTimeline(10)
	if err != nil {
		log.Fatal("Failed to get home timeline:", err)
	}
	fmt.Println("\nHome Timeline:")
	for _, t := range timeline.Timeline {
		fmt.Printf("Date: %s\n", t.Date)
		for _, card := range t.Cards {
			if card.UserNotification != nil {
				fmt.Printf("  - Notification: %+v\n", card.UserNotification.Parameters)
			}
			if card.HomeCard != nil {
				fmt.Printf("  - Card: %s (Valid: %s - %s)\n",
					card.HomeCard.ID,
					card.HomeCard.StartAt.Format("2006-01-02"),
					card.HomeCard.EndAt.Format("2006-01-02"),
				)
			}
		}
	}

	// Force update
	err = client.ForceUpdate()
	if err != nil {
		log.Fatal("Failed to force update:", err)
	}
	fmt.Println("\nForced update successful")

	// Get all transactions
	allTransactions, err := client.GetTransactions()
	if err != nil {
		log.Fatal("Failed to get all transactions:", err)
	}
	fmt.Println("\nRecommended Services:")
	for _, service := range allTransactions.EmptyState.RecommendedServices.Services {
		fmt.Printf("  - %s (%s) [Category: %s]\n",
			service.ServiceName,
			service.ServiceType,
			service.ServiceCategory.CategoryType,
		)
	}

	// Get specific user asset act
	if len(transactions.UserAssetActs) > 0 {
		act, err := client.GetUserAssetAct(transactions.UserAssetActs[0].ID.String())
		if err != nil {
			log.Fatal("Failed to get specific user asset act:", err)
		}
		fmt.Println("\nTransaction Details:")
		fmt.Printf("  Content: %s\n", act.UserAssetAct.Content)
		fmt.Printf("  Amount: ¥%.2f\n", act.UserAssetAct.Amount)
		fmt.Printf("  Date: %s\n", act.UserAssetAct.RecognizedAt.Format("2006-01-02"))
		fmt.Printf("  Account: %s (%s)\n",
			act.UserAssetAct.SubAccount.SubName,
			act.UserAssetAct.Account.Service.ServiceName,
		)
	}

	// Get specific account
	account, err := client.GetAccount(accounts.Accounts[0].ShowPath)
	if err != nil {
		log.Fatal("Failed to get specific account:", err)
	}
	fmt.Println("\nAccount Details:")
	fmt.Printf("  Name: %s\n", account.Account.DisplayName)
	fmt.Printf("  Total Assets: ¥%.2f\n", account.Account.TotalAsset)
	fmt.Printf("  Total Liabilities: ¥%.2f\n", account.Account.TotalLiability)
	fmt.Printf("  Last Updated: %s\n", account.Account.LastAggregatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("  Sub Accounts:")
	for _, sub := range account.Account.SubAccounts {
		fmt.Printf("    - %s (%s)\n", sub.SubName, sub.SubType)
		for _, summary := range sub.UserAssetDetSummaries {
			fmt.Printf("      %s: ¥%.2f\n", summary.AssetSubclassName, summary.JPYValue)
		}
	}
}
