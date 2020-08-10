package main

import (
  "fmt"
  "log"
  "io/ioutil"
  "net/http"
  "github.com/gorilla/mux"
  "html/template"
  "encoding/json"
  "os"
  "strconv"
)

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

//declaring slice with struct datatype
//syntax : var slice_name [] element_type
// here 'Products' is name of the struct(element_type) declared above
var product []Products

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
  file, err := os.OpenFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json", os.O_RDWR|os.O_APPEND, 0666)
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

//**************************homepage*******************************************
//Description: Renders the homepage.
//Connected files: static/homepage.html ; static/stylesheets/homepage.css

func homePage(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit: homePage")
  //to render html with net/http
  http.ServeFile(w, r, "static/homepage.html")
  //fmt.Fprintf(w, "Welcome to the HomePage!")
}
//to render using html/template package:
//func homePage(w http.ResponseWriter, r *http.Request){
// templ,err1 := template.ParseFiles("static/homepage.html")
// err1 = templ.ExecuteTemplate(w,"homepage.html",nil)
//   if err1 != nil{
//     fmt.Println(err1.Error())
//}
//}
//*****************************************************************************

//**************************view list of products******************************
//Description: json file is read & the data is stored in slice(product). This
//is passed to html page & the data is rendered in table form (by looping thro'
//slice).
//Connected files : static/view.html; static/stylesheets/view.css

func view(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit: View list")
  load_data()
	t, err := template.ParseFiles("static/view.html")
	checkError(err)
	t.Execute(w, product)
}

//alternate method of rendering json/string response are given below:

//to pass json as it is to server
//func view(w http.ResponseWriter, r *http.Request) {
// w.Header().Set("Content-Type", "application/json")
// fmt.Println("Endpoint Hit: View list")
// load_data()
//to return json response
// json.NewEncoder(w).Encode(product)
// }

//to render string response to html
//convert the byte array to string
//responseString := string(content)
//fmt.Fprint(w, responseString) ----> this is using net/http
//                         or
// templ,err1 := template.ParseFiles("static/view.html") --->using html/template
// err1 = templ.ExecuteTemplate(w,"view.html",responseString)
// if err1 != nil {
//   fmt.Println(err1.Error())
// }
//*****************************************************************************

//*************************CUD Ops Page****************************************
//Description: Render the crud.html page once the admin logins with the
//predefined credentials. contains the list of products along With their id,
//name,price & Availability. Also has the links to update,delete & add functions.
//Connected Files: static/crud.html ; static/stylesheets/crud.css
func crud(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit: CUD Ops page")
  load_data()
  fmt.Println(product)
  t, err := template.ParseFiles("static/crud.html")
  checkError(err)
  t.ExecuteTemplate(w, "crud.html", product)
  //http.ServeFile(w, r, "static/crud.html")
}
//*****************************************************************************

//***********************simple admin login************************************
//here the admin credentials have been hardcoded
//this will further be modified to (jwt/session management) if needed
//Description : crud page is rendered when logged in
//with hardcoded admin credentials
//Connected files : static/admin.html
func adlogin(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit: Admin Login hit")
  //renders an empty login Page
  //analogous to render_template in python flask framework
  if r.Method != "POST"{
    t,err := template.ParseFiles("static/admin.html")
    checkError(err)
    t.Execute(w,nil)
    return
  }
  //a more compact way of writing both these functions must be done
  if r.Method == "POST"{
    username := r.FormValue("name")
    password := r.FormValue("password")
    //credentials logic subject to change
    if username != "admin" && password !="adminpassword"{
      http.ServeFile(w, r, "static/admin.html")
      fmt.Fprintf(w, "Invalid Credentials.Admin Login only.")
    }else{
      redirect_url := "/crud"
      http.Redirect(w,r,redirect_url,302)
      //analogous to redirect_url in python flask framework
    }
  }
  // http.Redirect(w,r,"/login",302)
}
//*****************************************************************************

