package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"
)

type PlzLoc struct {
	Plz  int
	Lat  float64
	Lon  float64
	LatR float64 // This will be calculated 1time to speed the code up
	LonR float64
}

// Global Variables
var locations []PlzLoc // The Locations are stored in the RAM during runtime

const earthRadius = 6371.0

// In the init function the data is read from CSV. This will only run once so no need to speed up.
// Data can easily be udpated by changing the CSV content.
func init() {

	// Open the CSV file
	file, err := os.Open("zipcodes.de.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader from the file
	reader := csv.NewReader(file)

	// Read the header (and ignore it as we assume the columns are always in the same order)
	_, err = reader.Read()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Iterate through the rest of the records
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Convert the fields to integers
		plz, err1 := strconv.Atoi(record[1])
		lat, err2 := strconv.ParseFloat(record[9], 64)
		lon, err3 := strconv.ParseFloat(record[10], 64)

		// Do the radiant converstion once here, so we do not need it do to it later over again.
		latR := lat * math.Pi / 180.0
		lonR := lon * math.Pi / 180.0

		// If any conversions fail, we'll skip this record
		if err1 != nil || err2 != nil || err3 != nil {
			fmt.Printf("Error converting record: %+v\n", record)
			continue
		}

		locations = append(locations, PlzLoc{Plz: plz, Lat: lat, Lon: lon, LatR: latR, LonR: lonR})
	}

}

// Calculates a distance by assuming the earth is a perfect sphere
func haversine(lat1, lon1, lat2, lon2 float64) float64 {

	// Convert degrees to radians

	// Differences
	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	// Haversine formula
	// Complicated, probably slow
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

// Finds the start location by the PLZ the user entered.
func findeStartPunkt(plz int) (PlzLoc, error) {

	if plz == 99999 {
		return PlzLoc{Plz: plz, Lat: 0, Lon: 0}, errors.New("no Fallback found")
	}
	for _, location := range locations {
		if location.Plz == plz {
			return location, nil
		}
	}
	// Uses the next PLZ from the database if there is a number missing. Hopefully it is close ;-)
	return findeStartPunkt(plz + 1)
}

// Quick Check Function.
func isClose(lat1, lon1, lat2, lon2, distance float64) bool {

	latDiff := math.Abs(lat1-lat2) * 111.32 // Average values for Germany
	lonDiff := math.Abs(lon1-lon2) * 70.07

	totalDiff := latDiff + lonDiff

	return totalDiff < distance+10 // 10km Security added
}

func FindeOrte(plz int, radius int) ([]PlzLoc, error) {

	startpunkt, err := findeStartPunkt(plz)
	if err != nil {
		log.Fatal("Startpunkt Err:", err)
	}

	var orte []PlzLoc

	for _, location := range locations {
		if isClose(location.Lat, location.Lon, startpunkt.Lat, startpunkt.Lon, float64(radius)) { // quick check
			if haversine(location.LatR, location.LonR, startpunkt.LatR, startpunkt.LonR) <= float64(radius) { // exact check
				orte = append(orte, location)
			}
		}
	}

	return orte, nil

}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a postal code as an argument.")
		return
	}

	// Convert the argument to an integer
	plz, err := strconv.Atoi(os.Args[1]) // PLZ Mittelpunkt
	if err != nil {
		fmt.Printf("Error: Invalid postal code provided: %s\n", os.Args[1])
		return
	}

	radius, err := strconv.Atoi(os.Args[2]) // Radius in km
	if err != nil {
		fmt.Printf("Error: Invalid postal code provided: %s\n", os.Args[2])
		return
	}

	//fmt.Println(os.Args[0])
	fmt.Println(plz)

	timestart := time.Now()

	orte, err := FindeOrte(plz, radius)
	if err != nil {
		fmt.Println("Error findeOrte", err)
		return
	}

	fmt.Println("Search ended in :", (time.Since(timestart)))

	fmt.Println("Gefundene PLZ: ", len(orte))

}
