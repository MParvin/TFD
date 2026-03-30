# TFD (Telegram File Downloader)

[![CI](https://github.com/mparvin/tfd/actions/workflows/ci.yml/badge.svg)](https://github.com/mparvin/tfd/actions/workflows/ci.yml)
[![Release](https://github.com/mparvin/tfd/actions/workflows/release.yml/badge.svg)](https://github.com/mparvin/tfd/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mparvin/tfd)](https://goreportcard.com/report/github.com/mparvin/tfd)

TFD is a cross-platform CLI tool that runs a Telegram bot to automatically download and organize files sent to it. It supports photos, videos, documents, audio files, voice messages, and URLs, saving them to organized local directories.

## ✨ Features

- **File Type Support**: Photos, videos, documents, audio, voice messages, and URLs
- **Organized Storage**: Automatically saves files to categorized directories
- **User Authorization**: Whitelist-based access control for bot security
- **Proxy Support**: HTTP proxy support for restricted networks
- **Error Handling**: Robust error handling with admin notifications via Telegram
- **Cross-Platform**: Runs on Linux, macOS, and Windows
- **Path Security**: Sanitizes filenames to prevent path traversal attacks

## 🚀 Installation

### From GitHub Releases (Recommended)

1. Download the latest binary for your platform from [releases](https://github.com/mparvin/tfd/releases)
2. Make it executable (Linux/macOS):
   ```bash
   chmod +x tfd-*
   ```
3. Optionally, move to your PATH:
   ```bash
   sudo mv tfd-* /usr/local/bin/tfd
   ```

### Using Go Install

```bash
go install github.com/mparvin/tfd@latest
```

### From Source

```bash
git clone https://github.com/mparvin/tfd.git
cd tfd
go build -o tfd .
```

## ⚙️ Configuration

### 1. Create Configuration File

Copy the example configuration and edit it:

```bash
cp config.example.yaml ~/.tfd.yaml
nano ~/.tfd.yaml
```

### 2. Configuration Options

```yaml
# Telegram Bot Token (get from @BotFather)
token: "YOUR_BOT_TOKEN_HERE"

# Admin Chat ID (for error notifications)
admin_chat_id: YOUR_CHAT_ID

# List of allowed user Chat IDs
allowed_users:
  - CHAT_ID_1
  - CHAT_ID_2

# Support user username (for unauthorized user messages)
support_user: "YOUR_USERNAME"

# Directory paths for file storage (must end with /)
directories:
  photo: "/var/lib/tfd/Pictures/"
  video: "/var/lib/tfd/Videos/"
  music: "/var/lib/tfd/Music/"
  voice: "/var/lib/tfd/Voices/"
  document: "/var/lib/tfd/Documents/"
  other: "/var/lib/tfd/"

# Optional: Proxy settings (for regions where Telegram is blocked)
proxy: "http://proxy.example.com:8080"
```

### 3. Getting Bot Token and Chat IDs

1. **Bot Token**: Message [@BotFather](https://t.me/BotFather) on Telegram:
   - Send `/newbot`
   - Follow instructions to create your bot
   - Copy the token

2. **Your Chat ID**: Message [@userinfobot](https://t.me/userinfobot) to get your Chat ID

3. **Group Chat ID**: Add [@userinfobot](https://t.me/userinfobot) to a group to get the group Chat ID

## 🎯 Usage

### Basic Usage

```bash
./tfd
```

### With Custom Config File

```bash
./tfd --config /path/to/config.yaml
```

### With Proxy

```bash
./tfd --proxy http://proxy.example.com:8080
```

### Show Help

```bash
./tfd --help
```

## 🌐 Proxy Configuration

TFD supports HTTP and SOCKS5 proxies for regions where Telegram is blocked:
proxies for regions where Telegram is blocked:

### HTTP Proxy
```yaml
proxy: "http://proxy.example.com:8080"
```

### Via Command Line
```bash
./tfd --proxy http://proxy.example.com:8

## 📂 File Organization

TFD automatically organizes downloaded files:

- **Photos**: `directories.photo` (e.g., `/var/lib/tfd/Pictures/`)
- **Videos**: `directories.video` (e.g., `/var/lib/tfd/Videos/`)
- **Audio**: `directories.music` (e.g., `/var/lib/tfd/Music/`)
- **Voice**: `directories.voice` (e.g., `/var/lib/tfd/Voices/`)
- **Documents**: `directories.document` (e.g., `/var/lib/tfd/Documents/`)
- **URLs**: `directories.other` (e.g., `/var/lib/tfd/`)

Files are automatically timestamped: `photo_20240330_143022_image.jpg`

## 🔒 Security Features

- **User Whitelisting**: Only authorized users can interact with the bot
- **Path Traversal Protection**: Filenames are sanitized to prevent directory traversal
- **Admin Notifications**: Errors are automatically sent to the admin via Telegram
- **Token Security**: Configuration files are git-ignored to prevent token leaks

## 🛠️ Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes and test them
4. Commit your changes: `git commit -m 'Add amazing feature'`
5. Push to the branch: `git push origin feature/amazing-feature`
6. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/mparvin/tfd.git
cd tfd
go mod download
go build .
```

### Running Tests

```bash
go test ./...
go vet ./...
```

## 📋 Requirements

- Go 1.22+ (for building from source)
- Telegram Bot Token
- Network access to Telegram API (direct or via proxy)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/mparvin/tfd/issues)
- 💡 **Feature Requests**: [GitHub Issues](https://github.com/mparvin/tfd/issues)
- 📧 **General Questions**: Create a [GitHub Discussion](https://github.com/mparvin/tfd/discussions)

## ⭐ Star History

If this project helped you, please consider giving it a star! ⭐

---

**Made with ❤️ by [Mohammad Parvin](https://github.com/mparvin)**