//**********************add a product******************************************
//Description:renders product addition form, once filled the values are
//retrieved from the form and a new struct is defined. This is then appended
//to the slice & the json file is re-written with this new slice.
//Connected files: static/create.html
func create(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit: Adding a product")
  //renders an empty create form
  if r.Method != "POST"{
    t,err := template.ParseFiles("static/create.html")
    err = t.ExecuteTemplate(w,"create.html",nil)
    checkError(err)
    return
  }
  //a pointer to the predefined struct type is created here
  new_product := new(Products)

  //incase the data.json file is empty & this is the first product to be created
  //then the prod id is assumed to be 0(in sync with slice)
  //else the id-generator func is called
  if len(product) == 0 {
    new_product.Id = 0
  }else{
    new_product.Id = Idgenerator()+1
  }

  //other values(price,name & Availability) from the create form are retrieved
  getPrice,err := strconv.Atoi(r.FormValue("price"))
  checkError(err)
  new_product.Price = getPrice
  new_product.Product_Name = r.FormValue("product_name")
  new_product.Availability = r.FormValue("avail")
  // fmt.Println(new_product.Id,new_product.Price,new_product.Product_Name,new_product.Availability)

  // the below lines : append the product to the slice & rewrite the json file
  // & will get redirected to the crud ops page,where the new data will be displayed
  //open the json file
  file, err := os.OpenFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json", os.O_RDWR, 0644)
  // //fmt.Println(file)
  checkError(err)
  b, err := ioutil.ReadAll(file)
  json.Unmarshal(b,&product)
  checkError(err)
  // //stop immediate closing of file
  defer file.Close()
  product = append(product,*new_product)
  newdataBytes, err := json.MarshalIndent(&product, "", " ")
  // fmt.Println(newdataBytes)
	checkError(err)
	ioutil.WriteFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json", newdataBytes, 0666)
  file.Close()
	http.Redirect(w, r, "/crud", 301)
}
//******************************************************************************

//*************************Update a product*************************************
//Assuming its an ecom website, the admin is allowed to update only parameters
//namely :- Price & Availability of products
// Description : the particular product details(name,price & Availability) are
// displayed in the update form where the admin can change the price & Availability
//status.Once done based on the slice index of the product, the new details are
//updated  & the json file is re-written. Then it is redirected to crud page.
//Connected files : static/update.html
func update(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit : Update")
  // Get our params from the URL using Mux
  params := mux.Vars(r)
  // using this atoi method to parses the string into an integer
  requestId, _ := strconv.Atoi(params["id"])
  // Loop through collection of prods and find one with the id from the params
  for index, item := range product {
      if item.Id == requestId {
          t, err := template.ParseFiles("static/update.html")
          checkError(err)
          t.Execute(w, item)
          //json.NewEncoder(w).Encode(item)
          //return
          if r.Method == "POST"{
            fmt.Println("into post method")
            getPrice,err := strconv.Atoi(r.FormValue("price"))
            checkError(err)
            item.Price = getPrice
            item.Availability = r.FormValue("avail")
            item.Product_Name = r.FormValue("product_name")
            item.Id = requestId
            //fmt.Println(item.Price,item.Availability,item.Product_Name,item.Id)
            file, err := os.OpenFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json", os.O_RDWR, 0644)
            checkError(err)
            b, err := ioutil.ReadAll(file)
            json.Unmarshal(b,&product)
            checkError(err)
            // //stop immediate closing of file
            defer file.Close()
            product[index] = item
            newdataBytes, err := json.MarshalIndent(&product, "", " ")
            // fmt.Println(newdataBytes)
          	checkError(err)
          	ioutil.WriteFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json", newdataBytes, 0666)
            file.Close()
            // t, err := template.ParseFiles("static/crud.html")
            // checkError(err)
            // t.Execute(w,nil)
          }
          http.Redirect(w, r, "/crud", 301)
        }
      }
    }
//******************************************************************************

