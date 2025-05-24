# Archivist 💾
Archivist is a command-line tool written in Go for archiving and unarchiving files and directories. It supports multiple compression formats, including ZIP, TAR, TAR.GZ, TAR.BZ2, and TAR.XZ, with automatic format detection for unpacking.

## Installation ⬇️
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/archivist.git
   cd archivist
2. Install dependencies:
   ```bash
   go mod tidy
3. Build the binary:
   ```bash
   go build -o archivist 
4. (Recommended/Optional) Move the bivary to a directory in your PATH:
   ```bash
   mv archivist /usr/local/bin/

## Usage 📊
Archivist provides two main commands: pack to create archives and unpack to extract them.
### Packing Files or Directories
Create an archive from files or directory
```bash
archivist pack -m <format> <file>
```
-m: Compression format (zip, tar, targ.gz, tar.bz, tar.xz)

### Example
```bash
archivist pack -m zip my_folder
```
### Unpacking Archive
```bash
archivist unpack [--method <format>] <archive>
```
### Example
```bash
archivist unpack my_folder.zip
```


