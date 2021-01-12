package main 

import (
	"fmt"
	"os"
	"bufio"
	"log"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"time"
	"encoding/json"
	"math/rand"
)

/*
	Struct to hold all the different parts of the question
*/
type jsonObj struct {
	Rescode 						int `json:"response_code"`
	Question []struct {
		Category					string `json:"category"`
		Type							string `json:"type"`
		Difficulty				string `json:"difficulty"`
		Question 					string `json:"question"`
		Correct_Answers 	string `json:"correct_answer"`
		Incorrect_Answers []string `json:"incorrect_answers"`
	} `json:"results"`
}

/* 
	Function to check if the file for the scoreboard exists in user's directory
*/
func doesFileExist(fName string) bool {
	if _, err := os.Stat(fName); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}
/*
	Adapted from https://medium.com/@rgalus/sorting-algorithms-quick-sort-implementation-in-go-9ebfd91fe95f
*/
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

/*
	Function that takes in a list and reverses it
*/
func reverseList(scores [][]string) [][]string {
	for i,j := 0, len(scores)-1; i < j; i,j = i+1, j-1 {
		scores[i], scores[j] = scores[j], scores[i]
	} 
	return scores
}

func shuffleList(answers []string) []string {
	for i := range answers {
		j := rand.Intn(i + 1)
		answers[i], answers[j] = answers[j], answers[i]
	}

	return answers
}

/*
	Function that takes the fileName in as a string and returns the list of lists that contains the names and scores from the scoreboard
*/
func loadScoreboard(fileName string) [][]string {
	fileExist := doesFileExist(fileName)
	scores := make([][]string, 0, 10)
	
	if(fileExist) {
		file, err := os.Open(fileName)
		if(err != nil) {
			log.Fatal(err)
		}
		defer file.Close()
		
		scanner := bufio.NewScanner(file)
		// Read the file into the array
		for(scanner.Scan()) {
			line := strings.Split(scanner.Text(), ": ")
			scores = append(scores, line)
		}
		
		file.Close()
	}

	return scores
}

/*
	Function that prints the scoreboard
*/
func printScoreboard(scores [][]string) {
	sortList(scores, 0, len(scores)-1)

	fmt.Printf("SCOREBOARD\n--------------------------------------------------\n")
	scores = reverseList(scores)
	for i := 0; i < len(scores); i++ {
		fmt.Printf("%d. %s ----> %s\n", i+1, scores[i][0], scores[i][1])
	}
	fmt.Printf("--------------------------------------------------\n")
}

/*
	Function that loads the data from the api call into the struct
*/
func loadData() jsonObj {
	/*
		API: https://opentdb.com/api_config.php
	*/
	url := "https://opentdb.com/api.php?amount=10&type=multiple"
	client := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if(getErr != nil) {
		log.Fatal(getErr)
	}
	if(res.Body != nil) {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if(readErr != nil) {
		log.Fatal(readErr)
	}

	byteBody := []byte(body)

	var questions jsonObj
	err = json.Unmarshal(byteBody, &questions)

	if err == nil {
		fmt.Printf("Data loaded successfully!")
	} else {
		fmt.Printf("%s\n", err)
	}

	return questions
}

func giveTest(questions jsonObj) int {
	score := 0
	var choice string 
	allAnswers := make([]string, 4)
	
	for i := 0; i < len(questions.Question); i++ {
		allAnswers = nil
		choice = ""
		for j := 0; j < 3; j++ {
			allAnswers = append(allAnswers, questions.Question[i].Incorrect_Answers[j])
		}
		allAnswers = append(allAnswers, questions.Question[i].Correct_Answers)
		shuffledAnswers := shuffleList(allAnswers)

		/*
			Print the questions and the answers
		*/
		fmt.Printf("---------------------------------------\n")
		fmt.Printf("(%d/10) %s\n", i+1, questions.Question[i].Question)
		for j := 0; j < len(allAnswers); j++ {
			fmt.Printf("%d. %s\n", j + 1, shuffledAnswers[j])
		}
		fmt.Scanln(&choice)
		str_choice, err := strconv.Atoi(choice)
		if err != nil {
			fmt.Println("Can't turn string into int...")
			return -1 
		}
		if(shuffledAnswers[str_choice - 1] == questions.Question[i].Correct_Answers) {
			score++
			fmt.Printf("Correct! ---> Socre: %d\n", score)
		} else {
			fmt.Printf("Incorrect. The correct answer was:\n%s\nScore: %d\n", questions.Question[i].Correct_Answers, score)
		}
	}
	return score
}

func startQuiz(fileName string, questions jsonObj) {
	var choice string
	for ;; {
		fmt.Printf("Please choose an option:\n1. Start quiz\n2. See scoreboard\n3. Add a question\n4. Quit\n")
		fmt.Scanln(&choice)
		if(choice == "1") {
			score := giveTest(questions)
			switch score {
				case 0, 1, 2, 3, 4:
					fmt.Printf("Oof... a score of %d/10 is not that great. Try harder next time\n", score)
				case 5, 6, 7, 8:
					fmt.Printf("Not bad... not bad at all. Score: %d/10\n", score)
				case 9:
					fmt.Printf("Close but no cigar... Score: %d/10\n", score)
				case 10:
					fmt.Printf("Damn, you did it! %d/10!\n", score)
			}
			fmt.Printf("---------------------------------------\n")
		} else if(choice == "2") {
			scores := loadScoreboard(fileName)
			printScoreboard(scores)
		} else if(choice == "4") {
			fmt.Printf("Thanks for playing!\n")
			break
		}
	}
}

func main() {
	if(len(os.Args) < 2) {
		fmt.Printf("Usage: %s <filename>\n", os.Args[0])
		return
	} else {
		fName := os.Args[1]
		var questions jsonObj
		questions = loadData()
		startQuiz(fName, questions)
	}
}