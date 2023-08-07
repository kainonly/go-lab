package main

import (
	"encoding/base32"
	"fmt"
	"github.com/dgryski/dgoogauth"
	"net/url"
	"os"
	"rsc.io/qr"
)

const (
	qrFilename = "./qr.png"
)

func main() {
	// Example secret from here:
	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	secret := []byte{'H', 'e', 'l', 'l', 'o', '!', 0xDE, 0xAD, 0xBE, 0xEF}

	// Generate random secret instead of using the test value above.
	// secret := make([]byte, 10)
	// _, err := rand.Read(secret)
	// if err != nil {
	//	panic(err)
	// }

	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	account := "user@example.com"
	issuer := "NameOfMyService"

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		panic(err)
	}

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)

	params := url.Values{}
	params.Add("secret", secretBase32)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()
	fmt.Printf("URL is %s\n", URL.String())

	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(qrFilename, code.PNG(), 0600)
	if err != nil {
		panic(err)
	}

	fmt.Printf("QR code is in %s. Please scan it into Google Authenticator app.\n", qrFilename)

	// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
	// so allocate it once and reuse it for multiple calls.
	// *Caution*: if you have a scale-out service, this code won't detect replays to different
	// server replicas.
	otpc := &dgoogauth.OTPConfig{
		Secret:      secretBase32,
		WindowSize:  1,
		HotpCounter: 0,
		UTC:         true,
	}

	for {
		fmt.Printf("Please enter the token value (or q to quit): ")

		var token string
		fmt.Scanln(&token)
		if token == "q" {
			break
		}

		val, err := otpc.Authenticate(token)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !val {
			fmt.Println("Sorry, Not Authenticated")
			continue
		}

		fmt.Println("Authenticated!")
	}
}
