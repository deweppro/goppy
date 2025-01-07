/*
 *  Copyright (c) 2022-2025 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package env

type (
	// ENV type for environments (prod, dev, stage, etc)
	ENV string

	AppName        string
	AppVersion     string
	AppDescription string

	AppInfo struct {
		AppName        AppName
		AppVersion     AppVersion
		AppDescription AppDescription
	}
)

func NewAppInfo() AppInfo {
	return AppInfo{
		AppName:        "",
		AppVersion:     "",
		AppDescription: "",
	}
}
