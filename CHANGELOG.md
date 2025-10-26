# Changelog

All notable changes to the BLEF format and tools will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Web-based BLEF viewer and editor
- BLEF merge tool (combine multiple libraries)
- BLEF diff tool (compare libraries)
- Export to additional formats (Markdown, HTML)
- Platform-specific plugins (Calibre, browser extensions)
- OPDS feed support
- API specification for BLEF services

## [0.1.0] - 2025-10-26

### Added

#### Format Specification
- Initial BLEF format specification (RFC-style)
- Core data model with three main sections: books, entries, collections
- Support for multiple identifiers (ISBN13, ISBN10, UUID, OpenLibrary, Wikidata, Goodreads)
- Single reading status per book with multiple collection support
- Optional series and volume tracking
- Optional lending/loan management
- User data fields: ratings, reviews, private notes, tags, reading dates
- JSON Schema v0.1.0 for validation
- Media type registration: `application/vnd.blef+json`
- Recommended file extension: `.blef.json`

#### CLI Tool (blef-cli)
- **Validate Command**: Complete validation tool with
  - JSON Schema validation
  - Referential integrity checks
  - ISBN-13 check digit validation
  - Detailed statistics display
  - Color-coded output
- **Convert Command**: CSV to BLEF converter with
  - Auto-detection for Goodreads and Babelio formats
  - Interactive column mapping for custom CSV files
  - UUID generation for books without ISBN
  - Automatic validation post-conversion
  - Flags: `--output`, `--format`, `--no-validate`
- **View Command**: Interactive TUI viewer with
  - Tab navigation (Books/Collections/Stats)
  - Keyboard navigation (arrows, j/k, enter, esc)
  - Book detail view
  - Color-coded reading status
  - Search and filtering
- Test suite with comprehensive coverage
- Built with Go 1.21+ using cobra, bubbletea, lipgloss, survey

#### Examples
- Minimal valid BLEF file
- Complete BLEF file with all features
- French library example
- Series tracking example (Mistborn)
- All examples validated against schema

#### Documentation
- RFC-style specification document
- Complete JSON Schema with detailed validation rules
- README with quick start examples
- CLI tool documentation
- Schema validation guide
- Contributing guidelines
- Code of Conduct
- Security policy
- Platform compatibility matrix

#### Project Infrastructure
- MIT License for code
- CC0 1.0 for specification
- GitHub issue templates (bug, feature, question)
- Pull request template
- Makefile for CLI builds
- Comprehensive .gitignore

### Format Features
- **Books**: Unique bibliographic metadata per book
  - Required: id, title, authors, identifiers
  - Optional: subtitle, language, description, cover URL, edition, series, subjects
- **Entries**: User-specific data linking books to personal library
  - Required: book_id, collection_ids, user_data with status
  - Optional: ratings (0-5), review, private notes, tags, read dates, ownership info
- **Collections**: Custom shelves/lists with standardized types
  - Types: read, reading, to-read, wishlist, owned, custom
  - Each entry can belong to multiple collections
- **Identifiers**: ISBN13 or UUIDv4 as primary ID
  - Additional identifiers: ISBN10, ASIN, OpenLibrary, Wikidata, Goodreads
- **Status**: Enum-based reading states (read, reading, to-read, abandoned, wishlist)
- **Extensibility**: Optional metadata fields at all levels

### Technical Details
- Character encoding: UTF-8
- Timestamp format: ISO 8601
- Date format: ISO 8601 (YYYY-MM-DD)
- Language codes: ISO 639-1
- UUID: Version 4
- ISBN validation: Check digit verification

### Platform Compatibility
- Goodreads: CSV import/export supported
- Babelio: CSV import/export supported
- BookWyrm: Partial support via API
- Calibre: Plugin development possible
- Inventaire.io: Partial support (Wikidata IDs)

## Format Version History

### Version Numbering
- **Major** (1.x.x): Breaking changes to core structure
- **Minor** (x.1.x): New optional fields or features
- **Patch** (x.x.1): Clarifications, examples, bug fixes in schema

---

[Unreleased]: https://github.com/yoanbernabeu/BLEF/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/yoanbernabeu/BLEF/releases/tag/v0.1.0

