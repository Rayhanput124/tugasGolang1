package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Customer struct (Model) ...
type Customer struct {
	CustomerID   string `json:"CustomerID"`
	CompanyName  string `json:"CompanyName"`
	ContactName  string `json:"ContactName"`
	ContactTitle string `json:"ContactTitle"`
	Address      string `json:"Address"`
	City         string `json:"City"`
	Region       string `json:"Region"`
	PostalCode   string `json:"PostalCode"`
	Country      string `json:"Country"`
	Phone        string `json:"Phone"`
	Fax          string `json:"Fax"`
}

// Get all orders

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var customers []Customer

	sql := `SELECT
				CustomerID,
				IFNULL(CompanyName,''),
				IFNULL(ContactName,'') ContactName,
				IFNULL(ContactTitle,'') ContactTitle,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNull(Region,'') Region,
				IFNULL(PostalCode,'') PostalCode
				IFNULL(Country,'') Country,
				IFNULL(Phone,'') Phone ,
				IFNULL(Fax,'') Fax
			FROM customers`

	result, err := db.Query(sql)

	defer result.Close()

	if err != nil {
		panic(err.Error())
	}

	for result.Next() {

		var customer Customer
		err := result.Scan(&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
			&customer.ContactTitle, &customer.Address, &customer.City, &customer.Country, &customer.Region,
			&customer.Phone, &customer.PostalCode)

		if err != nil {
			panic(err.Error())
		}
		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)
}

func createCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		customerID := r.FormValue("CustomerID")
		companyName := r.FormValue("CompanyName")
		contactName := r.FormValue("ContactName")
		contactTitle := r.FormValue("ContactTitle")
		address := r.FormValue("Address")
		city := r.FormValue("City")
		country := r.FormValue("Country")
		region := r.FormValue("Region")
		phone := r.FormValue("Phone")
		postalCode := r.FormValue("PostalCode")
		fax := r.FormValue("Fax")

		stmt, err := db.Prepare("INSERT INTO customers (CustomerID,CompanyName,ContactName,ContactTitle,Address,City,Country,Region,Phone,PostalCode,Fax) VALUES (?,?,?,?,?,?,?,?,?,?)")

		_, err = stmt.Exec(customerID, companyName, contactName, contactTitle, address, city, country, region, phone, postalCode, fax)

		if err != nil {
			fmt.Fprintf(w, "Data Duplicate")
		} else {
			fmt.Fprintf(w, "Data Created")
		}

	}
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customers []Customer
	params := mux.Vars(r)

	sql := `SELECT
				CustomerID,
				IFNULL(CompanyName,''),
				IFNULL(ContactName,'') ContactName,
				IFNULL(ContactTitle,'') ContactTitle,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNULL(Country,'') Country,
				IFNULL(Region,'') Region,
				IFNULL(Phone,'') Phone ,
				IFNULL(PostalCode,'') PostalCode,
				IFNULL(Fax, '') Fax
			FROM customers WHERE CustomerID = ?`

	result, err := db.Query(sql, params["id"])

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var customer Customer

	for result.Next() {

		err := result.Scan(&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
			&customer.ContactTitle, &customer.Address, &customer.City, &customer.Country, &customer.Region,
			&customer.Phone, &customer.PostalCode, &customer.Fax)

		if err != nil {
			panic(err.Error())
		}

		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)
}

func updateCustomer(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		params := mux.Vars(r)

		newCompanyName := r.FormValue("CompanyName")
		newContactName := r.FormValue("ContactName")
		newContactTitle := r.FormValue("ContactTitle")
		newAddress := r.FormValue("Address")
		newCity := r.FormValue("City")
		newRegion := r.FormValue("Region")
		newCountry := r.FormValue("Country")
		newPostalCode := r.FormValue("PostalCode")
		newPhone := r.FormValue("Phone")
		newFax := r.FormValue("Fax")

		stmt, err := db.Prepare("UPDATE customers SET CompanyName = ?, ContactName = ?, ContactTitle = ?, Address = ?, City = ?, Region = ?, Country = ?, PostalCode = ?, Phone = ?, Fax = ?  WHERE CustomerID = ?")

		_, err = stmt.Exec(newCompanyName, newContactName, newContactTitle, newAddress, newCity, newRegion, newCountry, newPostalCode, newPhone, newFax, params["id"])

		if err != nil {
			fmt.Fprintf(w, "Data not found or Request error")
		}

		fmt.Fprintf(w, "Customer with CustomerID = %s was updated", params["id"])
	}
}

func deleteCustomer(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM customers WHERE CustomerID = ?")

	_, err = stmt.Exec(params["id"])

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Customer with ID = %s was deleted", params["id"])
}

func delCustomer(w http.ResponseWriter, r *http.Request) {

	customerID := r.FormValue("CustomerID")
	companyName := r.FormValue("CompanyName")

	stmt, err := db.Prepare("DELETE FROM customers WHERE CustomerID = ? AND CompanyName = ?")

	_, err = stmt.Exec(customerID, companyName)

	if err != nil {
		fmt.Fprintf(w, "delete failed")
	}

	fmt.Fprintf(w, "Customer with ID = %s was deleted", customerID)
}

func getPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var customers []Customer

	CustomerID := r.FormValue("CustomerID")
	CompanyName := r.FormValue("CompanyName")

	sql := `SELECT
				CustomerID,
				IFNULL(CompanyName,''),
				IFNULL(ContactName,'') ContactName,
				IFNULL(ContactTitle,'') ContactTitle,
				IFNULL(Address,'') Address,
				IFNULL(City,'') City,
				IFNULL(Country,'') Country,
				IFNULL(Phone,'') Phone ,
				IFNULL(PostalCode,'') PostalCode
			FROM customers WHERE CustomerID = ? AND CompanyName = ?`

	result, err := db.Query(sql, CustomerID, CompanyName)

	if err != nil {
		panic(err.Error())
	}

	defer result.Close()

	var customer Customer

	for result.Next() {

		err := result.Scan(&customer.CustomerID, &customer.CompanyName, &customer.ContactName,
			&customer.ContactTitle, &customer.Address, &customer.City, &customer.Country,
			&customer.Phone, &customer.PostalCode)

		if err != nil {
			panic(err.Error())
		}

		customers = append(customers, customer)
	}

	json.NewEncoder(w).Encode(customers)

}

// Main function
func main() {

	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/northwind")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/customers", getCustomers).Methods("GET")
	r.HandleFunc("/customers/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/customers", createCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", deleteCustomer).Methods("DELETE")

	// Delete Customers
	r.HandleFunc("/delcustomers", delCustomer).Methods("Post")

	//New
	r.HandleFunc("/getcustomer", getPost).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":4321", r))
}
