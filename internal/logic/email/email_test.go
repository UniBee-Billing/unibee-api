package email

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"unibee/api/bean"
	"unibee/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	_ "github.com/go-sql-driver/mysql"
)

func TestTemplateVariableReplacement(t *testing.T) {
	ctx := context.Background()
	now := gtime.New(gtime.Now())
	periodEnd := gtime.New(gtime.Now().AddDate(0, 1, 0))
	templateVariables := &bean.EmailTemplateVariable{
		InvoiceId:             "INV-2024-001",
		UserName:              "John Doe",
		MerchantProductName:   "Premium Subscription",
		MerchantCustomerEmail: "support@unibee.dev",
		MerchantName:          "Example Company",
		DateNow:               now,
		PeriodEnd:             periodEnd,
		PaymentAmount:         "99.99",
		RefundAmount:          "49.99",
		Currency:              "USD",
		TokenExpireMinute:     "30",
		CodeExpireMinute:      "15",
		Code:                  "123456",
		Link:                  "https://unibee.dev",
		HttpLink:              "https://unibee.dev",
		AccountHolder:         "UniBee Company Ltd",
		Address:               "123 Business St, City, Country",
		BIC:                   "EXAMPLBIC",
		IBAN:                  "GB29NWBK60161331926819",
		BankData:              "Bank of Example, Account: 12345678",
	}
	variableMap, _ := utility.ReflectTemplateStructToMap(templateVariables, "")
	var subject = "来自 Sliver plan 的银行转账信息"
	var content = "<p>您好，{User name}！</p><p>感谢您选择 {Merchant Product Name} 的套餐。您可以在账单后台的“发票（Invoices）”页面中查看您的账单。</p><p> 以下是我们用于银行转账的账户信息：</p><p>账户名：{Account Holder}</p><p> BIC（银行识别码）：<strong>TRWIBEB1XXX</strong></p><p> IBAN（国际银行账号）：<strong>BE57 9670 1926 1435</strong></p><p> Wise 收款地址：<strong>Avenue Louise 54, Room S52 Brussels 1050</strong></p><p>如有任何疑问，欢迎随时通过邮箱联系我们：{Merchant's customer support email address}。</p><p>请注意：此邮件为系统自动发送，请勿直接回复。如需帮助，请使用上方提供的联系方式。</p><p>{Merchant Name} 团队敬上</p>"
	utility.Assert(variableMap != nil, "template parse error")
	variableMap["CompanyName"] = "Test CompanyName"
	g.Log().Infof(ctx, "template variables:%v", utility.MarshalToJsonString(variableMap))
	for key, value := range variableMap {
		mapKey := "{{" + key + "}}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(subject) > 0 {
			subject = strings.Replace(subject, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, mapKey, htmlValue, -1)
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
	}
	for key, value := range variableMap {
		mapKey := "{" + key + "}"
		htmlKey := strings.Replace(mapKey, " ", "&nbsp;", 10)
		htmlValue := "<strong>" + value.(string) + "</strong>"
		if len(subject) > 0 {
			subject = strings.Replace(subject, mapKey, value.(string), -1)
		}
		if len(content) > 0 {
			content = strings.Replace(content, mapKey, htmlValue, -1)
			content = strings.Replace(content, htmlKey, htmlValue, -1)
		}
	}
	g.Log().Infof(ctx, "subject:%v,\ncontent:%v", subject, content)
}

func TestUpdateEmailTemplateVariables(t *testing.T) {
	// Database connection configuration - please modify according to your actual environment
	dbConfig := struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}{
		Host:     "localhost",
		Port:     3306,
		User:     "unibee",
		Password: "changeme",
		Database: "unibee",
	}

	// Create independent database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	// Use native sql connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Variable mapping relationships
	variableMappings := map[string]string{
		"{User name}":             "{UserName}",
		"{Merchant Product Name}": "{ProductName}",
		"{Merchant’s customer support email address}": "{SupportEmail}",
		"{Merchant's customer support email address}": "{SupportEmail}",
		"{Merchant Name}":  "{MerchantName}",
		"{Payment Amount}": "{PaymentAmount}",
		"{Refund Amount}":  "{RefundAmount}",
		"{Account Holder}": "{WireTransferAccountHolder}",
		"{Address}":        "{WireTransferAddress}",
		"{BIC}":            "{WireTransferBIC}",
		"{IBAN}":           "{WireTransferIBAN}",
		"{Bank Data}":      "{WireTransferBankData}",
		// HTML format variables (spaces replaced with &nbsp;)
		"{User&nbsp;name}":                  "{UserName}",
		"{Merchant&nbsp;Product&nbsp;Name}": "{ProductName}",
		"{Merchant’s&nbsp;customer&nbsp;support&nbsp;email&nbsp;address}": "{SupportEmail}",
		"{Merchant's&nbsp;customer&nbsp;support&nbsp;email&nbsp;address}": "{SupportEmail}",
		"{Merchant&nbsp;Name}":  "{MerchantName}",
		"{Payment&nbsp;Amount}": "{PaymentAmount}",
		"{Refund&nbsp;Amount}":  "{RefundAmount}",
		"{Account&nbsp;Holder}": "{WireTransferAccountHolder}",
		"{Bank&nbsp;Data}":      "{WireTransferBankData}",
	}

	// Tables to update
	tables := []string{"email_default_template", "merchant_email_template"}

	for _, tableName := range tables {
		t.Logf("Processing table: %s", tableName)

		// Query all records that need to be updated
		rows, err := db.Query(fmt.Sprintf("SELECT id, template_title, template_content FROM %s WHERE is_deleted = 0", tableName))
		if err != nil {
			t.Logf("Failed to query records from %s: %v", tableName, err)
			continue
		}
		defer rows.Close()

		updatedCount := 0
		for rows.Next() {
			var id int64
			var title, content string

			if err := rows.Scan(&id, &title, &content); err != nil {
				t.Logf("Failed to scan row from %s: %v", tableName, err)
				continue
			}

			// Update title and content
			newTitle := title
			newContent := content

			for oldVar, newVar := range variableMappings {
				newTitle = strings.ReplaceAll(newTitle, oldVar, newVar)
				newContent = strings.ReplaceAll(newContent, oldVar, newVar)
			}

			// Update database if content has changed
			if newTitle != title || newContent != content {
				_, err := db.Exec(fmt.Sprintf("UPDATE %s SET template_title = ?, template_content = ? WHERE id = ?", tableName),
					newTitle, newContent, id)
				if err != nil {
					t.Logf("Failed to update record %d in %s: %v", id, tableName, err)
					continue
				}
				updatedCount++
				t.Logf("Updated record %d in %s", id, tableName)
			}
		}

		if err := rows.Err(); err != nil {
			t.Logf("Error iterating rows in %s: %v", tableName, err)
		}

		t.Logf("Update completed for %s. Total records updated: %d", tableName, updatedCount)
	}

	t.Logf("All tables processed successfully")
}
