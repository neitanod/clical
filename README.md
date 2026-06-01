**English** · [Español](README.es.md)

# clical

**clical** is a multi-user calendar system with a command-line interface (CLI), designed to be AI-assisted.

## Main Features

- 🗓️ **Multi-user CLI calendar** - Full event management from the command line
- 📝 **Markdown + JSON storage** - Human-readable, manually editable files
- 🤖 **Optimized for AI assistance** - Reports designed for an AI to help proactively
- 📁 **Hierarchical organization** - Events organized by year/month/day
- 🔍 **Search and filtering** - Powerful search capabilities
- 📊 **Smart reports** - Daily, weekly and upcoming reports for AI consumption
- ⏰ **Cron-compatible** - Scheduled report execution

## Installation

### From source (Linux / macOS)

```bash
# Clone the repository
git clone https://github.com/neitanod/clical.git
cd clical

# Build
make build

# Install to /usr/local/bin (optional)
make install-system
```

### From source (Windows)

On Windows use the equivalent PowerShell scripts (no `make` required):

```powershell
git clone https://github.com/neitanod/clical.git
cd clical

# Build (produces .\clical.exe)
.\build.ps1

# Install to %USERPROFILE%\go\bin (optional)
.\install.ps1
```

For `clical` to be available in any terminal, make sure
`%USERPROFILE%\go\bin` is in your user `PATH`.

### AI-agent-assisted installation

If you use an agent with terminal access (Claude Code, Cursor, etc.), you can
install clical by pasting this prompt:

<https://github.com/neitanod/clical/blob/main/install_prompt.md>

The prompt guides the agent through OS detection, requirement checks, cloning,
building and installing on Linux/macOS or Windows.

### Requirements

- Go 1.23 or higher

## Quick Start

### 1. Create a user

```bash
clical user add --id=12345 --name="Your Name" --timezone="America/Argentina/Buenos_Aires"
```

### 2. Add an event

```bash
clical add --user=12345 \
  --datetime="2025-11-21 14:00" \
  --title="Client meeting" \
  --duration=60 \
  --location="Main Office" \
  --notes="Review Q4 proposal"
```

### 3. List events

```bash
# All events
clical list --user=12345

# Today's events
clical list --user=12345 --range=today

# This week's events
clical list --user=12345 --range=week
```

### 4. View daily report

```bash
clical daily-report --user=12345
```

## Available Commands

### User Management

```bash
# Create user
clical user add --id=ID --name="Name" --timezone="Timezone"

# List users
clical user list

# Show user details
clical user show --id=ID
```

### Event Management

```bash
# Add event
clical add --user=ID --datetime="YYYY-MM-DD HH:MM" --title="Title" [options]

# List events
clical list --user=ID [--from=DATE] [--to=DATE] [--range=RANGE] [--tags=TAG1,TAG2]

# Show event
clical show --user=ID --id=EVENT_ID

# Edit event
clical edit --user=ID --id=EVENT_ID [--title="New"] [--datetime="YYYY-MM-DD HH:MM"]

# Delete event
clical delete --user=ID --id=EVENT_ID [--force]
```

### AI Reports

```bash
# Full daily report
clical daily-report --user=ID [--date=YYYY-MM-DD]

# Tomorrow's report
clical tomorrow-report --user=ID

# Upcoming events
clical upcoming-report --user=ID --hours=2
clical upcoming-report --user=ID --count=5

# Weekly report
clical weekly-report --user=ID
```

### Other

```bash
# Version
clical version

# Help
clical --help
clical COMMAND --help
```

## Using with Cron

### Automatic reports

Edit your crontab:

```bash
crontab -e
```

Add lines like:

```bash
# Daily report at 7:00 AM
0 7 * * * /usr/local/bin/clical daily-report --user=12345 | mail -s "Today's Agenda" you@email.com

# Tomorrow's report at 8:00 PM
0 20 * * * /usr/local/bin/clical tomorrow-report --user=12345 | mail -s "Tomorrow's Agenda" you@email.com

# Hourly alerts during work hours
0 9-18 * * * /usr/local/bin/clical upcoming-report --user=12345 --hours=2 | mail -s "Upcoming Events" you@email.com

# Weekly report on Mondays at 7:00 AM
0 7 * * 1 /usr/local/bin/clical weekly-report --user=12345 | mail -s "Weekly Agenda" you@email.com
```

### Telegram integration (if you have a bot)

```bash
# Daily report via Telegram
0 7 * * * OUTPUT=$(/usr/local/bin/clical daily-report --user=12345); [ -n "$OUTPUT" ] && your-telegram-command "$OUTPUT"
```

## Storage Format

Data is stored under `~/.clical/data/` (configurable with `--data-dir`):

```
~/.clical/data/
└── users/
    └── 12345/
        ├── user.md              # User info (Markdown)
        ├── user.json            # User metadata
        ├── events/
        │   └── 2025/
        │       └── 11/
        │           └── 21/
        │               ├── 09-00-stand-up-meeting.md
        │               ├── 09-00-stand-up-meeting.json
        │               ├── 14-00-client-meeting.md
        │               └── 14-00-client-meeting.json
        └── .state/
            └── report-state.json
```

### Markdown file example

```markdown
# Client meeting

**Date:** 2025-11-21
**Time:** 14:00
**Duration:** 60 minutes
**Location:** Main Office
**Tags:** #work #client

## Notes

Review Q4 proposal and discuss timeline.

---

*Created: 2025-11-20 16:18*
*Updated: 2025-11-20 16:18*
*ID: e36e10014ea57372*
```

## Configuration

### Environment Variables

```bash
# Data directory (default: ~/.clical/data)
export CLICAL_DATA_DIR="/custom/path/data"

# Default user
export CLICAL_USER_ID="12345"
```

### Common Timezones

- `America/Argentina/Buenos_Aires`
- `America/Mexico_City`
- `America/New_York`
- `Europe/Madrid`
- `UTC`

Full list: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones

## Using with AI

clical is designed for an AI to assist you proactively. Example flow:

1. **07:00 AM** - Cron runs `daily-report`
   - AI receives today's agenda
   - AI greets you and presents the events
   - AI identifies pending tasks
   - AI suggests how to organize the day

2. **During the day** - Cron runs `upcoming-report` every hour
   - AI alerts you about imminent events
   - AI reminds you of needed preparation

3. **20:00 PM** - Cron runs `tomorrow-report`
   - AI presents tomorrow's view
   - AI suggests evening preparation

## Development

### Project Structure

```
clical/
├── cmd/clical/       # Entry point
├── pkg/              # Public packages
│   ├── calendar/     # Entry, Filter models
│   ├── storage/      # Filesystem storage
│   ├── user/         # User management
│   └── reporter/     # Report generation
├── internal/         # Private packages
│   ├── cli/          # Cobra commands
│   └── config/       # Configuration
├── ai/               # Development docs
│   ├── specs/        # Specifications
│   └── journal/      # Development journal
└── docs/             # Web documentation
```

### Build

```bash
make build
```

### Tests

```bash
make test
```

### Format code

```bash
make fmt
```

## Contributing

See [ai/specs/00_Overview.md](ai/specs/00_Overview.md) for full specifications.

## License

MIT

## Author

Developed by Sebastián Valencia with Claude (Anthropic) assistance.

---

**Version:** 0.1.0
**Status:** Working MVP - Under active development
