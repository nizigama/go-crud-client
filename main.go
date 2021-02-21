package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type person struct {
	ID         int    `json:"id"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Profession string `json:"profession"`
	Country    string `json:"country"`
}

var (
	selectedMenu int
)

func main() {

	fmt.Println("Welcome to the go client")
	fmt.Println("\n===App menu===:")

	fmt.Println("1. List all persons")
	fmt.Println("2. Create a new person")
	fmt.Println("3. View a specific person")
	fmt.Println("4. Update a person")
	fmt.Println("5. Delete a person")

	fmt.Print("\n=>")

	_, err := fmt.Scan(&selectedMenu)

	for err != nil {
		fmt.Println(err)
		fmt.Println("try again, or enter 20 to exit the app")
		fmt.Print("\n=>")
		_, err = fmt.Scan(&selectedMenu)
		if selectedMenu == 20 {
			return
		}
	}

	switch selectedMenu {
	case 1:
		displayPersons()
	case 2:
		recordData()
	case 3:
		retrievePerson()
	case 4:
		updateData()
	case 5:
		deletePerson()
	case 20:
		fmt.Println("Exiting app...")
		return
	default:
		fmt.Println("Invalid selection")
	}

}

func displayPersons() {
	persons, err := loadClients()
	if err != nil {
		fmt.Println("Error occurred:", err)
		return
	}
	for _, p := range persons {
		fmt.Printf("%d. %s %s\n", p.ID, p.FirstName, p.LastName)
	}
}

func loadClients() ([]person, error) {
	fmt.Println("Loading all persons...")
	response, err := http.Get("http://127.0.0.1:8000/api/person")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var persons []person

	err = json.Unmarshal(content, &persons)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return persons, nil
}

func recordData() {
	var newPerson person
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("First name: ")
	scanner.Scan()
	newPerson.FirstName = scanner.Text()
	fmt.Print("Last name: ")
	scanner.Scan()
	newPerson.LastName = scanner.Text()
	fmt.Print("Age: ")
	scanner.Scan()
	age, err := strconv.Atoi(scanner.Text())

	if err != nil {
		fmt.Println("Error occurred while reacording new person:", err)
		return
	}

	newPerson.Age = age
	var gender rune
	fmt.Print("Gender: ")
	fmt.Scanf("%c\n", &gender)

	newPerson.Gender = string(gender)

	fmt.Print("Profession: ")
	scanner.Scan()
	newPerson.Profession = scanner.Text()

	fmt.Print("Country of origin: ")
	scanner.Scan()
	newPerson.Country = scanner.Text()

	postRequest, err := json.Marshal(newPerson)

	if err != nil {
		fmt.Println("Error converting struct to json:", err)
		return
	}

	resp, err := http.Post("http://127.0.0.1:8000/api/person", "application/json", bytes.NewBuffer(postRequest))

	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	content := make([]byte, 10)

	for {
		newContent := make([]byte, 10)
		count, err := resp.Body.Read(newContent)
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			break
		}

		if count == 0 {
			break
		}

		content = append(content, newContent...)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("%s\n", content)
	} else {
		fmt.Println("New person created successfully")
	}
}

func retrievePerson() {
	fmt.Print("Person's ID:")
	var id int

	_, err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("Invalid entry:", err)
		return
	}

	url := fmt.Sprintf("http://127.0.0.1:8000/api/person/%d", id)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("error establishing connection", err)
		return
	}
	defer resp.Body.Close()

	var response []byte

	for {
		newContent := make([]byte, 10)
		count, err := resp.Body.Read(newContent)
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			break
		}

		if count == 0 {
			break
		}

		response = append(response, newContent...)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("%s\n", response)
	} else {
		retrievedPerson := person{}

		response = bytes.Trim(response, "\x00")
		err = json.Unmarshal(response, &retrievedPerson)

		if err != nil {
			fmt.Println("error unmarshalling", err)
			return
		}

		fmt.Printf("%+v\n", retrievedPerson)
	}
}

func updateData() {
	var newPerson person

	var id int

	fmt.Print("Person's id:")
	_, err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("Invalid entry:", err)
		return
	}

	url := fmt.Sprintf("http://127.0.0.1:8000/api/person/%d", id)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("First name: ")
	scanner.Scan()
	newPerson.FirstName = scanner.Text()
	fmt.Print("Last name: ")
	scanner.Scan()
	newPerson.LastName = scanner.Text()
	fmt.Print("Age: ")
	scanner.Scan()
	age, err := strconv.Atoi(scanner.Text())

	if err != nil {
		fmt.Println("Error occurred while reacording new person:", err)
		return
	}

	newPerson.Age = age
	var gender rune
	fmt.Print("Gender: ")
	fmt.Scanf("%c\n", &gender)

	newPerson.Gender = string(gender)

	fmt.Print("Profession: ")
	scanner.Scan()
	newPerson.Profession = scanner.Text()

	fmt.Print("Country of origin: ")
	scanner.Scan()
	newPerson.Country = scanner.Text()

	postRequest, err := json.Marshal(newPerson)

	if err != nil {
		fmt.Println("Error converting struct to json:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(postRequest))

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	content := make([]byte, 10)

	for {
		newContent := make([]byte, 10)
		count, err := resp.Body.Read(newContent)
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			break
		}

		if count == 0 {
			break
		}

		content = append(content, newContent...)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("%s\n", content)
	} else {
		fmt.Printf("Person with id: %d was updated successfully\n", id)
	}
}

func deletePerson() {
	fmt.Print("Person's ID:")
	var id int

	_, err := fmt.Scan(&id)
	if err != nil {
		fmt.Println("Invalid entry:", err)
		return
	}

	url := fmt.Sprintf("http://127.0.0.1:8000/api/person/%d", id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		fmt.Println("error establishing connection", err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	var response []byte

	for {
		newContent := make([]byte, 10)
		count, err := resp.Body.Read(newContent)
		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			break
		}

		if count == 0 {
			break
		}

		response = append(response, newContent...)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("%s\n", response)
	} else {
		fmt.Printf("Person with id: %d has been deleted successfully\n", id)
	}
}
