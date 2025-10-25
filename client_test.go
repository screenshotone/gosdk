package gosdk_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	screenshots "github.com/screenshotone/gosdk"
)

func ExampleClient_GenerateTakeURL() {
	client, err := screenshots.NewClient("IVmt2ghj9TG_jQ", "Sxt94yAj9aQSgg")
	if err != nil {
		// ...
	}

	options := screenshots.NewTakeOptions("https://scalabledeveloper.com").
		Format("png").
		FullPage(true).
		DeviceScaleFactor(2).
		BlockAds(true).
		BlockTrackers(true)

	u, err := client.GenerateTakeURL(options)
	if err != nil {
		// ...
	}

	fmt.Println(u.String())
	// Output: https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&block_ads=true&block_trackers=true&device_scale_factor=2&format=png&full_page=true&url=https%3A%2F%2Fscalabledeveloper.com&signature=85aabf7ac251563ec6158ef6839dd019bb79ce222cc85288a2e8cea0291a824e
}

func TestTakeURLGeneratesURL(t *testing.T) {
	testsCases := []struct {
		options     *screenshots.TakeOptions
		expectedURL string
	}{
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").FullPage(true).DeviceScaleFactor(1).ViewportHeight(1200).ViewportWidth(1200).Format("png").OmitBackground(true),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&device_scale_factor=1&format=png&full_page=true&omit_background=true&url=https%3A%2F%2Fscalabledeveloper.com&viewport_height=1200&viewport_width=1200&signature=3c0c5543599067322e8c84470702330e3687c6a08eef6b7311b71c32d04e1bd5",
		},
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").Format("jpg").ImageQuality(90),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&format=jpg&image_quality=90&url=https%3A%2F%2Fscalabledeveloper.com&signature=2e64625071a9686277aa2a3bdcc5b5ccb8d87a56ba15d5ab1a669bf878d3ef49",
		},
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").GeolocationLongitude(50.1234556).GeolocationLatitude(99.98765).GeolocationAccuracy(40),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&geolocation_accuracy=40&geolocation_latitude=99.98765&geolocation_longitude=50.1234556&url=https%3A%2F%2Fscalabledeveloper.com&signature=aea25992c369f3682084865a9e89b7a0b4f0e6618eb2657ac7c1b0e97b25323d",
		},
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").UserAgent("test").Authorization("auth").Headers("X-Header-1: val1", "X-Header-2: val2").Cookies("key=value", "key1=value1"),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&authorization=auth&cookies=key%3Dvalue&cookies=key1%3Dvalue1&headers=X-Header-1%3A+val1&headers=X-Header-2%3A+val2&url=https%3A%2F%2Fscalabledeveloper.com&user_agent=test&signature=9cd886f11f29473b92f47bcf09aec51963874ba4c49d5a0c2f38fd45d85f2d70",
		},
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").Delay(25).Timeout(15),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&delay=25&timeout=15&url=https%3A%2F%2Fscalabledeveloper.com&signature=947cc029d973160799896e59de3c0ea5d4eaa61b8fb233ededac87cd335f07ae",
		},
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").BlockAds(true).BlockTrackers(true).BlockRequests("*example*").BlockResources("image"),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&block_ads=true&block_requests=%2Aexample%2A&block_resources=image&block_trackers=true&url=https%3A%2F%2Fscalabledeveloper.com&signature=ea35f9f4b4bc68a7a779b9cb9ac532a50eaac9819b0267bdb4e328e103f1710a",
		},
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").TimeZone("Europe/Berlin"),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&time_zone=Europe%2FBerlin&url=https%3A%2F%2Fscalabledeveloper.com&signature=391e9aefa109a3f1af355e6bcb9de7714cdc1cb440bbeb2e5dac5a601d7cd58c",
		},
		{
			screenshots.NewTakeOptions("https://scalabledeveloper.com").Cache(true).CacheKey("test").CacheTTL(10000),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&cache=true&cache_key=test&cache_ttl=10000&url=https%3A%2F%2Fscalabledeveloper.com&signature=49dbe28dbc268f30528c72359bc2cfa41cd292d766e97043fe8fa3475e6a27b1",
		},
		{
			screenshots.NewTakeWithHTML("<h1>Hello, world!</h1>"),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&html=%3Ch1%3EHello%2C+world%21%3C%2Fh1%3E&signature=2e9559eaeb5ff8a6b0aa85ddeaaf2e65d8e7cf636741964488784864327e3901",
		},
		{
			screenshots.NewTakeOptions("https://example.com").
				PDFPrintBackground(true).
				PDFFitOnePage(true).
				PDFLandscape(true).
				PDFPaperFormat("a4"),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&pdf_fit_one_page=true&pdf_landscape=true&pdf_paper_format=a4&pdf_print_background=true&url=https%3A%2F%2Fexample.com&signature=2be4758936d2392a1776fdac4e0bb6dc4ad3aab1a380bfdb2c34eae6899ca556",
		},
		{
			screenshots.NewTakeOptions("https://example.com").
				ClipX(100).
				ClipY(200).
				ClipWidth(300).
				ClipHeight(400),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&clip_height=400&clip_width=300&clip_x=100&clip_y=200&url=https%3A%2F%2Fexample.com&signature=8c6815f9c6123177a65826b51f14aa774bd54e075a3e04156a0bee2e9b2650c7",
		},
		{
			screenshots.NewTakeOptions("https://example.com").
				FullPageAlgorithm("by_sections").
				SelectorScrollIntoView(true).
				IgnoreHostErrors(true),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&full_page_algorithm=by_sections&ignore_host_errors=true&selector_scroll_into_view=true&url=https%3A%2F%2Fexample.com&signature=8c7dd91d28a8289d75affde1aff70e6a73afa594a0557b55e1176fcabc321e26",
		},
		{
			screenshots.NewTakeOptions("https://example.com").
				StorageEndpoint("https://storage.example.com").
				StorageAccessKeyID("access123").
				StorageSecretAccessKey("secret456").
				StorageBucket("mybucket").
				StorageClass("standard"),
			"https://api.screenshotone.com/take?access_key=IVmt2ghj9TG_jQ&storage_access_key_id=access123&storage_bucket=mybucket&storage_class=standard&storage_endpoint=https%3A%2F%2Fstorage.example.com&storage_secret_access_key=secret456&url=https%3A%2F%2Fexample.com&signature=0b27223cf5ec9f47f43d902603e9b9578ca850fca9d5638e944eaef67c51d9d2",
		},
	}

	client, err := screenshots.NewClient("IVmt2ghj9TG_jQ", "Sxt94yAj9aQSgg")
	ok(t, err)
	for _, testCase := range testsCases {
		u, err := client.GenerateTakeURL(testCase.options)
		ok(t, err)

		equals(t, testCase.expectedURL, u.String())
	}
}

