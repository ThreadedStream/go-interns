package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"strconv"
)

func parseExcelFile(file io.Reader) []*Goods {
	wb, err := excelize.OpenReader(file)
	if err != nil {
		fmt.Errorf("ERR: %v", err)
		return nil
	}

	sheetName := "Sheet1"

	a := make([]*Goods, 0, len(wb.GetRows(sheetName)))

	for i := 1; i <= len(wb.GetRows(sheetName)); i++ {
		OfferId, _ := strconv.Atoi(wb.GetCellValue(sheetName, fmt.Sprintf("A%d", i)))
		Name := wb.GetCellValue(sheetName, fmt.Sprintf("B%d", i))
		Price, _ := strconv.ParseFloat(wb.GetCellValue(sheetName, fmt.Sprintf("C%d", i)), 32)
		Quantity, _ := strconv.Atoi(wb.GetCellValue(sheetName, fmt.Sprintf("D%d", i)))
		Available, _ := strconv.ParseBool(wb.GetCellValue(sheetName, fmt.Sprintf("E%d", i)))

		goods := Goods{
			OfferId:   OfferId,
			Name:      Name,
			Price:     Price,
			Quantity:  Quantity,
			Available: Available,
		}

		a = append(a, &goods)
	}
	return a
}
