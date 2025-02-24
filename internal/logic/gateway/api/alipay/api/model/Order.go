package model

type Store struct {
	ReferenceStoreId string   `json:"referenceStoreId,omitempty"`
	StoreName        string   `json:"storeName,omitempty"`
	StoreMCC         string   `json:"storeMCC,omitempty"`
	StoreDisplayName string   `json:"storeDisplayName,omitempty"`
	StoreTerminalId  string   `json:"storeTerminalId,omitempty"`
	StoreOperatorId  string   `json:"storeOperatorId,omitempty"`
	StoreAddress     *Address `json:"storeAddress,omitempty"`
	StorePhoneNo     string   `json:"storePhoneNo,omitempty"`
}

type Merchant struct {
	ReferenceMerchantId  string       `json:"referenceMerchantId,omitempty"`
	MerchantMCC          string       `json:"merchantMCC,omitempty"`
	MerchantName         string       `json:"merchantName,omitempty"`
	MerchantDisplayName  string       `json:"merchantDisplayName,omitempty"`
	MerchantAddress      *Address     `json:"merchantAddress,omitempty"`
	MerchantRegisterDate string       `json:"merchantRegisterDate,omitempty"`
	Store                *Store       `json:"store,omitempty"`
	MerchantType         MerchantType `json:"merchantType,omitempty"`
}

type Goods struct {
	ReferenceGoodsId   string             `json:"referenceGoodsId,omitempty"`
	GoodsName          string             `json:"goodsName,omitempty"`
	GoodsCategory      string             `json:"goodsCategory,omitempty"`
	GoodsBrand         string             `json:"goodsBrand,omitempty"`
	GoodsUnitAmount    *Amount            `json:"goodsUnitAmount,omitempty"`
	GoodsQuantity      string             `json:"goodsQuantity,omitempty"`
	GoodsSkuName       string             `json:"goodsSkuName,omitempty"`
	GoodsUrl           string             `json:"goodsUrl,omitempty"`
	DeliveryMethodType DeliveryMethodType `json:"deliveryMethodType,omitempty"`
	GoodsImageUrl      string             `json:"goodsImageUrl,omitempty"`
	PriceId            string             `json:"priceId,omitempty"`
}
type DeliveryEstimateInfo struct {
	Unit  string `json:"unit,omitempty"`
	Value string `json:"value,omitempty"`
}

type DeliveryEstimate struct {
	Minimum *DeliveryEstimateInfo `json:"minimum,omitempty"`
	Maximum *DeliveryEstimateInfo `json:"maximum,omitempty"`
}

type Shipping struct {
	ShippingName        *UserName         `json:"shippingName,omitempty"`
	ShippingAddress     *Address          `json:"shippingAddress,omitempty"`
	ShippingCarrier     string            `json:"shippingCarrier,omitempty"`
	ShippingPhoneNo     string            `json:"shippingPhoneNo,omitempty"`
	ShipToEmail         string            `json:"shipToEmail,omitempty"`
	ShippingFeeId       string            `json:"shippingFeeId,omitempty"`
	ShippingFee         *Amount           `json:"shippingFee,omitempty"`
	ShippingDescription string            `json:"shippingDescription,omitempty"`
	DeliveryEstimate    *DeliveryEstimate `json:"deliveryEstimate,omitempty"`
}

type Buyer struct {
	ReferenceBuyerId      string    `json:"referenceBuyerId,omitempty"`
	BuyerName             *UserName `json:"buyerName,omitempty"`
	BuyerPhoneNo          string    `json:"buyerPhoneNo,omitempty"`
	BuyerEmail            string    `json:"buyerEmail,omitempty"`
	BuyerRegistrationTime string    `json:"buyerRegistrationTime,omitempty"`
	IsAccountVerified     *bool     `json:"isAccountVerified,omitempty"`
	SuccessfulOrderCount  *int      `json:"successfulOrderCount,omitempty"`
}

type BrowserInfo struct {
	AcceptHeader      string `json:"acceptHeader,omitempty"`
	JavaEnabled       *bool  `json:"javaEnabled,omitempty"`
	JavaScriptEnabled *bool  `json:"javaScriptEnabled,omitempty"`
	Language          string `json:"language,omitempty"`
	UserAgent         string `json:"userAgent,omitempty"`
}

