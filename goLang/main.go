package main

import (
	"github.com/sattarfeizollahibarough/mygopkg/crawler"
	"github.com/sattarfeizollahibarough/mygopkg/mysqlDB"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	connection := mysqlDB.Initialize("appuser", "123456", "127.0.0.1", "3306", "TSEDB")
	columns := map[string]string{
		"symbolID":   "INT NOT NULL AUTO_INCREMENT PRIMARY KEY",
		"symbol":     "VARCHAR(150)",
		"symbolname": "VARCHAR(150)",
		"tseid":      "BIGINT",
	}
	mysqlDB.CreateTable(connection, "symbols", columns)
	htmlcode := crawler.ReadDynamicPage(`http://www.tsetmc.com/Loader.aspx?ParTree=15131F`)
	regx1 := regexp.MustCompile(`<a class="inst".*?</a>`)
	linkMatches := regx1.FindAllStringSubmatch(htmlcode, -1)
	var previd int64

	for _, linkGroup := range linkMatches {
		regx1 := regexp.MustCompile(`".*?"`)
		regx2 := regexp.MustCompile(`>.*?<`)
		replcr1 := strings.NewReplacer("\"", "")
		var strid, strsymbol string
		strid = regx1.FindAllString(linkGroup[0], -1)[2]
		strsymbol = regx2.FindAllString(linkGroup[0], -1)[0]
		strid = replcr1.Replace(strid)
		replcr1 = strings.NewReplacer(">", "")
		replcr2 := strings.NewReplacer("<", "")
		strsymbol = replcr1.Replace(strsymbol)
		strsymbol = replcr2.Replace(strsymbol)
		tseid, _ := strconv.ParseInt(strid, 0, 64)
		if previd != tseid {
			query := "INSERT INTO symbols(symbol,tseid) VALUES (\"" + strsymbol + "\"," + strid + ")"
			mysqlDB.ExecuteQuery(connection, query)
			previd = tseid
		} else {
			query := "UPDATE symbols SET symbolname=\"" + strsymbol + "\" WHERE tseid=" + strid
			mysqlDB.ExecuteQuery(connection, query)
		}

	}
	

}
