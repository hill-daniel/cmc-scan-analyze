package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"html/template"
	"io"
	"log"
	"net/mail"
	"text/tabwriter"
)

const cryptoTableTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cryptocurrency Table</title>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
        }

        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }

        th {
            background-color: #f2f2f2;
        }

        .positive {
            color: green;
            font-weight: bold;
        }

        .negative {
            color: red;
            font-weight: bold;
        }
    </style>
</head>
<body>
<h1>Cryptocurrency Data</h1>
<table>
    <thead>
    <tr>
        <th>CMC-Rank</th>
        <th>Name</th>
        <th>Symbol</th>
        <th>Change 24h</th>
        <th>Change 7d</th>
        <th>Change 30d</th>
        <th>Market Cap</th>
        <th>Price</th>
    </tr>
    </thead>
    <tbody>
    {{ range . }}
    <tr>
        <td>
            {{ if gt .RankChange 0 }}
            <span class="positive">▲{{ .RankChange }}</span>
            {{ else if lt .RankChange 0 }}
            <span class="negative">▼{{ .RankChange }}</span>
            {{ else }}
            {{ .RankChange }}
            {{ end }}
        </td>
        <td>{{ .Name }}</td>
        <td>{{ .Symbol }}</td>
        <td>
            {{ if gt .PercentChange24H 0.0 }}
            <span class="positive">{{ printf "%.2f%%" .PercentChange24H }}</span>
            {{ else if lt .PercentChange24H 0.0 }}
            <span class="negative">{{ printf "%.2f%%" .PercentChange24H }}</span>
            {{ else }}
            {{ printf "%.2f%%" .PercentChange24H }}
            {{ end }}
        </td>
        <td>
            {{ if gt .PercentChange7D 0.0 }}
            <span class="positive">{{ printf "%.2f%%" .PercentChange7D }}</span>
            {{ else if lt .PercentChange7D 0.0 }}
            <span class="negative">{{ printf "%.2f%%" .PercentChange7D }}</span>
            {{ else }}
            {{ printf "%.2f%%" .PercentChange7D }}
            {{ end }}
        </td>
        <td>
            {{ if gt .PercentChange30D 0.0 }}
            <span class="positive">{{ printf "%.2f%%" .PercentChange30D }}</span>
            {{ else if lt .PercentChange30D 0.0 }}
            <span class="negative">{{ printf "%.2f%%" .PercentChange30D }}</span>
            {{ else }}
            {{ printf "%.2f%%" .PercentChange30D }}
            {{ end }}
        </td>
        <td>{{ formatCurrency .MarketCap }}</td>
        <td>{{ formatTokenValue .Price }}</td>
    </tr>
    {{ end }}
    </tbody>
</table>
</body>
</html>
`

// EmailInput defines the expected input for the Lambda function
type EmailInput struct {
	Sender     string   `json:"sender"`
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	Body       string   `json:"body"`
	Html       string   `json:"html"`
}

func createHtml(changes []AssetQuote) (string, error) {
	p := message.NewPrinter(language.English)

	funcMap := template.FuncMap{
		"formatCurrency": func(value float64) string {
			return p.Sprintf("$%.2f", value)
		},
		"formatTokenValue": func(value float64) string {
			return p.Sprintf("$%.8f", value)
		},
	}
	tmpl, err := template.New("cryptoTable").Funcs(funcMap).Parse(cryptoTableTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, changes)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	// Output the generated HTML
	return buf.String(), nil
}

func writeChanges(output io.Writer, quotes []AssetQuote) {
	w := tabwriter.NewWriter(output, 1, 1, 1, ' ', 0)
	_, err := fmt.Fprintln(w, "#\t CMC-Rank\t Name\t Symbol\t Change 24h\t Change 7d\t Change 30d\t Market Cap\t Price")
	if err != nil {
		log.Fatalf("failed to format header: %v", err)
	}
	p := message.NewPrinter(language.English)

	for i, quote := range quotes {
		_, err = fmt.Fprintf(w, "#%d\t %d\t %s\t %s\t %.2f%%\t %.2f%%\t  %.2f%%\t %s\t %s \n", i+1, quote.RankChange, quote.Name, quote.Symbol, quote.PercentChange24H,
			quote.PercentChange7D, quote.PercentChange30D, p.Sprintf("%.2f", quote.MarketCap), p.Sprintf("%.8f", quote.Price))
		if err != nil {
			log.Fatalf("failed to format rankChange: %v", err)
		}
	}
	err = w.Flush()
	if err != nil {
		log.Fatalf("failed to flush writer: %v", err)
	}
}

func sendEmail(input EmailInput) error {
	// Validate email addresses
	if _, err := mail.ParseAddress(input.Sender); err != nil {
		return fmt.Errorf("invalid sender email: %v", err)
	}

	recipients := make([]*string, len(input.Recipients))
	for i, recipient := range input.Recipients {
		if _, err := mail.ParseAddress(recipient); err != nil {
			return fmt.Errorf("invalid recipient email: %v", err)
		}
		recipients[i] = &recipient
	}

	// Start AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create SES service client
	svc := ses.New(sess)

	// Compose the email
	inputMessage := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: recipients,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(input.Body),
				},
				Html: &ses.Content{
					Data: aws.String(input.Html)},
			},
			Subject: &ses.Content{
				Data: aws.String(input.Subject),
			},
		},
		Source: aws.String(input.Sender),
	}

	_, err = svc.SendEmail(inputMessage)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
