package gosdk_test

import (
	"fmt"
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
	}

	client, err := screenshots.NewClient("IVmt2ghj9TG_jQ", "Sxt94yAj9aQSgg")
	ok(t, err)
	for _, testCase := range testsCases {
		u, err := client.GenerateTakeURL(testCase.options)
		ok(t, err)

		equals(t, testCase.expectedURL, u.String())
	}
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
