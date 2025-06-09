# 🧹 Clearance

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square&logo=opensourceinitiative&logoColor=white)](LICENSE)
[![Windows](https://img.shields.io/badge/Windows-0078D6?style=flat-square&logo=windows&logoColor=white)](https://www.microsoft.com/windows)
[![PowerShell](https://img.shields.io/badge/PowerShell-5391FE?style=flat-square&logo=powershell&logoColor=white)](https://docs.microsoft.com/powershell)
[![GoReleaser](https://img.shields.io/badge/GoReleaser-00ADD8?style=flat-square&logo=go&logoColor=white)](https://goreleaser.com)
[![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white)](https://www.docker.com)
[![npm](https://img.shields.io/badge/npm-CB3837?style=flat-square&logo=npm&logoColor=white)](https://www.npmjs.com)
[![Yarn](https://img.shields.io/badge/Yarn-2C8EBB?style=flat-square&logo=yarn&logoColor=white)](https://yarnpkg.com)
[![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat-square&logo=github-actions&logoColor=white)](https://github.com/features/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/abdorrahmani/clearance?style=flat-square)](https://goreportcard.com/report/github.com/abdorrahmani/clearance)

A lightweight CLI tool to clean up development caches and free up disk space. Built with Go, Clearance helps you maintain a clean development environment by removing unnecessary cache files.

## 🚀 Features

- 🧹 Clean npm cache
- 🧶 Clean yarn cache
- 🐳 Clean Docker cache (optional)
- 🪟 Clean Windows WinSxS temp files
- 🔒 Safe, selectable cleanup operations
- 🎨 Beautiful CLI interface with color support
- ⚡ Fast and efficient cleanup operations

## 📋 Requirements

- Go 1.21 or higher
- Administrator privileges (required for system folders and Docker commands)
- npm, yarn, and/or Docker in PATH (if using respective cleanup options)

## 🛠️ Installation

### Using Go Install
```bash
go install github.com/abdorrahmani/clearance@latest
```

### Using Binary Release
1. Download the latest release from the [Releases page](https://github.com/abdorrahmani/clearance/releases)
2. Extract the ZIP file
3. Run `run-clearance.bat` or `clearance.exe` directly

## 💻 Usage

### Command Line Interface
```bash
clearance [--npm] [--yarn] [--docker] [--winsxs] [--all]
```

### Interactive Mode
Simply run:
```bash
clearance
```

### Examples

```bash
# Clean npm and yarn caches
clearance --npm --yarn

# Clean everything
clearance --all

# Clean only Docker cache
clearance --docker
```

## 🔧 Options

| Flag      | Description                    |
|-----------|--------------------------------|
| `--npm`   | Clean npm cache               |
| `--yarn`  | Clean yarn cache              |
| `--docker`| Clean Docker cache            |
| `--winsxs`| Clean WinSxS temp files       |
| `--all`   | Clean all caches              |

## ⚠️ Safety Notes

- 🔒 Always run with administrator privileges
- 🛡️ The tool only cleans known-safe locations
- 📁 For WinSxS, only the Temp directory is cleaned
- 🐳 Docker cleanup uses official Docker commands

## 🛠️ Development

### Building from Source
```bash
git clone https://github.com/abdorrahmani/clearance.git
cd clearance
go build -o clearance.exe
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m '✨ feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🌟 Support

If you need help or have questions, please visit [anophel.com/en/contact-us](https://anophel.com/en/contact-us)

---

Made with ❤️ by [Abdorrahmani](https://github.com/abdorrahmani) 