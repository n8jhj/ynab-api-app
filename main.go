package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const app_name = "ynab-api-app"
const cache_file_name = "api_key.txt"

func main() {
	fmt.Println("This is the YNAB API app!\n")

	// Get API key from user.
	var api_key string
	fmt.Println("Enter your API key:")
	fmt.Scanln(&api_key)

	// Store API key on local machine.
	// TODO: Look at https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	user_cache_dir, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("Can't find cache directory.")
		os.Exit(1)
	}
	app_cache_dir := filepath.Join(user_cache_dir, app_name)
	err = os.MkdirAll(app_cache_dir, 0700)
	if err != nil {
		fmt.Println("Can't make cache directory.")
		os.Exit(1)
	}
	cache_file_path := filepath.Join(app_cache_dir, cache_file_name)
	fmt.Printf("Cache file: %s\n", cache_file_path)

	file_contents := []byte(api_key)
	os.WriteFile(cache_file_path, file_contents, 0700)

	// Get API data.
	client := http.DefaultClient
	req, err := http.NewRequest("GET", "https://api.youneedabudget.com/v1/budgets", nil)
	if err != nil {
		fmt.Println("Request formulation failed.")
		os.Exit(1)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_key))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request failed.")
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Bad HTTP response.")
		os.Exit(1)
	}
	body_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body.")
		os.Exit(1)
	}
	body_string := string(body_bytes)
	fmt.Printf("BODY BYTES:\n%s\n", body_string)
}
