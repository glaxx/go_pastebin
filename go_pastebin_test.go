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
	"log"
	"testing"
)

const (
	testUser  = ""
	testPass  = ""
	apiDevKey = ""
)

func Test_PasteAnonymous(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	ret, err := pb.PasteAnonymous("wohoo", "api_test", "text", expire10Minutes, exposureUnlisted)
	if err != nil {
		log.Println(ret, err)
		t.Error("Paste failed, 1")
	} else {
		t.Log("PasteAnonymous passed with URL:", ret)
	}
}

func Test_PasteAnonymousSimple(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	ret, err := pb.PasteAnonymousSimple("wohoo")
	if err != nil {
		log.Println(ret, err)
		t.Error("Paste failed, 1")
	} else {
		t.Log("PasteAnonymousSimple passed with URL:", ret)
	}
}

func Test_GenerateUserSession(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	s, err := pb.GenerateUserSession(testUser, testPass)
	if err != nil {
		log.Println(err, s)
		t.Error("GenerateUserSession failed, 1")
	} else {
		t.Log("GenerateUserSession passed", s)
	}
}

func Test_Paste(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	s, err := pb.GenerateUserSession(testUser, testPass)
	if err != nil {
		t.Error("Paste failed at session creation")
	}

	ret, err := s.Paste(apiDevKey, "wohoo", "api_test", "text", expire10Minutes, exposureUnlisted)
	if err != nil {
		t.Error("Paste fialed at pasting")
	} else {
		t.Log("Paste passed with URL:", ret)
	}
}

func Test_ListPastes(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	s, err := pb.GenerateUserSession(testUser, testPass)
	if err != nil {
		t.Error("List failed at session creation")
	}
	pas, err := s.ListPastes(10)
	if err != nil {
		t.Error("List failed at gathering the list")
	} else {
		t.Log("List passed with gathered elements:", pas)
	}
}

func Test_ListTrendingPastes(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	p, err := pb.ListTrendingPastes()
	if err != nil {
		t.Log(err)
		t.Error("ListTrendingPastes failed")
	} else {
		t.Log("ListTrendingPastes passed with results:", p)
	}
}

func Test_DeletePaste(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	p, err := pb.GenerateUserSession(testUser, testPass)
	if err != nil {
		t.Error("DeletePaste failed at session creation")
	}
	pas, err := p.ListPastes(10)
	if err != nil {
		t.Error("DeletePaste failed at list fetch")
	}
	err = pb.DeletePaste(pas[0].PasteKey)
	if err != nil {
		t.Log(err)
		t.Error("DeletePaste failed at deleting")
	} else {
		t.Log("DeletePaste passed")
	}
}

func Test_PasteAnonymousBadAPIError(t *testing.T) {
	pb := NewPastebin(apiDevKey)
	ret, err := pb.PasteAnonymous("wohoo", "api_test", "text", expire10Minutes, "4")
	if err != nil && ret == nil {
		t.Error("PasteAnonymousBadAPIError failed")
	} else {
		t.Log("PasteAnonymousBadAPIError passed with err:", err)
	}
}
