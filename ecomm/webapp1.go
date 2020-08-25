package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var file_location = "D:/go1.14.3.windows-amd64/go/ecomm/data.json"
var log_location = "D:/go1.14.3.windows-amd64/go/ecomm/logs.txt"
var port string

//datatype for the app
//struct because we have different data types in our app
// syntax : type struct_name struct{}
type Products struct {
	Id           int
	Product_Name string
	Price        int
	Availability string
}

type Credential struct {
	UserName     string
	UserPassword string
}

//declaring slice with struct datatype
//syntax : var slice_name [] element_type
// here 'Products' is name of the struct(element_type) declared above
var product []Products
var cred []Credential

//log warnings
var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

//log function
func init() {
	file, err := os.OpenFile(log_location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

//error checking function
func checkError(err error) {
	if err != nil {
		//fmt.Println(err)
		ErrorLogger.Println(err)
	}
}

//id generation function
//the id generated is in sync with the index of the slice 'product'
func Idgenerator() int {
	maxId := product[0].Id
	for _, v := range product {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId
}

//**************************load json*******************************************
// Description: load the json data,unmarshal it & save to a slice (product)
func load_data() {
	//file==has memory location of file
	file, err := os.OpenFile(file_location, os.O_RDWR|os.O_APPEND, 0666)
	checkError(err)
	//file is read & stored in bytes form
	b, err := ioutil.ReadAll(file)
	//json bytes is parsed & stored as struct in slice-product
	json.Unmarshal(b, &product)
	checkError(err)
}

//****************************************************************************

//**************************view list of products******************************
//Description: json file is read & the data is stored in slice(product). It is
//then encoded & sent to server. Same func for /crud page also.

func view(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Endpoint Hit: View")
	load_data()
	//to return json response
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusOK) //StatusOK
		json.NewEncoder(w).Encode(product)
	} else {
		w.WriteHeader(http.StatusNotFound) //Status NotFound
	}
}

//*****************************************************************************

//***********************admin login/logout using cookie management************
//new cookie is created each time application is run. securecookie is a struct
// its inputs are haskey(compulsory),block-key(opti).
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func setSession(userName string, w http.ResponseWriter) {
	// map[key][value]
	value := map[string]string{
		"name": userName,
	}
	//:= 64bit(32bit("cookiename",value to be encoded))
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		InfoLogger.Println("Setting the cookie:", cookie)
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(cookie)

	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	InfoLogger.Println("Cookie cleared:", cookie)
	http.SetCookie(w, cookie)
	json.NewEncoder(w).Encode(cookie)
}

// login handler
func adlogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Endpoint Hit: Admin login")
	new_credentials := new(Credential)
	_ = json.NewDecoder(r.Body).Decode(&new_credentials)
	cred = append(cred, *new_credentials)
	Name := new_credentials.UserName
	Password := new_credentials.UserPassword	
	if Name == "admin" && Password == "adminpassword" {
		// .. check credentials ..
		InfoLogger.Println("Credentials checked")
		setSession(Name, w)
	} else {
		value := map[int]string{
			http.StatusUnauthorized: "Unauthorized user",
		}
		WarningLogger.Println(http.StatusUnauthorized, ":Unauthorized user")
		json.NewEncoder(w).Encode(value)
	}
}

// logout handler

func adlogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	InfoLogger.Println("Endpoint Hit: Admin logout")
	clearSession(w)
}

//*****************************************************************************

//**********************add a product******************************************
//Description:The json values are retrieved and decoded to a new struct. An Id
//is generated for json element. This is then added to the slice & the
//the json file is re-written with this new slice.

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Endpoint Hit: Add")
	//a pointer to the predefined struct type is created here
	new_product := new(Products)
	_ = json.NewDecoder(r.Body).Decode(&new_product)

	load_data()

	//ID generation:
	if len(product) == 0 {
		new_product.Id = 0
	} else {
		new_product.Id = Idgenerator() + 1
	}

	//condition for product name
	if new_product.Product_Name != " " {
		value := map[int]string{
			http.StatusOK: "Product added successfully",
		}
		product = append(product, *new_product)
		InfoLogger.Println("Product added successfully:", *new_product)
		json.NewEncoder(w).Encode(value)
	} else {
		value := map[int]string{
			http.StatusMethodNotAllowed: "Product Name not entered",
		}
		ErrorLogger.Println(http.StatusMethodNotAllowed, ":Product Name empty")
		json.NewEncoder(w).Encode(value)
	}
	newdataBytes, err := json.MarshalIndent(&product, "", " ")
	checkError(err)
	file, err := os.OpenFile(file_location, os.O_RDWR|os.O_APPEND, 0666)
	checkError(err)
	ioutil.WriteFile(file_location, newdataBytes, 0666)
	file.Close()
}

