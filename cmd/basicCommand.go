/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	model "simpl/models"
	database "simpl/utils"
	"strconv"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "To initiate or create a new user, merchant and transaction",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n Use simpl new -h or simpl new --help for more info\n")
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Helps to update the discount percent of merchant and credit limit of user",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n simpl update -h or --help \n")
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Helps to give the reports of merchant and user",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\n simpl report -h or --help for more info \n")
	},
}

var paybackCmd = &cobra.Command{
	Use:   "payback",
	Short: "For user to return back the amount borrowed",

	Run: func(cmd *cobra.Command, args []string) {
		sqlDB := database.DB
		var spents int
		Name, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("pass proper user name")
		}
		Amount, err := cmd.Flags().GetString("amount")
		if err != nil {
			fmt.Println("pass proper amount name")
		}
		amountInInt, _ := strconv.Atoi(Amount)
		sqlDB.Raw("select spent from users where user_name IN (?)", Name).Scan(&spents)
		fmt.Println(Name, Amount, amountInInt, spents)
		if amountInInt > spents {
			fmt.Println("Payback is more than spent. The current spent value is :", spents)
			os.Exit(1)
		} else {
			amountInInt = spents - amountInInt
			amountInMap := make(map[string]int)
			amountInMap["spent"] = amountInInt
			sqlDB.Model(model.Users{}).Where("user_name IN (?)", Name).Update("spent", amountInInt)
		}
		sqlDB.Raw("select spent from users where user_name IN (?)", Name).Scan(&spents)
		fmt.Printf("%v (dues: %d)", Name, spents)
	},
}

func init() {
	var Name string
	var Amount string
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(paybackCmd)
	paybackCmd.Flags().StringVarP(&Name, "username", "n", "", "To get the User Name")
	paybackCmd.Flags().StringVarP(&Amount, "amount", "c", "", "To get the money to be refunded")
	paybackCmd.MarkFlagRequired("username")
	paybackCmd.MarkFlagRequired("amount")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
