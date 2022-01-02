/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	model "simpl/models"
	database "simpl/utils"

	percentage "github.com/dariubs/percent"
	"github.com/spf13/cobra"
)

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// userCmd represents the user command
var userNewCmd = &cobra.Command{
	Use:   "user",
	Short: "\n A brief description of your usernewcmd",

	Run: func(cmd *cobra.Command, args []string) {
		sqlDB := database.DB
		Name, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("pass proper username")
		}
		Email, err := cmd.Flags().GetString("email")
		if err != nil {
			fmt.Println("pass proper Email")
		}
		value := isEmailValid(Email)
		if !value {
			fmt.Println("Enter proper mail address")
			os.Exit(1)
		}
		limit, err := cmd.Flags().GetString("creditlimit")
		if err != nil {
			fmt.Println("pass proper limit value")
		}
		credit, _ := strconv.Atoi(limit)
		users := model.Users{User_name: Name, Email: Email, Credit_limit: credit, Spent: 0}
		val := sqlDB.Select("user_name", "email", "credit_limit", "spent").Create(&users)
		if val.RowsAffected > 0 {
			fmt.Println(Name, "(", limit, ")")
		} else {
			fmt.Println(val.Error)
		}
	},
}

// merchantCmd represents the merchant command
var merchantNewCmd = &cobra.Command{
	Use:   "merchant",
	Short: "\n Creates new merchant",

	Run: func(cmd *cobra.Command, args []string) {
		sqlDB := database.DB
		Name, err := cmd.Flags().GetString("merchantname")
		if err != nil {
			fmt.Println("pass proper merchant name")
		}
		Email, err := cmd.Flags().GetString("email")
		if err != nil {
			fmt.Println("pass proper Email")
		}
		value := isEmailValid(Email)
		if !value {
			fmt.Println("Enter proper mail address")
			os.Exit(1)
		}
		Discount, err := cmd.Flags().GetString("discount")
		if err != nil {
			fmt.Println("pass proper merchant name")
		}

		disnt, _ := strconv.ParseFloat(Discount, 64)
		merchant := model.Merchants{Merchant_name: Name, Email: Email, Discount: disnt, Total_discount: 0}
		val := sqlDB.Select("merchant_name", "email", "discount", "total_discount").Create(&merchant)
		if val.RowsAffected > 0 {
			fmt.Printf("%v (%0.2F %%)", Name, disnt)
		} else {
			fmt.Println(val.Error)
		}
	},
}

// To update merchant
var merchantUpdateCmd = &cobra.Command{
	Use:   "merchant",
	Short: "\n Updates the merchant discount values",

	Run: func(cmd *cobra.Command, args []string) {
		sqlDB := database.DB
		Name, err := cmd.Flags().GetString("merchantname")
		if err != nil {
			fmt.Println("pass proper merchant name")
		}
		Discount, err := cmd.Flags().GetString("discount")
		if err != nil {
			fmt.Println("pass proper discount value")
		}

		disnt, _ := strconv.ParseFloat(Discount, 64)
		val := sqlDB.Model(model.Merchants{}).Where("merchant_name IN (?)", Name).Update("discount", disnt)
		if val.RowsAffected <= 0 {
			fmt.Println("Enter proper merchant name or updated discount")
		}
	},
}

// txnCmd represents the txn command
var txnCmd = &cobra.Command{
	Use:   "txn",
	Short: "\n A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		// Variable decleration
		var value model.User_spent
		var percent model.Merchant_discount
		//database connection
		sqlDB := database.DB
		// Getting the data from cmd
		UserName, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("pass proper merchant name")
		}
		merchantName, err := cmd.Flags().GetString("merchantname")
		if err != nil {
			fmt.Println("pass proper merchant name")
		}
		Amount, err := cmd.Flags().GetString("amount")
		if err != nil {
			fmt.Println("pass proper merchant name")
		}
		amountInInt, _ := strconv.Atoi(Amount)
		// Getting the amount details from user
		val := sqlDB.Raw("select credit_limit-spent as balance , spent  from users where user_name IN (?)", UserName).Scan(&value)
		if val.RowsAffected <= 0 {
			fmt.Println("Enter valid username")
			os.Exit(1)
		}
		if value.Balance < amountInInt {
			fmt.Println("rejected! (reason: credit limit)")
			os.Exit(1)
		} else {
			val = sqlDB.Raw("select discount , total_discount from merchants where merchant_name IN (?)", merchantName).Scan(&percent)
			if val.RowsAffected <= 0 {
				fmt.Println("Enter valid merchant name")
				os.Exit(1)
			}

			percent.Total_discount = percentage.PercentFloat(percent.Discount, float64(amountInInt)) + percent.Total_discount
			//updating merchant discount
			valMerchant := sqlDB.Model(model.Merchants{}).Where("merchant_name = (?)", merchantName).Update("total_discount", percent.Total_discount)
			if valMerchant.RowsAffected <= 0 {
				fmt.Println("merchant details not updated in merchant")
				os.Exit(1)
			}
			//Updating user spent
			val = sqlDB.Model(model.Users{}).Where("user_name IN (?)", UserName).Update("spent", value.Spent+amountInInt)
			if val.RowsAffected <= 0 {
				fmt.Println("user details not updated in user ", value.Balance+amountInInt)
				os.Exit(1)
			}
			fmt.Println("success!")
		}

	},
}

