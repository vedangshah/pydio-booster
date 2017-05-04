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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	str string
)

func init() {
	str = `{"OPTIONS":{"TYPE":"local"},"PATH":"\/private\/var\/tmp\/admin_1490284406tmpDownload.zip","UNLINK_AFTER_READ":true}`
}

func TestOptions(t *testing.T) {
	Convey("Testing encoding options via a channel", t, func() {

		c := make(chan interface{})

		go func() {
			c <- str
		}()
	})
}