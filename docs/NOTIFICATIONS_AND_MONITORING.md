# Notifications and Monitoring Dashboard

## Overview

Repo-lyzer now includes a comprehensive notification system and monitoring dashboard to track repository analyses, exports, and real-time repository changes.

## Features

### 1. Notifications Page

The notifications page displays a chronological list of all activities:

- **Repository Analyses**: Track when repositories are analyzed
- **Export Operations**: Monitor successful and failed exports (JSON, PDF, Markdown, CSV, HTML)
- **Monitoring Alerts**: View real-time updates from monitored repositories

#### Accessing Notifications

**From the Menu:**
- Navigate to "👀 Monitor Repository" in the main menu
- Press `m` as a keyboard shortcut

**From CLI:**
```bash
# View notifications directly
repo-lyzer notifications

# Alias commands
repo-lyzer notif
repo-lyzer notify
```

#### Keyboard Shortcuts

- `↑/↓` or `j/k`: Navigate through notifications
- `g`: Jump to first notification
- `G`: Jump to last notification
- `d`: Delete selected notification
- `c`: Clear all notifications
- `q/ESC`: Return to main menu

### 2. Monitoring Dashboard

The monitoring dashboard provides real-time tracking of repository changes.

#### Features

- **Real-time Updates**: Monitor commits, issues, PRs, and contributor changes
- **Configurable Intervals**: Set custom check intervals (default: 5 minutes)
- **Auto-scroll**: Automatically scroll to latest notifications
- **Persistent History**: Keep up to 50 recent notifications per session

