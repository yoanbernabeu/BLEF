# BLEF CLI Tool

A command-line tool for working with BLEF (Book Library Exchange Format) files.

## Features

- **Validate** BLEF files against the JSON schema
- **Convert** CSV files from Goodreads, Babelio, and other platforms to BLEF format
- **View** BLEF files in an interactive terminal UI

## Installation

### Quick Install (macOS/Linux)

Install the latest release with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/yoanbernabeu/BLEF/main/tools/blef-cli/install.sh | bash
```

This will:
- Detect your OS and architecture automatically
- Download the latest release
- Install to `~/.local/bin/blef-cli`
- Verify the installation

### Manual Installation

#### Download Pre-built Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/yoanbernabeu/BLEF/releases):

**macOS:**
```bash
# Intel Mac
curl -L https://github.com/yoanbernabeu/BLEF/releases/latest/download/blef-cli-latest-darwin-amd64.tar.gz | tar xz

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/yoanbernabeu/BLEF/releases/latest/download/blef-cli-latest-darwin-arm64.tar.gz | tar xz

# Move to PATH
sudo mv blef-cli-darwin-* /usr/local/bin/blef-cli
```

**Linux:**
```bash
# x86_64
curl -L https://github.com/yoanbernabeu/BLEF/releases/latest/download/blef-cli-latest-linux-amd64.tar.gz | tar xz

# ARM64
curl -L https://github.com/yoanbernabeu/BLEF/releases/latest/download/blef-cli-latest-linux-arm64.tar.gz | tar xz

# Move to PATH
sudo mv blef-cli-linux-* /usr/local/bin/blef-cli
```

**Windows:**

Download the appropriate `.tar.gz` from the [releases page](https://github.com/yoanbernabeu/BLEF/releases/latest), extract it, and add the binary to your PATH.

#### From Source

```bash
git clone https://github.com/yoanbernabeu/BLEF.git
cd BLEF/tools/blef-cli
go build -o blef-cli .
```

#### Using Go Install

```bash
go install github.com/yoanbernabeu/BLEF/tools/blef-cli@latest
```

## Commands

### Validate

Validate a BLEF file for correctness:

```bash
blef-cli validate my-library.blef.json
```

Features:
- JSON schema validation
- Referential integrity checks
- ISBN-13 check digit validation
- Statistics display

### Convert

Convert CSV files to BLEF format:

```bash
# Auto-detect format
blef-cli convert books.csv

# Specify output file
blef-cli convert books.csv -o my-library.blef.json

# Force specific format
blef-cli convert books.csv -f goodreads

# Skip validation
blef-cli convert books.csv --no-validate
```

Supported formats:
- **Goodreads** - Library export CSV
- **Babelio** - Library export CSV
- **Custom** - Interactive column mapping

Flags:
- `-o, --output` - Output file path (default: input.blef.json)
- `-f, --format` - Force format (goodreads, babelio)
- `--no-validate` - Skip validation after conversion

#### Interactive Mapping

If the CSV format isn't recognized, you'll be prompted to map each column:

```
📋 Column Mapping
Map column 'ISBN' to:
> ISBN-13
  ISBN-10
  Title
  Author
  ...
```

### View

Launch an interactive terminal viewer:

```bash
blef-cli view my-library.blef.json
```

Features:
- Browse books, collections, and statistics
- View detailed book information
- Color-coded reading status
- Fast keyboard navigation

Controls:
- `↑/↓` or `j/k` - Navigate
- `Tab` - Switch views (Books/Collections/Stats)
- `Enter` - View details
- `Esc` - Go back
- `q` - Quit

## CSV Format Requirements

### Goodreads Export

Export your library from Goodreads (My Books → Import/Export → Export Library).

Expected columns:
- Book Id
- Title
- Author
- ISBN13, ISBN
- My Rating
- Exclusive Shelf
- Date Read

### Babelio Export

Export your library from Babelio.

Expected columns:
- EAN
- Titre
- Auteur
- Note
- Étagère

### Custom CSV

Any CSV file with book data. You'll be prompted to map columns interactively.

Required fields:
- Title
- Author (recommended)
- Some identifier (ISBN-13, ISBN-10, or unique ID)

## Examples

### Convert Goodreads Export

```bash
blef-cli convert goodreads_library_export.csv -o my-books.blef.json
```

### Validate and View

```bash
blef-cli validate my-books.blef.json
blef-cli view my-books.blef.json
```

### Custom CSV with Interactive Mapping

```bash
blef-cli convert custom_books.csv
# Follow the prompts to map your columns
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Install Locally

```bash
make install
```

### Clean

```bash
make clean
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [survey](https://github.com/AlecAivazis/survey) - Interactive prompts
- [gojsonschema](https://github.com/xeipuuv/gojsonschema) - JSON Schema validation
- [uuid](https://github.com/google/uuid) - UUID generation

## Contributing

See the main [CONTRIBUTING.md](../../CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](../../LICENSE) for details.

## Links

- [BLEF Specification](../../docs/SPECIFICATION.md)
- [Examples](../../examples/)
- [JSON Schema](../../schema/)
- [GitHub Repository](https://github.com/yoanbernabeu/BLEF)

