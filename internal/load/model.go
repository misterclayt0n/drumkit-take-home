package load

type Load struct {
	ExternalTMSLoadID string         `json:"externalTMSLoadID"`
	FreightLoadID     string         `json:"freightLoadID"`
	Status            string         `json:"status"`
	Customer          Customer       `json:"customer"`
	BillTo            BillTo         `json:"billTo"`
	Pickup            Pickup         `json:"pickup"`
	Consignee         Consignee      `json:"consignee"`
	Carrier           Carrier        `json:"carrier"`
	RateData          RateData       `json:"rateData"`
	Specifications    Specifications `json:"specifications"`
	InPalletCount     int            `json:"inPalletCount"`
	OutPalletCount    int            `json:"outPalletCount"`
	NumCommodities    int            `json:"numCommodities"`
	TotalWeight       float64        `json:"totalWeight"`
	BillableWeight    float64        `json:"billableWeight"`
	PONums            string         `json:"poNums"`
	Operator          string         `json:"operator"`
	RouteMiles        float64        `json:"routeMiles"`
}

type Customer struct {
	ExternalTMSID string `json:"externalTMSId"`
	Name          string `json:"name"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zipcode       string `json:"zipcode"`
	Country       string `json:"country"`
	Contact       string `json:"contact"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	RefNumber     string `json:"refNumber"`
}

type BillTo struct {
	ExternalTMSID string `json:"externalTMSId"`
	Name          string `json:"name"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zipcode       string `json:"zipcode"`
	Country       string `json:"country"`
	Contact       string `json:"contact"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
}

type Pickup struct {
	ExternalTMSID string `json:"externalTMSId"`
	Name          string `json:"name"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zipcode       string `json:"zipcode"`
	Country       string `json:"country"`
	Contact       string `json:"contact"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	BusinessHours string `json:"businessHours"`
	RefNumber     string `json:"refNumber"`
	ReadyTime     string `json:"readyTime"`
	ApptTime      string `json:"apptTime"`
	ApptNote      string `json:"apptNote"`
	Timezone      string `json:"timezone"`
	WarehouseID   string `json:"warehouseId"`
}

type Consignee struct {
	ExternalTMSID string `json:"externalTMSId"`
	Name          string `json:"name"`
	AddressLine1  string `json:"addressLine1"`
	AddressLine2  string `json:"addressLine2"`
	City          string `json:"city"`
	State         string `json:"state"`
	Zipcode       string `json:"zipcode"`
	Country       string `json:"country"`
	Contact       string `json:"contact"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	BusinessHours string `json:"businessHours"`
	RefNumber     string `json:"refNumber"`
	MustDeliver   string `json:"mustDeliver"`
	ApptTime      string `json:"apptTime"`
	ApptNote      string `json:"apptNote"`
	Timezone      string `json:"timezone"`
	WarehouseID   string `json:"warehouseId"`
}

type Carrier struct {
	MCNumber                 string `json:"mcNumber"`
	DOTNumber                string `json:"dotNumber"`
	Name                     string `json:"name"`
	Phone                    string `json:"phone"`
	Dispatcher               string `json:"dispatcher"`
	SealNumber               string `json:"sealNumber"`
	SCAC                     string `json:"scac"`
	FirstDriverName          string `json:"firstDriverName"`
	FirstDriverPhone         string `json:"firstDriverPhone"`
	SecondDriverName         string `json:"secondDriverName"`
	SecondDriverPhone        string `json:"secondDriverPhone"`
	Email                    string `json:"email"`
	DispatchCity             string `json:"dispatchCity"`
	DispatchState            string `json:"dispatchState"`
	ExternalTMSTruckID       string `json:"externalTMSTruckId"`
	ExternalTMSTrailerID     string `json:"externalTMSTrailerId"`
	ConfirmationSentTime     string `json:"confirmationSentTime"`
	ConfirmationReceivedTime string `json:"confirmationReceivedTime"`
	DispatchedTime           string `json:"dispatchedTime"`
	ExpectedPickupTime       string `json:"expectedPickupTime"`
	PickupStart              string `json:"pickupStart"`
	PickupEnd                string `json:"pickupEnd"`
	ExpectedDeliveryTime     string `json:"expectedDeliveryTime"`
	DeliveryStart            string `json:"deliveryStart"`
	DeliveryEnd              string `json:"deliveryEnd"`
	SignedBy                 string `json:"signedBy"`
	ExternalTMSID            string `json:"externalTMSId"`
}

type RateData struct {
	CustomerRateType  string  `json:"customerRateType"`
	CustomerNumHours  float64 `json:"customerNumHours"`
	CustomerLHRateUSD float64 `json:"customerLhRateUsd"`
	FSCPercent        float64 `json:"fscPercent"`
	FSCPerMile        float64 `json:"fscPerMile"`
	CarrierRateType   string  `json:"carrierRateType"`
	CarrierNumHours   float64 `json:"carrierNumHours"`
	CarrierLHRateUSD  float64 `json:"carrierLhRateUsd"`
	CarrierMaxRate    float64 `json:"carrierMaxRate"`
	NetProfitUSD      float64 `json:"netProfitUsd"`
	ProfitPercent     float64 `json:"profitPercent"`
}

type Specifications struct {
	MinTempFahrenheit float64 `json:"minTempFahrenheit"`
	MaxTempFahrenheit float64 `json:"maxTempFahrenheit"`
	LiftgatePickup    bool    `json:"liftgatePickup"`
	LiftgateDelivery  bool    `json:"liftgateDelivery"`
	InsidePickup      bool    `json:"insidePickup"`
	InsideDelivery    bool    `json:"insideDelivery"`
	Tarps             bool    `json:"tarps"`
	Oversized         bool    `json:"oversized"`
	Hazmat            bool    `json:"hazmat"`
	Straps            bool    `json:"straps"`
	Permits           bool    `json:"permits"`
	Escorts           bool    `json:"escorts"`
	Seal              bool    `json:"seal"`
	CustomBonded      bool    `json:"customBonded"`
	Labor             bool    `json:"labor"`
}

type ListParams struct {
	Status               string
	CustomerID           string
	PickupDateSearchFrom string
	PickupDateSearchTo   string
	Page                 int
	Limit                int
}

type ListResponse struct {
	Data       []Load     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Total int `json:"total"`
	Pages int `json:"pages"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type CreateResponse struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"createdAt"`
}
