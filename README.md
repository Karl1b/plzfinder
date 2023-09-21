# plzfinder
Finds and returns surrounding PLZ (postal code) in GERMANY with a specific radius around a PLZ

```bash
go get github.com/karl1b/plzfinder
```

```go
// init: run this once for your service to load the data from csv into ram
// point to the csv file location.
plzfinder.LoadCSV("zipcodes.de.csv")

// tbis is how you search for all PLZ in 5km radius around 06110 
orte, err := plzfinder.FindeOrte("06110", 5)
if err != nil {
	// handle error
}
```

I got the DATA from HERE:
https://github.com/zauberware/postal-codes-json-xml-csv

Thanks matey! ;-)




