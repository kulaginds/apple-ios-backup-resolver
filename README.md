# apple-ios-backup-resolver
Convert backup-format files to human-readable file structure.

## What is iOS backups?
It's folder in path `~/Library/Application Support/MobileSync/Backup` os *nix system or `%appdata%\Apple\MobileSync\Backup` on windows.

Backup folder contains some folders called like UID (ex: 00008030-001950C90AE9802E). Each folder - one backup.

Each backup folder contains two-letters named folders and some files:

- Info.plist
- Manifest.db (*) - main file
- Manifest.plist
- Status.plist

Manifest.db contains real file names, which was on iPhone.

## How to run?
1. Build program from source code:
    ```bash
    go mod tidy
    make build 
    ```
    After this command in current directory was created executable file `apple-ios-backup-resolver`.
2. Run program:
   ```bash
   ./apple-ios-backup-resolver -src {SOURCE_DIR} -dst {DESTINATION_DIR}
   ```
   Where:
   - {SOURCE_DIR} - backup directory (ex: ~/Library/Application Support/MobileSync/Backup/00008030-001950C90AE9802E)
   - {DESTINATION_DIR} - folder with human-readable file structure (ex: ~/Downloads/iPhoneBackup)

In my case on 8GB backup on Macbook Pro program works approximately 5 minutes.
