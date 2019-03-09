package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	verbose *bool
)

func main() {
	port := flag.Int("p", 8080, "The Port the Server listens")
	verbose = flag.Bool("v", false, "increase verbosity")

	flag.Parse()

	http.HandleFunc("/", helpHandler)
	http.HandleFunc("/help", helpHandler)
	http.HandleFunc("/status/", statusHandler)
	http.HandleFunc("/delay/", delayHandler)
	http.HandleFunc("/ip", ipHandler)
	http.HandleFunc("/userAgent", userAgentHandler)
	http.HandleFunc("/basicAuth/", basicAuthHandler)
	http.HandleFunc("/file/", fileHandler)

	fmt.Printf("Listen on Port %d\n", *port)
	fmt.Printf("Open http://localhost:%d/\n", *port)

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		printVerbose(r.URL.String())
	}

	helpText := `<html>
	<head>
		<title>HTTPTestEndpoint</title>
		<style>
			body{
				background-color: #313131;
				color: #FFFFFF;
				font-family: 'Open Sans', sans-serif;
				font-weight: 300;
			}
			h1{
				font-weight: 600;
				text-align: center;
			}
			h2{
				font-weight: 600;
				cursor: pointer;
			}
			.method{
				border-bottom: 1px solid #FFFFFF;
				margin: auto;
				padding: 10px;
				max-width: 700px;
			}
			.down{
				background: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIyOCIgaGVpZ2h0PSIyOCIgdmlld0JveD0iMCAwIDI4IDI4Ij4KICA8cG9seWdvbiBmaWxsPSIjRkZGIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiIHBvaW50cz0iMjIgNC42NjcgOS41NTYgMTQgMjIgMjMuMzMzIDIyIDI2IDYgMTQgMjIgMiIgdHJhbnNmb3JtPSJyb3RhdGUoLTkwIDE0IDE0KSIvPgo8L3N2Zz4K);
				float: right;
				height: 25px;
				width: 25px;
			}
			table{
				border: 1px solid #FFFFFF;
				width: 100%;
			}
			th, td{
				padding: 5px;
			}
			ul{
				margin: 0;
			}
		</style>
		<script>
			window.addEventListener("load", function(){
				var headings = document.getElementsByTagName("h2");

				for(var i = 0; i < headings.length; i++){
					headings[i].addEventListener("click", details);
				}
			});
			
			function details(e){
				if(e.srcElement.nodeName == "H2"){
					var element = e.target.parentElement.querySelector(".details");
				} else{
					var element = e.target.parentElement.parentElement.querySelector(".details");
				}
			
				if(element.hidden == true){
					document.querySelectorAll(".details:not([hidden])").forEach(function(e){
						e.hidden = true;
					});

					element.hidden = false;
				} else{
					element.hidden = true;
				}
			}
		</script>
		<link href="https://fonts.googleapis.com/css?family=Open+Sans:300,600" rel="stylesheet">
	</head>
	<body>
		<h1>HTTPTestEndpoint</h1>
		<div class="method">
			<h2>/status/[code]<div class="down"></div></h2>
			<div class="details" hidden>
				<p>Answers with the passed HTTP status code. Optionally can the received HTTP body and content-type header also be in the answer.</p>
				<h3>Parameters</h3>
				<table rules="all">
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Description</th>
						<th>DataType</th>
					</tr>
					<tr>
						<td>code</th>
						<td>Path</th>
						<td>The HTTP Status Code</td>
						<td>Number</td>
					</tr>
					<tr>
						<td>returnBody</td>
						<td>Query</td>
						<td>Return the Body and Content-Type that was sent to the Server</td>
						<td>Bool</th>
					</tr>
				</table>
			</div>
		</div>
		<div class="method">
			<h2>/delay/[seconds]<div class="down"></div></h2>
			<div class="details" hidden>
				<p>Waits the passed seconds and answers 200 OK. Optionally can the received HTTP body and content-type header also be in the answer.</p>
				<h3>Parameters</h3>
				<table rules="all">
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Description</th>
						<th>DataType</th>
					</tr>
					<tr>
						<td>seconds</th>
						<td>Path</th>
						<td>The Seconds the Request should wait</td>
						<td>Number</td>
					</tr>
					<tr>
						<td>returnBody</td>
						<td>Query</td>
						<td>Return the Body and Content-Type that was sent to the Server</td>
						<td>Bool</th>
					</tr>
				</table>
			</div>
		</div>
		<div class="method">
			<h2>/ip<div class="down"></div></h2>
			<div class="details" hidden>
				<p>Returns the remote IP address in the HTTP body.</p>
				<h3>Parameters</h3>
				<table rules="all">
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Description</th>
						<th>DataType</th>
					</tr>
					<tr>
						<td>format</td>
						<td>Query</td>
						<td>
							The Return Format of the Data<br/>
							<ul>
								<li>text (default)</li>
								<li>json</li>
							</ul>
						</td>
						<td>Text</th>
					</tr>
				</table>
			</div>
		</div>
		<div class="method">
			<h2>/userAgent<div class="down"></div></h2>
			<div class="details" hidden>
				<p>Returns the user-agent string from the HTTP header in the HTTP body.</p>
				<h3>Parameters</h3>
				<table rules="all">
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Description</th>
						<th>DataType</th>
					</tr>
					<tr>
						<td>format</td>
						<td>Query</td>
						<td>
							The Return Format of the Data<br/>
							<ul>
								<li>text (default)</li>
								<li>json</li>
							</ul>
						</td>
						<td>Text</th>
					</tr>
				</table>
			</div>
		</div>
		<div class="method">
			<h2>/basicAuth/[username]/[password]<div class="down"></div></h2>
			<div class="details" hidden>
				<p>Checks if the username and password of the authorization are the same as the one passed in the URL path. If yes 200 OK and if not 401 unauthorized will be answered.</p>.
				<h3>Parameters</h3>
				<table rules="all">
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Description</th>
						<th>DataType</th>
					</tr>
					<tr>
						<td>username</td>
						<td>Path</td>
						<td></td>
						<td>Text</th>
					</tr>
					<tr>
						<td>password</td>
						<td>Path</td>
						<td></td>
						<td>Text</th>
					</tr>
				</table>
			</div>
		</div>
		<div class="method">
			<h2>/file/[filename]<div class="down"></div></h2>
			<div class="details" hidden>
				<p>Returns the content of the passed file in the HTTP body. The file must be in the same directory the HTTPTestEndpoint is.</p>
				<h3>Parameters</h3>
				<table rules="all">
					<tr>
						<th>Name</th>
						<th>Type</th>
						<th>Description</th>
						<th>DataType</th>
					</tr>
					<tr>
						<td>filename</td>
						<td>Path</td>
						<td></td>
						<td>Text</th>
					</tr>
				</table>
			</div>
		</div>
	</body>
</html>`

	fmt.Fprintf(w, "%s", helpText)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		printVerbose(r.URL.String())
	}

	path := strings.Split(r.URL.Path, "/")

	if path[2] == "" {
		http.Error(w, "URL must be /status/[code (number)]", http.StatusBadRequest)
		return
	}

	statusCode, err := strconv.Atoi(path[2])

	if err != nil {
		http.Error(w, "URL must be /status/[code (number)]", http.StatusBadRequest)
		return
	}

	returnBody := r.FormValue("returnBody")

	if returnBody == "true" {
		addHeader(w, r)
	}

	// must be after i set the Headers
	w.WriteHeader(statusCode)

	if returnBody == "true" {
		addBody(w, r)
	}
}

