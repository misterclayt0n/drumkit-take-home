package turvo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"drumkit-take-home/internal/load"
)

type createShipmentResponse struct {
	Status  string `json:"Status"`
	Details struct {
		ID            int    `json:"id"`
		Created       string `json:"created"`
		CreatedDate   string `json:"createdDate"`
		LastUpdatedOn string `json:"lastUpdatedOn"`
		StatusHistory []struct {
			LastUpdatedOn string `json:"lastUpdatedOn"`
		} `json:"statusHistory"`
	} `json:"details"`
}

type createShipmentRequest struct {
	LTLShipment    bool                          `json:"ltlShipment"`
	StartDate      shipmentDateTime              `json:"startDate"`
	EndDate        shipmentDateTime              `json:"endDate"`
	Status         *shipmentCreateStatus         `json:"status,omitempty"`
	Lane           shipmentLane                  `json:"lane"`
	GlobalRoute    []shipmentCreateStop          `json:"globalRoute,omitempty"`
	Transportation shipmentCreateTransport       `json:"transportation"`
	CustomerOrder  []shipmentCreateCustomerOrder `json:"customerOrder"`
}

type shipmentDateTime struct {
	Date     string `json:"date"`
	TimeZone string `json:"timeZone"`
}

type shipmentCreateStatus struct {
	Code shipmentCode `json:"code"`
}

type shipmentLane struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type shipmentCreateStop struct {
	Name            string                    `json:"name,omitempty"`
	StopType        shipmentCode              `json:"stopType"`
	Timezone        string                    `json:"timezone,omitempty"`
	Location        shipmentCreateLocation    `json:"location"`
	SegmentSequence int                       `json:"segmentSequence"`
	Sequence        int                       `json:"sequence"`
	State           string                    `json:"state"`
	Appointment     shipmentCreateAppointment `json:"appointment"`
	PONumbers       []string                  `json:"poNumbers,omitempty"`
	Notes           string                    `json:"notes,omitempty"`
}

type shipmentCreateLocation struct {
	ID int `json:"id"`
}

type shipmentCreateAppointment struct {
	Date     string `json:"date"`
	Timezone string `json:"timezone,omitempty"`
	Flex     int    `json:"flex"`
	HasTime  bool   `json:"hasTime"`
}

type shipmentCreateTransport struct {
	Mode        shipmentCode `json:"mode"`
	ServiceType shipmentCode `json:"serviceType"`
}

type shipmentCreateCustomerOrder struct {
	Customer    shipmentCreateCustomer     `json:"customer"`
	ExternalIDs []shipmentCreateExternalID `json:"externalIds,omitempty"`
}

type shipmentCreateCustomer struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type shipmentCreateExternalID struct {
	Type               shipmentCode `json:"type"`
	Value              string       `json:"value"`
	CopyToCarrierOrder bool         `json:"copyToCarrierOrder"`
}

func (c *Client) CreateLoad(ctx context.Context, input load.Load) (load.CreateResponse, error) {
	requestBody, err := buildCreateShipmentRequest(input)
	if err != nil {
		return load.CreateResponse{}, err
	}

	var response createShipmentResponse
	if err := c.postJSON(ctx, "/shipments?fullResponse=true", requestBody, &response); err != nil {
		return load.CreateResponse{}, err
	}

	return load.CreateResponse{
		ID:        response.Details.ID,
		CreatedAt: firstNonEmpty(response.Details.Created, response.Details.CreatedDate, firstStatusHistoryTimestamp(response.Details.StatusHistory), time.Now().UTC().Format(time.RFC3339Nano)),
	}, nil
}

