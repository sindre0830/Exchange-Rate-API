package fun

import (
	"errors"
	"strconv"
	"strings"
)

func ValidateDates(dates string) error {
	var err error
	//split date by '-' for format checking
	arrDate := strings.Split(dates, "-")
	//check if date has correct amount of elements
	if (len(arrDate) != 6) || (len(dates) != 21) {
		err = errors.New("date: not enough elements in input")
		return err
	}
	//check if start date is using correct format YYYY-MM-DD
	if (len(arrDate[0]) != 4) || (len(arrDate[1]) != 2) || (len(arrDate[2]) != 2) {
		err = errors.New("date: start date doesn't follow required date format YYYY-MM-DD")
		return err
	}
	//check if end date is using correct format YYYY-MM-DD
	if (len(arrDate[3]) != 4) || (len(arrDate[4]) != 2) || (len(arrDate[5]) != 2) {
		err = errors.New("date: end date doesn't follow required date format YYYY-MM-DD")
		return err
	}
	//check if all date elements are integers and larger than 0. 'hehe-01-00' == false
	for _, elemDate := range arrDate {
		elemDateNum, err := strconv.Atoi(elemDate)
		if err != nil || elemDateNum < 1 {
			err = errors.New("date: wrong type, should be numbers that are larger than 0")
			return err
		}
	}
	//check if end date is larger or equal than start date
	if dates[:10] > dates[11:] {
		err = errors.New("date: start date is larger than end date")
	}
	return err
}