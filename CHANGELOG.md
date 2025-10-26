# Changelog

All notable changes to the BLEF format and tools will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- CLI validator tool
- CSV to BLEF converter
- Reference implementations for popular platforms

## [0.1.0] - 2025-10-26

### Added
- Initial BLEF format specification
- Core data model with three main sections: books, entries, collections
- Support for multiple identifiers (ISBN13, ISBN10, UUID, OpenLibrary, Wikidata)
- Single reading status per book with multiple collection support
- Optional series and volume tracking
- Optional lending/loan management
- User data fields: ratings, notes, tags, dates
- JSON Schema for validation
- Example files
- Project documentation (README, CONTRIBUTING, CODE_OF_CONDUCT)
- GitHub issue and PR templates

### Format Features
- **Books**: Unique bibliographic metadata per book
- **Entries**: User-specific data linking books to personal library
- **Collections**: Custom shelves/lists with standardized types
- **Identifiers**: ISBN13 or UUIDv4 as primary ID
- **Status**: Enum-based reading states (read, reading, to-read, abandoned, wishlist)
- **Extensibility**: Optional fields for future enhancements

### Documentation
- Complete specification in English and French
- JSON Schema v0.1.0
- Quick start guide
- Contributing guidelines
- MIT License for code, CC0 for specification

## Format Version History

### Version Numbering
- **Major** (1.x.x): Breaking changes to core structure
- **Minor** (x.1.x): New optional fields or features
- **Patch** (x.x.1): Clarifications, examples, bug fixes in schema

---

[Unreleased]: https://github.com/yoanbernabeu/BLEF/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/yoanbernabeu/BLEF/releases/tag/v0.1.0

