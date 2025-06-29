# 🎮 creamapi-dlc-updater

A standalone Go application that automatically updates the `[dlc]` section of a `cream_api.ini` file for any Steam game by fetching the latest DLC information from the Steam API.

## ✨ Features

- ✅ Reads the `appid` from the `[steam]` section of `cream_api.ini`.
- ✅ Fetches DLC appids from the Steam Store API.
- ✅ Fetches DLC names from the Steam AppList API.
- ✅ Clears the old `[dlc]` section in `cream_api.ini`.
- ✅ Populates the `[dlc]` section with a sorted list of all current DLCs and their names.
- ✅ Creates a detailed log file (`creamapi-dlc-updater.log`) for each run.

## ⚙️ How it works

The application is a command-line executable. When you run it, it performs the following steps:

1.  Locates the `cream_api.ini` file in its own directory.
2.  Reads the base game's `appid`.
3.  Connects to the Steam APIs to get all official DLCs for that game.
4.  Updates the `cream_api.ini` file with the fetched data.

## 📋 Requirements

- 🖥️ Windows (the main target for CreamAPI).
- 🌐 An internet connection to reach the Steam API.
- 🐹 Go (the 1.21 version is recommended for building from source).

## 🛠️ Compilation

1.  Make sure you have [Go installed](https://go.dev/doc/install) and configured.
2.  Open a terminal or command prompt in the project's root directory.
3.  Initialize the Go module and fetch dependencies (only needs to be done once):
    ```sh
    go mod tidy
    ```
4.  Build the executable:
    ```sh
    go build
    ```
5.  This will create an executable file (e.g., `creamapi-dlc-updater.exe`) in the directory.

## 🚀 Usage

1.  Place the compiled executable (e.g., `creamapi-dlc-updater.exe`) in the same folder as the game's main executable and your `cream_api.ini` file.
2.  Run `creamapi-dlc-updater.exe`.
3.  The application will automatically update the `[dlc]` section in `cream_api.ini`. Check the `creamapi-dlc-updater.log` file for a detailed report of the operations.

## 🐛 Troubleshooting

- **"Failed to read INI file"**: Make sure `cream_api.ini` is in the same directory as the executable.
- **"Could not read appid from INI"**: Ensure your `cream_api.ini` has a `[steam]` section with a valid `appid` key (e.g., `appid = 1244460`).
- **Network Errors**: Check your internet connection and firewall settings. The application needs to access `store.steampowered.com` and `api.steampowered.com`.
- For any other issues, please check the `creamapi-dlc-updater.log` file for detailed error messages.

## 📄 License

This project is licensed under the MIT License.
