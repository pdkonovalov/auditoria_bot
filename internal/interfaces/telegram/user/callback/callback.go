package callback

import (
	"fmt"
	"strings"
)

const (
	Event = "event"
)

const (
	EventsByDate = "eventsbydate"
)

const (
	Booking = "booking"
)

const (
	EditBooking        = "editbooking"
	EditFormat         = "editformat"
	EditPayment        = "editpayment"
	EditAdditionalInfo = "editadditionalinfo"
)

const (
	DeleteBooking = "deletebooking"
)

const (
	ShowBooking = "showbooking"
)

const (
	EditContactInfo = "editcontactinfo"
)

func Encode(data map[string]string) string {
	parts := make([]string, 0)
	for key, value := range data {
		parts = append(parts,
			fmt.Sprintf("%s=%s", key, value),
		)
	}
	return strings.Join(parts, "&")
}

func Decode(b string) (map[string]string, error) {
	if b == "" {
		return make(map[string]string), nil
	}
	parts := strings.Split(b, "&")
	data := make(map[string]string)
	for _, part := range parts {
		keyvalue := strings.Split(part, "=")
		if len(keyvalue) != 2 {
			return nil, fmt.Errorf("Invalid callback data format")
		}
		data[keyvalue[0]] = keyvalue[1]
	}
	return data, nil
}
