package atome

import "time"

// Consumption represents the consumption for a specified period of time
type Consumption struct {
	// Total is the total consumed electricity (in Wh)
	Total int
	// Price is the current price (in euros)
	Price float64
	// Co2Impact represents the current CO2 impact of the consumption
	// It's in grams
	Co2Impact int
}

type authenticateRequest struct {
	Email             string `json:"email"`
	MobileInformation string `json:"mobileInformation"`
	PlainPassword     string `json:"plainPassword"`
}

type authenticateResponse struct {
	Firstname          string               `json:"firstname"`
	Lastname           string               `json:"lastname"`
	Phone              string               `json:"phone"`
	ID                 int                  `json:"id"`
	Address1           string               `json:"address1"`
	Address2           interface{}          `json:"address2"`
	Phone2             string               `json:"phone2"`
	Zipcode            string               `json:"zipcode"`
	City               string               `json:"city"`
	EmailNotification  string               `json:"emailNotification"`
	Subscriptions      []subscriptions      `json:"subscriptions"`
	NtfUserPreferences []ntfUserPreferences `json:"ntf_user_preferences"`
}

type lastConsumptionResponse struct {
	IsConnected       bool      `json:"isConnected"`
	Time              time.Time `json:"time"`
	Total             int       `json:"total"`
	Price             float64   `json:"price"`
	LogTargetStd      float64   `json:"logTargetStd"`
	StartPeriod       time.Time `json:"startPeriod"`
	EndPeriod         time.Time `json:"endPeriod"`
	PricingChange     int       `json:"pricingChange"`
	ImpactCo2         int       `json:"impactCo2"`
	ImpactVehicleType string    `json:"impactVehicleType"`
	ImpactVehicleKm   int       `json:"impactVehicleKm"`
	BasicPrice        float64   `json:"basicPrice"`
}

type probe struct {
	ID                string      `json:"id"`
	Prm               string      `json:"prm"`
	HardwareVersion   int         `json:"hardwareVersion"`
	FirmwareVersion   int         `json:"firmwareVersion"`
	Address1          string      `json:"address1"`
	Address2          interface{} `json:"address2"`
	Zipcode           string      `json:"zipcode"`
	City              string      `json:"city"`
	Country           string      `json:"country"`
	Area              int         `json:"area"`
	Building          interface{} `json:"building"`
	Housing           string      `json:"housing"`
	Level             int         `json:"level"`
	FirmwareUpdatedAt time.Time   `json:"firmwareUpdatedAt"`
	Rate              int         `json:"rate"`
	Wifi              bool        `json:"wifi"`
	MeterType         int         `json:"meterType"`
	Rssi              float64     `json:"rssi"`
}
type pricing struct {
	ID      string      `json:"id"`
	StartAt time.Time   `json:"startAt"`
	EndAt   interface{} `json:"endAt"`
	Name    string      `json:"name"`
	Cost    float64     `json:"cost"`
	Index   int         `json:"index"`
	Color   string      `json:"color"`
	Code    string      `json:"code"`
}
type pricePeriods struct {
	ID      string    `json:"id"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Pricing pricing   `json:"pricing"`
}

type day struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	PricePeriods []pricePeriods `json:"pricePeriods"`
}

type weekProfile struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Monday    day    `json:"monday"`
	Tuesday   day    `json:"tuesday"`
	Wednesday day    `json:"wednesday"`
	Thursday  day    `json:"thursday"`
	Friday    day    `json:"friday"`
	Saturday  day    `json:"saturday"`
	Sunday    day    `json:"sunday"`
}
type periods struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Start       time.Time   `json:"start"`
	End         time.Time   `json:"end"`
	WeekProfile weekProfile `json:"weekProfile"`
}

type subscriptions struct {
	ID                         string        `json:"id"`
	Reference                  string        `json:"reference"`
	Heater                     string        `json:"heater"`
	CustomerID                 string        `json:"customerId"`
	SerialNumber               string        `json:"serialNumber"`
	SubscriptionCommercialName string        `json:"subscriptionCommercialName"`
	QuantityPeople             int           `json:"quantityPeople"`
	ErlState                   int           `json:"erlState"`
	TutorialSeen               bool          `json:"tutorialSeen"`
	WizardPassed               bool          `json:"wizardPassed"`
	WizardUpToDate             bool          `json:"wizardUpToDate"`
	Goal                       int           `json:"goal"`
	Probe                      probe         `json:"probe"`
	OrderTrackingURL           string        `json:"orderTrackingURL"`
	LastEligibilityResult      string        `json:"lastEligibilityResult"`
	ErlPaired                  bool          `json:"erlPaired"`
	Periods                    []periods     `json:"periods"`
	Matricule                  string        `json:"matricule"`
	Type                       string        `json:"type"`
	PairingDate                time.Time     `json:"pairingDate"`
	OfferCode                  int           `json:"offerCode"`
	SpecialPricingList         []interface{} `json:"specialPricingList"`
	IsGreen                    bool          `json:"isGreen"`
	IsOnline                   bool          `json:"isOnline"`
}

type ntfUserPreferences struct {
	CustomMinPowerLimitEnabled   int  `json:"custom_min_power_limit_enabled"`
	CustomMaxPowerLimitEnabled   int  `json:"custom_max_power_limit_enabled"`
	CustomPowerLimitEnabled      int  `json:"custom_power_limit_enabled"`
	CustomPowerLimitMinW         int  `json:"custom_power_limit_min_w"`
	CustomPowerLimitMaxW         int  `json:"custom_power_limit_max_w"`
	PowerLimitEnabled            int  `json:"power_limit_enabled"`
	EnergyOverrunMinWh           int  `json:"energy_overrun_min_wh"`
	EnergyOverrunMaxWh           int  `json:"energy_overrun_max_wh"`
	EnergyOverrunWhEnabled       int  `json:"energy_overrun_wh_enabled"`
	EnergyOverrunMinEur          int  `json:"energy_overrun_min_eur"`
	EnergyOverrunMaxEur          int  `json:"energy_overrun_max_eur"`
	EnergyOverrunEurEnabled      int  `json:"energy_overrun_eur_enabled"`
	UnplugEnabled                int  `json:"unplug_enabled"`
	EnergyConsumedDailyEnabled   int  `json:"energy_consumed_daily_enabled"`
	EnergyConsumedWeeklyEnabled  int  `json:"energy_consumed_weekly_enabled"`
	EnergyConsumedMonthlyEnabled int  `json:"energy_consumed_monthly_enabled"`
	ExpertMode                   bool `json:"expert_mode"`
	Co2DailyNotifEnable          bool `json:"co2_daily_notif_enable"`
	Co2EmissionPeakNotifEnable   bool `json:"co2_emission_peak_notif_enable"`
}
