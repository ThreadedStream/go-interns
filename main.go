package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

//Used to associate task id with corresponding status of operation being executed by a server
// NOTE: Beware that there exist better alternatives to handle such type of problems
// I guess that using redis would do the trick
var taskIdsToStatus = make(map[int]interface{}, 0)

type App struct {
	Router *mux.Router
	Conn   *sql.DB
	Server *http.Server
}

type Goods struct {
	OfferId   int     `json:"offer_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Available bool    `json:"available"`
	SellerId  int     `json:"seller_id"`
}

type Stats struct {
	Deleted int `json:"deleted"`
	Created int `json:"created"`
	Updated int `json:"updated"`
	Errors  int `json:"errors"`
}

func (g *Goods) ToString() string {
	return fmt.Sprintf("Name: %d\n"+"OfferId: %s\n"+
		"Price: %f\n"+
		"Quantity: %d\n"+
		"Available: %t\n", g.OfferId, g.Name, g.Price, g.Quantity, g.Available)
}

func (g *Goods) getGoods(db *sql.DB) error {
	return errors.New("not implemented")
}

func (a *App) initRoutes() {
	a.Router.Path("/load_goods").HandlerFunc(a.loadGoods).Methods("POST")
	a.Router.Path("/retrieve_goods").HandlerFunc(a.retrieveGoods).Methods("GET")
	a.Router.Path("/get_status").HandlerFunc(a.getStatusOfOperation).Methods("GET")
}

func (a *App) initialize(user, password, dbname, addr string) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "192.168.1.40", 5432, user, password, dbname)

	var err error
	a.Conn, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
		return
	}
	a.Router = mux.NewRouter()

	a.Server = &http.Server{
		Addr:         addr,
		Handler:      a.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//Initialize router and routes
	a.initRoutes()
}

func (a *App) Run() {
	log.Printf("Starting server on %s\n", a.Server.Addr)
	log.Fatal(a.Server.ListenAndServe())
}

func main() {
	user := "postgres"
	password := "135797531"
	dbname := "intern"
	addr := os.Getenv("ADDR")
	a := App{}

	a.initialize(user, password, dbname, addr)
	a.initRoutes()
	//a.checkSellerOfferExistence(1,1)
	log.Println("hello")
	a.Run()
}