//************************Delete a product**************************************
// Description: Since slice index is in sync with product id, based on the id
// the product is deleted from the slice & then the slice is re-written into
// the data.json file. Then the crud page is displayed, with new list.
//Connected files : static/crud.html
func delete(w http.ResponseWriter, r*http.Request){
  w.Header().Set("Content-Type", "text/html")
  fmt.Println("Endpoint Hit : Delete")
  // Get our params from the URL using Mux
  params := mux.Vars(r)
  // using this atoi method to parses the string into an integer
  requestId, _ := strconv.Atoi(params["id"])
  // Loop through collection of prods and find one with the id from the params
  for index, item := range product {
      if item.Id == requestId {
        file, err := os.OpenFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json", os.O_RDWR, 0644)
        checkError(err)
        b, err := ioutil.ReadAll(file)
        json.Unmarshal(b,&product)
        checkError(err)
        // //stop immediate closing of file
        defer file.Close()
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
        checkError(err)
        ioutil.WriteFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json", newdataBytes, 0666)
        file.Close()
        t, err := template.ParseFiles("static/crud.html")
        checkError(err)
        t.Execute(w,nil)
      }
    }
  }
//******************************************************************************

            // fmt.Println("index: ",index)
            // fmt.Println("the item is: ",item)
            // fmt.Println("Before change:" , product[index])
            // product[index] = item
            // fmt.Println("After change: ",product[index])
            // fmt.Println("the whole product list is :" ,product)
            // fmt.Println(item)



            //product = append(product,*new_product)

    //  }

//  }
  // if r.Method == "POST"{
  //   getPrice,err := strconv.Atoi(r.FormValue("price"))
  //   checkError(err)
  //   item.Price = getPrice
  //   item.Availability = r.FormValue("avail")
  //   item.Product_Name = r.FormValue("product_name")
  //   item.Id = requestId
  //   product = append(product,item)
  //   fmt.Println(product)

  //}
  //json.NewEncoder(w).Encode(&Products{})
//}

//a function to contain all the routes
func handleRequests() {

    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/view", view)
    myRouter.HandleFunc("/login",adlogin)
    myRouter.HandleFunc("/crud",crud)
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

  // //json file is opened, read & converted to byte array & stored in content variable
  //   var content, err = ioutil.ReadFile("D:/go1.14.3.windows-amd64/go/ecomm/data.json")
  //   if err != nil {
  //       fmt.Println(err.Error())
  //   }
  //   //convert the byte array to string
  //   responseString := string(content)
  //   templ,err1 := template.ParseFiles("static/view.html")
  //   err1 = templ.ExecuteTemplate(w,"view.html",responseString)
  //   if err1 != nil {
  //     fmt.Println(err1.Error())
  //   }
  // }
//******************************************************************************

//   file, err := os.OpenFile("data.json", os.O_RDWR|os.O_APPEND, 0666)
//   if err != nil {
// 		fmt.Println(err)
// 	}
//   b, err := ioutil.ReadAll(file)
//   var prods Prod
// 	json.Unmarshal(b, &prods.product)
//   if err != nil {
// 		fmt.Println(err)
// 	}
//   au := &prods
//   fmt.Println(au)
//   t, err := template.ParseFiles("static/view.html")
//   if err != nil {
// 		fmt.Println(err)
// 	}
// 	t.ExecuteTemplate(w, au)
// }

  //fmt.Println("Endpoint Hit: Opened file")
  //product_json = json.NewEncoder(w).Encode(product)
  //load_data()
  //to return json response
  //json.NewEncoder(w).Encode(product)
  // templ,err1 := template.ParseFiles("static/view.html")
  // err1 = templ.ExecuteTemplate(w,"view.html",product)
  // if err1 != nil {
  //   fmt.Println(err1.Error())
  // }
  //}


//*****************************************************************************
//convert to json & write to data.json file
// b, err := json.Marshal(productss)
//   if err != nil {
//       http.Error(w, err.Error(), 500)
//       return
//   }
// f.Write(b)
// f.Close()
// //fmt.Fprintf(w, "Product added")
// http.Redirect(w, r, "/add", 301)
// templ,err1 := template.ParseFiles("static/create.html")
// err1 = templ.ExecuteTemplate(w,"create.html",nil)
// if err1 != nil {
//   fmt.Println(err1.Error())
// }
//}
