package turvo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"drumkit-take-home/internal/load"
)

type shipmentCode struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type shipmentStatus struct {
	Code shipmentCode `json:"code"`
}

type shipmentSummaryCustomer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type shipmentSummaryCustomerOrder struct {
	ID       int                     `json:"id"`
	Deleted  bool                    `json:"deleted"`
	Customer shipmentSummaryCustomer `json:"customer"`
}

type shipmentSummary struct {
	ID            int                            `json:"id"`
	CustomID      string                         `json:"customId"`
	Created       string                         `json:"created"`
	CreatedDate   string                         `json:"createdDate"`
	CustomerOrder []shipmentSummaryCustomerOrder `json:"customerOrder"`
	Status        shipmentStatus                 `json:"status"`
}

type shipmentListResponse struct {
	Status  string `json:"Status"`
	Details struct {
		Pagination Pagination        `json:"pagination"`
		Shipments  []shipmentSummary `json:"shipments"`
	} `json:"details"`
}

type shipmentEquipment struct {
	Deleted   bool    `json:"deleted"`
	Temp      float64 `json:"temp"`
	TempUnits struct {
		Value string `json:"value"`
	} `json:"tempUnits"`
}

type shipmentContributor struct {
	Deleted         bool `json:"deleted"`
	ContributorUser struct {
		Name string `json:"name"`
	} `json:"contributorUser"`
	Title struct {
		Value string `json:"value"`
	} `json:"title"`
}

type shipmentItem struct {
	Deleted     bool        `json:"deleted"`
	Qty         float64     `json:"qty"`
	Weight      float64     `json:"weight"`
	HandlingQty flexibleInt `json:"handlingQty"`
	Unit        struct {
		Value string `json:"value"`
	} `json:"unit"`
	HandlingUnit struct {
		Value string `json:"value"`
	} `json:"handlingUnit"`
	GrossWeight struct {
		Weight float64 `json:"weight"`
	} `json:"grossWeight"`
	NetWeight struct {
		Weight float64 `json:"weight"`
	} `json:"netWeight"`
}

type shipmentExternalID struct {
	Deleted bool `json:"deleted"`
	Type    struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"type"`
	Value string `json:"value"`
}

type shipmentEmail struct {
	Email string `json:"email"`
}

type shipmentBillTo struct {
	ID      string `json:"id"`
	BillTo  string `json:"billTo"`
	Address struct {
		Line1 string `json:"line1"`
		State struct {
			Name string `json:"name"`
		} `json:"state"`
		Zip     string `json:"zip"`
		Country struct {
			Name string `json:"name"`
		} `json:"country"`
		City struct {
			Name string `json:"name"`
		} `json:"city"`
	} `json:"address"`
	Emails []shipmentEmail `json:"emails"`
}

type shipmentFreightTerms struct {
	BillTo shipmentBillTo `json:"billTo"`
}

type shipmentRouteStop struct {
	ID            int    `json:"id"`
	Deleted       bool   `json:"deleted"`
	Phone         string `json:"phone"`
	Sequence      int    `json:"sequence"`
	AppointmentNo string `json:"appointmentNo"`
	Notes         string `json:"notes"`
	StopType      struct {
		Value string `json:"value"`
	} `json:"stopType"`
	Location struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"location"`
	Address struct {
		Line1   string `json:"line1"`
		City    string `json:"city"`
		State   string `json:"state"`
		Zip     string `json:"zip"`
		Country string `json:"country"`
	} `json:"address"`
	Appointment struct {
		Start    string `json:"start"`
		TimeZone string `json:"timeZone"`
	} `json:"appointment"`
}