func buildCreateShipmentRequest(input load.Load) (createShipmentRequest, error) {
	customerID, err := parseRequiredInt(input.Customer.ExternalTMSID, "customer.externalTMSId")
	if err != nil {
		return createShipmentRequest{}, err
	}

	pickupLocationID, err := parseRequiredInt(firstNonEmpty(input.Pickup.ExternalTMSID, input.Pickup.WarehouseID), "pickup.externalTMSId")
	if err != nil {
		return createShipmentRequest{}, err
	}

	deliveryLocationID, err := parseRequiredInt(firstNonEmpty(input.Consignee.ExternalTMSID, input.Consignee.WarehouseID), "consignee.externalTMSId")
	if err != nil {
		return createShipmentRequest{}, err
	}

	pickupTime, err := requireTimestamp(firstNonEmpty(input.Pickup.ApptTime, input.Pickup.ReadyTime), "pickup.apptTime")
	if err != nil {
		return createShipmentRequest{}, err
	}

	deliveryTime, err := requireTimestamp(firstNonEmpty(input.Consignee.ApptTime, input.Consignee.MustDeliver), "consignee.apptTime")
	if err != nil {
		return createShipmentRequest{}, err
	}

	pickupTimezone := firstNonEmpty(input.Pickup.Timezone, "America/New_York")
	deliveryTimezone := firstNonEmpty(input.Consignee.Timezone, pickupTimezone, "America/New_York")
	poNumbers := splitCSV(input.PONums)
	externalIDs := buildExternalIDs(input)
	status := buildShipmentStatus(input.Status)

	return createShipmentRequest{
		LTLShipment: false,
		StartDate: shipmentDateTime{
			Date:     pickupTime,
			TimeZone: pickupTimezone,
		},
		EndDate: shipmentDateTime{
			Date:     deliveryTime,
			TimeZone: deliveryTimezone,
		},
		Status: status,
		Lane: shipmentLane{
			Start: formatLane(input.Pickup.City, input.Pickup.State, input.Pickup.Country),
			End:   formatLane(input.Consignee.City, input.Consignee.State, input.Consignee.Country),
		},
		GlobalRoute: []shipmentCreateStop{
			{
				Name:            input.Pickup.Name,
				StopType:        shipmentCode{Key: "1500", Value: "Pickup"},
				Timezone:        pickupTimezone,
				Location:        shipmentCreateLocation{ID: pickupLocationID},
				SegmentSequence: 0,
				Sequence:        0,
				State:           "OPEN",
				Appointment: shipmentCreateAppointment{
					Date:     pickupTime,
					Timezone: pickupTimezone,
					Flex:     0,
					HasTime:  true,
				},
				PONumbers: poNumbers,
				Notes:     input.Pickup.ApptNote,
			},
			{
				Name:            input.Consignee.Name,
				StopType:        shipmentCode{Key: "1501", Value: "Delivery"},
				Timezone:        deliveryTimezone,
				Location:        shipmentCreateLocation{ID: deliveryLocationID},
				SegmentSequence: 0,
				Sequence:        1,
				State:           "OPEN",
				Appointment: shipmentCreateAppointment{
					Date:     deliveryTime,
					Timezone: deliveryTimezone,
					Flex:     0,
					HasTime:  true,
				},
				PONumbers: poNumbers,
				Notes:     input.Consignee.ApptNote,
			},
		},
		Transportation: shipmentCreateTransport{
			Mode:        shipmentCode{Key: "24105", Value: "TL"},
			ServiceType: shipmentCode{Key: "24304", Value: "Any"},
		},
		CustomerOrder: []shipmentCreateCustomerOrder{
			{
				Customer: shipmentCreateCustomer{
					ID:   customerID,
					Name: input.Customer.Name,
				},
				ExternalIDs: externalIDs,
			},
		},
	}, nil
}

func buildShipmentStatus(status string) *shipmentCreateStatus {
	trimmed := strings.TrimSpace(strings.ToLower(status))
	if trimmed == "" {
		return nil
	}

	code, ok := map[string]shipmentCode{
		"tendered": {Key: "2101", Value: "Tendered"},
		"covered":  {Key: "2102", Value: "Covered"},
	}[trimmed]
	if !ok {
		return nil
	}

	return &shipmentCreateStatus{Code: code}
}

func buildExternalIDs(input load.Load) []shipmentCreateExternalID {
	var externalIDs []shipmentCreateExternalID
	seen := map[string]struct{}{}

	add := func(key, value, raw string) {
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" {
			return
		}
		identity := key + ":" + trimmed
		if _, ok := seen[identity]; ok {
			return
		}
		seen[identity] = struct{}{}
		externalIDs = append(externalIDs, shipmentCreateExternalID{
			Type:               shipmentCode{Key: key, Value: value},
			Value:              trimmed,
			CopyToCarrierOrder: false,
		})
	}

	for _, poNumber := range splitCSV(input.PONums) {
		add("1400", "Purchase order #", poNumber)
	}
	for _, ref := range splitCSV(input.Customer.RefNumber) {
		add("1401", "Reference Number", ref)
	}
	add("1401", "Reference Number", input.ExternalTMSLoadID)
	add("1401", "Reference Number", input.FreightLoadID)

	return externalIDs
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	seen := map[string]struct{}{}
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		items = append(items, trimmed)
	}
	return items
}

func formatLane(city, state, country string) string {
	parts := make([]string, 0, 3)
	for _, part := range []string{city, state, country} {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return strings.Join(parts, ", ")
}

func parseRequiredInt(value, field string) (int, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, fmt.Errorf("%s is required", field)
	}

	parsed, err := strconv.Atoi(trimmed)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", field)
	}

	return parsed, nil
}

func requireTimestamp(value, field string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s is required", field)
	}

	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return "", fmt.Errorf("%s must be a valid RFC3339 timestamp", field)
	}

	return parsed.UTC().Format(time.RFC3339), nil
}

func firstStatusHistoryTimestamp(entries []struct {
	LastUpdatedOn string `json:"lastUpdatedOn"`
}) string {
	for _, entry := range entries {
		if strings.TrimSpace(entry.LastUpdatedOn) != "" {
			return entry.LastUpdatedOn
		}
	}
	return ""
}
