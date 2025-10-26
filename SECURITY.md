# Security Policy

## Supported Versions

As BLEF is an open format specification, security updates primarily concern the JSON Schema validation and any official tools provided.

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Security Considerations for BLEF

### Data Privacy

BLEF files may contain personal information including:

- 📚 **Reading history** — Personal reading preferences and habits
- 📝 **Reviews and notes** — Private thoughts and opinions
- 👤 **User information** — Name, email, identifiers
- 🏠 **Lending information** — Social connections (who borrowed books)

**Recommendations:**

1. **Sanitize before sharing** — Remove or anonymize personal data before publishing BLEF files
2. **Private fields** — Use `private_notes` for information you don't want to export
3. **Review exports** — Always review generated BLEF files before sharing
4. **Secure storage** — Store BLEF files securely, especially if they contain personal data

### For Tool Developers

When implementing BLEF import/export:

- ✅ **Provide privacy controls** — Let users choose what to export
- ✅ **Warn about personal data** — Alert users before exporting personal information
- ✅ **Secure transmission** — Use HTTPS for any network operations
- ✅ **Validate input** — Always validate imported BLEF files against the schema
- ✅ **Sanitize user input** — Prevent injection attacks in notes/reviews
- ⚠️ **Don't log sensitive data** — Avoid logging user data during processing
- ⚠️ **Handle errors safely** — Don't expose internal paths or system information

### JSON Schema Validation

Always validate BLEF files against the official JSON Schema to prevent:

- Malformed data
- Invalid identifiers
- Missing required fields
- Type mismatches

### No Executable Code

BLEF is a **data-only format**. It should never contain:

- ❌ JavaScript or executable code
- ❌ SQL queries or database commands
- ❌ System commands
- ❌ URLs that auto-execute on load

Implementations should treat BLEF files as **pure data**.

## Reporting a Vulnerability

### Scope

Please report vulnerabilities related to:

- Security issues in the JSON Schema
- Privacy concerns in the specification
- Security flaws in official tools (when available)
- Documentation that could lead to insecure implementations

### How to Report

**For security issues, please DO NOT open a public issue.**

Instead, report security vulnerabilities privately:

1. **Email**: contact@yoandev.co
   - Subject: `[SECURITY] Brief description`
   - Include detailed description of the vulnerability
   - Provide steps to reproduce (if applicable)
   - Suggest a fix (if you have one)

2. **Expected Response Time**:
   - **Initial response**: Within 48 hours
   - **Status update**: Within 7 days
   - **Fix timeline**: Depends on severity

3. **Severity Levels**:
   - 🔴 **Critical**: Data exposure, privacy breach — Fixed ASAP
   - 🟡 **High**: Potential security issue — Fixed within 30 days
   - 🟢 **Medium**: Minor concern — Fixed in next release
   - ⚪ **Low**: Theoretical issue — Addressed when possible

### What Happens Next

1. **Acknowledgment** — We'll confirm receipt of your report
2. **Assessment** — We'll evaluate the severity and impact
3. **Fix Development** — We'll work on a fix (privately if critical)
4. **Disclosure** — We'll coordinate public disclosure with you
5. **Credit** — You'll be credited in the security advisory (if desired)

### Disclosure Policy

- We follow **coordinated disclosure**
- Security fixes are released before public disclosure
- We'll work with you on the disclosure timeline
- Typical embargo period: 90 days (negotiable)

### Security Advisories

Published security advisories will be posted:
- GitHub Security Advisories
- CHANGELOG.md with `[SECURITY]` tag
- Project README (for critical issues)

## Security Best Practices for Users

### When Exporting

- [ ] Review what data you're exporting
- [ ] Remove sensitive private notes
- [ ] Anonymize lending information if needed
- [ ] Consider using generic collection names

### When Importing

- [ ] Only import BLEF files from trusted sources
- [ ] Validate files against the schema before importing
- [ ] Review the file contents before importing
- [ ] Backup your existing library before importing

### When Sharing

- [ ] Don't share raw BLEF files publicly unless intentional
- [ ] Sanitize personal information
- [ ] Consider who can see your reading history
- [ ] Be aware of what your reviews/notes reveal

## Third-Party Implementations

BLEF is an open standard. We cannot guarantee the security of third-party implementations. When using BLEF-compatible tools:

- ✅ Check if the tool is open-source
- ✅ Review the tool's privacy policy
- ✅ Verify what data is being exported
- ✅ Check if data is transmitted securely
- ✅ Look for security audits or reviews

## GDPR and Privacy Regulations

BLEF files may contain personal data subject to GDPR and similar regulations.

**For users**: You have the right to:
- Export your data (BLEF helps with this!)
- Delete your data
- Know what data is stored
- Correct inaccurate data

**For tool developers**: You must:
- Obtain consent before collecting data
- Provide easy export (BLEF is ideal for this)
- Allow users to delete their data
- Secure personal data appropriately
- Respect user privacy preferences

## Security Updates

Security-related changes to the specification or schema will be announced:

- GitHub Security Advisories
- GitHub Releases with `[SECURITY]` tag
- Project discussions
- CHANGELOG.md

Subscribe to repository notifications to stay informed.

## Questions?

For general security questions (not vulnerabilities):
- Open a [GitHub Discussion](https://github.com/yoanbernabeu/BLEF/discussions)
- Check existing security-related issues

---

**Thank you for helping keep BLEF and its community safe!** 🔒

*Last updated: October 2025*

