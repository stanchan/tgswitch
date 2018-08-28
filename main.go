package main

/*
* Version 0.3.0
* Compatible with Mac OS X ONLY
 */

/*** OPERATION WORKFLOW ***/
/*
* 1- Create /usr/local/terragrunt directory if does not exist
* 2- Download binary file from url to /usr/local/terragrunt
* 3- Rename the file from `terragrunt` to `terragrunt_version`
* 4- Read the existing symlink for terragrunt (Check if it's a homebrew symlink)
* 6- Remove that symlink (Check if it's a homebrew symlink)
* 7- Create new symlink to binary  `terragrunt_version`
 */

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/pborman/getopt"
	lib "github.com/warrensbox/terragrunt-switcher/lib"
)

const (
	terragruntURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases"
)

var version = "0.1.0\n"

func main() {

	versionFlag := getopt.BoolLong("version", 'v', "displays the version of tgshift")
	helpFlag := getopt.BoolLong("help", 'h', "displays help message")
	_ = versionFlag

	getopt.Parse()
	args := getopt.Args()

	if *versionFlag {
		fmt.Printf("\nVersion: %v\n", version)
	} else if *helpFlag {
		usageMessage()
	} else {

		if len(args) == 1 {

			semverRegex := regexp.MustCompile(`\A\d+(\.\d+){2}\z`)
			if semverRegex.MatchString(args[0]) {
				requestedVersion := args[0]

				//check if version exist before downloading it
				tflist, _ := lib.GetTGList(terragruntURL)
				exist := lib.VersionExist(requestedVersion, tflist)

				if exist {
					lib.AddRecent(requestedVersion) //add to recent file for faster lookup
					lib.Install(requestedVersion)
				} else {
					fmt.Println("Not a valid terragrunt version")
				}

			} else {
				fmt.Println("Not a valid terragrunt version")
				fmt.Println("Args must be a valid terragrunt version")
				usageMessage()
			}

		} else if len(args) == 0 {

			tglist, _ := lib.GetTGList(terragruntURL)
			recentVersions, _ := lib.GetRecentVersions() //get recent versions from RECENT file
			tglist = append(recentVersions, tglist...)   //append recent versions to the top of the list
			tglist = lib.RemoveDuplicateVersions(tglist) //remove duplicate version

			/* prompt user to select version of terragrunt */
			prompt := promptui.Select{
				Label: "Select terragrunt version",
				Items: tglist,
			}

			_, tgversion, errPrompt := prompt.Run()

			if errPrompt != nil {
				log.Printf("Prompt failed %v\n", errPrompt)
				os.Exit(1)
			}

			fmt.Printf("Terragrunt version %q selected\n", tgversion)
			lib.Install(tgversion)
			lib.AddRecent(tgversion) //add to recent file for faster lookup (cache)
			os.Exit(0)
		} else {
			usageMessage()
		}
	}

}

func usageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terragrunt version as an argument, or choose from a menu")
}
