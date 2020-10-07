//sources: https://www.goal.com/en-us/live-scores
	//https://github.com/PuerkitoBio/goquery
	//https://www.golangprograms.com/golang-html-parser.html
	//https://stackoverflow.com/questions/22891644/how-can-i-clear-the-terminal-screen-in-go
	//https://medium.com/@inhereat/terminal-color-rendering-tool-library-support-8-16-colors-256-colors-by-golang-a68fb8deee86

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
    "runtime"
	"time"
	"github.com/gookit/color"
	"github.com/PuerkitoBio/goquery"
)

var (
	resultsHomeOld []string
	resultsAwayOld []string
	resultsStatusOld []string
)

//Main Function
func main(){
	url := "https://www.goal.com/en-us/live-scores"
	var counter = 0

	//Run the scoreboard loop until the user stops or exits the program
	for {
		//Get the url
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scoreboard get url: %v\n", err)
			os.Exit(1)
		}
		//Close at the end
		defer resp.Body.Close()
		
		//Read through the returned html body
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "scoreboard goquery: %v\n", err)
			os.Exit(1)
		}

		var resultsHome []string
		var resultsAway []string
		var resultsStatus []string
		var index = 0
		//Get the home teams and corresponding goals
		doc.Find(".team-home").Each(func(i int, s *goquery.Selection) {
			title := s.Find("span").Text();
			for i, val := range title {
				if val == '>'{
					index = i
				}
			}
			homeScore := title[0:1] + title[index+2:]
			resultsHome = append(resultsHome, homeScore)
		})
		//Get the away teams and corresponding goals
		doc.Find(".team-away").Each(func(i int, s *goquery.Selection) {
			title := s.Find("span").Text();
			for i, val := range title {
				if val == '>'{
					index = i
				}
			}
			awayScore := title[0:1] + title[index+2:]
			resultsAway = append(resultsAway, awayScore)
		})
		//Get the match status
		doc.Find(".match-status").Each(func(i int, s *goquery.Selection) {
			title := s.Find("span").Text();
			status := title
			resultsStatus = append(resultsStatus, status)
		})

		//Clear the command window to remove old scores
		CallClear()
		//Display the time the program has been running for a reference
		counter++
		fmt.Fprintf(os.Stderr, "Scoreboard running for %v minutes...\n\n", counter)
		//Print the results if any were found
		if len(resultsHome) == 0 {
			fmt.Println("No games found")
		} else {
			for i, _ := range resultsHome {
				//Do not print games if they haven't happened, were cancelled, or were postponed 
				if resultsStatus[i] != "" && resultsStatus[i] != "CAN" && resultsStatus[i] != "POS"{
					//Compare the results to the old results and print changed results in red (do not consider a change in status)
					if counter != 1 && (resultsHome[i] != resultsHomeOld[i] || resultsAway[i] != resultsAwayOld[i]){
						color.Red.Printf("%s %s v %s\n", resultsStatus[i], resultsHome[i], resultsAway[i])
					} else {
						color.White.Printf("%s %s v %s\n", resultsStatus[i], resultsHome[i], resultsAway[i])
					}
				}
			}
		}
		//Update the saved results for the next comparison
		resultsStatusOld = resultsStatus
		resultsHomeOld = resultsHome
		resultsAwayOld = resultsAway
		//Wait a minute before getting the results again
		delayMin(1)
	}

}

//Delay function
func delayMin(t time.Duration) {
	time.Sleep(t * time.Minute)
}

//Create different clear commands based on os
var clear map[string]func()
func init() {
    clear = make(map[string]func())
    clear["linux"] = func() { 
        cmd := exec.Command("clear")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
    clear["windows"] = func() {
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

//Clear screen function
func CallClear() {
	//Clear the screen with commands based on os
    value, ok := clear[runtime.GOOS]
    if ok {
        value()
    } else {
        panic("Clear screen unavaliable for os")
    }
}
