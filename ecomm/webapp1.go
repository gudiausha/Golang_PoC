package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/securecookie"
  "encoding/json"
  "os"
  "strconv"
)

var file_location = "D:/go1.14.3.windows-amd64/go/ecomm/data.json"
var port string

//datatype for the app
//struct because we have different data types in our app
// syntax : type struct_name struct{}
type Products struct{
  Id     int
  Product_Name string
  Price   int
  Availability string
}

type Credential struct{
  UserName string
  UserPassword string
}

// type Id_No struct{
//   Id_No int
// }


//declaring slice with struct datatype
//syntax : var slice_name [] element_type
// here 'Products' is name of the struct(element_type) declared above
var product []Products
var cred []Credential
//var id_no []Id_No

//error checking function
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
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
func load_data(){
  file, err := os.OpenFile(file_location, os.O_RDWR|os.O_APPEND, 0666)
  //specify your file path ^
  //fmt.Println(file)
  checkError(err)
  //file is read & stored in bytes form
  b, err := ioutil.ReadAll(file)
  //json bytes r unmarshaled to readable form
  //& each json is stored as a struct in slice product
  json.Unmarshal(b,&product)
  checkError(err)
  //fmt.Println(product[1].Product_Name)
  //fmt.Println(product.Product_Name,product.Availability)
}
//****************************************************************************

//**************************view list of products******************************
//Description: json file is read & the data is stored in slice(product). It is
//then encoded & sent to server. Same func for /crud page also.

func view(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Headers","content-type")
  w.Header().Set("Content-Type", "application/json")
  fmt.Println("Endpoint Hit: View list")
  load_data()
  //to return json response
  json.NewEncoder(w).Encode(product)
}

//*****************************************************************************

//***********************admin login/logout using cookie management************
//new cookie is created each time application is run. securecookie is a struct
// its inputs are haskey(compulsory),block-key(opti).
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// func getUserName(r *http.Request) (userName string) {
// 	if cookie, err := r.Cookie("session"); err == nil {
// 		cookieValue := make(map[string]string)
// 		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
// 			userName = cookieValue["name"]
// 		}
// 	}
// 	return userName
// }

func setSession(userName string, w http.ResponseWriter) {
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
    fmt.Println("Setting the cookie:",cookie)
		http.SetCookie(w, cookie)
    json.NewEncoder(w).Encode(cookie)

	}
}
//
func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
  json.NewEncoder(w).Encode(cookie)
}

// login handler

func adlogin(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  w.Header().Set("Access-Control-Allow-Headers","content-type")
  w.Header().Set("Content-Type", "application/json")
  fmt.Println("Endpoint Hit: Admin login func")
  new_credentials := new(Credential)
  _ = json.NewDecoder(r.Body).Decode(&new_credentials)
  fmt.Println("Printing the credentials: ", new_credentials)
  cred = append(cred,*new_credentials)
  Name := new_credentials.UserName
  Password := new_credentials.UserPassword

  fmt.Println("Printing the username credentials:",Name)
  fmt.Println("Printing the password credentials:",Password)
  fmt.Println("Printing the credentials: ",*new_credentials)

	if Name == "admin" && Password == "adminpassword" {
		// .. check credentials ..
    //name:= cred.Name
    fmt.Println("going inside check credentials func")
		setSession(Name, w)
	}
  // once this condition satisfies, it must go to the crud page.
	// http.Redirect(response, request, redirectTarget, 302)
}

// logout handler

func adlogout(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  w.Header().Set("Access-Control-Allow-Headers","content-type")
  w.Header().Set("Content-Type", "application/json")
  fmt.Println("Endpoint Hit: Admin logout func")
	clearSession(w)
  //this must go to the view page after logout
	//http.Redirect(response, request, "/", 302)
}

//*****************************************************************************

//**********************add a product******************************************
//Description:The json values are retrieved and decoded to a new struct. An Id
//is generated for json element. This is then added to the slice & the
//the json file is re-written with this new slice.

func create(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Headers","content-type")
  w.Header().Set("Content-Type", "application/json")
  fmt.Println("Endpoint Hit: Adding a product")
  //a pointer to the predefined struct type is created here
  new_product := new(Products)
  _ = json.NewDecoder(r.Body).Decode(&new_product)
  fmt.Println(new_product)
  load_data()
  //defer file.Close()
  //fmt.Println("Before Change:",new_product.Id)
  //ID generation:
  if len(product) == 0 {
    new_product.Id = 0
  }else{
    new_product.Id = Idgenerator()+1
  }

  //fmt.Println("After Change: ",new_product.Id)
  //fmt.Println("Printing the struct value", &new_product)
  //fmt.Println("Printinsg the A value", a)
  if new_product.Product_Name != " " {
    product = append(product,*new_product)
  }
  newdataBytes, err := json.MarshalIndent(&product, "", " ")
  // fmt.Println(newdataBytes)
  checkError(err)
  file, err := os.OpenFile(file_location, os.O_RDWR|os.O_APPEND, 0666)
  checkError(err)
  ioutil.WriteFile(file_location, newdataBytes, 0666)
  file.Close()

  // product = append(product,*new_product)
  // newdataBytes, err := json.MarshalIndent(&product, "", " ")
  // // fmt.Println(newdataBytes)
	// checkError(err)
  // file, err := os.OpenFile(file_location, os.O_RDWR|os.O_APPEND, 0666)
  // checkError(err)
	// ioutil.WriteFile(file_location, newdataBytes, 0666)
  // file.Close()
  // fmt.Println("Printing the new struct",*new_product)
  // fmt.Println("Printing the whole product", product)
}
//******************************************************************************

