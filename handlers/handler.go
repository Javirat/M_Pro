package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Database struct {
	Db_Dates  []Dates  `json:"db_dates"`
	Db_Prayer []Prayer `json:"db_prayer"`
	Db_Status []Status `json:"db_status"`
}

type Status struct {
	Status_id   int    `json:"status_id"`
	Status_name string `json:"status_name"`
}

type Prayer struct {
	Prayer_id   int    `json:"prayer_id"`
	Prayer_name string `json:"prayer_name"`
}

type Record struct {
	Date   string `json:"date"`
	Prayer Prayer `json:"prayer"`
	Status Status `json:"status"`
}

type Dates struct {
	Day int64 `json:"day"`
	// Month int64 `json:"month"`
	// Year  int64 `json:"year"`
}

func DatesJson() ([]Record, error) {
	data, err := ioutil.ReadFile("dates.json")
	if os.IsNotExist(err) {
		return []Record{}, nil
	}
	if err != nil {
		return nil, err
	}
	var records []Record

	err = json.Unmarshal(data, &records)

	return records, err
}

func SaveD(records []Record) error {
	data, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("dates.json", data, 0644)
}

func PrJson() ([]Prayer, error) {
	data, err := ioutil.ReadFile("prayer.json")
	if err != nil {
		return nil, err
	}
	var prayers []Prayer

	err = json.Unmarshal(data, &prayers)

	return prayers, err
}

func StJson() ([]Status, error) {
	data, err := ioutil.ReadFile("status.json")
	if err != nil {
		return nil, err
	}
	var st []Status

	err = json.Unmarshal(data, &st)

	return st, err
}

func ChangeCH(c *gin.Context) {
	dateStr := c.Param("date")
	if _, err := time.Parse("02-01-2006", dateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be D-M"})
		return
	}

	prayers, err := PrJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var prayer *Prayer

	for i := range prayers {
		if strings.EqualFold(prayers[i].Prayer_name, c.Param("prayer")) {
			prayer = &prayers[i]
			break
		}
	}
	if prayer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "prayer not found"})
		return
	}

	var req struct {
		StatusID int `json:"status_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statuses, err := StJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var status *Status

	for i := range statuses {
		if statuses[i].Status_id == req.StatusID {
			status = &statuses[i]
			break
		}
	}

	if status == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "status not found"})
		return
	}

	records, err := DatesJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newRec := Record{Date: dateStr, Prayer: *prayer, Status: *status}

	found := false

	for i := range records {
		if records[i].Date == dateStr && records[i].Prayer.Prayer_id == prayer.Prayer_id {
			records[i] = newRec
			found = true
			break
		}
	}

	if !found {
		records = append(records, newRec)
	}

	if err := SaveD(records); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, newRec)
}

func GetHistory(c *gin.Context) {
	dateStr := c.Param("date")
	if _, err := time.Parse("02-01-2006", dateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be D-M"})
		return
	}

	date, err := DatesJson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"dates": date,
	})

}
