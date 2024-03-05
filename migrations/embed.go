// Copyright (c) 2024 Vasiliy Vasilyuk. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migrations

import "embed"

//go:embed *.up.sql
var FS embed.FS
