/**
 * http-redirect-tracking a poc about how to implement a tracking system using HTTP redirect headers.
 *
 * @see http://elie.im/blog/security/tracking-users-that-block-cookies-with-a-http-redirect/
 * @author {tom@0x101.com}
 */
package main
 
import (
	"fmt"
	"http"
	"template"
	"strings"
	"time"
	"strconv"
)

const portString = ":8080"
const id = "0123456789abcdefghijklmnopqrstuvwxyz"

// Used to represent the content of the main template
type Page struct {
	Body	string
	Title	string
	Time	string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	timestamp := strconv.Itoa64(time.Seconds())
	tracking := strings.Split(r.RawURL, "?id=", 2)

	var trackingId string = ""
	if len(tracking) > 1 {
		// We have already an id
		trackingId = tracking[1]
	} else {
		// Inject tracking id
		trackingId = id

		t := time.LocalTime()
		date := t.Format(time.RFC1123)

		t.Year = t.Year +1;
		expires := t.Format(time.RFC1123)
		maxAge :=  t.Format(time.RFC1123)

		// Write necessary headers
		w.Header().Set("Pragma",  "public");
		w.Header().Set("Last-Modified", date);
		w.Header().Set("Cache-Control", "maxage=" + maxAge);
		w.Header().Set("Pragma", "public");
		w.Header().Set("Expires", expires +" GMT\n");
		
		// Redirect
		http.Redirect(w, r, "?id="+trackingId, 301)
	}
	
	// Prepare the template
	page := Page{Body: trackingId, Title: "http-redirect-tracking", Time: timestamp}
	t, _ := template.ParseFile("./templates/index.html", nil)
	t.Execute(w, page)
}

func main() {
	fmt.Printf("Running server\n")
	http.HandleFunc("/view/", viewHandler)
	http.ListenAndServe(portString, nil)
}

