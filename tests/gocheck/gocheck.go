// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gocheck

import (
	"fmt"

	"360.cn/apollo/apollo/gocheck/aligncheck"
	"360.cn/apollo/apollo/gocheck/structcheck"
	"360.cn/apollo/apollo/gocheck/varcheck"
)

func Check(dirs map[string]string) map[string][]string {
	checks := make(map[string][]string)
	fmt.Println("Checking var...")
	for _, dir := range dirs {
		varchecks, err := varcheck.VarCheck(string(dir))
		if err != nil {
			fmt.Println(err)
		} else {
			checks[dir] = append(checks[dir], varchecks...)
		}
	}
	fmt.Println("Checking struct...")
	for _, dir := range dirs {
		structchecks, err := structcheck.StructCheck(string(dir))
		if err != nil {
			fmt.Println(err)
		} else {
			checks[dir] = append(checks[dir], structchecks...)
		}
	}
	fmt.Println("Checking align...")
	for _, dir := range dirs {
		alignchecks, err := aligncheck.AlignCheck(string(dir))
		if err != nil {
			fmt.Println(err)
		} else {
			checks[dir] = append(checks[dir], alignchecks...)
		}
	}

	return checks
}
