package main 
  
import (
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
	// "sort"
	"strconv"
)

func doesFileExist(fName string) bool {
	if _, err := os.Stat(fName); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

func sortLists(scores []int, names []string) {
	scoresPivot := scores[len(scores) - 1]
	namesPivot := names[len(scores) - 1]

	fmt.Println("End of scores:", scoresPivot)
	fmt.Println("End of names:", namesPivot)
}

func main() {

	if(len(os.Args) < 2) {
		fmt.Printf("Usage: %s <filename>\n", os.Args[0])
		return
	} else {
		fName := os.Args[1]

		fileExist := doesFileExist(fName)

		if(fileExist) {
			file, err := os.Open(fName)
			if(err != nil) {
				log.Fatal(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			names := make([]string, 0, 10)
			scores := make([]int, 0, 10)
			for(scanner.Scan()) {
				splits := strings.Split(scanner.Text(), ": ")
				// splits[1] give the score on each line
				names = append(names, splits[0])
				intScore, err := strconv.Atoi(splits[1])
				if(err != nil) {
					fmt.Printf("Could not convert string to int: %s", err)
					return
				} else {
					scores = append(scores, intScore)
				}
			}

			for i := 0; i < len(names); i++ {
				fmt.Printf("%d. %s %d\n", i+1, names[i], scores[i])
			}

			sortLists(scores, names)
		} else {
			fmt.Println("There is no scoreboard! You're the first player!")
		}
	}

	// var firstName, lastName string
	
	// fmt.Println("Please input your first name:")
	// fmt.Scanln(&firstName)
	
	// fmt.Println("Please enter your last name:")
	// fmt.Scanln(&lastName)

	// fullName := firstName + " " + lastName
	
	// l, error := f.WriteString(fullName)
	// if(error != nil) {
	// 	fmt.Println(error)
	// 	f.Close()
	// 	return
	// }

	// fmt.Println(l, "bytes written successfully")
	// error = f.Close()
	// if(error != nil) {
	// 	fmt.Println(error)
	// 	fmt.Println("File could not close")
	// 	return
	// }

}