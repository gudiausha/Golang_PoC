# Golang_PoC - An E-comm Website
### Description
E-comm Website is a basic web application containing all the CRUD actions. The users can view the products present &amp; an authorized user can create,update and delete a product.
#### Technologies :
* Backend : Golang
* Frontend : Vue-js

### Pre-requisties 
* Go initial setup
* Mux package in Go. To download type `go get github.com/gorilla/mux` in cmd
* Secure cookie package in Go. To download type `go get github.com/securecookie` in cmd

### Folder Structure
The folder structure along with their descriptions is given below.

    ecomm/
     ├── webapp1.go     #main program - contains views for CRUD ops
     ├── data.json      # json file which is used as database 
     └── log.txt        # logs all the details    
                 
### Deployment
1. Download this repo in the go root folder
2. Download the required packages
3. Run this in cmd `go run webapp1.go`
4. Go to a brower and type `http//localhost:10002`, the app runs in the local server

## API Routes

Path | Method | Required JSON | Header | Description
---|---|---|---|---
/view | GET | -- | Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type:"application/json" | Displays the list of products
/add | POST | Product Info |Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type:"application/json" | Creates a new products by adding the info to JSON file
/update/{id} | POST |Updated Product Info | Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type:"application/json" | Update the product with new info
/delete/{id} | DELETE |--| Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type:"application/json" | Deletes the product at the particular index
/adminlogin | POST |Admin Credentials| Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type:"application/json" | Login into the website
/logout | POST |--| Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type:"application/json" | Logout from the website

### Snapshots 
To view the images of the webapp click here:
[Snapshots](SNAPSHOTS.md)

### Tasks
* Completed:
    - [x] Reading data from database
    - [x] Appending data to database
    - [x] Updating data in database
    - [x] Deleting data from database
    - [x] Simple admin login (using if condition)
    - [x] Rendering the homepage/view/admin/crud pages to html (this was done using basic html & css)
    
* To Do:
    - [x] Authorization using either JWT/Session Cookies in golang
         - [x] Writing login & logout views
    - [x] Integrating with Vue-js
         - [x] Fix cross site errors
         - [x] Display the product list in tabular form
         - [x] After updation redirect to CRUD page
         - [x] After deletion redirect to CRUD page
         - [x] Design login & logout pages (if needed can revamp the whole site)
         
   





                
          
