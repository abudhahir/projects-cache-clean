# 📸 Generated Screenshots

This directory contains generated screenshots for Cache Remover documentation.

## 🚧 Current Status

The screenshot files currently exist as placeholders. To generate actual GIF screenshots:

1. **Install dependencies**:
   ```bash
   go install github.com/charmbracelet/vhs@latest
   brew install ttyd
   ```

2. **Generate screenshots**:
   ```bash
   cd screenshots-utility
   ./setup-test-data.sh
   ./generate.sh
   ```

3. **Or generate individual screenshots**:
   ```bash
   vhs tapes/dry-run.tape
   vhs tapes/basic-usage.tape
   vhs tapes/verbose.tape
   # etc.
   ```

## 📋 Screenshot Files

| File | Purpose | Documentation |
|------|---------|---------------|
| `dry-run.gif` | Safe preview mode demo | QUICKSTART.md, USAGE.md |
| `basic-usage.gif` | Standard cleanup workflow | QUICKSTART.md |
| `verbose.gif` | Detailed output with optimizations | USAGE.md |
| `interactive.gif` | Per-project confirmation | USAGE.md |
| `ui-demo.gif` | TUI interface demonstration | README.md, QUICKSTART.md |
| `performance.gif` | Optimization comparison | README.md, USAGE.md |
| `quickstart.gif` | Complete workflow demo | README.md |

## 🎯 Integration

These screenshots are already referenced in:
- **README.md** - Visual demo section
- **QUICKSTART.md** - Visual examples section  
- **USAGE.md** - Visual examples section

The documentation will show the screenshots once they are generated.

## 💡 Notes

- Screenshots are generated using VHS (Video Tape Simulator)
- Each screenshot corresponds to a `.tape` file in the `tapes/` directory
- Test data is automatically created by `setup-test-data.sh`
- Screenshots show actual terminal sessions with the cache remover