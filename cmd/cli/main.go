/*
Copyright (C) 2022 - now  Rat Technology Co., Ltd.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"os"
	_ "time/tzdata"
)

func main() {
	if os.Geteuid() != 0 {
		panic("panel must run as root")
	}

	cli, err := initCli()
	if err != nil {
		panic(err)
	}

	if err = cli.Run(); err != nil {
		panic(err)
	}
}
