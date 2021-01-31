package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

func (a *App) checkSellerOfferExistence(sellerId int, offerId int) bool {
	query := fmt.Sprintf("SELECT COUNT(*) AS count FROM goods WHERE seller_id=%d AND offer_id=%d;", sellerId, offerId)
	rows, err := a.Conn.Query(query)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return getRowCount(rows) != 0
}

func getRowCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}
	return count
}

func (a *App) insertGoods(g Goods, sellerId int) error {
	query := fmt.Sprintf("INSERT INTO goods(offer_id, name, price, quantity, available, seller_id) VALUES (%d,'%s',%f,%d,%t,%d)", g.OfferId, g.Name, g.Price, g.Quantity, g.Available, sellerId)

	_, err := a.Conn.Query(query)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) getGoods(sellerId, offerId, query string) (error, []map[string]interface{}) {
	mapping := make([]map[string]interface{}, 0)
	var sellerIdInt int
	var offerIdInt int
	if sellerId != "" {
		var err error
		sellerIdInt, err = strconv.Atoi(sellerId)
		if err != nil {
			return err, nil
		}
	} else {
		sellerIdInt = -2
	}

	if offerId != "" {
		var err error
		offerIdInt, err = strconv.Atoi(offerId)
		if err != nil {
			return err, nil
		}
	} else {
		offerIdInt = -2
	}
	databaseQuery := buildQuery(offerIdInt, sellerIdInt, query)
	rows, err := a.Conn.Query(databaseQuery)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	mapping = mapify(rows)

	return nil, mapping
}

func buildQuery(offerIdInt int, sellerIdInt int, query string) string {
	var databaseQuery string
	if query != "" && offerIdInt != -2 && sellerIdInt != -2 {
		databaseQuery = fmt.Sprintf("SELECT * FROM goods WHERE seller_id=%d AND offer_id=%d AND name='%s'", sellerIdInt, offerIdInt, query)
	} else if query != "" && (offerIdInt != -2 || sellerIdInt != -2) {
		databaseQuery = fmt.Sprintf("SELECT offer_id, name, price, quantity, available, seller_id, REGEXP_MATCHES(name, '%s') FROM goods WHERE offer_id=%d OR seller_id=%d", query, offerIdInt, sellerIdInt)
	} else {
		databaseQuery = fmt.Sprintf("SELECT offer_id, name, price, quantity, available, seller_id, REGEXP_MATCHES(name, '%s') AS title FROM goods", query)
	}

	return databaseQuery
}

func mapify(rows *sql.Rows) []map[string]interface{} {
	mapping := make([]map[string]interface{}, 0)
	var offerId int
	var name string
	var price float64
	var quantity int
	var available bool
	var sellerId int
	var title string
	for rows.Next() {
		err := rows.Scan(&offerId, &name, &price, &quantity, &available, &sellerId, &title)
		if err != nil {
			fmt.Println(err)
		}
		elem := map[string]interface{}{"OfferId": offerId, "Name": name, "Price": price, "Quantity": quantity, "Available": available, "SellerId": sellerId}
		mapping = append(mapping, elem)
	}
	return mapping
}

func (a *App) updateGoods(g Goods) (bool, error) {
	if !g.Available {
		_, err := a.Conn.Query("DELETE FROM goods WHERE offer_id=$1 AND seller_id=$2", g.OfferId, g.SellerId)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	_, err := a.Conn.Query("UPDATE goods SET name = $1, price = $2, quantity = $3", g.Name, g.Price, g.Quantity)
	if err != nil {
		return false, nil
	}

	return false, nil
}
