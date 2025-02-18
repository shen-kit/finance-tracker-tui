package backend

import (
	"database/sql"
	"fmt"
	"time"
)

type DataRow interface {
	// spread to a slice of strings, used to display as a table row
	SpreadToStrings() []string
}

type Record struct {
	Id    int
	Date  time.Time
	CatId int
	Desc  string
	Amt   int
}

func (rec Record) Spread() (int, time.Time, string, int, int) {
	return rec.Id, rec.Date, rec.Desc, rec.Amt, rec.CatId
}

func (rec Record) SpreadToStrings() []string {
	id, date, desc, amt, catId := rec.Spread()
	return []string{
		fmt.Sprint(id),
		date.Format("2006-01-02"),
		GetCategoryNameFromId(catId),
		desc,
		fmt.Sprintf("$%.2f", float32(amt)/100),
	}
}

type Category struct {
	Id       int
	Name     string
	IsIncome bool
	Desc     string
}

func (cat Category) Spread() (int, string, bool, string) {
	return cat.Id, cat.Name, cat.IsIncome, cat.Desc
}

func (c Category) SpreadToStrings() []string {
	if c.IsIncome {
		return []string{fmt.Sprint(c.Id), c.Name, "Income", c.Desc}
	} else {
		return []string{fmt.Sprint(c.Id), c.Name, "Expenditure", c.Desc}
	}
}

type Investment struct {
	Id        int
	Date      time.Time
	Code      string
	Unitprice int
	Qty       float32
}

func (inv Investment) Spread() (int, time.Time, string, int, float32) {
	return inv.Id, inv.Date, inv.Code, inv.Unitprice, inv.Qty
}

func (inv Investment) SpreadToStrings() []string {
	return []string{
		fmt.Sprint(inv.Id),
		inv.Date.Format("2006-01-02"),
		inv.Code,
		fmt.Sprintf("$%.2f", float32(inv.Unitprice)/100),
		fmt.Sprintf("%.1f", inv.Qty),
		fmt.Sprintf("$%.2f", float32(inv.Unitprice)/100*inv.Qty),
	}
}

func dbRowsToInvestments(rows *sql.Rows) []DataRow {
	var investments []DataRow

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var inv Investment
		if err := rows.Scan(&inv.Id, &inv.Date, &inv.Code, &inv.Unitprice, &inv.Qty); err != nil {
			panic(err)
		}
		investments = append(investments, inv)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return investments
}

func dbRowsToRecords(rows *sql.Rows) []DataRow {
	var records []DataRow

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.Id, &rec.Date, &rec.Desc, &rec.Amt, &rec.CatId); err != nil {
			panic(err)
		}
		records = append(records, rec)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return records
}

func dbRowsToCategories(rows *sql.Rows) []DataRow {
	var categories []DataRow

	// for each row, assign column data to struct fields and append struct to slice
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.Desc, &cat.IsIncome); err != nil {
			panic(err)
		}
		categories = append(categories, cat)
	}

	// check for errors then return
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return categories
}

type CategoryYear struct {
	CatId     int
	MonthSums [12]int // sum of records for this category for each month
}

func (cy CategoryYear) SpreadToStrings() []string {
	var res = make([]string, 13, 13)
	res[0] = GetCategoryNameFromId(cy.CatId)
	for i, val := range cy.MonthSums {
		res[i+1] = fmt.Sprintf("$%.0f", float32(val)/100)
	}
	return res
}

type FilterOpts struct {
	minCost   float32
	maxCost   float32
	startDate time.Time
	endDate   time.Time
	catIds    []int
	code      string
}

func NewFilterOpts() FilterOpts {
	/*
	  Set default options for filters, allow functions to be passed to modify these
	*/
	startDate, _ := makeDate(2000, 1, 1)
	endDate, _ := makeDate(3000, 1, 1)

	opts := &FilterOpts{
		minCost:   -10000000,
		maxCost:   10000000,
		startDate: startDate,
		endDate:   endDate,
		catIds:    []int{},
		code:      "",
	}

	return *opts
}

func (opts FilterOpts) WithMinCost(val float32) FilterOpts {
	opts.minCost = val
	return opts
}

func (opts FilterOpts) WithMaxCost(val float32) FilterOpts {
	opts.maxCost = val
	return opts
}

func (opts FilterOpts) WithStartDate(val time.Time) FilterOpts {
	opts.startDate = val
	return opts
}

func (opts FilterOpts) WithEndDate(val time.Time) FilterOpts {
	opts.endDate = val
	return opts
}

func (opts FilterOpts) WithCatId(val []int) FilterOpts {
	opts.catIds = val
	return opts
}

func (opts FilterOpts) WithCode(val string) FilterOpts {
	opts.code = val
	return opts
}
