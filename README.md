# Encrypted Chunk Mailer

**Encrypted Chunk Mailer** is a Go-based CLI tool that allows users to **split, encrypt, send, and retrieve files via email** using MIME formatting and Gmail/IMAP integration. It is designed to support secure file sharing over common email services, with support for encryption and chunked delivery.

---

## Project Status

> ⚠️ This project is currently under development.

- Attachment sending speed is being optimized.
- Automatic retrieval of attachments from email is in progress.
- Google Drive API integration may be introduced as an alternative to email-based delivery.
- A `--help` flag and improved CLI UI&UX are planned.

---

## Features

- ✅ File splitting and AES-GCM encryption
- ✅ MIME-formatted email attachments with chunked delivery
- ✅ Optional ZIP compression for file chunks
- ✅ IMAP-based attachment reading
- ✅ Secure Gmail transmission using App Passwords

---

## Prerequisites

- Go 1.18 or newer
- Enable IMAP in your Gmail settings
- Create a [Gmail App Password](https://support.google.com/accounts/answer/185833?hl=en) if 2FA is enabled

---

## Installation

```bash
git clone https://github.com/yourusername/encrypted-chunk-mailer.git
cd encrypted-chunk-mailer
go build
```
## Usage
### 1. Add a File
```
go run . add filename.ext
```
- Splits and encrypts the provided file into multiple chunks.

- Adds the resulting file and key metadata to MetaData.json, which acts as the internal metadata store.

### 2. Send Encrypted Chunks via Email
```
go run . push <fileID> recipient@example.com
```
- Sends the encrypted chunks as MIME attachments to the provided email address.

- Note: This process may be slow for large files. Optimization or Google Drive upload support is in development.

- Recommended for smaller files such as images or documents.

### 3. Manually Reconstruct File from Chunks
```
go run . pull chunkFile1 chunkFile2 ... <decryptionKey>
```
- Decrypts and reconstructs the original file from its encrypted chunks using the provided key.

- Currently a manual process—automatic retrieval is planned.

# Important Notes

    This is the developer’s first Go project and is actively evolving.

    Use caution when sharing or storing decryption keys.

    Contributions, feedback, and bug reports are welcome.
