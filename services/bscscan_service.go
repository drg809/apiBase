package services

import (
	"bytes"
	"fmt"
	"os"

	"github.com/valyala/fasthttp"
)

func ExampleGetGzippedJsonWithFastHttp() {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	S3_ACCESS_KEY := os.Getenv("BSCSCAN_KEY")
	req.SetRequestURI("https://api.bscscan.com/api?module=stats&action=bnbprice&apikey=" + S3_ACCESS_KEY)
	// fasthttp does not automatically request a gzipped response.
	// We must explicitly ask for it.
	req.Header.Set("Accept-Encoding", "gzip")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Perform the request
	err := fasthttp.Do(req, resp)
	if err != nil {
		fmt.Printf("Client get failed: %s\n", err)
		return
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		fmt.Printf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode())
		return
	}

	// Verify the content type
	contentType := resp.Header.Peek("Content-Type")
	if bytes.Index(contentType, []byte("application/json")) != 0 {
		fmt.Printf("Expected content type application/json but got %s\n", contentType)
		return
	}

	// Do we need to decompress the response?
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		fmt.Println("Unzipping...")
		body, _ = resp.BodyGunzip()
	} else {
		body = resp.Body()
	}

	fmt.Printf("Response body is: %s", body)
}