type shipmentGlobalStop struct {
	ID       int  `json:"id"`
	Deleted  bool `json:"deleted"`
	StopType struct {
		Value string `json:"value"`
	} `json:"stopType"`
	Location struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"location"`
	Address struct {
		Line1       string `json:"line1"`
		City        string `json:"city"`
		State       string `json:"state"`
		CountryCode string `json:"countryCode"`
	} `json:"address"`
	Timezone    string `json:"timezone"`
	Appointment struct {
		Date     string `json:"date"`
		TimeZone string `json:"timeZone"`
	} `json:"appointment"`
	AppointmentNo string `json:"appointmentNo"`
	Notes         string `json:"notes"`
}

type shipmentCustomerOrder struct {
	ID           int                     `json:"id"`
	Deleted      bool                    `json:"deleted"`
	TotalMiles   float64                 `json:"totalMiles"`
	Customer     shipmentSummaryCustomer `json:"customer"`
	Items        []shipmentItem          `json:"items"`
	Route        []shipmentRouteStop     `json:"route"`
	ExternalIDs  []shipmentExternalID    `json:"externalIds"`
	FreightTerms shipmentFreightTerms    `json:"freightTerms"`
}

type shipmentDetail struct {
	ID               int                     `json:"id"`
	CustomID         string                  `json:"customId"`
	NetCustomerCosts float64                 `json:"netCustomerCosts"`
	NetCarrierCosts  float64                 `json:"netCarrierCosts"`
	NetRevenue       float64                 `json:"netRevenue"`
	Status           shipmentStatus          `json:"status"`
	Equipment        []shipmentEquipment     `json:"equipment"`
	Contributors     []shipmentContributor   `json:"contributors"`
	CustomerOrder    []shipmentCustomerOrder `json:"customerOrder"`
	GlobalRoute      []shipmentGlobalStop    `json:"globalRoute"`
}

type shipmentDetailResponse struct {
	Status  string         `json:"Status"`
	Details shipmentDetail `json:"details"`
}

type flexibleInt int

func (v *flexibleInt) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		*v = 0
		return nil
	}

	if strings.HasPrefix(trimmed, `"`) {
		var asString string
		if err := json.Unmarshal(data, &asString); err != nil {
			return err
		}
		trimmed = strings.TrimSpace(asString)
		if trimmed == "" {
			*v = 0
			return nil
		}
	}

	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return fmt.Errorf("parse flexibleInt %q: %w", trimmed, err)
	}
	if math.Mod(parsed, 1) != 0 {
		return fmt.Errorf("parse flexibleInt %q: value must be a whole number", trimmed)
	}

	*v = flexibleInt(int(parsed))
	return nil
}

func (c *Client) ListLoads(ctx context.Context, params load.ListParams) (load.ListResponse, error) {
	page := params.Page
	if page < 1 {
		page = 1
	}

	limit := params.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	start := (page - 1) * limit
	response, err := c.listShipmentSummaryPage(ctx, start, limit, params)
	if err != nil {
		return load.ListResponse{}, err
	}

	details, err := c.fetchShipmentDetails(ctx, response.Details.Shipments)
	if err != nil {
		return load.ListResponse{}, err
	}

	loads := make([]load.Load, 0, len(details))
	for _, detail := range details {
		loads = append(loads, mapShipmentToLoad(detail))
	}

	pagination := paginationFromSummaryPage(page, limit, start, len(response.Details.Shipments), response.Details.Pagination.MoreAvailable)

	return load.ListResponse{
		Data:       loads,
		Pagination: pagination,
	}, nil
}

func (c *Client) listAllShipmentSummaries(ctx context.Context, params load.ListParams) ([]shipmentSummary, error) {
	var (
		all       []shipmentSummary
		start     int
		pageCount int
	)

	for {
		page, err := c.listShipmentSummaryPage(ctx, start, 100, params)
		if err != nil {
			return nil, err
		}

		all = append(all, page.Details.Shipments...)
		pageCount++

		if !page.Details.Pagination.MoreAvailable {
			break
		}

		start += len(page.Details.Shipments)
		if len(page.Details.Shipments) == 0 {
			break
		}
		if pageCount > 500 {
			return nil, fmt.Errorf("aborting pagination after %d pages", pageCount)
		}
	}

	return all, nil
}

func (c *Client) listShipmentSummaryPage(ctx context.Context, start int, pageSize int, params load.ListParams) (shipmentListResponse, error) {
	values := url.Values{}
	if start > 0 {
		values.Set("start", strconv.Itoa(start))
	}
	if pageSize > 0 {
		values.Set("pageSize", strconv.Itoa(pageSize))
	}
	if params.Status != "" {
		if statusCode := turvoStatusFilterCode(params.Status); statusCode != "" {
			values.Set("status[eq]", statusCode)
		}
	}
	if strings.TrimSpace(params.CustomerID) != "" {
		values.Set("customerId[eq]", strings.TrimSpace(params.CustomerID))
	}
	if params.PickupDateSearchFrom != "" {
		values.Set("pickupDate[gte]", formatDateFilter(params.PickupDateSearchFrom, false))
	}
	if params.PickupDateSearchTo != "" {
		values.Set("pickupDate[lte]", formatDateFilter(params.PickupDateSearchTo, true))
	}

	path := "/shipments/list"
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var response shipmentListResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return shipmentListResponse{}, err
	}

	return response, nil
}

