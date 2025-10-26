# BLEF JSON Schema

This directory contains the official JSON Schema for validating BLEF files.

## Current Version

**[blef-schema-v0.1.0.json](./blef-schema-v0.1.0.json)** — Version 0.1.0

## What is JSON Schema?

JSON Schema is a vocabulary that allows you to annotate and validate JSON documents. It provides:

- ✅ **Validation** — Ensure BLEF files are correctly formatted
- 📚 **Documentation** — Self-documenting structure
- 🔧 **Tooling** — Enable auto-completion in IDEs
- 🛡️ **Type Safety** — Catch errors before processing

## Using the Schema

### Online Validation

Quick validation using web tools:

1. Copy your BLEF file content
2. Visit [JSONSchema.net Validator](https://www.jsonschemavalidator.net/)
3. Paste the schema and your BLEF file
4. Check for validation errors

### Command Line Validation

#### Using AJV (Node.js)

```bash
# Install AJV CLI
npm install -g ajv-cli

# Validate a BLEF file
ajv validate -s blef-schema-v0.1.0.json -d your-library.blef.json

# Validate with detailed errors
ajv validate -s blef-schema-v0.1.0.json -d your-library.blef.json --errors=text
```

#### Using check-jsonschema (Python)

```bash
# Install check-jsonschema
pip install check-jsonschema

# Validate a BLEF file
check-jsonschema --schemafile blef-schema-v0.1.0.json your-library.blef.json
```

#### Using jsonschema (Python)

```python
import json
import jsonschema

# Load schema
with open('blef-schema-v0.1.0.json') as f:
    schema = json.load(f)

# Load BLEF file
with open('your-library.blef.json') as f:
    blef_data = json.load(f)

# Validate
try:
    jsonschema.validate(instance=blef_data, schema=schema)
    print("✅ Valid BLEF file!")
except jsonschema.ValidationError as e:
    print(f"❌ Invalid BLEF file: {e.message}")
```

### IDE Integration

#### VS Code

1. Install the **YAML** extension (it also supports JSON Schema)
2. Add to your settings.json:

```json
{
  "json.schemas": [
    {
      "fileMatch": ["*.blef.json"],
      "url": "./schema/blef-schema-v0.1.0.json"
    }
  ]
}
```

Now `.blef.json` files will have auto-completion and inline validation!

#### JetBrains IDEs (IntelliJ, WebStorm, etc.)

1. Go to **Settings** → **Languages & Frameworks** → **Schemas and DTDs** → **JSON Schema Mappings**
2. Click **+** to add a new mapping
3. Set **Schema file**: `blef-schema-v0.1.0.json`
4. Set **Schema version**: `JSON Schema version 7`
5. Add file pattern: `*.blef.json`

## Schema Structure

The BLEF schema defines:

### Root Object

```json
{
  "format": "BLEF",           // Required: Must be "BLEF"
  "version": "0.1.0",         // Required: Semantic version
  "exported_at": "datetime",  // Required: ISO 8601 timestamp
  "user": {},                 // Optional: User information
  "books": [],                // Required: Array of books
  "collections": [],          // Required: Array of collections (min 1)
  "entries": []               // Required: Array of entries
}
```

### Key Definitions

- **`book`** — Unique book with metadata
- **`author`** — Author information
- **`collection`** — Shelf/list definition
- **`entry`** — User data linking books to collections
- **`readDate`** — Reading date range with progress

## Validation Rules

### Required Fields

#### Root Level
- `format`, `version`, `exported_at`
- `books`, `collections`, `entries` (arrays can be empty except collections)

#### Book
- `id` (ISBN13 or UUIDv4)
- `title`, `authors`, `identifiers`

#### Entry
- `book_id`, `collection_ids`, `user_data`
- `user_data.status` (one of: read, reading, to-read, abandoned, wishlist)

### Constraints

- **Book ID**: Must be valid ISBN13 or UUIDv4
- **ISBN13**: Pattern `^97[89]\d{10}$`
- **UUIDv4**: Standard UUID v4 format
- **Language**: ISO 639-1 code (e.g., `en`, `fr`, `en-US`)
- **Status**: Must be one of the enum values
- **Rating**: Number between 0 and 5
- **Collections**: Minimum 1 collection required
- **Entry collection_ids**: Minimum 1 reference required

### Format Validation

- **date-time**: ISO 8601 format (e.g., `2025-10-26T14:00:00Z`)
- **date**: ISO 8601 date (e.g., `2025-10-26`)
- **email**: Valid email format
- **uri**: Valid URI format

## Common Validation Errors

### ❌ "Missing required property"

**Problem**: A required field is missing

**Solution**: Add the missing field:
```json
{
  "format": "BLEF",           // Required
  "version": "0.1.0",         // Required
  "exported_at": "2025-10-26T14:00:00Z"  // Required
}
```

### ❌ "Does not match pattern"

**Problem**: ID format is invalid

**Solution**: Use valid ISBN13 or UUIDv4:
```json
{
  "id": "9780156013987"  // ✅ Valid ISBN13
  // or
  "id": "550e8400-e29b-41d4-a716-446655440000"  // ✅ Valid UUIDv4
}
```

### ❌ "Must be one of [enum values]"

**Problem**: Invalid enum value (e.g., status)

**Solution**: Use valid enum value:
```json
{
  "status": "read"  // ✅ Valid: read, reading, to-read, abandoned, wishlist
}
```

### ❌ "Array must contain at least 1 item"

**Problem**: Empty required array (e.g., collections)

**Solution**: Add at least one item:
```json
{
  "collections": [
    { "id": "read", "name": "Read", "type": "read" }
  ]
}
```

## Schema Versioning

The schema follows the same versioning as the BLEF format:

- **Major** (1.x.x): Breaking changes
- **Minor** (x.1.x): New optional fields
- **Patch** (x.x.1): Bug fixes, clarifications

### Version Compatibility

| BLEF Version | Schema File                  | Status      |
|--------------|------------------------------|-------------|
| 0.1.0        | blef-schema-v0.1.0.json      | ✅ Current   |

## Contributing

Found an issue with the schema? Want to suggest improvements?

- **Bug**: Open an issue with the `bug` label
- **Enhancement**: Open an issue with the `enhancement` label
- **Question**: Use GitHub Discussions

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## Resources

### JSON Schema Resources

- [JSON Schema Official Site](https://json-schema.org/)
- [Understanding JSON Schema](https://json-schema.org/understanding-json-schema/)
- [JSON Schema Validator](https://www.jsonschemavalidator.net/)

### BLEF Resources

- [BLEF Specification](../blef_specification.md)
- [Example Files](../examples/)
- [Contributing Guide](../CONTRIBUTING.md)

---

*Need help? Open a [discussion](https://github.com/yoanbernabeu/BLEF/discussions) or [issue](https://github.com/yoanbernabeu/BLEF/issues).*

