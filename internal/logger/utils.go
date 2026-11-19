// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package logger

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime"

	"github.com/core-laboratories/maxio/internal/color"
)

var ansiRE = regexp.MustCompile("(\x1b[^m]*m)")

// Print ANSI Control escape
func ansiEscape(format string, args ...any) {
	Esc := "\x1b"
	fmt.Printf("%s%s", Esc, fmt.Sprintf(format, args...))
}

func ansiMoveRight(n int) {
	if runtime.GOOS == "windows" {
		return
	}
	if color.IsTerminal() {
		ansiEscape("[%dC", n)
	}
}

func ansiSaveAttributes() {
	if runtime.GOOS == "windows" {
		return
	}
	if color.IsTerminal() {
		ansiEscape("7")
	}
}

func ansiRestoreAttributes() {
	if runtime.GOOS == "windows" {
		return
	}
	if color.IsTerminal() {
		ansiEscape("8")
	}
}

// logIgnoreError if true,the error will ignore.
func logIgnoreError(err error) bool {
	return err == nil || errors.Is(err, context.Canceled) || errors.Is(err, http.ErrServerClosed) || err.Error() == "drive not found"
}
