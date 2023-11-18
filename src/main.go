package main

import (
	"fmt"
	"os"
	"github.com/akamensky/argparse"
	"github.com/heaveless/bds-downloader/utils"
)

const (
	ColorRed    = "\033[91m"
	ColorGreen  = "\033[92m"
	ColorYellow = "\033[93m"
	ColorBlue   = "\033[94m"
	ColorReset  = "\033[0m"
)

func main() {
	parser := argparse.NewParser("bds-down", "Download and install BDS.")
	skipAgreePtr := parser.Flag("y", "yes", &argparse.Options{Required: false, Help: "Skip the agreement" })
	parser.Parse(os.Args)

	skipAgree := *skipAgreePtr

	fmt.Print(`Before using this software, please read: 
	- Minecraft End User License Agreement   https://minecraft.net/terms
	- Microsoft Privacy Policy               https://go.microsoft.com/fwlink/?LinkId=521839
	Please enter y if you agree with the above terms: `)

	var agree string
	if skipAgree {
		agree = "y"
	} else {
		fmt.Scanln(&agree)
	}

	if agree != "y" {
		fmt.Println(ColorYellow + "You must agree with the above terms to use this software." + ColorReset)
		return
	}

	var ver, err = utils.GetLatestVersion()
	if err != nil {
		fmt.Println(ColorRed+"ERROR:", err)
	}

	fmt.Println("Latest version: " + ColorBlue + ver + ColorReset)
	err = utils.Install(ver)
	if err != nil {
		fmt.Println(ColorRed+"ERROR:", err, ColorReset)
	}

	fmt.Println(ColorGreen + "Install complete." + ColorReset)
	return
}