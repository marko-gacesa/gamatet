// Copyright (c) 2025 by Marko Gaćeša
// Licensed under the GNU GPL v3 or later. See the LICENSE file for details.

package i18n

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/marko-gacesa/gamatet/logic/lang"
)

//go:embed translations/*.lang
var translationsFS embed.FS

type Key string

func (key Key) String() string {
	return lang.Str(string(key))
}

func T(key Key) string {
	return lang.Str(string(key))
}

func Tf(key Key, args ...any) string {
	return fmt.Sprintf(lang.Str(string(key)), args...)
}

var reFileName = regexp.MustCompile(`^(\w+)\.lang$`)
var reLine = regexp.MustCompile(`^(\w+)\s*=\s*(.*)\s*$`)

var escaper = strings.NewReplacer(
	"\\n", "\n",
	"\\t", "\t")

func ParseEmbeddedLanguages(logger *slog.Logger) {
	const directory = "translations"

	translations, err := translationsFS.ReadDir(directory)
	if err != nil {
		logger.Error("Failed to list embedded files", "error", err)
		return
	}

	for _, translation := range translations {
		if translation.IsDir() {
			continue
		}

		fileName := translation.Name()

		partsFileName := reFileName.FindStringSubmatch(fileName)
		if len(partsFileName) != 2 {
			continue
		}

		language := lang.Lang(partsFileName[1])

		data, err := translationsFS.ReadFile(directory + "/" + fileName)
		if err != nil {
			logger.Warn("Failed to open embedded file", "file_name", fileName, "error", err)
			continue
		}

		m := make(map[string]string)

		scan := bufio.NewScanner(bytes.NewReader(data))
		lineNumber := 0
		for scan.Scan() {
			line := scan.Text()
			lineNumber++

			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := reLine.FindStringSubmatch(line)
			if len(parts) != 3 {
				logger.Warn("Invalid line",
					"language", language, "line_number", lineNumber, "line", line)
				continue
			}

			key := parts[1]
			value := parts[2]

			value = escaper.Replace(value)

			if _, ok := m[key]; ok {
				logger.Warn("Duplicate entry",
					"language", language, "line_number", lineNumber, "key", key, "value", value)
				continue
			}

			m[key] = value
		}

		name := m[KeyLanguageName]
		if name == "" {
			logger.Warn("Language name is missing", "language", language)
			continue
		}

		lang.Define(language, m)
	}
}