//******************************************************************************

//*************************Update a product*************************************
//Assuming its an ecom website, the admin is allowed to update only parameters
//namely :- Price & Availability of products
// Description : the particular product details(name,price & Availability) are
// displayed in the update form where the admin can change the price & Availability
//status.Once done based on the slice index of the product, the new details are
//updated  & the json file is re-written. Then it is redirected to crud page.

func update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Endpoint Hit : Update")
	var flag bool //default false

	new_id := new(Products)
	_ = json.NewDecoder(r.Body).Decode(&new_id)
	// Get our params from the URL using Mux- returns map of the req
	//id:"value in url"
	params := mux.Vars(r)
	// using this atoi method to parses the string into an integer
	requestId, _ := strconv.Atoi(params["id"])
	// Loop through collection of prods and find one with the id from the params
	load_data()
	for index, item := range product {
		if item.Id == requestId {
			item.Price = new_id.Price
			item.Availability = new_id.Availability
			item.Product_Name = new_id.Product_Name
			item.Id = requestId
			product[index] = item
			InfoLogger.Println("Updated:", item)
			newdataBytes, err := json.MarshalIndent(&product, "", " ")
			checkError(err)
			file, err := os.OpenFile(file_location, os.O_RDWR, 0644)
			checkError(err)
			ioutil.WriteFile(file_location, newdataBytes, 0666)
			file.Close()
			flag = true
			value := map[int]string{
				http.StatusOK: "Product updated successfully",
			}
			json.NewEncoder(w).Encode(value)
		}
	}
	if !flag {
		value := map[int]string{
			http.StatusNotFound: "Request ID not found",
		}
		ErrorLogger.Println(http.StatusNotFound, ":Request ID not found")
		json.NewEncoder(w).Encode(value)
	}
}

//******************************************************************************

//************************Delete a product**************************************
// Description: Since slice index is in sync with product id, based on the id
// the product is deleted from the slice & then the slice is re-written into
// the data.json file. Then the crud page is displayed, with new list.

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Content-Type", "application/json")
	InfoLogger.Println("Endpoint Hit : Delete")
	var delflag bool
	// Get our params from the URL using Mux
	params := mux.Vars(r)
	requestId, _ := strconv.Atoi(params["id"])
	load_data()
	// Loop through collection of prods and find one with the id from the params
	for index, item := range product {
		if item.Id == requestId {
			delflag = true
			//This delete method is costly...as the input increases the time taken
			// also increases.
			// Remove the element at index i from product.
			copy(product[index:], product[index+1:]) // Shift [i+1:] left one index.
			product = product[:len(product)-1]       // Truncate slice.
			if len(product) != 0 {
				for index1, element := range product {
					element.Id = index1
					product[index1] = element
				}
			}
			InfoLogger.Println("Deleted:", product) // [{item.id,item.name,item.price,item.avail}]
			newdataBytes, err := json.MarshalIndent(&product, "", " ")
			file, err := os.OpenFile(file_location, os.O_RDWR, 0644)
			checkError(err)
			ioutil.WriteFile(file_location, newdataBytes, 0666)
			file.Close()
			value := map[int]string{
				http.StatusOK: "Product deleted successfully",
			}
			json.NewEncoder(w).Encode(value)
		}
	}
	if !delflag {
		value := map[int]string{
			http.StatusNotFound: "Request ID not found",
		}
		ErrorLogger.Println(http.StatusNotFound, ":Request ID not found")
		json.NewEncoder(w).Encode(value)
	}
}

//*****************************************************************************

//a function to contain all the routes
func handleRequests() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/view", view)
	myRouter.HandleFunc("/adminlogin", adlogin)
	myRouter.HandleFunc("/logout", adlogout)
	myRouter.HandleFunc("/crud", view)
	myRouter.HandleFunc("/add", create)
	myRouter.HandleFunc("/update/{id}", update)
	myRouter.HandleFunc("/delete/{id}", delete)
	port = ":10002"
	log.Fatal(http.ListenAndServe(port, myRouter))

}

func main() {
	port = ":10002"
	fmt.Println("Listening And Serving on " + "http://localhost" + port)
	InfoLogger.Println("Listening And Serving on " + "http://localhost" + port)
	handleRequests()
}
