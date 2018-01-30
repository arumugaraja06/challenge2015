package main

import (
	"fmt"
	"./Structs"
	"./Utility"
	"os"
	"encoding/json"
	"time"
)

var processData structs.ProcessData
var num int
var start, end time.Time

func main() {
	var ResultDatas []structs.ResultData
	if len(os.Args) == 3 {
		start = time.Now()
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
			personDetails := FetchDataFromMovieBuff(person)

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
			for _, personMovies := range personDetails.Movies {
				if processData.Visited[personMovies.URL] {
					continue
				}

				processData.Visited[personMovies.URL] = true

				personMovieDetails := FetchDataFromMovieBuff(personMovies.URL)

				//Processing Casts
				for _, movieCast := range personMovieDetails.Cast {
					if processData.Visited[movieCast.URL] {
						continue
					}

					processData.Visited[movieCast.URL] = true
					processData.TotalPersons = append(processData.TotalPersons, movieCast.URL)
					processData.Result[movieCast.URL] = structs.ResultData{personMovies.Name, personMovies.Role, personDetails.Name, movieCast.Role, movieCast.Name}
				}

			}
		}
	}
	return ResultData
}

func FetchDataForPersons(person1Url, person2Url string) {
	
	person1Data := FetchDataFromMovieBuff(person1Url)

	person2Data := FetchDataFromMovieBuff(person2Url)

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

func FetchDataFromMovieBuff(partUrl string) structs.MainData {
	//fmt.Println("Going to fetch data from URL ", partUrl)
	var movieData structs.MainData
	num++
	//Fetching Data from URL
	var retData = make(chan []byte)
	go utility.HttpGet(partUrl, retData)
	byteData := <- retData
	if retData != nil {
		err := json.Unmarshal(byteData, &movieData)
		if err != nil {
			//fmt.Println("Error in unmarshaling. Error is %s", err)
		}
		return movieData
	} else {
		return movieData
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
		end = time.Now()
		fmt.Println("\n\n\nTotal Requests --> ", num)
		fmt.Println("Started At --> ", start)
		fmt.Println("Ends At --> ", end, "\n\n")
	}
}