func TestGenerateUnsignedTakeURL(t *testing.T) {
	client, err := screenshots.NewClient("test-key", "")
	ok(t, err)

	options := screenshots.NewTakeOptions("https://example.com")
	u, err := client.GenerateUnsignedTakeURL(options)
	ok(t, err)

	expected := "https://api.screenshotone.com/take?access_key=test-key&url=https%3A%2F%2Fexample.com"
	equals(t, expected, u.String())
}

func TestGenerateTakeURLRequiresSecretKey(t *testing.T) {
	client, err := screenshots.NewClient("test-key", "")
	ok(t, err)

	options := screenshots.NewTakeOptions("https://example.com")
	_, err = client.GenerateTakeURL(options)
	errorred(t, err, "secret key is required")
}

func TestTakeAcceptsOKStatusCode(t *testing.T) {
	mockClient := &http.Client{
		Transport: &mockRoundTripper{
			statusCode: http.StatusOK,
			body:       []byte("test image data"),
		},
	}

	client, err := screenshots.NewClientWithHTTPClient("test-key", "test-secret", mockClient)
	ok(t, err)

	options := screenshots.NewTakeOptions("https://example.com")
	image, _, err := client.Take(context.Background(), options)
	ok(t, err)

	equals(t, "test image data", string(image))
}

func TestTakeAcceptsCreatedStatusCode(t *testing.T) {
	mockClient := &http.Client{
		Transport: &mockRoundTripper{
			statusCode: http.StatusCreated,
			body:       []byte(""),
		},
	}

	client, err := screenshots.NewClientWithHTTPClient("test-key", "test-secret", mockClient)
	ok(t, err)

	options := screenshots.NewTakeOptions("https://example.com")
	image, _, err := client.Take(context.Background(), options)
	ok(t, err)

	equals(t, "", string(image))
}

func TestTakeRejectsOtherStatusCodes(t *testing.T) {
	mockClient := &http.Client{
		Transport: &mockRoundTripper{
			statusCode: http.StatusBadRequest,
			body:       []byte("bad request"),
		},
	}

	client, err := screenshots.NewClientWithHTTPClient("test-key", "test-secret", mockClient)
	ok(t, err)

	options := screenshots.NewTakeOptions("https://example.com")
	_, _, err = client.Take(context.Background(), options)
	errorred(t, err, "the server returned a response: 400 Bad Request")
}

// errorred fails the test if an err is nil or message is not found in the message string.
func errorred(tb testing.TB, err error, message string) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: expected error, but got nil\033[39m\n\n", filepath.Base(file), line)
		tb.FailNow()
		return
	}

	if !strings.Contains(err.Error(), message) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: \"%s\" not found in \"%s\"\033[39m\n\n", filepath.Base(file), line, message, err.Error())
		tb.FailNow()
		return
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type mockRoundTripper struct {
	statusCode int
	body       []byte
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.statusCode,
		Status:     http.StatusText(m.statusCode),
		Body:       io.NopCloser(bytes.NewReader(m.body)),
		Header:     make(http.Header),
	}, nil
}
