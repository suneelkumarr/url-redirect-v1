package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	Id           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

// generateShortURL generates a short URL for the provided original URL.
// It uses the MD5 hashing algorithm to create a unique, shortened version of the original URL.
// The first 8 characters of the hashed data are returned as the short URL.
//
// Parameters:
//   - OriginalURL (string) - The original URL to be shortened.
//
// Returns:
//   - string - A shortened version of the original URL.
func genrateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))
	data := hasher.Sum(nil)
	hash := hex.EncodeToString(data)
	return hash[:8]
}

// createUrl creates a short URL for the provided original URL and stores the mapping of the short URL to the original URL in the urlDB map.
// It uses the genrateShortURL function to generate a shortened version of the original URL.
// The function then returns the short URL as a string.
//
// Parameters:
//   - OriginalURL (string) - The original URL to be shortened.
//
// Returns:
//   - string - A shortened version of the original URL.
func createUrl(OriginalURL string) string {
	shorturl := genrateShortURL(OriginalURL)
	id := shorturl
	urlDB[id] = URL{Id: id, OriginalURL: OriginalURL, ShortURL: shorturl, CreationDate: time.Now()}
	return shorturl
}

func getUrl(id string) (URL, error) {

	if url, ok := urlDB[id]; ok {
		return url, nil
	} else {
		return URL{}, fmt.Errorf("url not found")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Method GET: ", r.Method)
	fmt.Fprintf(w, "Hello World")
}

/**
 * @summary Creates a short URL for the provided original URL and returns it as a JSON response.
 * @description This function takes an original URL as input, generates a short URL using a hashing algorithm,
 *  and stores the mapping of the short URL to the original URL in the urlDB map.
 *  It then returns the short URL as a JSON response with the "short_url" key.
 * @param w http.ResponseWriter - The HTTP response writer to write the JSON response.
 * @param r *http.Request - The HTTP request containing the original URL in the JSON body.
 * @return http.ResponseWriter - The HTTP response writer with the JSON response written to it.
 * @return error - An error if there is an issue decoding the JSON body or creating the short URL.
 */
func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shorturl := createUrl(data.URL)
	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shorturl}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

/**
 * redirectUrlHandler redirects the user to the original URL based on the provided short URL id.
 *
 * @param w http.ResponseWriter - The HTTP response writer to write the redirect response.
 * @param r *http.Request - The HTTP request containing the short URL id in the path.
 */
func redirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getUrl(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	fmt.Println("Url shortners...")
	OriginalURL := "https://www.google.com"
	genrateShortURL(OriginalURL)

	http.HandleFunc("/", handler)
	http.HandleFunc("/shorter", shortUrlHandler)
	http.HandleFunc("/redirect/", redirectUrlHandler)
	fmt.Println("Server is Running on the port :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error on Starting the server", err)
	}
}
