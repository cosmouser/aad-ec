package actions

// AccessResponse holds the json response data from the Graph API
type AccessResponse struct {
	TokenType    string `json:"token_type"`     // "Bearer"
	ExpiresIn    string `json:"expires_in"`     // "3600",
	ExtExpiresIn string `json:"ext_expires_in"` // "10800",
	ExpiresOn    string `json:"expires_on"`     // "1488429872",
	NotBefore    string `json:"not_before"`     // "1488425972",
	Resource     string `json:"resource"`       // "https://management.core.windows.net/",
	AccessToken  string `json:"access_token"`   // "eyJ0eBAi3n..."
}

// APResponse is the response that graph gives for the getPlans call
type APResponse struct {
	Value []AssignedPlan `json:"value"`
}

// AssignedLicense represents a license assigned to a user
type AssignedLicense struct {
	DisabledPlans []string `json:"disabledPlans"`
	SkuID         string   `json:"skuId"`
}

// AssignedPlan represents a plan assigned to a user
type AssignedPlan struct {
	AssignedDateTime string `json:"assignedDateTime"`
	CapabilityStatus string `json:"capabilityStatus"`
	Service          string `json:"service"`
	ServicePlanID    string `json:"servicePlanId"`
}

// GetPlansResponse is the json response for the getPlans call
type GetPlansResponse struct {
	AssignedPlans []AssignedPlan `json:"assignedPlans"`
}

// ErrorResponse is the root of the json error response
type ErrorResponse struct {
	Error Error `json:"error"`
}

// Error contains information about the error
type Error struct {
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	InnerError InnerError `json:"innerError"`
}

// InnerError contains error metadata
type InnerError struct {
	RequestID string `json:"requestId"`
	Date      string `json:"date"`
}
