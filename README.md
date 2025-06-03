# Encrypted Chunk Mailer

A secure, Go-based CLI tool for splitting, encrypting, and transmitting files via email with MIME formatting and Gmail integration.

## Overview

Encrypted Chunk Mailer enables secure file sharing through email by automatically splitting files into encrypted chunks, sending them as email attachments, and providing tools for reconstruction. This approach allows you to securely share files of any size and also get the file anywhere by using your gmail id like a cloud too, while maintaining privacy and bypassing traditional file size limitations of email and cloud services.

## Key Features

- **AES-GCM Encryption**: Military-grade encryption for all file chunks
- **Email Integration**: Seamless Gmail support with App Password authentication
- **Intelligent Chunking**: Automatic file splitting for optimal transmission
- **Privacy First**: Files are unrecognizable without the decryption tool
- **Used as Cloud**: Use your gmail id anywhere to get those chunks
- **Cross-Platform Sharing**: Share chunks via email, messaging apps, or any platform

## Installation

### Prerequisites

- **Go 1.18+** - [Download Go](https://golang.org/dl/)
- **App Password** (required if 2FA is enabled) - [Setup Guide](https://support.google.com/accounts/answer/185833?hl=en)

### Setup

```bash
git clone https://github.com/SureshS03/P2PFile.git
cd P2PFile
go build -o sender
```

## Usage

### Adding Files for Encryption

Split and encrypt a file into secure chunks:

```bash
./sender add filename.ext
```

This command:
- Splits the file into manageable chunks
- Encrypts each chunk using AES-GCM
- Stores metadata in `MetaData.json`
- Generates a unique file ID and decryption key

### Sending via Email

Automatically send encrypted chunks to a recipient:

```bash
./sender push <fileID> recipient@gmail.com
```

**Recommended for:**
- Images and documents (< 25MB)
- Quick automated delivery
- Trusted recipients with reliable email

### Manual File Reconstruction

Decrypt and reconstruct the original file:

```bash
./sender pull chunkFile1 chunkFile2 chunkFile3 <decryptionKey>
```

**Use when:**
- Chunks were shared manually
- Working with locally stored chunks

### Clearing Metadata

Reset the metadata storage:

```bash
./sender clear
```

## Workflow Examples

### Scenario 1: Secure Document Sharing
```bash
# Encrypt a confidential document
./sender add contract.pdf

# Send to colleague via email
./sender push <FileId> colleague@gamil.com

# Colleague receives chunks and key separately
# They reconstruct: ./sender pull chunk1 chunk2 chunk3 <key>
```

### Scenario 2: Large File via Multiple Channels
```bash
# Encrypt a large presentation
./sender add presentation.pptx

# Share chunks via different platforms (WhatsApp, Telegram, etc.)
# Share decryption key through secure channel
# Recipient reconstructs locally
```

## Security Considerations

- **Key Management**: Store decryption keys separately from chunks
- **Email Security**: Use App Passwords, never regular passwords
- **Transmission**: Consider sharing keys through different channels than chunks
- **Local Storage**: Securely delete original files after encryption if needed

## Current Limitations & Roadmap

### Known Issues
- Email transmission speed optimization in progress
- Large files (>50MB) may required more time

### Upcoming Features
- [ ] Automatic email attachment retrieval using IMAP
- [ ] Google Drive API integration
- [ ] Comprehensive `--help` documentation
- [ ] Progress indicators and improved UX
- [ ] Batch processing support

## Configuration

### Gmail Setup
1. Generate an App Password (if 2FA enabled)
2. Use App Password instead of regular password

### Metadata Structure
The tool maintains a `MetaData.json` file containing:
- File IDs and original names
- Chunk information and locations
- Encryption metadata (keys stored separately for security)

## Contributing

This project welcomes contributions! Areas where help is needed:

### Getting Started
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT License - See LICENSE file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/SureshS03/P2PFile/issues)
- **Email**: suthamani51@gmail.com

---

> **Security Notice**: This tool is designed for legitimate file sharing purposes. Always comply with your organization's security policies and local regulations when transmitting sensitive data.
