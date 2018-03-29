package main

import (
	"errors"
	"strconv"
	"time"
)

// UrssafToPeriod convertit le format de période urssaf en type Period
func UrssafToPeriod(urssaf string) (Periode, error) {
	// format en 4 ou 6 caractère YYQM ou YYYYQM
	// si YY < 50 alors YYYY = 20YY sinon YYYY = 19YY
	// si QM == 62 alors période annuelle sur YYYY
	// si M == 0 alors période trimestrielle sur le trimestre Q de YYYY
	// si 0 < M < 4 alors mois M du trimestre Q

	period := Periode{}

	if len(urssaf) == 4 {
		if urssaf[0:2] < "50" {
			urssaf = "20" + urssaf
		} else {
			urssaf = "19" + urssaf
		}
	}

	if len(urssaf) != 6 {
		return period, errors.New("Valeur non autorisée")
	}

	year, err := strconv.Atoi(urssaf[0:4])
	if err != nil {
		return period, err
	}

	if urssaf[4:6] == "62" {
		period.Start = time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		period.End = time.Date(year+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	} else {
		quarter, err := strconv.Atoi(urssaf[4:5])
		if err != nil {
			return period, err
		}
		monthOfQuarter, err := strconv.Atoi(urssaf[5:6])
		if err != nil {
			return period, err
		}
		if monthOfQuarter == 0 {
			period.Start = time.Date(year, time.Month((quarter-1)*3+1), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+4), 1, 0, 0, 0, 0, time.UTC)
		} else {
			period.Start = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter+1), 1, 0, 0, 0, 0, time.UTC)
		}
	}
	return period, nil
}

// UrssafToDate Convertit le format de date urssaf en type Date
func UrssafToDate(urssaf string) (time.Time, error) {
	// Date au format YYYMMJJ
	// YYY = YYYY - 1900

	intUrsaff, err := strconv.Atoi(urssaf)
	if err != nil {
		return time.Time{}, errors.New("Valeur non autorisée")
	}
	strDate := strconv.Itoa(intUrsaff + 19000000)
	date, err := time.Parse("20060102", strDate)
	if err != nil {
		return time.Time{}, errors.New("Valeur non autorisée")
	}

	return date, nil
}
