package state

import (
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
)

const (
	Init = entity.UserStateInit
)

const (
	EditContactInfoWaitInput = "editcontactinfo/waitinput"
)

const (
	BookingWaitInputContactInfo    = "booking/waitinput/contactinfo"
	BookingWaitInputFormat         = "booking/waitinput/format"
	BookingWaitInputPayment        = "booking/waitinput/payment"
	BookingWaitInputAdditionalInfo = "booking/waitinput/additionalinfo"
)

const (
	EditBookingWaitInputFormat         = "editbooking/waitinput/format"
	EditBookingWaitInputPayment        = "editbooking/waitinput/payment"
	EditBookingWaitInputAdditionalInfo = "editbooking/waitinput/additionalinfo"
)

const (
	DeleteBookingWaitInputAreYouSure = "deletebooking/waitinput/areyousure"
)
