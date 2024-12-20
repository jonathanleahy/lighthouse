package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	queue         []string
	mu            sync.Mutex
	channels             = make(map[string]chan string)
	urlMap               = make(map[string][]string)
	urlCounter    uint64 = 0
	submissionLog *os.File
	processingLog *os.File
	allowedOrigin = "http://localhost:3000"
)

type Job struct {
	url     string
	hashtag string
}

func init() {
	if origin := os.Getenv("ALLOWED_ORIGIN"); origin != "" {
		allowedOrigin = origin
	}

	var err error
	submissionLog, err = os.OpenFile("submissions.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open submissions log: %v", err)
	}

	processingLog, err = os.OpenFile("processing.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open processing log: %v", err)
	}
}

func main() {
	defer submissionLog.Close()
	defer processingLog.Close()

	rand.Seed(time.Now().UnixNano())

	log.Println("Starting the web server on :8080")
	log.Printf("CORS enabled for origin: %s", allowedOrigin)

	http.HandleFunc("/", serveHomePage)
	http.HandleFunc("/submit", handleURL)
	http.HandleFunc("/queue-length", getQueueLength)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	if r.Method == "OPTIONS" {
		return
	}

	html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>URL Submitter</title>
        <script>
            function submitURLJS() {
                var url = document.getElementById("url").value;
                var hashtag = document.getElementById("hashtag").value;
                fetch('/submit?url=' + encodeURIComponent(url) + '&hashtag=' + encodeURIComponent(hashtag))
                    .then(response => response.text())
                    .then(data => {
                        document.getElementById("response").innerText = data;
                        setTimeout(() => {
                            updateQueueLength();
                            setInterval(updateQueueLength, 5000);
                        }, 500);
                    });
            }

            function updateQueueLength() {
                fetch('/queue-length')
                    .then(response => response.json())
                    .then(data => {
                        var queueInfo = "";
                        for (var hashtag in data.channels) {
                            queueInfo += "<h2>Hashtag: " + hashtag + "</h2>";
                            queueInfo += "<p>Items: " + data.channels[hashtag].join(", ") + "</p>";
                        }
                        document.getElementById("queue-info").innerHTML = queueInfo;
                    });
            }

            setInterval(updateQueueLength, 5000);
        </script>
    </head>
    <body onload="updateQueueLength()">
        <h1>URL Submitter</h1>
        <input type="text" id="url" placeholder="Enter URL" />
        <input type="text" id="hashtag" placeholder="Enter Hashtag" />
        <button onclick="submitURLJS()">Submit</button>
        <p id="response"></p>
        <div id="queue-info"></div>
    </body>
    </html>
    `
	fmt.Fprint(w, html)
}

func generateSequentialURL() string {
	count := atomic.AddUint64(&urlCounter, 1)
	return fmt.Sprintf("http://example.com/url_%d", count)
}

func handleURL(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	if r.Method == "OPTIONS" {
		return
	}

	url := r.URL.Query().Get("url")
	hashtag := r.URL.Query().Get("hashtag")
	log.Printf("Received URL: %s, Hashtag: %s", url, hashtag)

	if url == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}
	if hashtag == "" {
		http.Error(w, "Hashtag parameter is missing", http.StatusBadRequest)
		return
	}

	// Check for mock case first
	if url == "mock" && hashtag == "mock" {
		log.Println("Generating random URLs distributed across channels")
		go generateRandomHashtagsAndSubmit()
		fmt.Fprintf(w, "Generating mock data")
		return
	}

	mu.Lock()
	// Ensure channel exists
	if _, exists := channels[hashtag]; !exists {
		log.Printf("Creating new channel for hashtag: %s", hashtag)
		channels[hashtag] = make(chan string, 10000)
		go processHashtagChannel(hashtag, channels[hashtag])
	}

	// Initialize or update urlMap
	if _, exists := urlMap[hashtag]; !exists {
		urlMap[hashtag] = make([]string, 0)
	}
	urlMap[hashtag] = append(urlMap[hashtag], url)
	queue = append(queue, url)

	// Get channel reference before unlocking
	ch := channels[hashtag]
	mu.Unlock()

	// Send to hashtag channel
	go func() {
		select {
		case ch <- url:
			log.Printf("Successfully sent URL to hashtag channel: %s", hashtag)
		case <-time.After(100 * time.Millisecond):
			log.Printf("Channel %s is full, retrying in background", hashtag)
			ch <- url // This will block until space is available
		}
	}()

	fmt.Fprintf(w, "URL received: %s with hashtag: %s", url, hashtag)
}

func getQueueLength(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	if r.Method == "OPTIONS" {
		return
	}

	mu.Lock()
	channelsCopy := make(map[string][]string)
	for hashtag, urls := range urlMap {
		if len(urls) > 0 {
			channelsCopy[hashtag] = append([]string{}, urls...)
		}
	}
	mu.Unlock()

	response := struct {
		Channels map[string][]string `json:"channels"`
	}{
		Channels: channelsCopy,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func processHashtagChannel(hashtag string, ch chan string) {
	var delaySeconds time.Duration = 5
	var numWorkers int = 1

	if strings.HasPrefix(hashtag, "channel_") {
		if channelNum, err := strconv.Atoi(strings.TrimPrefix(hashtag, "channel_")); err == nil {
			delaySeconds = time.Duration(channelNum)
			numWorkers = channelNum + 1
			log.Printf("Channel %s will process with %d second delay using %d concurrent workers",
				hashtag, channelNum, numWorkers)
		}
	}

	// Create a buffered channel for jobs
	jobs := make(chan Job, 10000)

	// Start the workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for job := range jobs {
				processJob(workerId, job, delaySeconds)
			}
		}(i)
	}

	// Process incoming URLs
	for url := range ch {
		log.Printf("Received new URL in channel %s: %s", hashtag, url)
		jobs <- Job{url: url, hashtag: hashtag}
	}

	// Cleanup
	close(jobs)
	wg.Wait()
}

func processJob(workerId int, job Job, delaySeconds time.Duration) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	log.Printf("Worker %d starting to process URL: %s for hashtag: %s",
		workerId, job.url, job.hashtag)

	// Log start of processing
	fmt.Fprintf(processingLog, "%s - PROCESSING - Worker: %d, URL: %s, Hashtag: %s\n",
		timestamp, workerId, job.url, job.hashtag)

	minDelay := 1 * time.Second
	if delaySeconds < minDelay {
		time.Sleep(minDelay)
	} else {
		time.Sleep(delaySeconds * time.Second)
	}

	timestamp = time.Now().Format("2006-01-02 15:04:05.000")
	fmt.Fprintf(processingLog, "%s - COMPLETED - Worker: %d, URL: %s, Hashtag: %s\n",
		timestamp, workerId, job.url, job.hashtag)

	log.Printf("Worker %d completed URL: %s for hashtag: %s\n",
		workerId, job.url, job.hashtag)

	mu.Lock()
	if urls, exists := urlMap[job.hashtag]; exists {
		for i, u := range urls {
			if u == job.url {
				urlMap[job.hashtag] = append(urls[:i], urls[i+1:]...)
				break
			}
		}
		// Only remove from urlMap if empty, but keep the channel
		if len(urlMap[job.hashtag]) == 0 {
			delete(urlMap, job.hashtag)
			log.Printf("Removed empty urlMap entry for hashtag: %s", job.hashtag)
		}
	}
	mu.Unlock()
}

func submitURL(urlStr, hashtag string) {
	maxRetries := 3
	baseDelay := 50 * time.Millisecond
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	fmt.Fprintf(submissionLog, "%s - SUBMITTED - URL: %s, Hashtag: %s\n",
		timestamp, urlStr, hashtag)

	for retry := 0; retry < maxRetries; retry++ {
		if retry > 0 {
			delay := baseDelay * time.Duration(1<<uint(retry))
			time.Sleep(delay)
		}

		log.Printf("Submitting URL: %s with hashtag: %s (attempt %d)", urlStr, hashtag, retry+1)
		submitURL := fmt.Sprintf("http://localhost:8080/submit?url=%s&hashtag=%s",
			url.QueryEscape(urlStr), url.QueryEscape(hashtag))

		resp, err := http.Get(submitURL)
		if err != nil {
			log.Printf("Failed to submit URL: %s with hashtag: %s, error: %v",
				urlStr, hashtag, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to submit URL: %s with hashtag: %s, status code: %d",
				urlStr, hashtag, resp.StatusCode)
			continue
		}

		return
	}
	log.Printf("Failed to submit URL after %d attempts: %s with hashtag: %s",
		maxRetries, urlStr, hashtag)
}

func generateRandomHashtagsAndSubmit() {
	numChannels := rand.Intn(5) + 1
	channelNames := make([]string, numChannels+1)

	// Add channel_0 first
	channelNames[0] = "channel_0"

	// Generate remaining channel names
	for i := 1; i <= numChannels; i++ {
		channelNames[i] = fmt.Sprintf("channel_%d", i)
	}

	log.Printf("Created %d channels with different processing concurrency:", len(channelNames))
	for _, name := range channelNames {
		num := strings.TrimPrefix(name, "channel_")
		workers, _ := strconv.Atoi(num)
		log.Printf("- %s will process with %d concurrent workers", name, workers+1)

		mu.Lock()
		if _, exists := channels[name]; !exists {
			channels[name] = make(chan string, 10000)
			go processHashtagChannel(name, channels[name])
		}
		if _, exists := urlMap[name]; !exists {
			urlMap[name] = make([]string, 0)
		}
		mu.Unlock()
	}

	// Create a rate limiter for submissions
	rateLimiter := time.NewTicker(100 * time.Millisecond)
	defer rateLimiter.Stop()

	// Submit URLs to channel_0
	log.Printf("Submitting URLs to channel_0")
	for i := 0; i < 50; i++ {
		<-rateLimiter.C
		url := generateSequentialURL()
		submitURL(url, "channel_0")
	}

	// Generate and submit URLs to other channels
	log.Printf("Submitting URLs distributed across other channels")
	for i := 0; i < 50; i++ {
		<-rateLimiter.C
		randomChannel := channelNames[rand.Intn(numChannels)+1]
		url := generateSequentialURL()
		submitURL(url, randomChannel)
	}
}
