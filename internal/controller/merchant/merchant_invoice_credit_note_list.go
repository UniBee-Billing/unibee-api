package merchant

import (
	"context"
	"encoding/csv"
	"io"
	"strings"
	"unibee/api/merchant/invoice"
	_interface "unibee/internal/interface/context"
	"unibee/internal/logic/invoice/service"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerInvoice) CreditNoteList(ctx context.Context, req *invoice.CreditNoteListReq) (res *invoice.CreditNoteListRes, err error) {
	emails := make([]string, 0)

	// 1. Process directly passed emails parameter
	if len(req.Emails) > 0 {
		cleanedEmails := strings.ReplaceAll(req.Emails, ";", ",")
		emails = strings.Split(cleanedEmails, ",")
		// Clean each email, remove spaces
		for i, email := range emails {
			emails[i] = strings.TrimSpace(email)
		}
		// Filter empty strings
		var filteredEmails []string
		for _, email := range emails {
			if email != "" {
				filteredEmails = append(filteredEmails, email)
			}
		}
		emails = filteredEmails
	}

	// 2. If CSV file is uploaded, read emails from file
	if req.File != nil && req.File.Size > 0 {
		// Check file size limit (2MB = 2 * 1024 * 1024 bytes)
		const maxFileSize = 2 * 1024 * 1024 // 2MB
		utility.Assert(req.File.Size <= maxFileSize, "CSV file size exceeds 2MB limit")

		// Open uploaded file
		file, fileErr := req.File.Open()
		if fileErr != nil {
			g.Log().Errorf(ctx, "CreditNoteList CSV file open error: %s", fileErr.Error())
		} else {
			defer file.Close()

			// Create CSV reader
			reader := csv.NewReader(file)

			// Skip header row if exists (optional)
			firstRow, firstRowErr := reader.Read()
			if firstRowErr != nil && firstRowErr != io.EOF {
				g.Log().Errorf(ctx, "CreditNoteList CSV first row read error: %s", firstRowErr.Error())
			} else if firstRowErr == nil && len(firstRow) > 0 {
				// Check if first row is header (contains "email" or similar)
				firstField := strings.ToLower(strings.TrimSpace(firstRow[0]))
				if firstField == "email" || firstField == "e-mail" || firstField == "mail" {
					g.Log().Infof(ctx, "CreditNoteList CSV header detected, skipping first row")
				} else {
					// First row is data, process it
					if len(firstRow) > 0 {
						email := strings.TrimSpace(firstRow[0])
						if email != "" && utility.IsEmailValid(email) {
							emails = append(emails, email)
						}
					}
				}
			}

			// Read remaining CSV content
			for {
				record, readErr := reader.Read()
				if readErr == io.EOF {
					break
				}
				if readErr != nil {
					g.Log().Errorf(ctx, "CreditNoteList CSV read error: %s", readErr.Error())
					continue // Skip error lines, continue reading
				}

				// Validate record has at least one column
				if len(record) == 0 {
					g.Log().Infof(ctx, "CreditNoteList CSV empty row detected, skipping")
					continue
				}

				// Get email from first column
				email := strings.TrimSpace(record[0])
				if email != "" && utility.IsEmailValid(email) {
					// Check if already exists, avoid duplicates
					exists := false
					for _, existingEmail := range emails {
						if strings.EqualFold(existingEmail, email) {
							exists = true
							break
						}
					}
					if !exists {
						emails = append(emails, email)
					}
				} else if email != "" {
					g.Log().Infof(ctx, "CreditNoteList CSV invalid email format: %s", email)
				}
			}

			g.Log().Infof(ctx, "CreditNoteList CSV processing completed, found %d valid emails", len(emails))
		}
	}

	// Check total emails count limit
	const maxEmailsCount = 1000
	utility.Assert(len(emails) <= maxEmailsCount, "Total emails filter count exceeds 1000 limit")

	internalResult, err := service.CreditNoteList(ctx, &service.CreditNoteListInternalReq{
		MerchantId:      _interface.GetMerchantId(ctx),
		SearchKey:       req.SearchKey,
		Emails:          emails,
		Status:          req.Status,
		GatewayIds:      req.GatewayIds,
		PlanIds:         req.PlanIds,
		Currency:        req.Currency,
		SortField:       req.SortField,
		SortType:        req.SortType,
		Page:            req.Page,
		Count:           req.Count,
		CreateTimeStart: req.CreateTimeStart,
		CreateTimeEnd:   req.CreateTimeEnd,
		SkipTotal:       false,
	})
	if err != nil {
		return nil, err
	}
	return &invoice.CreditNoteListRes{CreditNotes: internalResult.CreditNotes, Total: internalResult.Total}, nil
}
