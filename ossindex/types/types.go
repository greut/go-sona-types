//
// Copyright 2018-present Sonatype Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package types

import (
	"fmt"

	"github.com/shopspring/decimal"
)

const (
	OssIndexDirName        = ".ossindex"
	OssIndexConfigFileName = ".oss-index-config"
	IQServerDirName        = ".iqserver"
	IQServerConfigFileName = ".iq-server-config"
)

type Configuration struct {
	Version    string
	CleanCache bool
	Username   string
	Token      string
}

type Coordinate struct {
	Coordinates     string
	Reference       string
	Vulnerabilities []Vulnerability
	InvalidSemVer   bool
}

type Vulnerability struct {
	ID          string
	Title       string
	Description string
	CvssScore   decimal.Decimal
	CvssVector  string
	Cve         string
	Reference   string
	Excluded    bool
}

func (c Coordinate) IsVulnerable() bool {
	for _, v := range c.Vulnerabilities {
		if !v.Excluded {
			return true
		}
	}
	return false
}

//Mark Excluded=true for all Vulnerabilities of the given Coordinate if their Title is in the list of exclusions
func (c *Coordinate) ExcludeVulnerabilities(exclusions []string) {
	for i := range c.Vulnerabilities {
		c.Vulnerabilities[i].maybeExcludeVulnerability(exclusions)
	}
}

//Mark the given vulnerability as excluded if it appears in the exclusion list
func (v *Vulnerability) maybeExcludeVulnerability(exclusions []string) {
	for _, ex := range exclusions {
		if v.Cve == ex || v.ID == ex {
			v.Excluded = true
		}
	}
}

type AuditRequest struct {
	Coordinates []string `json:"coordinates"`
}

// OSSIndexRateLimitError is a custom error implementation to allow us to return a better error response to the user
// as well as check the type of the error so we can surface this information.
type OSSIndexRateLimitError struct {
}

func (o *OSSIndexRateLimitError) Error() string {
	return `You have been rate limited by OSS Index.
If you do not have a OSS Index account, please visit https://ossindex.sonatype.org/user/register to register an account.
After registering and verifying your account, you can retrieve your username (Email Address), and API Token
at https://ossindex.sonatype.org/user/settings. Upon retrieving those, run 'nancy config', set your OSS Index
settings, and rerun Nancy.`
}

type OSSIndexError struct {
	Err     error
	Message string
}

func (o *OSSIndexError) Error() string {
	return fmt.Sprintf("An error occurred: %s, err: %e", o.Message, o.Err)
}