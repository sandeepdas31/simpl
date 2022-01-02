/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"simpl/cmd"
	database "simpl/utils"
)

func main() {
	database.GetInstancemysql()
	cmd.Execute()
}
