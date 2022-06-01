package gosdk

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const baseURL = "https://api.screenshotone.com"
const takePath = "/take"

// Client API client for the ScreenshotOne.com API.
type Client struct {
	accessKey, secretKey string

	httpClient *http.Client
}

// NewClient returns new API client for the ScreenshotOne.com API.
func NewClient(accessKey, secretKey string) (*Client, error) {
	client := &Client{accessKey, secretKey, &http.Client{}}

	return client, nil
}

// GenerateTakeURL generates URL for taking screenshots, but does not send any request.
func (client *Client) GenerateTakeURL(options *TakeOptions) (*url.URL, error) {
	// generate query
	query := options.query
	query.Set("access_key", client.accessKey)
	queryString := query.Encode()

	// sign the query string and append the signature
	hash := hmac.New(sha256.New, []byte(client.secretKey))
	_, err := hash.Write([]byte(queryString))
	if err != nil {
		return nil, fmt.Errorf("failed to sign the query string: %w", err)
	}
	signature := hex.EncodeToString(hash.Sum(nil))
	queryString += "&signature=" + signature

	u, err := url.Parse(baseURL + takePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL \"%s\": %w", baseURL+takePath, err)
	}
	u.RawQuery = queryString

	return u, nil
}

// Take takes screenshot and returns image or error if the request failed.
func (client *Client) Take(ctx context.Context, options *TakeOptions) ([]byte, *http.Response, error) {
	u, err := client.GenerateTakeURL(options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate URL: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to instantiate HTTP request: %w", err)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, response, fmt.Errorf("the server returned a response: %d %s", response.StatusCode, response.Status)
	}

	defer response.Body.Close()
	image, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read the image data from HTTP response: %w", err)
	}

	return image, nil, nil
}

// TakeOptions for the ScreenshotOne.com API take method.
type TakeOptions struct {
	query url.Values
}

// Returns options for the ScreenshotOne.com API take method.
func NewTakeOptions(pageURL string) *TakeOptions {
	query := url.Values{}
	query.Add("url", pageURL)

	return &TakeOptions{query: query}
}

// FullPage renders the full page.
func (o *TakeOptions) FullPage(fullPage bool) *TakeOptions {
	o.query.Add("full_page", strconv.FormatBool(fullPage))

	return o
}

// Format sets response format, one of: "png", "jpeg", "webp" or "jpg".
func (o *TakeOptions) Format(format string) *TakeOptions {
	o.query.Add("format", format)

	return o
}

// ImageQuality renders image with the specified quality. Available for the next formats: "jpeg" ("jpg"), "webp".
func (o *TakeOptions) ImageQuality(imageQuality int) *TakeOptions {
	o.query.Add("image_quality", strconv.Itoa(imageQuality))

	return o
}

// OmitBackground renders a transparent background for the image. Works only if the site has not defined background color.
// Available for the following response formats: "png", "webp".
func (o *TakeOptions) OmitBackground(omitBackground bool) *TakeOptions {
	o.query.Add("omit_background", strconv.FormatBool(omitBackground))

	return o
}

// ViewportWidth sets the width of the browser viewport (pixels).
func (o *TakeOptions) ViewportWidth(viewportWidth int) *TakeOptions {
	o.query.Add("viewport_width", strconv.Itoa(viewportWidth))

	return o
}

// ViewportWidth sets the height of the browser viewport (pixels).
func (o *TakeOptions) ViewportHeight(viewportHeight int) *TakeOptions {
	o.query.Add("viewport_height", strconv.Itoa(viewportHeight))

	return o
}

// DeviceScaleFactor sets the device scale factor. Acceptable value is one of: 1, 2 or 3.
func (o *TakeOptions) DeviceScaleFactor(viewportHeight int) *TakeOptions {
	o.query.Add("device_scale_factor", strconv.Itoa(viewportHeight))

	return o
}

// GeolocationLatitude sets geolocation latitude for the request.
// Both latitude and longitude are required if one of them is set.
func (o *TakeOptions) GeolocationLatitude(latitude float64) *TakeOptions {
	o.query.Add("geolocation_latitude", strconv.FormatFloat(latitude, byte('f'), -1, 64))

	return o
}

