/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"time"

	"github.com/Stefasaurus/llmsger/cmd"
)

func main() {
	start := time.Now()

	cmd.Execute()

	elapsed := time.Since(start)
	log.Printf("Program took %s", elapsed)
}
