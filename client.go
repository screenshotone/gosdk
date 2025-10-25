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

// NewClientWithHTTPClient returns new API client for the ScreenshotOne.com API with a custom HTTP client.
func NewClientWithHTTPClient(accessKey, secretKey string, httpClient *http.Client) (*Client, error) {
	client := &Client{accessKey, secretKey, httpClient}

	return client, nil
}

// GenerateTakeURL generates URL for taking screenshots with request signing.
func (client *Client) GenerateTakeURL(options *TakeOptions) (*url.URL, error) {
	if client.secretKey == "" {
		return nil, fmt.Errorf("secret key is required for signed URLs")
	}

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

// GenerateUnsignedTakeURL generates URL for taking screenshots without signing the request.
func (client *Client) GenerateUnsignedTakeURL(options *TakeOptions) (*url.URL, error) {
	// generate query
	query := options.query
	query.Set("access_key", client.accessKey)
	queryString := query.Encode()

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

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
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

	return NewTakeWithURL(pageURL)
}

// Returns options for the ScreenshotOne.com API take method.
func NewTakeWithURL(pageURL string) *TakeOptions {
	query := url.Values{}
	query.Add("url", pageURL)

	return &TakeOptions{query: query}
}

// Returns options for the ScreenshotOne.com API take method.
func NewTakeWithHTML(html string) *TakeOptions {
	query := url.Values{}
	query.Add("html", html)

	return &TakeOptions{query: query}
}

// Returns options for the ScreenshotOne.com API take method.
func NewTakeWithMarkdown(markdown string) *TakeOptions {
	query := url.Values{}
	query.Add("markdown", markdown)

	return &TakeOptions{query: query}
}

// Selector is a CSS-like selector of the element to take a screenshot of.
func (o *TakeOptions) Selector(selector string) *TakeOptions {
	o.query.Add("selector", selector)

	return o
}

// SelectorAlgorithm sets the algorithm for finding selectors.
func (o *TakeOptions) SelectorAlgorithm(algorithm string) *TakeOptions {
	o.query.Add("selector_algorithm", algorithm)

	return o
}

// ErrorOnSelectorNotFound determines the behavior of what to do when selector is not found.
func (o *TakeOptions) ErrorOnSelectorNotFound(errorOn bool) *TakeOptions {
	o.query.Add("error_on_selector_not_found", strconv.FormatBool(errorOn))

	return o
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

// Styles specifies custom CSS styles for the page.
func (o *TakeOptions) Styles(styles string) *TakeOptions {
	o.query.Add("styles", styles)

	return o
}

// Scripts specifies custom scripts for the page.
func (o *TakeOptions) Scripts(scripts string) *TakeOptions {
	o.query.Add("scripts", scripts)

	return o
}

// ScriptsWaitUntil sets when to wait for scripts to complete.
func (o *TakeOptions) ScriptsWaitUntil(waitUntil string) *TakeOptions {
	o.query.Add("scripts_wait_until", waitUntil)

	return o
}

// ImageQuality renders image with the specified quality. Available for the next formats: "jpeg" ("jpg"), "webp".
func (o *TakeOptions) ImageQuality(imageQuality int) *TakeOptions {
	o.query.Add("image_quality", strconv.Itoa(imageQuality))

	return o
}

// ImageHeight sets the height of the resulting image (pixels).
func (o *TakeOptions) ImageHeight(imageHeight int) *TakeOptions {
	o.query.Add("image_height", strconv.Itoa(imageHeight))

	return o
}

// ImageWidth sets the width of the resulting image (pixels).
func (o *TakeOptions) ImageWidth(imageWidth int) *TakeOptions {
	o.query.Add("image_width", strconv.Itoa(imageWidth))

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
func (o *TakeOptions) DeviceScaleFactor(deviceScaleFactor int) *TakeOptions {
	o.query.Add("device_scale_factor", strconv.Itoa(deviceScaleFactor))

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
func (o *TakeOptions) BlockResources(blockResources ...string) *TakeOptions {
	for _, blockResource := range blockResources {
		o.query.Add("block_resources", blockResource)
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

// WaitUntil waits until an event occurred before taking a screenshot or rendering HTML or PDF.
func (o *TakeOptions) WaitUntil(events ...string) *TakeOptions {
	for _, event := range events {
		o.query.Add("wait_until", event)
	}
	return o
}

// WaitForSelector waits until the element appears in DOM.
func (o *TakeOptions) WaitForSelector(selector string) *TakeOptions {
	o.query.Add("wait_for_selector", selector)
	return o
}

// WaitForSelectorAlgorithm sets the algorithm for waiting for selectors.
func (o *TakeOptions) WaitForSelectorAlgorithm(algorithm string) *TakeOptions {
	o.query.Add("wait_for_selector_algorithm", algorithm)
	return o
}

// NavigationTimeout sets the timeout for navigation.
func (o *TakeOptions) NavigationTimeout(timeout int) *TakeOptions {
	o.query.Add("navigation_timeout", strconv.Itoa(timeout))
	return o
}

// CaptureBeyondViewport controls whether to capture beyond the viewport.
func (o *TakeOptions) CaptureBeyondViewport(capture bool) *TakeOptions {
	o.query.Add("capture_beyond_viewport", strconv.FormatBool(capture))
	return o
}

// FullPageScroll controls whether to scroll to the bottom of the page and back to the top.
func (o *TakeOptions) FullPageScroll(scroll bool) *TakeOptions {
	o.query.Add("full_page_scroll", strconv.FormatBool(scroll))

	return o
}

// FullPageScrollDelay sets the delay for full page scrolling.
func (o *TakeOptions) FullPageScrollDelay(delay int) *TakeOptions {
	o.query.Add("full_page_scroll_delay", strconv.Itoa(delay))

	return o
}

// FullPageScrollBy sets how much to scroll by for full page screenshots.
func (o *TakeOptions) FullPageScrollBy(pixels int) *TakeOptions {
	o.query.Add("full_page_scroll_by", strconv.Itoa(pixels))

	return o
}

// FullPageMaxHeight sets the maximum height for full page screenshots.
func (o *TakeOptions) FullPageMaxHeight(height int) *TakeOptions {
	o.query.Add("full_page_max_height", strconv.Itoa(height))

	return o
}

// HideSelectors hides elements matching the given selectors.
func (o *TakeOptions) HideSelectors(selectors ...string) *TakeOptions {
	for _, selector := range selectors {
		o.query.Add("hide_selectors", selector)
	}

	return o
}

// Click specifies a selector to click before taking the screenshot.
func (o *TakeOptions) Click(selector string) *TakeOptions {
	o.query.Add("click", selector)

	return o
}

// ScrollIntoView scrolls the page to ensure the given selector is in view.
func (o *TakeOptions) ScrollIntoView(selector string) *TakeOptions {
	o.query.Add("scroll_into_view", selector)

	return o
}

// ScrollIntoViewAdjustTop adjusts the top position after scrolling into view.
func (o *TakeOptions) ScrollIntoViewAdjustTop(pixels int) *TakeOptions {
	o.query.Add("scroll_into_view_adjust_top", strconv.Itoa(pixels))

	return o
}

// DarkMode sets the dark mode for the screenshot.
func (o *TakeOptions) DarkMode(enabled bool) *TakeOptions {
	o.query.Add("dark_mode", strconv.FormatBool(enabled))

	return o
}

// ReducedMotion sets the reduced motion mode for the screenshot.
func (o *TakeOptions) ReducedMotion(enabled bool) *TakeOptions {
	o.query.Add("reduced_motion", strconv.FormatBool(enabled))

	return o
}

// MediaType sets the media type for the screenshot.
func (o *TakeOptions) MediaType(mediaType string) *TakeOptions {
	o.query.Add("media_type", mediaType)

	return o
}

// ViewportMobile sets whether the meta viewport tag is taken into account.
func (o *TakeOptions) ViewportMobile(mobile bool) *TakeOptions {
	o.query.Add("viewport_mobile", strconv.FormatBool(mobile))

	return o
}

// ViewportHasTouch sets whether the viewport supports touch events.
func (o *TakeOptions) ViewportHasTouch(hasTouch bool) *TakeOptions {
	o.query.Add("viewport_has_touch", strconv.FormatBool(hasTouch))

	return o
}

// ViewportLandscape sets whether the viewport is in landscape mode.
func (o *TakeOptions) ViewportLandscape(landscape bool) *TakeOptions {
	o.query.Add("viewport_landscape", strconv.FormatBool(landscape))

	return o
}

// ViewportDevice sets the device for viewport emulation.
func (o *TakeOptions) ViewportDevice(device string) *TakeOptions {
	o.query.Add("viewport_device", device)

	return o
}

// BlockCookieBanners blocks cookie banners and privacy notices.
func (o *TakeOptions) BlockCookieBanners(block bool) *TakeOptions {
	o.query.Add("block_cookie_banners", strconv.FormatBool(block))

	return o
}

// BlockBannersByHeuristics blocks banners using heuristics.
func (o *TakeOptions) BlockBannersByHeuristics(block bool) *TakeOptions {
	o.query.Add("block_banners_by_heuristics", strconv.FormatBool(block))

	return o
}

// BlockChats blocks chat widgets.
func (o *TakeOptions) BlockChats(block bool) *TakeOptions {
	o.query.Add("block_chats", strconv.FormatBool(block))

	return o
}

// BypassCSP bypasses Content Security Policy.
func (o *TakeOptions) BypassCSP(bypass bool) *TakeOptions {
	o.query.Add("bypass_csp", strconv.FormatBool(bypass))

	return o
}

// Proxy sets a custom proxy for the request.
func (o *TakeOptions) Proxy(proxyURL string) *TakeOptions {
	o.query.Add("proxy", proxyURL)

	return o
}

// IPCountryCode sets the country for IP-based geolocation.
func (o *TakeOptions) IPCountryCode(countryCode string) *TakeOptions {
	o.query.Add("ip_country_code", countryCode)

	return o
}

// ResponseType sets the type of response to return.
func (o *TakeOptions) ResponseType(responseType string) *TakeOptions {
	o.query.Add("response_type", responseType)

	return o
}

// Store enables storing the screenshot in S3-compatible storage.
func (o *TakeOptions) Store(store bool) *TakeOptions {
	o.query.Add("store", strconv.FormatBool(store))

	return o
}

// StoragePath sets the storage path for the screenshot.
func (o *TakeOptions) StoragePath(path string) *TakeOptions {
	o.query.Add("storage_path", path)

	return o
}

// StorageACL sets the ACL for the stored screenshot.
func (o *TakeOptions) StorageACL(acl string) *TakeOptions {
	o.query.Add("storage_acl", acl)

	return o
}

// StorageReturnLocation enables returning the storage location.
func (o *TakeOptions) StorageReturnLocation(returnLocation bool) *TakeOptions {
	o.query.Add("storage_return_location", strconv.FormatBool(returnLocation))

	return o
}

// Async enables asynchronous execution of the request.
func (o *TakeOptions) Async(async bool) *TakeOptions {
	o.query.Add("async", strconv.FormatBool(async))

	return o
}

// WebhookURL sets the URL for the webhook.
func (o *TakeOptions) WebhookURL(url string) *TakeOptions {
	o.query.Add("webhook_url", url)

	return o
}

// WebhookSign controls whether to sign the webhook request body.
func (o *TakeOptions) WebhookSign(sign bool) *TakeOptions {
	o.query.Add("webhook_sign", strconv.FormatBool(sign))

	return o
}

// WebhookErrors enables error details in webhook requests.
func (o *TakeOptions) WebhookErrors(enable bool) *TakeOptions {
	o.query.Add("webhook_errors", strconv.FormatBool(enable))

	return o
}

// RequestGPURendering requests GPU rendering for the screenshot.
func (o *TakeOptions) RequestGPURendering(request bool) *TakeOptions {
	o.query.Add("request_gpu_rendering", strconv.FormatBool(request))

	return o
}

// IncludeShadowDOM includes shadow DOM elements in the content.
func (o *TakeOptions) IncludeShadowDOM(include bool) *TakeOptions {
	o.query.Add("include_shadow_dom", strconv.FormatBool(include))

	return o
}

// AttachmentName sets the attachment name for the response.
func (o *TakeOptions) AttachmentName(name string) *TakeOptions {
	o.query.Add("attachment_name", name)

	return o
}

// ExternalIdentifier sets an external identifier for the request.
func (o *TakeOptions) ExternalIdentifier(identifier string) *TakeOptions {
	o.query.Add("external_identifier", identifier)

	return o
}

// FailIfGPURenderingFails forces the request to fail if GPU rendering fails.
func (o *TakeOptions) FailIfGPURenderingFails(fail bool) *TakeOptions {
	o.query.Add("fail_if_gpu_rendering_fails", strconv.FormatBool(fail))

	return o
}

// MetadataImageSize enables returning the actual image size metadata.
func (o *TakeOptions) MetadataImageSize(enable bool) *TakeOptions {
	o.query.Add("metadata_image_size", strconv.FormatBool(enable))

	return o
}

// MetadataFonts enables returning the fonts used by the website.
func (o *TakeOptions) MetadataFonts(enable bool) *TakeOptions {
	o.query.Add("metadata_fonts", strconv.FormatBool(enable))

	return o
}

// MetadataOpenGraph enables returning the Open Graph metadata.
func (o *TakeOptions) MetadataOpenGraph(enable bool) *TakeOptions {
	o.query.Add("metadata_open_graph", strconv.FormatBool(enable))

	return o
}

// MetadataPageTitle enables returning the page title.
func (o *TakeOptions) MetadataPageTitle(enable bool) *TakeOptions {
	o.query.Add("metadata_page_title", strconv.FormatBool(enable))

	return o
}

// MetadataHTTPResponseHeaders enables returning the HTTP response headers.
func (o *TakeOptions) MetadataHTTPResponseHeaders(enable bool) *TakeOptions {
	o.query.Add("metadata_http_response_headers", strconv.FormatBool(enable))

	return o
}

// MetadataHTTPResponseStatusCode enables returning the HTTP response status code.
func (o *TakeOptions) MetadataHTTPResponseStatusCode(enable bool) *TakeOptions {
	o.query.Add("metadata_http_response_status_code", strconv.FormatBool(enable))

	return o
}

// MetadataContent enables returning the content of the page.
func (o *TakeOptions) MetadataContent(enable bool) *TakeOptions {
	o.query.Add("metadata_content", strconv.FormatBool(enable))

	return o
}

// OpenAIAPIKey sets the OpenAI API key for vision integration.
func (o *TakeOptions) OpenAIAPIKey(key string) *TakeOptions {
	o.query.Add("openai_api_key", key)

	return o
}

// VisionPrompt sets the prompt for OpenAI vision integration.
func (o *TakeOptions) VisionPrompt(prompt string) *TakeOptions {
	o.query.Add("vision_prompt", prompt)

	return o
}

// VisionMaxTokens sets the maximum number of tokens for OpenAI vision integration.
func (o *TakeOptions) VisionMaxTokens(tokens int) *TakeOptions {
	o.query.Add("vision_max_tokens", strconv.Itoa(tokens))

	return o
}

// FailIfContentContains forces the request to fail if the specified text is matched on the page.
func (o *TakeOptions) FailIfContentContains(text string) *TakeOptions {
	o.query.Add("fail_if_content_contains", text)

	return o
}

// FailIfContentMissing forces the request to fail if the specified text is missing on the page.
func (o *TakeOptions) FailIfContentMissing(text string) *TakeOptions {
	o.query.Add("fail_if_content_missing", text)

	return o
}

// FailIfRequestFailed forces the request to fail if any network request fails during page loading.
func (o *TakeOptions) FailIfRequestFailed(pattern string) *TakeOptions {
	o.query.Add("fail_if_request_failed", pattern)

	return o
}

// PDFPrintBackground sets whether to print background graphics in PDF.
func (o *TakeOptions) PDFPrintBackground(pdfPrintBackground bool) *TakeOptions {
	o.query.Add("pdf_print_background", strconv.FormatBool(pdfPrintBackground))
	return o
}

// PDFFitOnePage tries to fit the website on one page for PDF output.
func (o *TakeOptions) PDFFitOnePage(pdfFitOnePage bool) *TakeOptions {
	o.query.Add("pdf_fit_one_page", strconv.FormatBool(pdfFitOnePage))
	return o
}

// PDFLandscape sets PDF orientation to landscape.
func (o *TakeOptions) PDFLandscape(pdfLandscape bool) *TakeOptions {
	o.query.Add("pdf_landscape", strconv.FormatBool(pdfLandscape))
	return o
}

// PDFPaperFormat specifies the paper format for PDF output.
func (o *TakeOptions) PDFPaperFormat(format string) *TakeOptions {
	o.query.Add("pdf_paper_format", format)
	return o
}

// PDFMargin sets the margin for PDF output.
func (o *TakeOptions) PDFMargin(margin string) *TakeOptions {
	o.query.Add("pdf_margin", margin)
	return o
}

// PDFMarginTop sets the top margin for PDF output.
func (o *TakeOptions) PDFMarginTop(margin string) *TakeOptions {
	o.query.Add("pdf_margin_top", margin)
	return o
}

// PDFMarginRight sets the right margin for PDF output.
func (o *TakeOptions) PDFMarginRight(margin string) *TakeOptions {
	o.query.Add("pdf_margin_right", margin)
	return o
}

// PDFMarginBottom sets the bottom margin for PDF output.
func (o *TakeOptions) PDFMarginBottom(margin string) *TakeOptions {
	o.query.Add("pdf_margin_bottom", margin)
	return o
}

// PDFMarginLeft sets the left margin for PDF output.
func (o *TakeOptions) PDFMarginLeft(margin string) *TakeOptions {
	o.query.Add("pdf_margin_left", margin)
	return o
}

// ClipX sets the x coordinate of the area to clip.
func (o *TakeOptions) ClipX(x int) *TakeOptions {
	o.query.Add("clip_x", strconv.Itoa(x))
	return o
}

// ClipY sets the y coordinate of the area to clip.
func (o *TakeOptions) ClipY(y int) *TakeOptions {
	o.query.Add("clip_y", strconv.Itoa(y))
	return o
}

// ClipWidth sets the width of the area to clip.
func (o *TakeOptions) ClipWidth(width int) *TakeOptions {
	o.query.Add("clip_width", strconv.Itoa(width))
	return o
}

// ClipHeight sets the height of the area to clip.
func (o *TakeOptions) ClipHeight(height int) *TakeOptions {
	o.query.Add("clip_height", strconv.Itoa(height))
	return o
}

// FullPageAlgorithm sets the algorithm for full page screenshots.
func (o *TakeOptions) FullPageAlgorithm(algorithm string) *TakeOptions {
	o.query.Add("full_page_algorithm", algorithm)
	return o
}

// SelectorScrollIntoView controls scrolling behavior when taking element screenshots.
func (o *TakeOptions) SelectorScrollIntoView(enable bool) *TakeOptions {
	o.query.Add("selector_scroll_into_view", strconv.FormatBool(enable))
	return o
}

// IgnoreHostErrors allows taking screenshots even when site returns error status codes.
func (o *TakeOptions) IgnoreHostErrors(ignore bool) *TakeOptions {
	o.query.Add("ignore_host_errors", strconv.FormatBool(ignore))
	return o
}

// ErrorOnClickSelectorNotFound controls error behavior when click selector is not found.
func (o *TakeOptions) ErrorOnClickSelectorNotFound(errorOn bool) *TakeOptions {
	o.query.Add("error_on_click_selector_not_found", strconv.FormatBool(errorOn))
	return o
}

// StorageEndpoint sets custom S3-compatible storage endpoint.
func (o *TakeOptions) StorageEndpoint(endpoint string) *TakeOptions {
	o.query.Add("storage_endpoint", endpoint)
	return o
}

// StorageAccessKeyID sets storage access key ID.
func (o *TakeOptions) StorageAccessKeyID(keyID string) *TakeOptions {
	o.query.Add("storage_access_key_id", keyID)
	return o
}

// StorageSecretAccessKey sets storage secret access key.
func (o *TakeOptions) StorageSecretAccessKey(key string) *TakeOptions {
	o.query.Add("storage_secret_access_key", key)
	return o
}

// StorageBucket sets storage bucket name.
func (o *TakeOptions) StorageBucket(bucket string) *TakeOptions {
	o.query.Add("storage_bucket", bucket)
	return o
}

// StorageClass sets storage class for the object.
func (o *TakeOptions) StorageClass(class string) *TakeOptions {
	o.query.Add("storage_class", class)
	return o
}

// MetadataIcon enables returning the favicon metadata.
func (o *TakeOptions) MetadataIcon(enable bool) *TakeOptions {
	o.query.Add("metadata_icon", strconv.FormatBool(enable))
	return o
}
