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
package pydhttp

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToken(t *testing.T) {

	token := NewToken("vXqzNCtQ5R7odtjIzVqB8OMW", "Ysz4npNH1KoOfRVfPH2J12Ia")

	Convey("Sending a simple request", t, func() {

		uri := "/api/ajxp_conf/scheduler_runAll"
		auth := token.GetQueryArgs(uri)

		So(auth, ShouldNotBeNil)
		So(auth.Token, ShouldNotBeNil)
		So(auth.Hash, ShouldNotBeNil)

		fmt.Printf("\n\n-----------\n\n")
		fmt.Println("Results")
		fmt.Printf("=========\n\n")
		fmt.Printf("URI : %v\n", uri)
		fmt.Printf("auth_token : %v\n", auth.Token)
		fmt.Printf("auth_hash: %v\n", auth.Hash)
	})
}
