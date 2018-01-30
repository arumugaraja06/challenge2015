package structs

type MainData struct {
	URL string `json:"url"`
	Type	string	`json:"type"`
	Name	string	`json:"name"`
	Cast	[]SubData	`json:"cast,omitempty"`
	Movies	[]SubData	`json:"movies,omitempty"`
}

type SubData struct {
	Name	string	`json:"name"`
	Role	string	`json:"role"`
	URL		string	`json:"url"`
}

type ResultData struct {
	MovieName	string
	FirstRole	string
	FirstName	string
	SecondRole	string
	SecondName	string
}

type ProcessData struct {
	Initial string
	Target	string
	InitialPerson	MainData
	TargetPerson	MainData
	TargetMovies	map[string]SubData
	ReachedPerson	map[string]bool
	TotalPersons	[]string
	Visited			map[string]bool
	Result			map[string]ResultData
}
