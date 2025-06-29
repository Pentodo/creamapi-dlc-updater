package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"gopkg.in/ini.v1"
)

// AppDetailsResponse matches the structure of the JSON response from /api/appdetails
type AppDetailsResponse map[string]struct {
	Success bool `json:"success"`
	Data    struct {
		DLC []int `json:"dlc"`
	} `json:"data"`
}

// AppListResponse matches the structure of the JSON response from /ISteamApps/GetAppList/v2/
type AppListResponse struct {
	Applist struct {
		Apps []struct {
			Appid int    `json:"appid"`
			Name  string `json:"name"`
		} `json:"apps"`
	} `json:"applist"`
}

// setupLogging configures the logger to write to a file and the console.
func setupLogging(logPath string) {
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Create a multi-writer to log to both file and console
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}

func main() {
	// Get the executable's directory
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	// Set up logging
	logPath := filepath.Join(exeDir, "creamapi-dlc-updater.log")
	setupLogging(logPath)

	// Set paths
	iniPath := filepath.Join(exeDir, "cream_api.ini")
	log.Printf("[INFO] INI path: %s", iniPath)

	// Load the INI file
	cfg, err := ini.Load(iniPath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read INI file: %v", err)
	}

	// Read the appid
	appid := cfg.Section("steam").Key("appid").String()
	if appid == "" {
		log.Fatal("[ERROR] Could not read appid from INI")
	}
	log.Printf("[INFO] Read appid: %s", appid)

	// --- Step 1: Get DLC AppIDs ---
	appDetailsURL := fmt.Sprintf("https://store.steampowered.com/api/appdetails?appids=%s", appid)
	log.Printf("[INFO] Requesting appdetails: %s", appDetailsURL)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(appDetailsURL)
	if err != nil {
		log.Fatalf("[ERROR] Failed to get app details: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read response body: %v", err)
	}
	log.Printf("[INFO] Received appdetails JSON, size: %d", len(body))

	var appDetails AppDetailsResponse
	if err := json.Unmarshal(body, &appDetails); err != nil {
		log.Fatalf("[ERROR] Failed to parse appdetails JSON: %v", err)
	}

	if details, ok := appDetails[appid]; !ok || !details.Success {
		log.Fatal("[ERROR] Appdetails request was not successful")
	}

	dlcIDs := appDetails[appid].Data.DLC
	if len(dlcIDs) == 0 {
		log.Println("[WARN] No DLCs found for this appid")
		log.Println("[INFO] Process finished")
		return
	}
	log.Printf("[INFO] Found %d DLC appids", len(dlcIDs))

	// --- Step 2: Get DLC names from the full app list ---
	appListURL := "https://api.steampowered.com/ISteamApps/GetAppList/v2/"
	log.Println("[INFO] Requesting app list")

	resp, err = client.Get(appListURL)
	if err != nil {
		log.Fatalf("[ERROR] Failed to get app list: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read app list response body: %v", err)
	}
	log.Printf("[INFO] Received app list JSON, size: %d", len(body))

	var appList AppListResponse
	if err := json.Unmarshal(body, &appList); err != nil {
		log.Fatalf("[ERROR] Failed to parse app list JSON: %v", err)
	}

	// Create a map for quick lookup of DLC names
	dlcIDSet := make(map[int]bool)
	for _, id := range dlcIDs {
		dlcIDSet[id] = true
	}
	dlcNames := make(map[int]string)
	for _, app := range appList.Applist.Apps {
		if _, found := dlcIDSet[app.Appid]; found {
			dlcNames[app.Appid] = app.Name
		}
	}

	// --- Step 3: Write DLCs to the INI file ---
	log.Println("[INFO] Writing DLC section to INI")

	// Delete the old section to ensure it's clean
	cfg.DeleteSection("dlc")
	dlcSection, err := cfg.NewSection("dlc")
	if err != nil {
		log.Fatalf("[ERROR] Failed to create new DLC section: %v", err)
	}

	// Sort DLCs by ID for consistent output
	sort.Ints(dlcIDs)

	for _, id := range dlcIDs {
		if name, found := dlcNames[id]; found {
			dlcSection.Key(strconv.Itoa(id)).SetValue(name)
			log.Printf("[INFO] DLC %d = %s", id, name)
		}
	}

	// Save changes back to the INI file
	if err := cfg.SaveTo(iniPath); err != nil {
		log.Fatalf("[ERROR] Failed to save INI file: %v", err)
	}

	log.Println("[INFO] Process finished successfully")
}
