# BLEF — Book Library Exchange Format

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/version-0.1.0-blue.svg)](https://github.com/yoanbernabeu/BLEF)

> **An open, interoperable standard for exchanging personal book library data**

BLEF (Book Library Exchange Format) is a simple, extensible JSON format designed to enable seamless data exchange between reading platforms, library management tools, and online book services.

## 🎯 Why BLEF?

Current book tracking platforms (Goodreads, Babelio, LibraryThing, etc.) use proprietary formats, making it difficult to:
- **Switch platforms** without losing your reading history
- **Backup** your personal library data
- **Share** your reading lists across tools
- **Own** your data in a standardized way

BLEF solves this by providing a **vendor-neutral**, **human-readable** format that anyone can implement.

## ✨ Key Features

- ✅ **Multiple collections** — Books can belong to several lists simultaneously
- ✅ **Single reading status** — Clear, unambiguous reading state per book
- ✅ **Normalized identifiers** — ISBN13, ISBN10, UUID, OpenLibrary, Wikidata
- ✅ **Series support** — Track book series and volumes
- ✅ **User annotations** — Notes, ratings, tags, reading dates
- ✅ **Optional lending** — Simple loan tracking
- ✅ **Extensible** — Easy to add custom fields

## 📦 Quick Example

```json
{
  "format": "BLEF",
  "version": "0.1.0",
  "exported_at": "2025-10-26T14:00:00Z",
  "books": [
    {
      "id": "9780156013987",
      "title": "The Little Prince",
      "authors": [{ "name": "Antoine de Saint-Exupéry" }],
      "identifiers": { "isbn13": "9780156013987" },
      "language": "en"
    }
  ],
  "collections": [
    { "id": "read", "name": "Read", "type": "read" }
  ],
  "entries": [
    {
      "book_id": "9780156013987",
      "collection_ids": ["read"],
      "user_data": { 
        "status": "read",
        "rating": 5,
        "read_dates": [{ "started": "2025-10-01", "finished": "2025-10-05" }]
      }
    }
  ]
}
```

## 🚀 Getting Started

### For Users

1. **Export** your library from your current platform (CSV/API)
2. **Convert** to BLEF format using the [CLI tool](./tools/blef-cli/)
   ```bash
   blef-cli convert goodreads_export.csv -o my-library.blef.json
   ```
3. **Validate** your BLEF file
   ```bash
   blef-cli validate my-library.blef.json
   ```
4. **View** your library interactively
   ```bash
   blef-cli view my-library.blef.json
   ```

### For Developers

1. Read the [specification](./docs/SPECIFICATION.md)
2. Implement import/export in your application
3. Validate files against the [JSON Schema](./schema/blef-schema-v0.1.0.json)
4. Use the [CLI tool](./tools/blef-cli/) for testing
5. Share your implementation with the community

## 📚 Documentation

- **[Full Specification](./docs/SPECIFICATION.md)** — Complete format documentation
- **[JSON Schema](./schema/blef-schema-v0.1.0.json)** — Validation schema
- **[Examples](./examples/)** — Sample BLEF files
- **[Converters](./tools/)** — Import/export tools (coming soon)

## 🔄 Platform Compatibility

| Platform      | Import | Export | Notes                    |
|---------------|--------|--------|--------------------------|
| Goodreads     | ✅     | ✅     | Via CSV conversion       |
| Babelio       | ✅     | ✅     | EAN → ISBN13 supported   |
| BookWyrm      | 🔶     | ✅     | API integration possible |
| Calibre       | ✅     | ✅     | Via plugin               |
| Inventaire.io | 🔶     | 🔶     | Wikidata IDs native      |

## 🛠️ Tools & Implementations

### Official Tools

#### BLEF CLI

Complete command-line tool for working with BLEF files:

**Quick Install (macOS/Linux):**
```bash
curl -fsSL https://raw.githubusercontent.com/yoanbernabeu/BLEF/main/tools/blef-cli/install.sh | bash
```

**Features:**
- ✅ Validator (JSON schema + integrity checks)
- ✅ CSV to BLEF converter (Goodreads, Babelio, custom)
- ✅ Interactive TUI viewer

**[Full Documentation →](./tools/blef-cli/)**

### Community Implementations
*Submit your implementation via PR!*

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](./CONTRIBUTING.md) for details.

Ways to contribute:
- 🐛 Report bugs or suggest features
- 📝 Improve documentation
- 🔧 Build converters and tools
- 💡 Share implementation feedback
- 🌍 Translate the specification

## 📄 License

- **Code & Tools**: [MIT License](./LICENSE)
- **Specification**: CC0 1.0 (Public Domain)

This means:
- ✅ Use BLEF in commercial projects
- ✅ Modify and extend the format
- ✅ No attribution required (though appreciated!)

## 🌟 Support the Project

If BLEF is useful to you:
- ⭐ Star this repository
- 📢 Share with your community
- 🔗 Implement in your tools
- 💬 Join the discussion

## 📞 Contact & Community

- **Issues**: [GitHub Issues](https://github.com/yoanbernabeu/BLEF/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yoanbernabeu/BLEF/discussions)
- **Email**: contact@yoandev.co

---

**Made with 📚 by the open reading community**

*Current version: 0.1.0 — Last updated: October 2025*

