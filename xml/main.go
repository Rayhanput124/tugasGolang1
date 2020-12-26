package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

// Root struct
type Root struct {
	XMLName   xml.Name `xml:"Root"`
	Text      string   `xml:",chardata"`
	Xmlns     string   `xml:"xmlns,attr"`
	Customers struct {
		Text     string `xml:",chardata"`
		Customer []struct {
			Text         string `xml:",chardata"`
			CustomerID   string `xml:"CustomerID,attr"`
			CompanyName  string `xml:"CompanyName"`
			ContactName  string `xml:"ContactName"`
			ContactTitle string `xml:"ContactTitle"`
			Phone        string `xml:"Phone"`
			FullAddress  struct {
				Text       string `xml:",chardata"`
				Address    string `xml:"Address"`
				City       string `xml:"City"`
				Region     string `xml:"Region"`
				PostalCode string `xml:"PostalCode"`
				Country    string `xml:"Country"`
			} `xml:"FullAddress"`
			Fax string `xml:"Fax"`
		} `xml:"Customer"`
	} `xml:"Customers"`
	Orders struct {
		Text  string `xml:",chardata"`
		Order []struct {
			Text         string `xml:",chardata"`
			CustomerID   string `xml:"CustomerID"`
			EmployeeID   string `xml:"EmployeeID"`
			OrderDate    string `xml:"OrderDate"`
			RequiredDate string `xml:"RequiredDate"`
			ShipInfo     struct {
				Text           string `xml:",chardata"`
				ShippedDate    string `xml:"ShippedDate,attr"`
				ShipVia        string `xml:"ShipVia"`
				Freight        string `xml:"Freight"`
				ShipName       string `xml:"ShipName"`
				ShipAddress    string `xml:"ShipAddress"`
				ShipCity       string `xml:"ShipCity"`
				ShipRegion     string `xml:"ShipRegion"`
				ShipPostalCode string `xml:"ShipPostalCode"`
				ShipCountry    string `xml:"ShipCountry"`
			} `xml:"ShipInfo"`
		} `xml:"Order"`
	} `xml:"Orders"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Root

	if r.Method == "POST" {
		if err = xml.Unmarshal(body, &request); err != nil {
			fmt.Fprintln(w, "Failed decoding json message")
		} else {
			sql := `INSERT INTO customers (CustomerID, CompanyName) VALUES (?, ?)`
			stmt, err := db.Prepare(sql)

			for _, each := range request.Customers.Customer {
				_, err = stmt.Exec(each.CustomerID, each.CompanyName)

				if err != nil {
					fmt.Fprintln(w, "Data duplicate")
				} else {
					fmt.Fprintln(w, "Data inserted", each.CustomerID)
				}
			}
		}
	}
}

func getOrders(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	var request Root

	if r.Method == "POST" {
		if err = xml.Unmarshal(body, &request); err != nil {
			fmt.Fprintln(w, "Failed decoding json message")
		} else {
			sql := "INSERT INTO orders (CustomerID, EmployeeID) VALUES(?, ?)"
			stmt, err := db.Prepare(sql)

			for _, each := range request.Orders.Order {
				_, err = stmt.Exec(each.CustomerID, each.EmployeeID)

				if err != nil {
					fmt.Fprintln(w, "Data duplicate")
				} else {
					fmt.Fprintln(w, "Data inserted", each.CustomerID)
				}
			}
		}
	}
}

func main() {

	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306)/northwind")
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	// Init router
	r := mux.NewRouter()

	fmt.Println("Server on :8181")

	// Route handles & endpoints
	r.HandleFunc("/customers", getCustomers).Methods("POST")
	r.HandleFunc("/orders", getOrders).Methods("POST")

	// Start server
	log.Fatal(http.ListenAndServe(":8181", r))

}
