// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	humanize "github.com/dustin/go-humanize"
	"github.com/core-laboratories/maxio/internal/color"
)

// prepareUpdateMessage - prepares the update message, only if a
// newer version is available.
func prepareUpdateMessage(downloadURL string, older time.Duration) string {
	if downloadURL == "" || older <= 0 {
		return ""
	}

	// Compute friendly duration string to indicate time
	// difference between newer and current release.
	t := time.Time{}
	newerThan := humanize.RelTime(t, t.Add(older), "before the latest release", "")

	if globalServerCtxt.JSON {
		return fmt.Sprintf("You are running an older version of MaxIO released %s, update: %s", newerThan, downloadURL)
	}

	// Return the nicely colored and formatted update message.
	return colorizeUpdateMessage(downloadURL, newerThan)
}

// colorizeUpdateMessage - inspired from Yeoman project npm package https://github.com/yeoman/update-notifier
func colorizeUpdateMessage(updateString string, newerThan string) string {
	msgLine1Fmt := " You are running an older version of MaxIO released %s "
	msgLine2Fmt := " Update: %s "

	// Calculate length *without* color coding: with ANSI terminal
	// color characters, the result is incorrect.
	line1Length := len(fmt.Sprintf(msgLine1Fmt, newerThan))
	line2Length := len(fmt.Sprintf(msgLine2Fmt, updateString))

	// Populate lines with color coding.
	line1InColor := fmt.Sprintf(msgLine1Fmt, color.YellowBold(newerThan))
	line2InColor := fmt.Sprintf(msgLine2Fmt, color.CyanBold(updateString))

	// calculate the rectangular box size.
	maxContentWidth := max(line1Length, line2Length)

	// termWidth is set to a default one to use when we are
	// not able to calculate terminal width via OS syscalls
	termWidth := 25
	if width, err := pb.GetTerminalWidth(); err == nil {
		termWidth = width
	}

	// Box cannot be printed if terminal width is small than maxContentWidth
	if maxContentWidth > termWidth {
		return "\n" + line1InColor + "\n" + line2InColor + "\n\n"
	}

	topLeftChar := "┏"
	topRightChar := "┓"
	bottomLeftChar := "┗"
	bottomRightChar := "┛"
	horizBarChar := "━"
	vertBarChar := "┃"
	// on windows terminal turn off unicode characters.
	if runtime.GOOS == globalWindowsOSName {
		topLeftChar = "+"
		topRightChar = "+"
		bottomLeftChar = "+"
		bottomRightChar = "+"
		horizBarChar = "-"
		vertBarChar = "|"
	}

	lines := []string{
		color.YellowBold(topLeftChar + strings.Repeat(horizBarChar, maxContentWidth) + topRightChar),
		color.YellowBold(vertBarChar) + line1InColor + strings.Repeat(" ", maxContentWidth-line1Length) + color.YellowBold(vertBarChar),
		color.YellowBold(vertBarChar) + line2InColor + strings.Repeat(" ", maxContentWidth-line2Length) + color.YellowBold(vertBarChar),
		color.YellowBold(bottomLeftChar + strings.Repeat(horizBarChar, maxContentWidth) + bottomRightChar),
	}
	return "\n" + strings.Join(lines, "\n") + "\n"
}
