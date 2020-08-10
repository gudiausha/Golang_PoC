# Golang_PoC - An E-comm Website
### Description
E-comm Website is a basic web application containing all the CRUD actions. The users can view the products present &amp; an authorized user can create,update and delete a product.
#### Technologies :
* Backend : Golang
* Frontend : Vue-js

### Pre-requisties 
* Go initial setup
* Mux package in Go. To download type `go get github.com/gorilla/mux` in cmd

### Folder Structure
The folder structure along with their descriptions is given below.

    ecomm/
     ├── webapp1.go     #main program - contains views for CRUD ops
     ├── data.json      # json file which is used as database 
     └── static/          
           ├── admin.html      
           ├── create.html
           ├── crud.html        
           ├── homepage.html
           ├── update.html
           ├── view.html
           ├── stylesheets/     
                 │       ├── homepage.css     
                 │       ├── view.css
                 │       ├── crud.css     
                 
### Deployment
1. Download this repo in the go root folder
2. Download the required packages
3. Run this in cmd `go run webapp1.go`
4. Go to a brower and type `http//localhost:10002`, the app runs in the local server

### Snapshots 
Homepage
![Homepage](images/homepage.JPG)

View list of Products
![View](images/view.JPG)

Admin Login
<img src="images/admin login.JPG" />

Addition of Product
<img src="images/add form.JPG" />

JSON file before addition 
    <img src="images/json file before add.JPG" />

JSON file after addition
    <img src="images/json file after add.JPG" />
    
UI after addition
<img src="images/form after add.JPG" />

Updation of Product
<img src="images/update prod1 form.JPG" />

UI before updation
<img src="images/form before update.JPG" />

UI after updation
<img src="images/after updation.JPG" />

JSON file after updation
<img src="images/file after updation.JPG" />
    
   





                
          
