package state

import (
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
)

const (
	Init = entity.UserStateInit
)

const (
	NewEventWaitInputFormat         = "newevent/waitinput/format"
	NewEventWaitInputPaid           = "newevent/waitinput/paid"
	NewEventWaitInputTitle          = "newevent/waitinput/title"
	NewEventWaitInputTime           = "newevent/waitinput/time"
	NewEventWaitInputPaymentDetails = "newevent/waitinput/paymentdetails"
	NewEventWaitInputPhotoText      = "newevent/waitinput/phototext"
)

const (
	EditEventWaitInputFormat         = "editevent/waitinput/format"
	EditEventWaitInputPaid           = "editevent/waitinput/paid"
	EditEventWaitInputTitle          = "editevent/waitinput/title"
	EditEventWaitInputTime           = "editevent/waitinput/time"
	EditEventWaitInputPaymentDetails = "editevent/waitinput/paymentdetails"
	EditEventWaitInputPhotoText      = "editevent/waitinput/phototext"
)

const (
	SendNotificationWaitInputPhotoText = "sendnotification/waitinput/phototext"
)

const (
	DeleteEventWaitInputAreYouSure = "deleteevent/waitinput/areyousure"
)
