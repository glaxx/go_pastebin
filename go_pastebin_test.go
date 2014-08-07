/*
 *    Copyright (C) 2014 Stefan Luecke
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors: Stefan Luecke <glaxx@glaxx.net>
 */

package go_pastebin

import (
 	"testing"
 	"log"
)

const (
	test_user = ""
	test_pw = ""
)

func Test_PasteAnonymous(t *testing.T) {
	ret, err := PasteAnonymous("wohoo", "api_test", "text", expire_10Minutes, private_unlisted)
	if err != nil {
		log.Println(ret, err)
		t.Error("Paste failed, 1")
	} else {
		t.Log("PasteAnonymous passed with URL:", ret)
	}
}

func Test_GenerateUserSession(t *testing.T) {
	s, err := GenerateUserSession(test_user, test_pw)
	if err != nil {
		log.Println(err, s)
		t.Error("GenerateUserSession failed, 1")
	} else {
		t.Log("GenerateUserSession passed", s)
	}
}

func Test_Paste(t *testing.T) {
	s, err := GenerateUserSession(test_user, test_pw)
	if err != nil {
		t.Error("Paste failed at session creation")
	}

	ret, err := s.Paste("wohoo", "api_test", "text", expire_10Minutes, private_unlisted)
	if err != nil {
		t.Error("Paste fialed at pasting")
	} else {
		t.Log("Paste passed with URL:", ret)
	}
}

func Test_ListPastes(t *testing.T) {
	s, err := GenerateUserSession(test_user, test_pw)
	if err != nil {
		t.Error("List failed at session creation")
	}
	pas, errr := s.ListPastes(10)
	if errr != nil {
		t.Error("List failed at gathering the list")
	} else {
		t.Log("List passed with gathered elements:", pas)
	}

}