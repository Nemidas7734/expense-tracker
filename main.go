package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"strconv"
)

type Expense struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

var expenses []Expense
var nextID int
var totalExpense int
var budget int

func loadExpensesFromFile() {
	file, err := os.Open("expenses.json")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error reading file", err)
		}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&expenses)
	if err != nil && err.Error() != "EOF" {
		fmt.Println("Error decoding file", err)
	}
	for _, expense := range expenses {
		if expense.ID >= nextID {
			nextID = expense.ID + 1
		}
		totalExpense += expense.Amount
	}
}

func writeExpenseToFile() {
	file, err := os.Create("expenses.json")
	if err != nil {
		fmt.Println("Error creating file", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(expenses)
	if err != nil {
		fmt.Println("Error encoding file", err)
		return
	}
}

func addExpense(description string, amount int) {
	expense := Expense{
		ID:          nextID,
		Description: description,
		Amount:      amount,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expenses = append(expenses, expense)
	nextID++
	totalExpense += amount
	writeExpenseToFile()
	fmt.Printf("Expense added successfully (ID: %d)\n", expense.ID)
}

func updateExpense(id int, description string, amount int) {
	for i := range expenses {
		if expenses[i].ID == id {
			totalExpense -= expenses[i].Amount
			expenses[i].Amount = amount
			expenses[i].Description = description
			expenses[i].UpdatedAt = time.Now()

			writeExpenseToFile()
			totalExpense += amount
			fmt.Printf("Expense updated successfully (ID: %d)\n", id)
			return
		}
	}
	fmt.Println("Expense not found")
}

func deleteExpense(id int) {
	for i := range expenses {
		if expenses[i].ID == id {
			totalExpense -= expenses[i].Amount
			expenses = append(expenses[:i], expenses[i+1:]...)
			writeExpenseToFile()
			fmt.Printf("Expense deleted successfully (ID: %d)\n", id)
			return
		}
	}
	fmt.Println("Expense not found")
}

func listExpenses() {
	if len(expenses) == 0 {
		fmt.Println("No expenses found.")
		return
	}

	fmt.Println("ID  Date       Description  Amount")
	for _, expense := range expenses {
		fmt.Printf("%d   %s  %s  $%d\n", expense.ID, expense.CreatedAt.Format("2006-01-02"), expense.Description, expense.Amount)
	}
}

func showSummary() {
	fmt.Printf("Total expenses: $%d\n", totalExpense)
	if budget > 0 && totalExpense > budget {
		fmt.Printf("Warning: You have exceeded your monthly budget of $%d\n", budget)
	}
}

func showMonthlySummary(month int) {
	var monthTotal int
	for _, expense := range expenses {
		if expense.CreatedAt.Month() == time.Month(month) && expense.CreatedAt.Year() == time.Now().Year() {
			monthTotal += expense.Amount
		}
	}
	fmt.Printf("Total expenses for %s: $%d\n", time.Month(month).String(), monthTotal)
}

func setBudget(amount int) {
	budget = amount
	fmt.Printf("Monthly budget set to: $%d\n", budget)
}

func showMenu() {
	fmt.Println("\n--- Expense Tracker ---")
	fmt.Println("1. Add an expense")
	fmt.Println("2. Update an expense")
	fmt.Println("3. Delete an expense")
	fmt.Println("4. List all expenses")
	fmt.Println("5. View summary of expenses")
	fmt.Println("6. View summary of expenses by month")
	fmt.Println("7. Set a monthly budget")
	fmt.Println("8. Exit")
	fmt.Print("Choose an option (1-8): ")
}

func main() {
	loadExpensesFromFile()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		showMenu()
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1": 
			fmt.Print("Enter description: ")
			scanner.Scan()
			description := scanner.Text()

			fmt.Print("Enter amount: ")
			scanner.Scan()
			amount, err := strconv.Atoi(scanner.Text())
			if err != nil || amount <= 0 {
				fmt.Println("Invalid amount.")
				continue
			}

			addExpense(description, amount)

		case "2": 
			fmt.Print("Enter ID of the expense to update: ")
			scanner.Scan()
			id, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Invalid ID.")
				continue
			}

			fmt.Print("Enter new description: ")
			scanner.Scan()
			description := scanner.Text()

			fmt.Print("Enter new amount: ")
			scanner.Scan()
			amount, err := strconv.Atoi(scanner.Text())
			if err != nil || amount <= 0 {
				fmt.Println("Invalid amount.")
				continue
			}

			updateExpense(id, description, amount)

		case "3":
			fmt.Print("Enter ID of the expense to delete: ")
			scanner.Scan()
			id, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Invalid ID.")
				continue
			}

			deleteExpense(id)

		case "4": 
			listExpenses()

		case "5": 
			showSummary()

		case "6": 
			fmt.Print("Enter month (1-12): ")
			scanner.Scan()
			month, err := strconv.Atoi(scanner.Text())
			if err != nil || month < 1 || month > 12 {
				fmt.Println("Invalid month.")
				continue
			}
			showMonthlySummary(month)

		case "7": 
			fmt.Print("Enter monthly budget: ")
			scanner.Scan()
			budgetAmount, err := strconv.Atoi(scanner.Text())
			if err != nil || budgetAmount <= 0 {
				fmt.Println("Invalid budget amount.")
				continue
			}
			setBudget(budgetAmount)

		case "8": 
			fmt.Println("Exiting application...")
			return

		default:
			fmt.Println("Invalid option. Please choose a number between 1 and 8.")
		}
	}
}
