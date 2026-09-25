package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Checks struct {
	Plus  int    `json:"plus"`
	Minus int    `json:"minus"`
	J     string `json:"jamoat"`
	Q     string `json:"qazo"`
}

type Dates struct {
	Dates  int    `json:"dates"`
	Months string `json:"month"`
}

type Database struct {
	Db_Checks []Checks `json:"db_checks"`
	Db_Dates  []Dates  `json:"db_dates"`
}

type help struct {
	Checks []Checks `json:"checks"`
	Dates  []Dates  `json:"dates"`
}

func CheckJson() (*Database, error) {
	data, err := ioutil.ReadFile("checks.json")
	if err != nil {
		return nil, err
	}

	var db Database

	if err := json.Unmarshal(data, &db.Db_Checks); err == nil {
		return &db, nil
	}

	if err := json.Unmarshal(data, &db); err != nil {
		return &db, nil
	}

	return &db, nil
}

func DatesJson() (*Database, error) {
	data, err := ioutil.ReadFile("dates.json")
	if err != nil {
		return nil, err
	}

	var db Database

	if err := json.Unmarshal(data, &db.Db_Checks); err == nil {
		return &db, nil
	}

	if err := json.Unmarshal(data, &db); err != nil {
		return &db, nil
	}

	return &db, nil
}

func PostCh(c *gin.Context) {
	var request help

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	chec, err := CheckJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	date, err := DatesJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// dates := request.Dates

	NewDate := 1

	for _, w := range date.Db_Dates {
		if w.Dates >= NewDate {
			NewDate = w.Dates + 1
		}
	}

	c.IndentedJSON(http.StatusCreated, gin.H{
		"date":   date,
		"change": chec,
	})
}

func GetHistory(c *gin.Context) {
	db, err := CheckJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"History": db,
	})
}