// report command
var discount = &cobra.Command{
	Use:   "discount",
	Short: "\n What is the total discount provided by a merchant",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("discount called")
		sqlDB := database.DB
		var value float64
		MerchantName, err := cmd.Flags().GetString("merchantname")
		if err != nil {
			fmt.Println("pass proper user name")
		}
		fmt.Println(MerchantName)
		val := sqlDB.Raw("select total_discount from merchants where merchant_name IN (?)", MerchantName).Scan(&value)
		if val.RowsAffected <= 0 {
			fmt.Println("Enter valid merchant name")
			os.Exit(1)
		}
		fmt.Println(value)
	},
}

var dues = &cobra.Command{
	Use:   "dues",
	Short: "\n Gives the report of a particular user",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dues called")
		sqlDB := database.DB
		var value int
		UserName, err := cmd.Flags().GetString("username")
		if err != nil {
			fmt.Println("pass proper user name")
		}
		fmt.Println(UserName)
		val := sqlDB.Raw("select spent from users where user_name IN (?)", UserName).Scan(&value)
		if val.RowsAffected <= 0 {
			fmt.Println("enter valid username")
			os.Exit(1)
		}
		fmt.Printf("Dues for %v is: %d", UserName, value)

	},
}

var usersAtLimit = &cobra.Command{
	Use:   "users-at-credit-limit",
	Short: "\n Users who have no credit left",

	Run: func(cmd *cobra.Command, args []string) {
		sqlDB := database.DB
		var value []string
		sqlDB.Raw("select user_name from users where credit_limit-spent = 0").Scan(&value)
		for i := 0; i < len(value); i++ {
			fmt.Println(value[i])
		}
	},
}

var totalDues = &cobra.Command{
	Use:   "total-dues",
	Short: "\n Gives list of all the user with dues",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("total dues called")
		sqlDB := database.DB
		var value []model.User_report
		var sum int
		sqlDB.Raw("select user_name, spent from users where spent > 0").Scan(&value)
		for i := 0; i < len(value); i++ {
			sum = sum + value[i].Spent
			fmt.Printf("%v : %d \n", value[i].User_name, value[i].Spent)
		}
		fmt.Println("Total :", sum)
	},
}

func init() {

	var Name string
	var Email string
	var Limit string
	var Discount string
	var MerchantName string
	var Amount string
	// creating new user
	newCmd.AddCommand(userNewCmd)
	userNewCmd.Flags().StringVarP(&Name, "username", "n", "", "Getting the user name")
	userNewCmd.Flags().StringVarP(&Email, "email", "m", "", "Getting the email address")
	userNewCmd.Flags().StringVarP(&Limit, "creditlimit", "c", "", "Getting the credit limit")
	userNewCmd.MarkFlagRequired("username")
	userNewCmd.MarkFlagRequired("email")
	userNewCmd.MarkFlagRequired("creditlimit")

	//merchant command
	newCmd.AddCommand(merchantNewCmd)
	merchantNewCmd.Flags().StringVarP(&Name, "merchantname", "n", "", "Getting the merchant name")
	merchantNewCmd.Flags().StringVarP(&Email, "email", "m", "", "Getting the email address")
	merchantNewCmd.Flags().StringVarP(&Discount, "discount", "d", "", "Getting the discount percent")
	merchantNewCmd.MarkFlagRequired("merchantname")
	merchantNewCmd.MarkFlagRequired("email")
	merchantNewCmd.MarkFlagRequired("discount")

	//Update merchant
	updateCmd.AddCommand(merchantUpdateCmd)
	merchantUpdateCmd.Flags().StringVarP(&Name, "merchantname", "n", "", "Getting the merchant name")
	merchantUpdateCmd.Flags().StringVarP(&Discount, "discount", "d", "", "Getting the update discount percent")
	merchantUpdateCmd.MarkFlagRequired("merchantname")
	merchantUpdateCmd.MarkFlagRequired("discount")

	// txm command
	newCmd.AddCommand(txnCmd)
	txnCmd.Flags().StringVarP(&Name, "username", "n", "", "Getting the user name")
	txnCmd.Flags().StringVarP(&MerchantName, "merchantname", "m", "", "Getting the email address")
	txnCmd.Flags().StringVarP(&Amount, "amount", "c", "", "Getting the credit limit")
	txnCmd.MarkFlagRequired("username")
	txnCmd.MarkFlagRequired("merchantname")
	txnCmd.MarkFlagRequired("amount")

	//Report command
	// Discount provided by merchant
	reportCmd.AddCommand(discount)
	discount.Flags().StringVarP(&Name, "merchantname", "n", "", "Getting the merchant name")
	discount.MarkFlagRequired("merchantname")
	//Dues pending for user
	reportCmd.AddCommand(dues)
	dues.Flags().StringVarP(&Name, "username", "n", "", "Getting the user name")
	dues.MarkFlagRequired("username")
	// Users who are at limit
	reportCmd.AddCommand(usersAtLimit)
	// All the dues of user
	reportCmd.AddCommand(totalDues)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
