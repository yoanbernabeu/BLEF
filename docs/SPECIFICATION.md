```
Network Working Group                                      Y. Bernabeu
Request for Comments: BLEF-001                                 YoanDev
Category: Standards Track                                 October 2025
```

# BLEF: Book Library Exchange Format

## Abstract

This document specifies the Book Library Exchange Format (BLEF), a JSON-based data interchange format for personal book library information. BLEF enables interoperability between reading platforms, library management tools, and online book services by providing a vendor-neutral, human-readable standard for representing book collections, reading status, user annotations, and lending information.

## Status of This Memo

This document specifies a standards track protocol for the Internet community, and requests discussion and suggestions for improvements. Distribution of this memo is unlimited.

## Copyright Notice

Copyright (c) 2025 BLEF Contributors. This document is subject to the rights, licenses and restrictions contained in MIT License for code implementations. The specification itself is released under CC0 1.0 Universal (Public Domain Dedication).

## Table of Contents

1. [Introduction](#1-introduction)
2. [Conventions and Terminology](#2-conventions-and-terminology)
3. [Format Overview](#3-format-overview)
4. [Data Model](#4-data-model)
5. [Identification and Uniqueness](#5-identification-and-uniqueness)
6. [Validation Rules](#6-validation-rules)
7. [JSON Schema](#7-json-schema)
8. [Examples](#8-examples)
9. [Security Considerations](#9-security-considerations)
10. [IANA Considerations](#10-iana-considerations)
11. [References](#11-references)
12. [Acknowledgments](#12-acknowledgments)

---

## 1. Introduction

### 1.1. Purpose

The Book Library Exchange Format (BLEF) addresses the problem of data portability and vendor lock-in in personal book library management. Current book tracking platforms employ proprietary data formats, making it difficult for users to:

- Migrate between platforms without data loss
- Maintain control over their reading history
- Create backups in a standardized format
- Share library data across multiple tools

BLEF provides a standardized, JSON-based format that focuses on reader usage patterns rather than exhaustive bibliographic metadata.

### 1.2. Scope

BLEF is designed to represent:

- Personal book collections and shelves
- Reading status and progress tracking
- User annotations (ratings, reviews, notes, tags)
- Book ownership and lending information
- Multi-list membership for books

BLEF explicitly does NOT cover:

- Network synchronization protocols
- Real-time collaborative features
- Commercial or recommendation algorithms
- Complex lending workflows with reminders
- Full bibliographic cataloging metadata

### 1.3. Goals

The primary goals of BLEF are:

1. **Simplicity**: Easy to understand and implement
2. **Readability**: Human-readable JSON format
3. **Extensibility**: Support for optional fields and future enhancements
4. **Interoperability**: Vendor-neutral design
5. **Data Ownership**: Enable user control over personal data

## 2. Conventions and Terminology

### 2.1. Requirements Language

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", and "OPTIONAL" in this document are to be interpreted as described in RFC 2119 [RFC2119].

### 2.2. Definitions

**Book**: A unique bibliographic work identified by ISBN-13 or UUID.

**Entry**: The association between a book and user-specific data.

**Collection**: A user-defined list or shelf containing zero or more books.

**Status**: An enumerated value representing the reading state of a book.

**Identifier**: A standardized code used to uniquely identify books (ISBN-13, ISBN-10, OpenLibrary ID, Wikidata QID, etc.).

### 2.3. Notational Conventions

JSON structures are represented using standard JSON notation [RFC8259]. Optional fields are indicated in accompanying text.

## 3. Format Overview

### 3.1. General Structure

A BLEF document is a JSON object containing the following top-level members:

```json
{
  "format": "BLEF",
  "version": "0.1.0",
  "exported_at": "2025-10-26T14:00:00Z",
  "user": { },
  "books": [ ],
  "collections": [ ],
  "entries": [ ]
}
```

### 3.2. Character Encoding

BLEF documents MUST be encoded in UTF-8 [RFC3629].

### 3.3. Media Type

The media type for BLEF documents is `application/vnd.blef+json`.

## 4. Data Model

### 4.1. Root Object

The root object MUST contain the following members:

- **format** (string, REQUIRED): MUST be the literal string "BLEF"
- **version** (string, REQUIRED): Version string following semantic versioning (e.g., "0.1.0")
- **exported_at** (string, REQUIRED): ISO 8601 [ISO8601] timestamp indicating export time
- **books** (array, REQUIRED): Array of book objects (MAY be empty)
- **collections** (array, REQUIRED): Array of collection objects (MUST contain at least one element)
- **entries** (array, REQUIRED): Array of entry objects (MAY be empty)

The root object MAY contain the following member:

- **user** (object, OPTIONAL): User information object

### 4.2. Book Object

A book object represents unique bibliographic metadata. It MUST contain:

- **id** (string, REQUIRED): Primary identifier, either:
  - ISBN-13: 13-digit number matching pattern `^97[89]\d{10}$`
  - UUID v4: Standard UUID format
- **title** (string, REQUIRED): Book title
- **authors** (array, REQUIRED): Array of author objects (MUST contain at least one)
- **identifiers** (object, REQUIRED): Object containing at least one identifier

A book object MAY contain:

- **subtitle** (string, OPTIONAL)
- **language** (string, OPTIONAL): ISO 639-1 language code
- **description** (string, OPTIONAL)
- **cover_url** (string, OPTIONAL): URI to cover image
- **edition** (object, OPTIONAL): Edition information
- **series** (object, OPTIONAL): Series information
- **subjects** (array, OPTIONAL): Array of subject strings
- **metadata** (object, OPTIONAL): Additional custom metadata

#### 4.2.1. Author Object

An author object MUST contain:

- **name** (string, REQUIRED): Author name

An author object MAY contain:

- **role** (string, OPTIONAL): One of: "author", "editor", "translator", "illustrator", "contributor"
- **identifiers** (object, OPTIONAL): Author identifier object

#### 4.2.2. Edition Object

An edition object MAY contain:

- **publisher** (string, OPTIONAL)
- **published_date** (string, OPTIONAL): ISO 8601 date or year
- **format** (string, OPTIONAL): One of: "hardcover", "paperback", "ebook", "audiobook", "other"
- **pages** (integer, OPTIONAL): Page count (MUST be >= 1)
- **edition_number** (string, OPTIONAL)

#### 4.2.3. Series Object

A series object MUST contain:

- **name** (string, REQUIRED): Series name

A series object MAY contain:

- **volume** (number or string, OPTIONAL): Volume number

### 4.3. Collection Object

A collection object MUST contain:

- **id** (string, REQUIRED): Unique collection identifier
- **name** (string, REQUIRED): Display name
- **type** (string, REQUIRED): One of: "read", "reading", "to-read", "wishlist", "owned", "custom"

A collection object MAY contain:

- **description** (string, OPTIONAL)
- **is_public** (boolean, OPTIONAL): Default true
- **created_at** (string, OPTIONAL): ISO 8601 timestamp
- **metadata** (object, OPTIONAL)

### 4.4. Entry Object

An entry object links a book to user-specific data. It MUST contain:

- **book_id** (string, REQUIRED): Reference to a book.id
- **collection_ids** (array, REQUIRED): Array of collection IDs (MUST contain at least one)
- **user_data** (object, REQUIRED): User data object

An entry object MAY contain:

- **ownership** (object, OPTIONAL): Ownership information
- **metadata** (object, OPTIONAL): Additional custom metadata

#### 4.4.1. User Data Object

A user_data object MUST contain:

- **status** (string, REQUIRED): One of: "read", "reading", "to-read", "abandoned", "wishlist"

A user_data object MAY contain:

- **rating** (number, OPTIONAL): Rating from 0 to 5
- **review** (string, OPTIONAL): User review text
- **private_notes** (string, OPTIONAL): Private notes
- **tags** (array, OPTIONAL): Array of tag strings
- **favorite** (boolean, OPTIONAL): Default false
- **read_dates** (array, OPTIONAL): Array of read date objects
- **added_at** (string, OPTIONAL): ISO 8601 timestamp

#### 4.4.2. Ownership Object

An ownership object MAY contain:

- **owned** (boolean, OPTIONAL): Whether user owns the book
- **loaned** (object, OPTIONAL): Loan information object

A loaned object MUST contain:

- **status** (boolean, REQUIRED): Currently loaned out

A loaned object MAY contain:

- **to** (string, OPTIONAL): Borrower name
- **date** (string, OPTIONAL): Loan date (ISO 8601)
- **notes** (string, OPTIONAL): Loan notes

#### 4.4.3. Read Date Object

A read date object MAY contain:

- **started** (string, OPTIONAL): Start date (ISO 8601 date format)
- **finished** (string, OPTIONAL): Finish date (ISO 8601 date format)
- **progress** (integer, OPTIONAL): Progress percentage (0-100)

## 5. Identification and Uniqueness

### 5.1. Book Identification

Each book MUST have a unique `id` field within a BLEF document. The id MUST be either:

1. A valid ISBN-13 (preferred for published books)
2. A valid UUID v4 (for books without ISBN or user-added content)

### 5.2. Reference Integrity

All `book_id` values in entry objects MUST reference an existing `books[].id` value in the same document.

All values in `collection_ids` arrays MUST reference existing `collections[].id` values in the same document.

### 5.3. Uniqueness Constraints

- `books[].id` values MUST be unique within the document
- `collections[].id` values MUST be unique within the document
- An entry object MAY NOT appear multiple times with the same `book_id`

## 6. Validation Rules

### 6.1. Required Fields

Implementations MUST validate that all REQUIRED fields are present.

### 6.2. Data Types

Implementations MUST validate that field values match their specified types.

### 6.3. Enumerated Values

Implementations MUST validate that enumerated string values (status, collection type, format, author role) match one of the specified allowed values.

### 6.4. Format Validation

Implementations SHOULD validate:

- ISBN-13 format and check digit
- UUID v4 format
- ISO 8601 datetime and date formats
- ISO 639-1 language codes
- URI format for URLs

### 6.5. Range Validation

Implementations MUST validate:

- Rating values are between 0 and 5 (inclusive)
- Progress values are between 0 and 100 (inclusive)
- Page counts are positive integers

## 7. JSON Schema

The normative JSON Schema definition for BLEF v0.1.0 is provided separately in the file `blef-schema-v0.1.0.json`.

The schema follows JSON Schema Draft 2020-12 [JSON-SCHEMA].

Implementations SHOULD validate BLEF documents against the official schema.

## 8. Examples

### 8.1. Minimal Valid Document

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
      "user_data": { "status": "read" }
    }
  ]
}
```

### 8.2. Invalid Document

The following document is invalid because it lacks required fields:

```json
{
  "format": "BLEF",
  "version": "0.1.0",
  "books": []
}
```

Reasons for invalidity:
- Missing `exported_at` (REQUIRED)
- Missing `collections` (REQUIRED)
- Missing `entries` (REQUIRED)

### 8.3. Complete Example

For a comprehensive example demonstrating all features, see the file `examples/complete.blef.json` in the specification repository.

## 9. Security Considerations

### 9.1. Privacy

BLEF documents MAY contain personal information including:

- Reading habits and preferences
- Personal notes and reviews
- Social relationships (via lending information)
- User identification data

Implementations SHOULD:

- Warn users before exporting personal data
- Provide options to exclude sensitive fields
- Use secure transport (HTTPS) when transmitting BLEF data
- Respect user privacy preferences

### 9.2. Data Validation

Implementations MUST validate all input to prevent:

- Injection attacks in note/review fields
- Path traversal via file references
- Excessive resource consumption via large documents

### 9.3. No Executable Content

BLEF is a data-only format. Implementations MUST NOT:

- Execute code contained in BLEF documents
- Evaluate strings as code
- Follow URLs automatically without user consent

### 9.4. GDPR Compliance

BLEF documents containing personal data are subject to GDPR and similar regulations. Implementations SHOULD:

- Obtain user consent before data export
- Provide mechanisms for data deletion
- Honor user data portability rights
- Implement appropriate security measures

## 10. IANA Considerations

### 10.1. Media Type Registration

This specification requests registration of the media type `application/vnd.blef+json`:

- Type name: application
- Subtype name: vnd.blef+json
- Required parameters: None
- Optional parameters: None
- Encoding considerations: UTF-8
- Security considerations: See Section 9
- Interoperability considerations: None
- Published specification: This document
- Applications that use this media type: Book library management systems, reading platforms, data portability tools

### 10.2. File Extension

The recommended file extension for BLEF documents is `.blef.json`.

## 11. References

### 11.1. Normative References

**[RFC2119]** Bradner, S., "Key words for use in RFCs to Indicate Requirement Levels", BCP 14, RFC 2119, March 1997.

**[RFC3629]** Yergeau, F., "UTF-8, a transformation format of ISO 10646", STD 63, RFC 3629, November 2003.

**[RFC8259]** Bray, T., "The JavaScript Object Notation (JSON) Data Interchange Format", STD 90, RFC 8259, December 2017.

**[ISO8601]** ISO 8601:2004, "Data elements and interchange formats - Information interchange - Representation of dates and times".

**[JSON-SCHEMA]** "JSON Schema: A Media Type for Describing JSON Documents", draft-bhutton-json-schema-01, December 2020.

### 11.2. Informative References

**[ISBN]** ISO 2108:2017, "Information and documentation - International Standard Book Number (ISBN)".

**[UUID]** RFC 4122, "A Universally Unique IDentifier (UUID) URN Namespace", July 2005.

**[ISO639]** ISO 639-1:2002, "Codes for the representation of names of languages - Part 1: Alpha-2 code".

### 11.3. External Resources

- JSON Schema: https://json-schema.org/
- BLEF Repository: https://github.com/yoanbernabeu/BLEF
- Examples: https://github.com/yoanbernabeu/BLEF/tree/main/examples

## 12. Acknowledgments

BLEF was designed to address real-world interoperability challenges in personal book library management. The author thanks the open reading community for their input and feedback.

Special acknowledgment to existing platforms (Goodreads, Babelio, LibraryThing, Calibre, BookWyrm, Inventaire.io) whose feature sets informed the design of this specification.

---

## Appendix A. Platform Compatibility

This section is informative.

| Platform      | Import Support | Export Support | Notes                           |
|---------------|----------------|----------------|---------------------------------|
| Goodreads     | Via CSV        | Via CSV        | CSV conversion required         |
| Babelio       | Via CSV        | Via CSV        | EAN to ISBN-13 mapping          |
| BookWyrm      | Custom API     | Supported      | ActivityPub Book alignment      |
| Calibre       | Via Plugin     | Via Plugin     | Custom/UUID ID support          |
| Inventaire.io | Via API        | Partial        | Native Wikidata ID support      |

## Appendix B. Version History

### Version 0.1.0 (October 2025)

Initial specification release including:
- Core data model (books, entries, collections)
- Multiple identifier support
- Single status with multi-collection membership
- Optional series and lending features
- JSON Schema validation

---

## Author's Address

Yoan Bernabeu  
Email: contact@yoandev.co  
GitHub: https://github.com/yoanbernabeu/BLEF

---

**Document Status**: Standards Track  
**Version**: 0.1.0  
**Last Updated**: October 26, 2025
