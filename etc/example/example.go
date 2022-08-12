package entity

type Account struct {
	ID string
}

type User struct {
	ID string

	AccountID string `dddoc:"belongsTo, Account"`
}

type Billing struct {
	ID string

	AccountID string `dddoc:"belongsTo, Account"`

	Items []*BillingItem `dddoc:"hasMany, BillingItem, hasMany(0..n)"`
}

type BillingItem struct {
	ID string
}

type AccountFactory interface{}

type UserFactory interface{}

type BillingFactory interface{}

type AccountRepository interface{}

type UserRepository interface{}

type BillingRepository interface{}
