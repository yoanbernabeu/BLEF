# BLEF Examples

This directory contains example BLEF files demonstrating various features and use cases of the format.

## 📁 Example Files

### [`minimal.blef.json`](./minimal.blef.json)
**The simplest valid BLEF file**

- Single book (The Little Prince)
- One collection (Read)
- Minimal required fields only
- Perfect for understanding the basic structure

**Use this for:**
- Learning the format basics
- Testing validators
- Starting your own BLEF file

---

### [`complete.blef.json`](./complete.blef.json)
**Comprehensive example showcasing all features**

- Multiple books with full metadata
- Various collections (read, reading, to-read, custom)
- User information
- Multiple readings of the same book
- Ratings, reviews, tags
- Ownership and lending information
- Series information
- Different book formats

**Use this for:**
- Understanding all available fields
- Reference implementation
- Testing advanced features

---

### [`french-library.blef.json`](./french-library.blef.json)
**French book library example**

- French language books
- French collection names
- French ISBN identifiers
- Classic literature examples
- Demonstrates multilingual support

**Use this for:**
- Non-English libraries
- Understanding language support
- French platform integration

---

### [`series-tracking.blef.json`](./series-tracking.blef.json)
**Series and saga tracking**

- Multiple books from the same series (Mistborn)
- Series name and volume numbers
- Progress tracking across series
- Custom series collection
- Reading order tracking

**Use this for:**
- Managing book series
- Tracking reading progress in sagas
- Series-focused applications

---

## 🧪 Validation

All examples in this directory are valid BLEF files that pass schema validation.

To validate these files:

```bash
# Using ajv-cli (requires installation)
ajv validate -s ../schema/blef-schema-v0.1.0.json -d minimal.blef.json
ajv validate -s ../schema/blef-schema-v0.1.0.json -d complete.blef.json
ajv validate -s ../schema/blef-schema-v0.1.0.json -d french-library.blef.json
ajv validate -s ../schema/blef-schema-v0.1.0.json -d series-tracking.blef.json
```

Or using online validators:
- [JSONSchema.net](https://www.jsonschemavalidator.net/)
- [JSON Schema Validator](https://json-schema-validator.herokuapp.com/)

## 💡 Usage Tips

### For Developers

1. **Start with minimal** — Understand the basic structure first
2. **Review complete** — See all possible fields and their usage
3. **Check language examples** — Understand i18n support
4. **Examine series tracking** — Learn about series management

### For Users

These examples show what your exported library data might look like. When choosing or building a BLEF-compatible tool, check if it supports the features shown in these examples that matter to you.

### For Tool Creators

Use these examples as:
- **Test fixtures** for your importer/exporter
- **Documentation references** for your users
- **Validation targets** for your implementation

## 🔧 Modifying Examples

Feel free to:
- ✅ Copy and modify for your own use
- ✅ Use as templates for new BLEF files
- ✅ Test with different values

Remember to:
- ⚠️ Validate against the schema after changes
- ⚠️ Maintain required fields
- ⚠️ Use valid ISBNs or UUIDs for book IDs

## 📝 Adding More Examples

Have an interesting use case? We welcome example contributions!

Good examples to add:
- Audiobook-focused library
- Academic/textbook library
- Manga/comic collection
- Multilingual mixed library
- Large library (100+ books)
- Migration from specific platforms

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

---

## 📚 Quick Reference

### Minimal Required Fields

```json
{
  "format": "BLEF",
  "version": "0.1.0",
  "exported_at": "ISO-8601-datetime",
  "books": [ /* at least one book */ ],
  "collections": [ /* at least one collection */ ],
  "entries": [ /* links books to collections */ ]
}
```

### Book Required Fields

```json
{
  "id": "ISBN13 or UUIDv4",
  "title": "string",
  "authors": [{ "name": "string" }],
  "identifiers": { /* at least one */ }
}
```

### Entry Required Fields

```json
{
  "book_id": "reference to book.id",
  "collection_ids": ["at least one"],
  "user_data": {
    "status": "read|reading|to-read|abandoned|wishlist"
  }
}
```

---

*Need help? Check the [full specification](../docs/SPECIFICATION.md) or [open an issue](https://github.com/yoanbernabeu/BLEF/issues).*

