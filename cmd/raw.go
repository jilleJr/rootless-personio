// SPDX-FileCopyrightText: 2023 Kalle Fagerberg
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jilleJr/rootless-personio/pkg/personio"
	"github.com/jilleJr/rootless-personio/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rawFlags = struct {
	method   string
	data     string
	headers  []string
	formData []string
	noLogin  bool
}{}

var rawCmd = &cobra.Command{
	Use:   "raw </url/path?query=param>",
	Short: "Send a raw HTTP request to the API",
	Long: `Send a raw HTTP request to the API
as a logged in user, and print the resulting JSON data.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		urlArg := args[0]
		baseURL, err := getBaseURL(urlArg)
		if err != nil {
			return err
		}
		client, err := personio.New(baseURL)
		if err != nil {
			return err
		}
		log.Debug().Str("baseUrl", client.BaseURL).Msg("Created valid client.")

		if !rawFlags.noLogin {
			var missingCredentials bool
			if cfg.Auth.Email == "" {
				missingCredentials = true
				log.Error().Msg("Missing email! Must set auth.email config or PERSONIO_AUTH_EMAIL env var.")
			}
			if cfg.Auth.Password == "" {
				missingCredentials = true
				log.Error().Msg("Missing password! Must set auth.password config or PERSONIO_AUTH_PASSWORD env var.")
			}
			if missingCredentials {
				return errors.New("missing credentials")
			}
			if err := client.Login(cfg.Auth.Email, cfg.Auth.Password); err != nil {
				return err
			}
			log.Info().Int("employeeId", client.EmployeeID).
				Msg("Successfully logged in.")
		}

		method := http.MethodGet

		body, err := getDataFlagReader(rawFlags.data)
		if err != nil {
			return err
		}
		if body != nil {
			method = http.MethodPost
			defer body.Close()
		}

		if rawFlags.method != "" {
			method = rawFlags.method
		}

		req, err := http.NewRequest(method, urlArg, body)
		if err != nil {
			return err
		}
		respModel, err := client.Raw(req)
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(respModel, "", "  ")
		if err != nil {
			return err
		}
		if colorized, err := util.ColorizeJSON(b); err == nil {
			fmt.Println(string(colorized))
		} else {
			fmt.Println(string(b))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rawCmd)

	rawCmd.Flags().StringVarP(&rawFlags.method, "method", "X", rawFlags.method, `Request method to use (default "POST" if with --data flag, otherwise "GET")`)
	rawCmd.Flags().StringVarP(&rawFlags.data, "data", "d", rawFlags.data, `Request body ("@filename" for reading from file, or "@-" from STDIN)`)
	rawCmd.Flags().StringArrayVarP(&rawFlags.headers, "header", "H", nil, `Add custom headers to request (format "Key: value")`)
	rawCmd.Flags().StringArrayVarP(&rawFlags.formData, "form", "F", nil, `Add multipart MIME data (format "key=value")`)
	rawCmd.Flags().BoolVar(&rawFlags.noLogin, "no-login", false, `Skip logging in before the request`)
}

func getBaseURL(urlArg string) (string, error) {
	if cfg.BaseURL != "" {
		return cfg.BaseURL, nil
	}
	u, err := url.Parse(urlArg)
	if err != nil {
		return "", fmt.Errorf("parse positional arg as URL: %w", err)
	}
	if u.Host == "" {
		return "", errors.New("must provide url config or hostname in positional arg")
	}
	u.Path = ""
	return u.String(), nil
}

func getDataFlagReader(dataFlag string) (io.ReadCloser, error) {
	if dataFlag == "" {
		return nil, nil
	}
	if dataFlag[0] == '@' {
		path := dataFlag[1:]

		if path == "-" {
			return os.Stdin, nil
		}

		file, err := os.Open(path)
		return file, err
	}
	strReader := strings.NewReader(dataFlag)
	return io.NopCloser(strReader), nil
}
