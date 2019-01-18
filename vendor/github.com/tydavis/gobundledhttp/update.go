// Updates certificates by downloading from https://curl.haxx.se/ca/cacert.pem
// Requires a machine with exisitng certificate pool in order to run

// +build ignore

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

var cacertsTemplate = template.Must(template.New("").Parse(`package certificates

// PemCerts is the default pool of certificates from the Mozilla certificate
// collection. This collection is also used by cURL as indicated here:
// https://curl.haxx.se/docs/caextract.html
var PemCerts = []byte({{ .  }})
`))

func main() {
	urlString := "https://curl.haxx.se/ca/cacert.pem"
	// Using the basic http client because we can't load certs ourselves,
	// wget in busybox can't contact https sites either.
	downloadresp, _ := http.Get(urlString)
	certFile, err := ioutil.ReadAll(downloadresp.Body)
	defer downloadresp.Body.Close()
	if err != nil {
		log.Fatalf("Failed to read response into byte array: %v", err)
	}

	// regex matching is nice
	reg := regexp.MustCompile(`-----BEGIN CERTIFICATE-----[\n|\S]+-----END CERTIFICATE-----`)
	matches := reg.FindAll(certFile, -1) // Find all matches

	// Byte slice to construct the go file must be zero cap so there are no
	// leading null bytes in resulting file
	certsgo := make([]byte, 0)
	// Add backtick (`) because byte array won't process correctly without it
	certsgo = append(certsgo, byte('`'))
	for _, b := range matches {
		x := append([]byte{'\n'}, b...)
		certsgo = append(certsgo, x...)
	}
	certsgo = append(certsgo, byte('`'))

	// Write out the resulting go file
	fh, _ := os.Create("certificates/certificates.go")
	defer fh.Close()
	cacertsTemplate.Execute(fh, string(certsgo))
}
