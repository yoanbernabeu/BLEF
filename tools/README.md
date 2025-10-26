# BLEF Tools

This directory contains official tools for working with BLEF (Book Library Exchange Format) files.

## Available Tools

### [blef-cli](./blef-cli/)

A comprehensive command-line interface for BLEF file operations.

**Features:**
- ✅ **Validate** - Validate BLEF files against JSON schema with integrity checks
- ✅ **Convert** - Convert CSV files from Goodreads, Babelio, and custom formats to BLEF
- ✅ **View** - Interactive terminal UI for browsing BLEF files

**Quick Start:**
```bash
cd blef-cli
make build
./blef-cli validate ../../examples/minimal.blef.json
```

**Documentation:** See [blef-cli/README.md](./blef-cli/README.md)

## Future Tools

### Planned
- **blef-web** - Web-based BLEF viewer and editor
- **blef-merge** - Tool to merge multiple BLEF files
- **blef-diff** - Compare two BLEF files and show differences
- **blef-export** - Export BLEF to various formats (CSV, Markdown, HTML)

## Contributing

Want to create a new tool? See our [Contributing Guidelines](../CONTRIBUTING.md).

Tools should:
- Follow the BLEF specification
- Include comprehensive tests
- Provide good documentation
- Use the JSON Schema for validation

## License

All tools in this directory are licensed under the MIT License - see [../LICENSE](../LICENSE) for details.

