package structures

//data entry from the energyData.csv file
type DataEntry struct {
	Country     string  `json:"name"`
	CountryCode string  `json:"isoCode"`
	Year        int     `json:"year"`
	Percentage  float64 `json:"percentage"`
}

//countyr struct for the country collecting third party service
type Country struct {
	Borders     []string `json:"borders"`
	CountryCode string   `json:"cca2"`
	Name        struct {
		Common string `json:"common"`
	} `json:"name"`
}

//to store information
type Info struct {
	RESTStatus  int     `json:"countries_api"`
	NotifStatus int     `json:"notification_db"`
	Webhooks    int     `json:"webhooks"`
	Version     string  `json:"version"`
	Uptime      float64 `json:"uptime"`
}

//content of a webhook
type Webhook struct {
	URL     string `json:"url"`
	Country string `json:"country"`
	Calls   int    `json:"calls"`
}

//id given to each webhook
type WebhookRegistration struct {
	ID      string `json:"webhook_id"`
	Webhook Webhook
}

//structure of the binary search three
type BSTNode struct {
	Data   []DataEntry
	Letter rune
	Left   *BSTNode
	Right  *BSTNode
}
