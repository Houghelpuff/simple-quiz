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

func sortList(scores [][]string, start, end int) {
	if (end - start) < 1 {
		return
	}

	pivot := scores[end]
	splitIdx := start

	for i := start; i < end; i++ {
		// fmt.Println(strconv.Atoi(pivot[1]))

		intPivot, err := strconv.Atoi(pivot[1])
		intScore, err := strconv.Atoi(scores[i][1])

		if err != nil {
			fmt.Println("Can't turn string into int...")
			return
		}

		if intScore < intPivot {
			temp := scores[splitIdx]

			scores[splitIdx] = scores[i]
			scores[i] = temp

			splitIdx++
		}
	}

	scores[end] = scores[splitIdx]
	scores[splitIdx] = pivot

	sortList(scores, start, splitIdx-1)
	sortList(scores, splitIdx+1, end)
}

func reverseList(scores [][]string) [][]string {
	for i,j := 0, len(scores)-1; i < j; i,j = i+1, j-1 {
		scores[i], scores[j] = scores[j], scores[i]
	} 
	return scores
}

func printScoreboard(scores [][]string) {
	fmt.Println("SCOREBOARD\n--------------------------------------------------")
	scores = reverseList(scores)
		for i := 0; i < len(scores); i++ {
			fmt.Printf("%d. %s ----> %s\n", i+1, scores[i][0], scores[i][1])
		}
	fmt.Println("--------------------------------------------------")
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
			scores := make([][]string, 0, 10)

			// Read the file into the array
			for(scanner.Scan()) {
				line := strings.Split(scanner.Text(), ": ")
				scores = append(scores, line)
			}

			file.Close()

			sortList(scores, 0, len(scores)-1)
			printScoreboard(scores)


		} else {
			fmt.Println("Uh oh, there is no scoreboard! You're the first player!")
		}
	}
}