// GeolocationLatitude sets geolocation longitude for the request.
// Both latitude and longitude are required if one of them is set.
func (o *TakeOptions) GeolocationLongitude(longitude float64) *TakeOptions {
	o.query.Add("geolocation_longitude", strconv.FormatFloat(longitude, byte('f'), -1, 64))

	return o
}

// GeolocationAccuracy sets the geolocation accuracy in meters.
func (o *TakeOptions) GeolocationAccuracy(accuracy int) *TakeOptions {
	o.query.Add("geolocation_accuracy", strconv.Itoa(accuracy))

	return o
}

// BlockAds blocks ads.
func (o *TakeOptions) BlockAds(blockAds bool) *TakeOptions {
	o.query.Add("block_ads", strconv.FormatBool(blockAds))

	return o
}

// BlockTrackers blocks trackers.
func (o *TakeOptions) BlockTrackers(blockTrackers bool) *TakeOptions {
	o.query.Add("block_trackers", strconv.FormatBool(blockTrackers))

	return o
}

// BlockRequests blocks requests by specifying URL, domain, or even a simple pattern.
func (o *TakeOptions) BlockRequests(blockRequests ...string) *TakeOptions {
	for _, blockRequest := range blockRequests {
		o.query.Add("block_requests", blockRequest)
	}

	return o
}

// BlockResources blocks loading resources by type.
// Available resource types are: "document", "stylesheet", "image", "media", "font", "script", "texttrack", "xhr", "fetch", "eventsource", "websocket", "manifest", "other".
func (o *TakeOptions) BlockResources(blockRequests ...string) *TakeOptions {
	for _, blockRequest := range blockRequests {
		o.query.Add("block_resources", blockRequest)
	}

	return o
}

// Cache allows to cache the screenshot.
func (o *TakeOptions) Cache(cache bool) *TakeOptions {
	o.query.Add("cache", strconv.FormatBool(cache))

	return o
}

// CacheTTL sets cache TTL.
func (o *TakeOptions) CacheTTL(cacheTTL int) *TakeOptions {
	o.query.Add("cache_ttl", strconv.Itoa(cacheTTL))

	return o
}

// CacheTTL sets cache key.
func (o *TakeOptions) CacheKey(cacheKey string) *TakeOptions {
	o.query.Add("cache_key", cacheKey)

	return o
}

// UserAgent sets a user agent for the request.
func (o *TakeOptions) UserAgent(userAgent string) *TakeOptions {
	o.query.Add("user_agent", userAgent)

	return o
}

// Authorization sets an authorization header for the request.
func (o *TakeOptions) Authorization(authorization string) *TakeOptions {
	o.query.Add("authorization", authorization)

	return o
}

// Cookies set cookies for the request.
func (o *TakeOptions) Cookies(cookies ...string) *TakeOptions {
	for _, cookie := range cookies {
		o.query.Add("cookies", cookie)
	}

	return o
}

// Headers sets extra headers for the request.
func (o *TakeOptions) Headers(headers ...string) *TakeOptions {
	for _, header := range headers {
		o.query.Add("headers", header)
	}

	return o
}

// TimeZone sets time zone for the request.
// Available time zones are: "America/Santiago", "Asia/Shanghai", "Europe/Berlin", "America/Guayaquil", "Europe/Madrid", "Pacific/Majuro", "Asia/Kuala_Lumpur", "Pacific/Auckland", "Europe/Lisbon", "Europe/Kiev", "Asia/Tashkent", "Europe/London".
func (o *TakeOptions) TimeZone(timeZone string) *TakeOptions {
	o.query.Add("time_zone", timeZone)

	return o
}

// Delay sets delay.
func (o *TakeOptions) Delay(delay int) *TakeOptions {
	o.query.Add("delay", strconv.Itoa(delay))

	return o
}

// Timeout sets timeout.
func (o *TakeOptions) Timeout(timeout int) *TakeOptions {
	o.query.Add("timeout", strconv.Itoa(timeout))

	return o
}
