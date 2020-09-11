package main

import (
    "fmt"
    "os"
    "strconv"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("    go " + os.Args[0] + "my5g-rantester CMD [OPTIONS]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("    AttachGnb")
	fmt.Println("    MultiAttachUesInQueue NUM_OF_UES")
	fmt.Println("    MultiAttachUesInConcurrencyWithGNBs")
	fmt.Println("    MultiAttachUesInConcurrencyWithTNLAs NUM_OF_UES")
	fmt.Println("    MultiAttachGnbInQueue NUM_OF_GNBS")
	fmt.Println("    MultiAttachGnbInConcurrency NUM_OF_GNBS")
}

func main() {
	if len(os.Args) > 1 {
		switch cmd := os.Args[1]; cmd {
			case "AttachGnb":
				fmt.Println(testAttachGnb())
			case "MultiAttachUesInQueue":
				if len(os.Args) > 2 {
					numOfUes, err := strconv.Atoi(os.Args[2])
					if err == nil {
						fmt.Println(testMultiAttachUesInQueue(numOfUes))
					}
				} else {
					usage()
				}
			case "MultiAttachUesInConcurrencyWithGNBs":
				fmt.Println(testMultiAttachUesInConcurrencyWithGNBs())
			case "MultiAttachUesInConcurrencyWithTNLAs":
				if len(os.Args) > 2 {
					numOfUes, err := strconv.Atoi(os.Args[2])
					if err == nil {
						fmt.Println(testMultiAttachUesInConcurrencyWithTNLAs(numOfUes))
					}
				} else {
					usage()
				}
			case "MultiAttachGnbInQueue":
				if len(os.Args) > 2 {
					numOfGnbs, err := strconv.Atoi(os.Args[2])
					if err == nil {
						fmt.Println(testMultiAttachGnbInQueue(numOfGnbs))
					}
				} else {
					usage()
				}
			case "MultiAttachGnbInConcurrency":
				if len(os.Args) > 2 {
					numOfGnbs, err := strconv.Atoi(os.Args[2])
					if err == nil {
						fmt.Println(testMultiAttachGnbInConcurrency(numOfGnbs))
					}
				} else {
					usage()
				}
			default:
				usage()

		}
	}
}
