package contacts

type Contact struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	FName        string `json:"fname"`
	LName        string `json:"lname"`
	EmailAddress string `json:"emailaddress"`
	Subscribed   bool   `json:"subscribed"`
	Unsubscribed bool   `json:"unsubscribed"`
}

type ContactParams struct {
	FName        string
	LName        string
	EmailAddress string
	Subscribed   bool
	Unsubscribed bool
}

func NewContact(params *ContactParams) *Contact {
	return &Contact{
		FName:        params.FName,
		LName:        params.LName,
		EmailAddress: params.EmailAddress,
		Subscribed:   params.Subscribed,
		Unsubscribed: params.Unsubscribed,
	}
}

type ContactStore interface {
	CreateContact(fname, lname, email string, subscribed, unsubscribed bool) error
	GetContact(email string) (*Contact, error)
	SearchContacts(email string) (*[]Contact, error)
	GetSubscribersList() (*[]Contact, error)
}
