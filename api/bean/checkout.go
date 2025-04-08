package bean

type CheckoutSignIn struct {
	Redirect              bool   `json:"redirect" description:"should redirect to sign in page"`
	Url                   string `json:"url" description:"redirect url"`
	DuplicateSubscription bool   `json:"DuplicateSubscription" description:"whether contain active or incomplete subscription or not"`
}
