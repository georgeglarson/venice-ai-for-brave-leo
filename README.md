# Leo Venice.AI Configuration Tool

> #BraveAI #VeniceAI #LeoAI #AIHelper #LLM #LLaMA3 #AIConfiguration #BraveExtension

**📱 Website: [georgeglarson.github.io/venice-ai-for-brave-leo](https://georgeglarson.github.io/venice-ai-for-brave-leo)**

A universal tool to seamlessly configure Leo AI helper in Brave browser to use Venice.AI's powerful LLaMA 3.3 70B model. This tool automates the process of setting up your Brave browser to use Venice.AI as your default AI assistant.

[![Website Preview](https://img.shields.io/badge/View_Website-blue?style=for-the-badge&logo=github)](https://georgeglarson.github.io/venice-ai-for-brave-leo)
[![GitHub Repo](https://img.shields.io/badge/GitHub_Repo-black?style=for-the-badge&logo=github)](https://github.com/georgeglarson/venice-ai-for-brave-leo)

## ⚡ Quickstart

1. **Download the executable** for your platform:
   - [Windows (leo_venice_config.exe)](https://github.com/georgeglarson/venice-ai-for-brave-leo/releases/latest/download/leo_venice_config.exe)
   - [macOS (leo_venice_config_mac)](https://github.com/georgeglarson/venice-ai-for-brave-leo/releases/latest/download/leo_venice_config_mac)
   - [Linux (leo_venice_config_linux)](https://github.com/georgeglarson/venice-ai-for-brave-leo/releases/latest/download/leo_venice_config_linux)

2. **Run the executable**:
   - **Windows**: Double-click the downloaded file or run it from the command line
   - **macOS**: Open Terminal, navigate to the download location, and run:
     ```bash
     chmod +x ./leo_venice_config_mac && ./leo_venice_config_mac
     ```
   - **Linux**: Open Terminal, navigate to the download location, and run:
     ```bash
     chmod +x ./leo_venice_config_linux && ./leo_venice_config_linux
     ```
   - The tool will prompt you for your Venice.AI API key

3. **Completely Close and Restart Brave browser**:
   - **IMPORTANT**: You must *completely* close Brave browser (including any background processes) and restart it for changes to take effect
   - On Windows, you may need to check Task Manager and end any Brave processes before restarting
   - On macOS, right-click the Brave icon in the dock and select "Quit"
   - On Linux, ensure all Brave processes are terminated with `pkill -9 -i brave`

That's it! Your Brave browser is now configured to use Venice.AI with Leo AI helper.

## 🌐 Website

Visit our [project website](https://georgeglarson.github.io/venice-ai-for-brave-leo) for:
- Easy downloads for all platforms
- Interactive quickstart guide
- Detailed documentation
- Contact information

The website provides a user-friendly interface for downloading and learning about the tool.

## 🚀 Features

- **One-Click Configuration**: Automatically configures Leo AI helper in Brave browser
- **Cross-Platform Support**: Works on Windows, macOS, and Linux
- **Backup Creation**: Creates a backup of your existing configuration
- **User-Friendly Interface**: Simple GUI on Windows, CLI on other platforms
- **Secure**: Your API key is stored only in your local Brave configuration

## 🛠️ Installation

### Prerequisites

- [Brave Browser](https://brave.com/download/) with Leo AI helper
- [Venice.AI API Key](https://venice.ai/settings/api)

### Option 1: Download Pre-built Executable (Recommended)

Simply download the executable for your platform from the [Releases page](https://github.com/georgeglarson/venice-ai-for-brave-leo/releases/latest).

### Option 2: Building from Source (Advanced)

If you prefer to build from source, you'll need:
- [Go](https://golang.org/dl/) 1.22 or higher (project uses Go 1.22.2)

Then follow these steps:

1. Clone this repository:
   ```bash
   git clone https://github.com/georgeglarson/venice-ai-for-brave-leo.git
   cd venice-ai-for-brave-leo
   ```

2. Initialize the Go module and get dependencies:
   ```bash
   go get github.com/google/uuid
   ```

3. Build for your platform using the included build scripts:

   **Linux/macOS:**
   ```bash
   ./build.sh
   ```

   **Windows:**
   ```bash
   build.bat
   ```

   These scripts will build executables for all supported platforms (Windows, macOS, and Linux) and place them in the `build` directory.

   **Manual Build (if needed):**
   ```bash
   # For Windows
   GOOS=windows GOARCH=amd64 go build -o leo_venice_config.exe main.go types.go ui.go browser.go preferences.go messagebox_windows.go

   # For macOS
   GOOS=darwin GOARCH=amd64 go build -o leo_venice_config_mac main.go types.go ui.go browser.go preferences.go messagebox_stub.go

   # For Linux
   GOOS=linux GOARCH=amd64 go build -o leo_venice_config_linux main.go types.go ui.go browser.go preferences.go messagebox_stub.go
   ```

## 📋 Usage

### Easiest Method (Windows)

Simply **double-click** the `leo_venice_config.exe` file. A dialog will appear prompting you for your Venice.AI API key.

### Command Line Method

You can also run the tool from the command line:

**With API key as argument:**

**Windows:**
```bash
leo_venice_config.exe -key=YOUR_API_KEY
```

**macOS/Linux:**
```bash
./leo_venice_config -key=YOUR_API_KEY
```

**With interactive prompt:**

**Windows:**
```bash
leo_venice_config.exe
```

**macOS/Linux:**
```bash
./leo_venice_config
```

## ⚙️ Configuration Details

The tool will automatically:

1. Find Brave's Preferences file on your system
2. Create a backup of your current configuration
3. Configure Leo to use Venice.AI with:
   - Label: Venice.AI
   - Model: llama-3.3-70b
   - Endpoint: https://api.venice.ai/api/v1/chat/completions
   - Context size: 4000
4. Set Venice.AI as the default model for Leo AI helper

## 🔒 Security

- Your API key is stored only in your local Brave browser configuration
- The tool creates a backup of your configuration before making any changes
- No data is sent to any server except when you use Leo AI helper with Venice.AI

### Windows SmartScreen Warning

When running the executable on Windows, you may see a "Windows protected your PC" message from SmartScreen. This happens because the executable isn't digitally signed with a certificate.

To run the application:
1. Click "More info" on the warning dialog
2. Click "Run anyway" to proceed

This is a standard security feature in Windows for applications downloaded from the internet. The source code is fully open and available for inspection on GitHub.

## 🤝 Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## 🚀 Development and Release Process

### Development

The codebase is available in two formats:

1. **Modular Structure** (Recommended for development):
   - **types.go**: Type definitions and constants
   - **ui.go**: UI/dialog functions
   - **browser.go**: Browser-related functions
   - **preferences.go**: Preferences file manipulation
   - **main.go**: Main program logic
   - **messagebox_windows.go**: Windows-specific MessageBox implementation
   - **messagebox_stub.go**: Stub implementation for non-Windows platforms

2. **Monolithic File**:
   - **leo_venice_config.go**: A single file containing all functionality, provided for reference and compatibility with older build processes

The modular structure is recommended for development as it provides better organization and maintainability. The build scripts use the modular files to create the executables.

### Building Locally

To build the application locally, use the provided build scripts:

**Linux/macOS:**
```bash
./build.sh
```

**Windows:**
```batch
build.bat
```

### Creating a Release

To create a new release with automatic binary builds:

1. Make sure all your changes are committed and pushed to the repository

2. Run the release script with a version number:

   **Linux/macOS:**
   ```bash
   ./release.sh v1.0.2
   ```

   **Windows:**
   ```batch
   release.bat v1.0.2
   ```

3. The script will create and push a git tag, which will trigger the GitHub Actions workflow
   
4. GitHub Actions will automatically build the binaries for all platforms and attach them to the release

5. You can check the status of the build at: https://github.com/georgeglarson/venice-ai-for-brave-leo/actions

## � Contact

For questions or support, contact:
- Email: george.g.larson@gmail.com
- Twitter: [@g3ologic](https://twitter.com/g3ologic)
- GitHub: [georgeglarson](https://github.com/georgeglarson)

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

#AITools #BraveExtensions #VeniceAI #LeoAI #AIConfiguration #LLaMA #OpenSource #CrossPlatform #GoLang