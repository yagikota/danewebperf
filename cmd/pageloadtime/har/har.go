package har

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// HTTP Archive (HAR) 1.2 Spec
// http://www.softwareishard.com/blog/har-12-spec/

type Har struct {
	Log Log `json:"log"`
}

func (h *Har) DropEachResponseContent() *Har {
	for i := range h.Log.Entries {
		h.Log.Entries[i].Response.Content.Text = ""
	}
	return h
}

func (h *Har) ExistPages() bool {
	return len(h.Log.Pages) > 0
}

func (h *Har) ValidPageLoadTime() bool {
	if !h.ExistPages() {
		return false
	}
	return h.Log.Pages[0].PageTimings.OnLoad > 0
}

func (h *Har) StartedDateTimeOfFirstPage() time.Time {
	if !h.ExistPages() {
		return time.Time{}
	}
	return h.Log.Pages[0].StartedDateTime
}

func (h *Har) OnContentLoadOfFirstPage() int {
	if !h.ExistPages() {
		return -1
	}
	return h.Log.Pages[0].PageTimings.OnContentLoad
}

func (h *Har) OnLoadOfFirstPage() int {
	if !h.ExistPages() {
		return -1
	}
	return h.Log.Pages[0].PageTimings.OnLoad
}

func (h *Har) Entries() []Entry {
	return h.Log.Entries
}

