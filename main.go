package main

import (
	"fmt"
	"./Structs"
	"./Utility"
	"os"
	"sync"
	"time"
	"encoding/json"
)

var processData structs.ProcessData
var wg sync.WaitGroup

func main() {
	var ResultDatas []structs.ResultData
	if len(os.Args) == 3 {
		var person1Url, person2Url string
		person1Url = os.Args[1]
		person2Url = os.Args[2]

		if os.Args[1] == os.Args[2] {
			PrintResultData(ResultDatas)
		} else {
			//Initialising Variables
			processData.TargetMovies, processData.Visited, processData.Result, processData.ReachedPerson = make(map[string]structs.SubData), make(map[string]bool), make(map[string]structs.ResultData), make(map[string]bool)

			//Fetching Data for the URL
			FetchDataForPersons(person1Url, person2Url)
		
			ResultDatas = CalculateSeparation()
			PrintResultData(ResultDatas)
		}
	} else {
		fmt.Println("Please enter neccassary input data")
	}
}

func CalculateSeparation() []structs.ResultData {
	var ResultData []structs.ResultData
	for true {
		for _, person := range processData.TotalPersons {
			var personChan = make(chan structs.MainData)
			wg.Add(1)
			go FetchDataFromMovieBuff(person, personChan, &wg)
			personDetails := <- personChan
			wg.Wait()
			go CloseChannel(personChan)

			//Checing Matches
			for _, personMovies := range personDetails.Movies {
				if processData.TargetMovies[personMovies.URL].URL == personMovies.URL { //Person Matched
					if _, exist := processData.Result[personDetails.URL]; exist {
						ResultData = append(ResultData, processData.Result[personDetails.URL], structs.ResultData{personMovies.Name, personMovies.Role, personDetails.Name, processData.TargetMovies[personMovies.URL].Role, processData.TargetPerson.Name})
					} else {
						ResultData = append(ResultData, structs.ResultData{personMovies.Name, personMovies.Role, personDetails.Name, processData.TargetMovies[personMovies.URL].Role, processData.TargetPerson.Name})
					}
					return ResultData
				}
			}

			//Processing Movies
			var personMovieChan = make(chan structs.MainData)
			for _, personMovies := range personDetails.Movies {
				if processData.Visited[personMovies.URL] {
					continue
				}

				processData.Visited[personMovies.URL] = true
				wg.Add(1)
				go FetchDataFromMovieBuff(personMovies.URL, personMovieChan, &wg)
			}
			go func () {
				for personMovieDetails := range personMovieChan {
					
					//Processing Casts
					for _, movieCast := range personMovieDetails.Cast {
						if processData.Visited[movieCast.URL] {
							continue
						}
	
						processData.Visited[movieCast.URL] = true
						processData.TotalPersons = append(processData.TotalPersons, movieCast.URL)
						processData.Result[movieCast.URL] = structs.ResultData{personMovieDetails.Name, personMovieDetails.Type, personDetails.Name, movieCast.Role, movieCast.Name}
					}
				}
			}()
			wg.Wait()
			go CloseChannel(personMovieChan)
		}
	}
	return ResultData
}

func FetchDataForPersons(person1Url, person2Url string) {
	var person1Chan = make(chan structs.MainData)
	wg.Add(1)
	go FetchDataFromMovieBuff(person1Url, person1Chan, &wg)
	person1Data := <- person1Chan
	wg.Wait()
	go CloseChannel(person1Chan)

	var person2Chan = make(chan structs.MainData)
	wg.Add(1)
	go FetchDataFromMovieBuff(person2Url, person2Chan, &wg)
	person2Data := <- person2Chan
	wg.Wait()
	go CloseChannel(person2Chan)

	if len(person1Data.Movies) > len(person2Data.Movies) {
		processData.Initial, processData.Target = person2Url, person1Url
		processData.InitialPerson, processData.TargetPerson = person2Data, person1Data
	} else {
		processData.Initial, processData.Target = person1Url, person2Url
		processData.InitialPerson, processData.TargetPerson = person1Data, person2Data
	}

	for i := range processData.TargetPerson.Movies {
		processData.TargetMovies[processData.TargetPerson.Movies[i].URL] = processData.TargetPerson.Movies[i]
	}

	processData.TotalPersons = append(processData.TotalPersons, processData.Initial)
	processData.Visited[processData.Initial] = true
}

func FetchDataFromMovieBuff(partUrl string, movieData chan structs.MainData, wg *sync.WaitGroup) {
	//fmt.Println("Going to fetch data from URL ", partUrl)
	//Fetching Data from URL
	//defer close(movieData)
	defer wg.Done()
	time.Sleep(1 * time.Millisecond)
	var retData structs.MainData
	byteData := utility.HttpGet(partUrl)
	if byteData != nil {
		err := json.Unmarshal(byteData, &retData)
		if err != nil {
			//fmt.Println("Error in unmarshaling. Error is %s", err)
		}
		movieData <- retData
	} else {
		movieData <- retData
	}
}

func PrintResultData(ResultDatas []structs.ResultData) {
	fmt.Printf("\n\nDegrees of Separation: \t%d", len(ResultDatas))

	if len(ResultDatas) > 0 {

		for i := range ResultDatas {
			var movieData = ResultDatas[i]
			fmt.Printf("\n\n%d. Movie: \t%s", i+1, movieData.MovieName)
			fmt.Printf("\n%s:\t%s", movieData.FirstRole, movieData.FirstName)
			fmt.Printf("\n%s:\t%s\n", movieData.SecondRole, movieData.SecondName)
		}
	}
}

func CloseChannel(closeChan chan structs.MainData) {
	close(closeChan)
}