type Env struct {
	TerminalType             TerminalType `json:"terminalType,omitempty"`
	OsType                   OsType       `json:"osType,omitempty"`
	UserAgent                string       `json:"userAgent,omitempty"`
	DeviceTokenId            string       `json:"deviceTokenId,omitempty"`
	ClientIp                 string       `json:"clientIp,omitempty"`
	CookieId                 string       `json:"cookieId,omitempty"`
	ExtendInfo               string       `json:"extendInfo,omitempty"`
	StoreTerminalId          string       `json:"storeTerminalId,omitempty"`
	StoreTerminalRequestTime string       `json:"storeTerminalRequestTime,omitempty"`
	BrowserInfo              *BrowserInfo `json:"browserInfo,omitempty"`
	ColorDepth               string       `json:"colorDepth,omitempty"`
	ScreenHeight             string       `json:"screenHeight,omitempty"`
	ScreenWidth              string       `json:"screenWidth,omitempty"`
	TimeZoneOffset           int          `json:"timeZoneOffset,omitempty"`
	DeviceBrand              string       `json:"deviceBrand,omitempty"`
	DeviceModel              string       `json:"deviceModel,omitempty"`
	DeviceLanguage           string       `json:"deviceLanguage,omitempty"`
	DeviceId                 string       `json:"deviceId,omitempty"`
}

type Leg struct {
	DepartureTime        string    `json:"departureTime,omitempty"`
	ArrivalTime          string    `json:"arrivalTime,omitempty"`
	DepartureAddress     Address   `json:"departureAddress,omitempty"`
	ArrivalAddress       Address   `json:"arrivalAddress,omitempty"`
	CarrierName          string    `json:"carrierName,omitempty"`
	CarrierNo            string    `json:"carrierNo,omitempty"`
	ClassType            ClassType `json:"classType,omitempty"`
	DepartureAirportCode string    `json:"departureAirportCode,omitempty"`
	ArrivalAirportCode   string    `json:"arrivalAirportCode,omitempty"`
}

type Passenger struct {
	PassengerName    UserName `json:"passengerName,omitempty"`
	PassengerEmail   string   `json:"passengerEmail,omitempty"`
	PassengerPhoneNo string   `json:"passengerPhoneNo,omitempty"`
}
type Transit struct {
	TransitType TransitType `json:"transitType,omitempty"`
	Legs        []Leg       `json:"legs,omitempty"`
	Passengers  []Passenger `json:"passengers,omitempty"`
}
type Address struct {
	Region   string `json:"region,omitempty"`
	State    string `json:"state,omitempty"`
	City     string `json:"city,omitempty"`
	Address1 string `json:"address1,omitempty"`
	Address2 string `json:"address2,omitempty"`
	ZipCode  string `json:"zipCode,omitempty"`
	Label    string `json:"label,omitempty"`
}

type UserName struct {
	FirstName  string `json:"firstName,omitempty"`
	MiddleName string `json:"middleName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	FullName   string `json:"fullName,omitempty"`
}

type Lodging struct {
	HotelName      string      `json:"hotelName,omitempty"`
	HotelAddress   *Address    `json:"hotelAddress,omitempty"`
	CheckInDate    string      `json:"checkInDate,omitempty"`
	CheckOutDate   string      `json:"checkOutDate,omitempty"`
	NumberOfNights int         `json:"numberOfNights,omitempty"`
	NumberOfRooms  int         `json:"numberOfRooms,omitempty"`
	GuestNames     []*UserName `json:"guestNames,omitempty"`
}

type Gaming struct {
	GameName        string `json:"gameName,omitempty"`
	ToppedUpUser    string `json:"toppedUpUser,omitempty"`
	ToppedUpEmail   string `json:"toppedUpEmail,omitempty"`
	ToppedUpPhoneNo string `json:"toppedUpPhoneNo,omitempty"`
}

type Order struct {
	ReferenceOrderId string    `json:"referenceOrderId,omitempty"`
	OrderDescription string    `json:"orderDescription,omitempty"`
	OrderAmount      *Amount   `json:"orderAmount,omitempty"`
	Merchant         *Merchant `json:"merchant,omitempty"`
	Goods            []Goods   `json:"goods,omitempty"`
	Shipping         *Shipping `json:"shipping,omitempty"`
	Buyer            *Buyer    `json:"buyer,omitempty"`
	Env              *Env      `json:"env,omitempty"`
	ExtendInfo       string    `json:"extendInfo,omitempty"`
	Transit          *Transit  `json:"transit,omitempty"`
	Lodging          *Lodging  `json:"lodging,omitempty"`
	Gaming           *Gaming   `json:"gaming,omitempty"`
	OrderCreatedTime string    `json:"orderCreatedTime,omitempty"`
	NeedDeclaration  bool      `json:"needDeclaration,omitempty"`
}