func paginationFromSummaryPage(page int, limit int, start int, count int, moreAvailable bool) load.Pagination {
	total := start + count
	if moreAvailable {
		total++
	}

	pages := 0
	if total > 0 {
		pages = (total + limit - 1) / limit
	}
	if moreAvailable && pages < page+1 {
		pages = page + 1
	}

	return load.Pagination{
		Total:       total,
		Pages:       pages,
		Page:        page,
		Limit:       limit,
		Approximate: moreAvailable,
	}
}

func turvoStatusFilterCode(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "2101", "tendered":
		return "2101"
	case "2102", "covered":
		return "2102"
	default:
		return ""
	}
}
func (c *Client) fetchShipmentDetails(ctx context.Context, summaries []shipmentSummary) ([]shipmentDetail, error) {
	if len(summaries) == 0 {
		return nil, nil
	}

	type result struct {
		index   int
		detail  shipmentDetail
		err     error
		skipped bool
	}

	results := make([]shipmentDetail, len(summaries))
	present := make([]bool, len(summaries))
	jobs := make(chan int)
	resultCh := make(chan result, len(summaries))

	workerCount := len(summaries)
	if workerCount > 8 {
		workerCount = 8
	}

	var wg sync.WaitGroup
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				detail, skipped, err := c.getShipmentDetailWithRetry(ctx, summaries[idx].ID)
				resultCh <- result{index: idx, detail: detail, err: err, skipped: skipped}
			}
		}()
	}

	for idx := range summaries {
		jobs <- idx
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for result := range resultCh {
		if result.err != nil {
			return nil, result.err
		}
		if result.skipped {
			continue
		}
		results[result.index] = result.detail
		present[result.index] = true
	}

	filtered := make([]shipmentDetail, 0, len(results))
	for idx, detail := range results {
		if present[idx] {
			filtered = append(filtered, detail)
		}
	}

	return filtered, nil
}

func (c *Client) getShipmentDetail(ctx context.Context, shipmentID int) (shipmentDetail, error) {
	var response shipmentDetailResponse
	if err := c.getJSON(ctx, "/shipments/"+strconv.Itoa(shipmentID), &response); err != nil {
		return shipmentDetail{}, err
	}
	return response.Details, nil
}

func (c *Client) getShipmentDetailWithRetry(ctx context.Context, shipmentID int) (shipmentDetail, bool, error) {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		detail, err := c.getShipmentDetail(ctx, shipmentID)
		if err == nil {
			return detail, false, nil
		}

		var statusErr *apiStatusError
		if !errors.As(err, &statusErr) || !shouldRetryShipmentDetail(statusErr.Status) {
			return shipmentDetail{}, false, err
		}

		log.Printf("shipment detail retry shipment_id=%d attempt=%d status=%d", shipmentID, attempt, statusErr.Status)
		if attempt == maxAttempts {
			log.Printf("shipment detail skipped shipment_id=%d retries=%d status=%d", shipmentID, attempt, statusErr.Status)
			return shipmentDetail{}, true, nil
		}

		if err := sleepWithContext(ctx, time.Duration(attempt)*250*time.Millisecond); err != nil {
			return shipmentDetail{}, false, err
		}
	}

	return shipmentDetail{}, true, nil
}

