// Package pydio contains all objects needed by the Pydio system
/*
 * Copyright 2007-2016 Abstrium <contact (at) pydio.com>
 * This file is part of Pydio.
 *
 * Pydio is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com/>.
 */
package pydio

import (
	"encoding/json"
	"testing"

	"github.com/pydio/pydio-booster/encoding/path"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	fakeRepo *Repo
)

func init() {
	fakeRepo = &Repo{
		ID: "repo",
	}
}

func TestUnmarshalRepo(t *testing.T) {
	Convey("Testing unmarshal a Repo path format", t, func() {
		var urlBlob = []byte("/repo/dir1/file1.txt")
		var repo Repo
		err := path.Unmarshal(urlBlob, &repo)
		So(err, ShouldBeNil)

		So(repo, ShouldResemble, *fakeRepo)
	})

	Convey("Testing decoding ", t, func() {
		var repo Repo

		dec := json.NewDecoder(fakeRepo)

		err := dec.Decode(&repo)

		So(err, ShouldBeNil)

		So(&repo, ShouldResemble, fakeRepo)

	})
}