//*************************Update a product*************************************
//Assuming its an ecom website, the admin is allowed to update only parameters
//namely :- Price & Availability of products
// Description : the particular product details(name,price & Availability) are
// displayed in the update form where the admin can change the price & Availability
//status.Once done based on the slice index of the product, the new details are
//updated  & the json file is re-written. Then it is redirected to crud page.

func update(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Headers","content-type")
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit : Update")

  new_id := new(Products)
  _ = json.NewDecoder(r.Body).Decode(&new_id)
  // Get our params from the URL using Mux
  params := mux.Vars(r)
  // using this atoi method to parses the string into an integer
  requestId, _ := strconv.Atoi(params["id"])
  fmt.Println(*new_id)
  //requestId := new_id.Id
  // Loop through collection of prods and find one with the id from the params
  load_data()
  for index, item := range product {
      if item.Id == requestId {
          item.Price = new_id.Price
          item.Availability = new_id.Availability
          item.Product_Name = new_id.Product_Name
          item.Id = requestId
          //fmt.Println(item.Price,item.Availability,item.Product_Name,item.Id)
          product[index] = item
          newdataBytes, err := json.MarshalIndent(&product, "", " ")
          // fmt.Println(newdataBytes)
          checkError(err)
          file, err := os.OpenFile(file_location, os.O_RDWR, 0644)
          checkError(err)
          ioutil.WriteFile(file_location, newdataBytes, 0666)
          file.Close()
          }
          //Once updated it must get redirected to crud page
          //http.Redirect(w, r, "/crud", 301)
        }
      }
//******************************************************************************

//************************Delete a product**************************************
// Description: Since slice index is in sync with product id, based on the id
// the product is deleted from the slice & then the slice is re-written into
// the data.json file. Then the crud page is displayed, with new list.

func delete(w http.ResponseWriter, r*http.Request){
  w.Header().Set("Access-Control-Allow-Origin","*")
  w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  w.Header().Set("Access-Control-Allow-Headers","content-type")
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit : Delete")
  // Get our params from the URL using Mux
  params := mux.Vars(r)
  requestId, _ := strconv.Atoi(params["id"])
  // id_delete := new(Id_No)
  // _ = json.NewDecoder(r.Body).Decode(&id_delete)
  // requestId := id_delete.Id_No
  load_data()
  // Loop through collection of prods and find one with the id from the params
  for index, item := range product {
      if item.Id == requestId {
        //This delete method is costly...as the input increases the time taken
        // also increases.
        // Remove the element at index i from product.
        copy(product[index:], product[index+1:]) // Shift [i+1:] left one index.
        product = product[:len(product)-1]     // Truncate slice.
        if len(product)!=0{
		        for index1,element := range product{
			           //fmt.Println(element)
			           element.Id = index1
			           //fmt.Println(element)
			           product[index1] = element
		             }
	           }
        fmt.Println("after change",product)
        fmt.Println(product) // [{item.id,item.name,item.price,item.avail}]
        newdataBytes, err := json.MarshalIndent(&product, "", " ")
        // fmt.Println(newdataBytes)
        file, err := os.OpenFile(file_location, os.O_RDWR, 0644)
        checkError(err)
        ioutil.WriteFile(file_location, newdataBytes, 0666)
        file.Close()
      }
    }
  }
//*****************************************************************************

//a function to contain all the routes
func handleRequests() {

    myRouter := mux.NewRouter().StrictSlash(true)
    //myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/view", view)
    myRouter.HandleFunc("/adminlogin",adlogin)
    myRouter.HandleFunc("/logout",adlogout)
    myRouter.HandleFunc("/crud",view)
    myRouter.HandleFunc("/add",create)
    myRouter.HandleFunc("/update/{id}",update)
    myRouter.HandleFunc("/delete/{id}",delete)
    port = ":10002"
    log.Fatal(http.ListenAndServe(port, myRouter))

  }

func main() {
      port = ":10002"
      fmt.Println("Listening And Serving on " + "http://localhost"+port)
      handleRequests()
    }