func delayHandler(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		printVerbose(r.URL.String())
	}

	path := strings.Split(r.URL.Path, "/")

	if path[2] == "" {
		http.Error(w, "URL must be /status/[code (number)]", http.StatusBadRequest)
		return
	}

	timeSleep, err := strconv.Atoi(path[2])

	if err != nil {
		http.Error(w, "URL must be /status/[code (number)]", http.StatusBadRequest)
		return
	}

	returnBody := r.FormValue("returnBody")

	time.Sleep(time.Duration(timeSleep) * time.Second)

	if returnBody == "true" {
		addHeader(w, r)
		addBody(w, r)
	}
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		printVerbose(r.URL.String())
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		fmt.Println(err)

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	userIP := net.ParseIP(ip)

	if err != nil {
		fmt.Println(err)

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var body string

	format := r.FormValue("format")

	if format != "" {
		if format == "json" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			body = fmt.Sprintf("{\"ip\":\"%s\"}", userIP)
		} else {
			body = fmt.Sprintf("IP: %s", userIP)
		}
	} else {
		body = fmt.Sprintf("IP: %s", userIP)
	}

	fmt.Fprintf(w, "%s", body)
}

func userAgentHandler(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		printVerbose(r.URL.String())
	}

	userAgentString := r.Header.Get("User-Agent")

	var body string

	format := r.FormValue("format")

	if format != "" {
		if format == "json" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			body = fmt.Sprintf("{\"userAgent\":\"%s\"}", userAgentString)
		} else {
			body = fmt.Sprintf("UserAgent: %s", userAgentString)
		}
	} else {
		body = fmt.Sprintf("UserAgent: %s", userAgentString)
	}

	fmt.Fprintf(w, "%s", body)
}

func basicAuthHandler(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		printVerbose(r.URL.String())
	}

	path := strings.Split(r.URL.Path, "/")

	if path[2] == "" || path[3] == "" {
		http.Error(w, "URL must be /basicAuth/[UserName]/[Password]", http.StatusBadRequest)
		return
	}

	userName := path[2]
	password := path[3]

	auth := strings.Split(r.Header.Get("Authorization"), " ")

	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "Authorization Method must be Basic", http.StatusUnauthorized)
		return
	}

	authPayload, err := base64.StdEncoding.DecodeString(auth[1])

	if err != nil {
		fmt.Println(err)

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	authPair := strings.Split(string(authPayload), ":")

	if len(authPair) != 2 || authPair[0] != userName || authPair[1] != password {
		http.Error(w, "UserName or Password is not correct", http.StatusUnauthorized)
		return
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if *verbose {
		printVerbose(r.URL.String())
	}

	path := strings.Split(r.URL.Path, "/")

	if path[2] == "" {
		http.Error(w, "URL must be /file/[FileName]", http.StatusBadRequest)
		return
	}

	fileData, err := ioutil.ReadFile(path[2])

	if err != nil {
		fmt.Println(err)

		http.Error(w, "Could not Open File", http.StatusInternalServerError)
		return
	}

	w.Write(fileData)
}

func printVerbose(url string) {
	currentTime := time.Now()

	fmt.Printf("%s %s\n", currentTime.Format("2006.01.02 15:04:05"), url)
}

func addHeader(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
}

func addBody(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
	}

	w.Write(body)
}