func (h *Har) Informational() []Entry {
	var entries []Entry
	for _, entry := range h.Log.Entries {
		statusCode := entry.Response.Status
		if statusCode >= 200 && statusCode < 300 {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (h *Har) SuccessEntries() []Entry {
	var entries []Entry
	for _, entry := range h.Log.Entries {
		statusCode := entry.Response.Status
		if statusCode >= 200 && statusCode < 300 {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (h *Har) RedirectionEntries() []Entry {
	var entries []Entry
	for _, entry := range h.Log.Entries {
		statusCode := entry.Response.Status
		if statusCode >= 300 && statusCode < 400 {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (h *Har) ClientErrorEntries() []Entry {
	var entries []Entry
	for _, entry := range h.Log.Entries {
		statusCode := entry.Response.Status
		if statusCode >= 400 && statusCode < 500 {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (h *Har) ServerErrorEntries() []Entry {
	var entries []Entry
	for _, entry := range h.Log.Entries {
		statusCode := entry.Response.Status
		if statusCode >= 500 && statusCode < 600 {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (h *Har) Save(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := json.Marshal(h)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, out, 0644); err != nil {
		return err
	}
	return nil
}

type Log struct {
	// version [string] - Version number of the format. If empty, string "1.1" is assumed by default.
	Version string `json:"version"`
	// creator [object] - Name and version info of the log creator application.
	Creator Creator `json:"creator"`
	// browser [object, optional] - Name and version info of used browser.
	Browser Browser `json:"browser"`
	// pages [array, optional] - List of all exported (tracked) pages. Leave out this field if the application does not support grouping by pages.
	Pages []Page `json:"pages"`
	// entries [array] - List of all exported (tracked) requests.
	Entries []Entry `json:"entries"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Creator struct {
	// name [string] - Name of the application/browser used to export the log.
	Name string `json:"name"`
	// version [string] - Version of the application/browser used to export the log.
	Version string `json:"version"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Browser struct {
	// name [string] - Name of the application/browser used to export the log.
	Name string `json:"name"`
	// version [string] - Version of the application/browser used to export the log.
	Version string `json:"version"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Page struct {
	// startedDateTime [string] - Date and time stamp for the beginning of the page load (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD, e.g. 2009-07-24T19:20:30.45+01:00).
	StartedDateTime time.Time `json:"startedDateTime"`
	// id [string] - Unique identifier of a page within the <log>. Entries use it to refer the parent page.
	ID string `json:"id"`
	// title [string] - Page title.
	Title string `json:"title"`
	// pageTimings[object] - Detailed timing info about page load.
	PageTimings PageTimings `json:"pageTimings"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type PageTimings struct {
	// onContentLoad [number, optional] - Content of the page loaded. Number of milliseconds since page load started (page.startedDateTime). Use -1 if the timing does not apply to the current request.
	OnContentLoad int `json:"onContentLoad"`
	// onLoad [number,optional] - Page is loaded (onLoad event fired). Number of milliseconds since page load started (page.startedDateTime). Use -1 if the timing does not apply to the current request.
	OnLoad int `json:"onLoad"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Entry struct {
	// pageref [string, unique, optional] - Reference to the parent page. Leave out this field if the application does not support grouping by pages.
	Pageref string `json:"pageref"`
	// startedDateTime [string] - Date and time stamp of the request start (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD).
	StartedDateTime time.Time `json:"startedDateTime"`
	// time [number] - Total elapsed time of the request in milliseconds. This is the sum of all timings available in the timings object (i.e. not including -1 values) .
	//
	// entry.timings.blocked + entry.timings.dns +entry.timings.connect + entry.timings.send + entry.timings.wait +entry.timings.receive
	Time int `json:"time"`
	// request [object] - Detailed info about the request.
	Request Request `json:"request"`
	// response [object] - Detailed info about the response.
	Response Response `json:"response"`
	// cache [object] - Info about cache usage.
	Cache Cache `json:"cache"`
	// timings [object] - Detailed timing info about request/response round trip.
	Timings       Timings `json:"timings"`
	SecurityState string  `json:"_securityState"`
	// serverIPAddress [string, optional] (new in 1.2) - IP address of the server that was connected (result of DNS resolution).
	ServerIPAddress string `json:"serverIPAddress"`
	// connection [string, optional] (new in 1.2) - Unique ID of the parent TCP/IP connection, can be the client or server port number. Note that a port number doesn't have to be unique identifier in cases where the port is shared for more connections. If the port isn't available for the application, any other unique connection ID can be used instead (e.g. connection index). Leave out this field if the application doesn't support this info.
	Connection string `json:"connection"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Request struct {
	// method [string] - Request method (GET, POST, ...).
	Method string `json:"method"`
	// url [string] - Absolute URL of the request (fragments are not included).
	URL string `json:"url"`
	// httpVersion [string] - Request HTTP Version.
	HTTPVersion string `json:"httpVersion"`
	// cookies [array] - List of cookie objects.
	Cookies []Cookie `json:"cookies"`
	// headers [array] - List of header objects.
	Headers []Header `json:"headers"`
	// queryString [array] - List of query parameter objects.
	QueryString []QueryParam `json:"queryString"`
	// postData [object, optional] - Posted data info.
	PostData PostData `json:"postData"`
	// headersSize [number] - Total number of bytes from the start of the HTTP request message until (and including) the double CRLF before the body. Set to -1 if the info is not available.
	HeadersSize int `json:"headersSize"`
	// bodySize [number] - Size of the request body (POST data payload) in bytes. Set to -1 if the info is not available.
	BodySize int `json:"bodySize"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Cookie struct {
	// name [string] - The name of the cookie.
	Name string `json:"name"`
	// value [string] - The cookie value.
	Value string `json:"value"`
	// path [string, optional] - The path pertaining to the cookie.
	Path string `json:"path"`
	// domain [string, optional] - The host of the cookie.
	Domain string `json:"domain"`
	// expires [string, optional] - Cookie expiration time. (ISO 8601 - YYYY-MM-DDThh:mm:ss.sTZD, e.g. 2009-07-24T19:20:30.123+02:00).
	Expires time.Time `json:"expires"`
	// httpOnly [boolean, optional] - Set to true if the cookie is HTTP only, false otherwise.
	HTTPOnly bool `json:"httpOnly"`
	// secure [boolean, optional] (new in 1.2) - True if the cookie was transmitted over ssl, false otherwise.
	Secure bool `json:"secure"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Header struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

type QueryParam struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

type PostData struct {
	// mimeType [string] - Mime type of posted data.
	MimeType string `json:"mimeType"`
	// params [array] - List of posted parameters (in case of URL encoded parameters).
	Params []Param `json:"params"`
	// text [string] - Plain text posted data
	Text string `json:"text"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Param struct {
	// name [string] - name of a posted parameter.
	Name string `json:"name"`
	// value [string, optional] - value of a posted parameter or content of a posted file.
	Value string `json:"value"`
	// fileName [string, optional] - name of a posted file.
	FileName string `json:"fileName"`
	// contentType [string, optional] - content type of a posted file.
	ContentType string `json:"contentType"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Response struct {
	// status [number] - Response status.
	Status int `json:"status"`

	// statusText [string] - Response status description.
	StatusText string `json:"statusText"`

	// httpVersion [string] - Response HTTP Version.
	HTTPVersion string `json:"httpVersion"`

	// cookies [array] - List of cookie objects.
	Cookies []any `json:"cookies"`

	// headers [array] - List of header objects.
	Headers []Header `json:"headers"`

	// content [object] - Details about the response body.
	Content Content `json:"content"`

	// redirectURL [string] - Redirection target URL from the Location response header.
	RedirectURL string `json:"redirectURL"`

	// headersSize [number]* - Total number of bytes from the start of the HTTP response message until (and including) the double CRLF before the body. Set to -1 if the info is not available.
	// bodySize [number] - Size of the received response body in bytes. Set to zero in case of responses coming from the cache (304). Set to -1 if the info is not available.
	//
	// The size of received response-headers is computed only from headers that are really received from the server. Additional headers appended by the browser are not included in this number, but they appear in the list of header objects.
	// The total response size received can be computed as follows (if both values are available):
	// totalSize = entry.response.headersSize + entry.response.bodySize
	HeadersSize int `json:"headersSize"`
	BodySize    int `json:"bodySize"`

	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Content struct {
	// size [number] - Length of the returned content in bytes. Should be equal to response.bodySize if there is no compression and bigger when the content has been compressed.
	Size int `json:"size"`
	// compression [number, optional] - Number of bytes saved. Leave out this field if the information is not available.
	Compression int `json:"compression"`
	// mimeType [string] - MIME type of the response text (value of the Content-Type response header). The charset attribute of the MIME type is included (if available).
	MimeType string `json:"mimeType"`
	// text [string, optional] - Response body sent from the server or loaded from the browser cache. This field is populated with textual content only. The text field is either HTTP decoded text or a encoded (e.g. "base64") representation of the response body. Leave out this field if the information is not available.
	Text string `json:"text"`
	// encoding [string, optional] (new in 1.2) - Encoding used for response text field e.g "base64". Leave out this field if the text field is HTTP decoded (decompressed & unchunked), than trans-coded from its original character set into UTF-8.
	Encoding string `json:"encoding"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Cache struct {
	// beforeRequest [object, optional] - State of a cache entry before the request. Leave out this field if the information is not available.
	BeforeRequest CacheRequest `json:"beforeRequest"`
	// afterRequest [object, optional] - State of a cache entry after the request. Leave out this field if the information is not available.
	AfterRequest CacheRequest `json:"afterRequest"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type CacheRequest struct {
	// expires [string, optional] - Expiration time of the cache entry.
	Expires time.Time `json:"expires"`
	// lastAccess [string] - The last time the cache entry was opened.
	LastAccess time.Time `json:"lastAccess"`
	// eTag [string] - Etag
	ETag string `json:"eTag"`
	// hitCount [number] - The number of times the cache entry has been opened.
	HitCount int `json:"hitCount"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type Timings struct {
	// 	blocked [number, optional] - Time spent in a queue waiting for a network connection. Use -1 if the timing does not apply to the current request.
	Blocked int `json:"blocked"`
	// dns [number, optional] - DNS resolution time. The time required to resolve a host name. Use -1 if the timing does not apply to the current request.
	DNS int `json:"dns"`
	// connect [number, optional] - Time required to create TCP connection. Use -1 if the timing does not apply to the current request.
	Connect int `json:"connect"`
	// send [number] - Time required to send HTTP request to the server.
	Send int `json:"send"`
	// wait [number] - Waiting for a response from the server.
	Wait int `json:"wait"`
	// receive [number] - Time required to read entire response from the server (or cache).
	Receive int `json:"receive"`
	// ssl [number, optional] (new in 1.2) - Time required for SSL/TLS negotiation. If this field is defined then the time is also included in the connect field (to ensure backward compatibility with HAR 1.1). Use -1 if the timing does not apply to the current request.
	Ssl int `json:"ssl"`
	// comment [string, optional] (new in 1.2) - A comment provided by the user or the application.
	Comment string `json:"comment"`
}

type CSVFormat struct {
	Records []Record
}

type Record struct {
	Status string
	Method string
	Domain string
	File   string
	// Initiator        string // ignore this field
	MIMEType                string
	CompressedSize          string
	UnCompressedSize        string
	PageLoadStartedDateTime string
	Transaction             Transaction
}

type Transaction struct {
	StartedDateTime string
	Queued          string
	Started         string
	Downloaded      string
	Blocked         string
	DNSResolution   string
	Connecting      string
	TLSSetup        string
	Sending         string
	Waiting         string
	Receiving       string
}

func (h *Har) ConvertCSVFormat() CSVFormat {
	var records []Record
	for _, entry := range h.Log.Entries {
		u, err := url.Parse(entry.Request.URL)
		if err != nil {
			log.Println(err)
			continue
		}
		domain := u.Hostname()

		pathSegments := strings.Split(u.Path, "/")
		file := pathSegments[len(pathSegments)-1]
		if file == "" {
			file = "/"
		}

		queued := int(entry.StartedDateTime.Sub(h.StartedDateTimeOfFirstPage()).Milliseconds())

		started := queued + entry.Timings.Blocked

		downloaded := queued + entry.Time

		record := Record{
			Status:                  strconv.Itoa(entry.Response.Status),
			Method:                  entry.Request.Method,
			Domain:                  domain,
			File:                    file,
			MIMEType:                entry.Response.Content.MimeType,
			CompressedSize:          strconv.Itoa(entry.Response.BodySize),
			UnCompressedSize:        strconv.Itoa(entry.Response.Content.Size),
			PageLoadStartedDateTime: h.StartedDateTimeOfFirstPage().Format("2006-01-02 15:04:05.000"),
			Transaction: Transaction{
				StartedDateTime: entry.StartedDateTime.Format("2006-01-02 15:04:05.000"),
				Queued:          strconv.Itoa(queued),
				Started:         strconv.Itoa(started),
				Downloaded:      strconv.Itoa(downloaded),
				Blocked:         strconv.Itoa(entry.Timings.Blocked),
				DNSResolution:   strconv.Itoa(entry.Timings.DNS),
				Connecting:      strconv.Itoa(entry.Timings.Connect),
				TLSSetup:        strconv.Itoa(entry.Timings.Ssl),
				Sending:         strconv.Itoa(entry.Timings.Send),
				Waiting:         strconv.Itoa(entry.Timings.Wait),
				Receiving:       strconv.Itoa(entry.Timings.Receive),
			},
		}
		records = append(records, record)
	}
	return CSVFormat{Records: records}
}

func (har CSVFormat) SaveAsCSV(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		"Status",
		"Method",
		"Domain",
		"File",
		"MIMEType",
		"CompressedSize(B)",
		"UnCompressedSize(B)",
		"PageLoadStartedDateTime",
		"StartedDateTime",
		"Queued(ms)",
		"Started(ms)",
		"Downloaded(ms)",
		"Blocked(ms)",
		"DNSResolution(ms)",
		"Connecting(ms)",
		"TLSSetup(ms)",
		"Sending(ms)",
		"Waiting(ms)",
		"Receiving(ms)",
	}); err != nil {
		return err
	}

	for _, record := range har.Records {
		record := []string{
			record.Status,
			record.Method,
			record.Domain,
			record.File,
			record.MIMEType,
			record.CompressedSize,
			record.UnCompressedSize,
			record.PageLoadStartedDateTime,
			record.Transaction.StartedDateTime,
			record.Transaction.Queued,
			record.Transaction.Started,
			record.Transaction.Downloaded,
			record.Transaction.Blocked,
			record.Transaction.DNSResolution,
			record.Transaction.Connecting,
			record.Transaction.TLSSetup,
			record.Transaction.Sending,
			record.Transaction.Waiting,
			record.Transaction.Receiving,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}