func shouldRetryShipmentDetail(status int) bool {
	return status == http.StatusBadGateway || status == http.StatusServiceUnavailable || status == http.StatusGatewayTimeout
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func filterSummaries(summaries []shipmentSummary, params load.ListParams) []shipmentSummary {
	filtered := make([]shipmentSummary, 0, len(summaries))

	for _, summary := range summaries {
		if params.Status != "" && !strings.EqualFold(strings.TrimSpace(summary.Status.Code.Value), strings.TrimSpace(params.Status)) {
			continue
		}
		if params.CustomerID != "" && !summaryHasCustomer(summary, params.CustomerID) {
			continue
		}
		filtered = append(filtered, summary)
	}

	return filtered
}

func summaryHasCustomer(summary shipmentSummary, customerID string) bool {
	for _, customerOrder := range summary.CustomerOrder {
		if customerOrder.Deleted {
			continue
		}
		if strconv.Itoa(customerOrder.Customer.ID) == customerID {
			return true
		}
	}
	return false
}

func mapShipmentToLoad(detail shipmentDetail) load.Load {
	customerOrder := firstActiveCustomerOrder(detail)
	pickupStop := firstPickupStop(customerOrder.Route, detail.GlobalRoute)
	deliveryStop := firstDeliveryStop(customerOrder.Route, detail.GlobalRoute)
	poNumbers := collectExternalIDs(customerOrder.ExternalIDs, "1400", "Purchase order #")
	refNumbers := collectExternalIDs(customerOrder.ExternalIDs, "1401", "Reference #")
	totalWeight, billableWeight, numCommodities, inPalletCount, outPalletCount := summarizeItems(customerOrder.Items)
	customerRate := detail.NetCustomerCosts
	carrierRate := detail.NetCarrierCosts
	netProfit := detail.NetRevenue
	profitPercent := 0.0
	if customerRate != 0 {
		profitPercent = (netProfit / customerRate) * 100
	}

	return load.Load{
		ExternalTMSLoadID: strconv.Itoa(detail.ID),
		FreightLoadID:     detail.CustomID,
		Status:            detail.Status.Code.Value,
		Customer: load.Customer{
			ExternalTMSID: strconv.Itoa(customerOrder.Customer.ID),
			Name:          customerOrder.Customer.Name,
			RefNumber:     strings.Join(refNumbers, ", "),
		},
		BillTo: load.BillTo{
			ExternalTMSID: customerOrder.FreightTerms.BillTo.ID,
			Name:          customerOrder.FreightTerms.BillTo.BillTo,
			AddressLine1:  customerOrder.FreightTerms.BillTo.Address.Line1,
			City:          customerOrder.FreightTerms.BillTo.Address.City.Name,
			State:         customerOrder.FreightTerms.BillTo.Address.State.Name,
			Zipcode:       customerOrder.FreightTerms.BillTo.Address.Zip,
			Country:       customerOrder.FreightTerms.BillTo.Address.Country.Name,
			Email:         firstEmail(customerOrder.FreightTerms.BillTo.Emails),
		},
		Pickup: load.Pickup{
			ExternalTMSID: pickupStop.externalID,
			Name:          pickupStop.name,
			AddressLine1:  pickupStop.addressLine1,
			City:          pickupStop.city,
			State:         pickupStop.state,
			Zipcode:       pickupStop.zipcode,
			Country:       pickupStop.country,
			Phone:         pickupStop.phone,
			RefNumber:     strings.Join(refNumbers, ", "),
			ReadyTime:     pickupStop.appointmentTime,
			ApptTime:      pickupStop.appointmentTime,
			ApptNote:      pickupStop.note,
			Timezone:      pickupStop.timezone,
			WarehouseID:   pickupStop.externalID,
		},
		Consignee: load.Consignee{
			ExternalTMSID: deliveryStop.externalID,
			Name:          deliveryStop.name,
			AddressLine1:  deliveryStop.addressLine1,
			City:          deliveryStop.city,
			State:         deliveryStop.state,
			Zipcode:       deliveryStop.zipcode,
			Country:       deliveryStop.country,
			Phone:         deliveryStop.phone,
			RefNumber:     strings.Join(refNumbers, ", "),
			MustDeliver:   deliveryStop.appointmentTime,
			ApptTime:      deliveryStop.appointmentTime,
			ApptNote:      deliveryStop.note,
			Timezone:      deliveryStop.timezone,
			WarehouseID:   deliveryStop.externalID,
		},
		Carrier: load.Carrier{},
		RateData: load.RateData{
			CustomerLHRateUSD: customerRate,
			CarrierLHRateUSD:  carrierRate,
			NetProfitUSD:      netProfit,
			ProfitPercent:     profitPercent,
		},
		Specifications: mapSpecifications(detail.Equipment),
		InPalletCount:  inPalletCount,
		OutPalletCount: outPalletCount,
		NumCommodities: numCommodities,
		TotalWeight:    totalWeight,
		BillableWeight: billableWeight,
		PONums:         strings.Join(poNumbers, ", "),
		Operator:       firstOperator(detail.Contributors),
		RouteMiles:     customerOrder.TotalMiles,
	}
}

type mappedStop struct {
	externalID      string
	name            string
	addressLine1    string
	city            string
	state           string
	zipcode         string
	country         string
	phone           string
	appointmentTime string
	timezone        string
	note            string
}

func firstActiveCustomerOrder(detail shipmentDetail) shipmentCustomerOrder {
	for _, customerOrder := range detail.CustomerOrder {
		if !customerOrder.Deleted {
			return customerOrder
		}
	}
	return shipmentCustomerOrder{}
}

func firstPickupStop(route []shipmentRouteStop, globalRoute []shipmentGlobalStop) mappedStop {
	for _, stop := range route {
		if stop.Deleted || !strings.EqualFold(stop.StopType.Value, "Pickup") {
			continue
		}
		return mappedStop{
			externalID:      strconv.Itoa(stop.Location.ID),
			name:            stop.Location.Name,
			addressLine1:    stop.Address.Line1,
			city:            stop.Address.City,
			state:           stop.Address.State,
			zipcode:         stop.Address.Zip,
			country:         stop.Address.Country,
			phone:           stop.Phone,
			appointmentTime: stop.Appointment.Start,
			timezone:        stop.Appointment.TimeZone,
			note:            stop.Notes,
		}
	}

	for _, stop := range globalRoute {
		if stop.Deleted || !strings.EqualFold(stop.StopType.Value, "Pickup") {
			continue
		}
		return mappedStop{
			externalID:      strconv.Itoa(stop.Location.ID),
			name:            stop.Location.Name,
			addressLine1:    stop.Address.Line1,
			city:            stop.Address.City,
			state:           stop.Address.State,
			country:         stop.Address.CountryCode,
			appointmentTime: stop.Appointment.Date,
			timezone:        firstNonEmpty(stop.Appointment.TimeZone, stop.Timezone),
			note:            stop.Notes,
		}
	}

	return mappedStop{}
}

func firstDeliveryStop(route []shipmentRouteStop, globalRoute []shipmentGlobalStop) mappedStop {
	for _, stop := range route {
		if stop.Deleted || !strings.EqualFold(stop.StopType.Value, "Delivery") {
			continue
		}
		return mappedStop{
			externalID:      strconv.Itoa(stop.Location.ID),
			name:            stop.Location.Name,
			addressLine1:    stop.Address.Line1,
			city:            stop.Address.City,
			state:           stop.Address.State,
			zipcode:         stop.Address.Zip,
			country:         stop.Address.Country,
			phone:           stop.Phone,
			appointmentTime: stop.Appointment.Start,
			timezone:        stop.Appointment.TimeZone,
			note:            stop.Notes,
		}
	}

	for _, stop := range globalRoute {
		if stop.Deleted || !strings.EqualFold(stop.StopType.Value, "Delivery") {
			continue
		}
		return mappedStop{
			externalID:      strconv.Itoa(stop.Location.ID),
			name:            stop.Location.Name,
			addressLine1:    stop.Address.Line1,
			city:            stop.Address.City,
			state:           stop.Address.State,
			country:         stop.Address.CountryCode,
			appointmentTime: stop.Appointment.Date,
			timezone:        firstNonEmpty(stop.Appointment.TimeZone, stop.Timezone),
			note:            stop.Notes,
		}
	}

	return mappedStop{}
}

func collectExternalIDs(externalIDs []shipmentExternalID, key, value string) []string {
	collected := make([]string, 0)
	seen := map[string]struct{}{}

	for _, externalID := range externalIDs {
		if externalID.Deleted {
			continue
		}
		if externalID.Type.Key != key && !strings.EqualFold(externalID.Type.Value, value) {
			continue
		}
		trimmed := strings.TrimSpace(externalID.Value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		collected = append(collected, trimmed)
	}

	return collected
}

func summarizeItems(items []shipmentItem) (totalWeight float64, billableWeight float64, numCommodities int, inPalletCount int, outPalletCount int) {
	for _, item := range items {
		if item.Deleted {
			continue
		}
		numCommodities++
		weight := firstPositive(item.GrossWeight.Weight, item.NetWeight.Weight, item.Weight)
		totalWeight += weight
		billableWeight += weight

		unitName := strings.ToLower(firstNonEmpty(item.HandlingUnit.Value, item.Unit.Value))
		if strings.Contains(unitName, "pallet") {
			inPalletCount += int(item.HandlingQty)
			outPalletCount += int(item.HandlingQty)
		}
	}

	return totalWeight, billableWeight, numCommodities, inPalletCount, outPalletCount
}

func mapSpecifications(equipment []shipmentEquipment) load.Specifications {
	for _, item := range equipment {
		if item.Deleted {
			continue
		}
		temp := item.Temp
		if strings.Contains(strings.ToUpper(item.TempUnits.Value), "C") {
			temp = (temp * 9 / 5) + 32
		}
		return load.Specifications{
			MinTempFahrenheit: temp,
			MaxTempFahrenheit: temp,
		}
	}
	return load.Specifications{}
}

func firstOperator(contributors []shipmentContributor) string {
	for _, contributor := range contributors {
		if contributor.Deleted {
			continue
		}
		return contributor.ContributorUser.Name
	}
	return ""
}

func firstEmail(emails []shipmentEmail) string {
	for _, email := range emails {
		if strings.TrimSpace(email.Email) != "" {
			return email.Email
		}
	}
	return ""
}

func formatDateFilter(value string, endOfDay bool) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.Contains(trimmed, "T") {
		return trimmed
	}
	if endOfDay {
		return trimmed + "T23:59:59Z"
	}
	return trimmed + "T00:00:00Z"
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func firstPositive(values ...float64) float64 {
	for _, value := range values {
		if value > 0 {
			return value
		}
	}
	return 0
}
