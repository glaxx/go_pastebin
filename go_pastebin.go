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
// This package offers functions to interact with the Pastebin API, for
// further information check http://pastebin.com/api .
package go_pastebin

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	api_dev_key = ""
	login_url   = "http://pastebin.com/api/api_login.php"
	post_url    = "http://pastebin.com/api/api_post.php"

	expire_never     = "N"
	expire_10Minutes = "10M"
	expire_1Hour     = "1H"
	expire_1Day      = "1D"
	expire_1Week     = "1W"
	expire_2Weeks    = "2W"
	expire_1Month    = "1M"

	private_public   = "0"
	private_unlisted = "1"
	private_private  = "2"
)

type Paste struct {
	Paste_key          string
	Paste_date         time.Time
	Paste_title        string
	Paste_size         int
	Paste_expire_date  time.Time
	Paste_private      int
	Paste_format_long  string
	Paste_format_short string
	Paste_url          *url.URL
	Paste_hits         int
}

type Session struct {
	api_user_key string
}

func PasteAnonymous(paste, title, format, expire, private string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}

	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)
	pasteOptions.Set("api_paste_name", title)
	pasteOptions.Set("api_paste_format", format)
	pasteOptions.Set("api_paste_expire_date", expire)
	pasteOptions.Set("api_paste_private", private)

	return pasteRequest(post_url, pasteOptions)
}

func PasteAnonymousSimple(paste string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)

	return pasteRequest(post_url, pasteOptions)
}

func ListTrendingPastes() (pastes []Paste, err error) {
	listOptions := url.Values{}
	listOptions.Set("api_option", "trends")

	p, err := listRequest(post_url, listOptions)
	if err != nil {
		return nil, err
	} else {
		return p, err
	}
}

// This function request (and returns) a Session key object.
func GenerateUserSession(username, password string) (se *Session, err error) {
	var s Session
	userOptions := url.Values{}
	userOptions.Set("api_user_name", username)
	userOptions.Set("api_user_password", password)
	userOptions.Set("api_dev_key", api_dev_key)

	resp, err := http.PostForm(login_url, userOptions)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// TODO: Catch BAD API request-answers

	s.api_user_key = string(body)

	return &s, nil
}

func (s *Session) ListPastes(result_limit int) (pastes []Paste, err error) {
	listOptions := url.Values{}
	listOptions.Set("api_user_key", s.api_user_key)
	listOptions.Set("api_option", "list")
	listOptions.Set("api_result_limit", strconv.Itoa(result_limit))
	p, err := listRequest(post_url, listOptions)
	if err != nil {
		return nil, err
	} else {
		return p, err
	}

}

func (s *Session) DeletePaste(paste_key string) (err error) {
	qryOptions := url.Values{}
	qryOptions.Set("api_paste_key", paste_key)
	qryOptions.Set("api_user_key", s.api_user_key)
	qryOptions.Set("api_dev_key", api_dev_key)
	qryOptions.Set("api_option", "delete")

	resp, err := http.PostForm(post_url, qryOptions)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "Paste Removed" {
		return errors.New(string(body))
	} else {
		return nil
	}
}

func (s *Session) Paste(paste, title, format, expire, private string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}

	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)
	pasteOptions.Set("api_paste_name", title)
	pasteOptions.Set("api_paste_format", format)
	pasteOptions.Set("api_paste_expire_date", expire)
	pasteOptions.Set("api_paste_private", private)

	pasteOptions.Set("api_user_key", s.api_user_key)

	return pasteRequest(post_url, pasteOptions)
}

func (s *Session) PasteSimple(paste string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)

	return pasteRequest(post_url, pasteOptions)
}

func pasteRequest(req_url string, options url.Values) (pasteURL *url.URL, err error) {
	options.Set("api_dev_key", api_dev_key)

	resp, err := http.PostForm(req_url, options)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(body), "Bad API request") {
		return nil, errors.New(string(body))
	}

	return url.Parse(string(body))
}

func listRequest(req_url string, options url.Values) (pastes []Paste, err error) {
	options.Set("api_dev_key", api_dev_key)

	resp, err := http.PostForm(req_url, options)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := "<fckn_invalid_xml>" + string(body) + "</fckn_invalid_xml>"

	type internal_paste struct {
		Paste_key          string `xml:"paste_key"`
		Paste_date         int64  `xml:"paste_date"`
		Paste_title        string `xml:"paste_title"`
		Paste_size         int    `xml:"paste_size"`
		Paste_expire_date  int64  `xml:"paste_expire_date"`
		Paste_private      int    `xml:"paste_private"`
		Paste_format_long  string `xml:"paste_format_long"`
		Paste_format_short string `xml:"paste_format_short"`
		Paste_url          string `xml:"paste_url"`
		Paste_hits         int    `xml:"paste_hits"`
	}
	type wrapper_paste struct {
		XMLName xml.Name         `xml:"fckn_invalid_xml"`
		Pastes  []internal_paste `xml:"paste"`
	}
	p := wrapper_paste{}

	err = xml.Unmarshal([]byte(result), &p)
	if err != nil {
		return nil, err
	}

	res_pastes := make([]Paste, len(p.Pastes))
	for i, p := range p.Pastes {
		res_pastes[i].Paste_key = p.Paste_key
		res_pastes[i].Paste_date = time.Unix(p.Paste_date, 0)
		res_pastes[i].Paste_title = p.Paste_title
		res_pastes[i].Paste_size = p.Paste_size
		if p.Paste_expire_date == 0 {
			res_pastes[i].Paste_expire_date = time.Date(9999, time.Month(5), 23, 5, 23, 2, 3, time.UTC) // the law of fives is never wrong ;)
		} else {
			res_pastes[i].Paste_expire_date = time.Unix(p.Paste_expire_date, 0)
		}
		res_pastes[i].Paste_private = p.Paste_private
		res_pastes[i].Paste_format_long = p.Paste_format_long
		res_pastes[i].Paste_format_short = p.Paste_format_short
		res_pastes[i].Paste_url, err = url.Parse(p.Paste_url)
		if err != nil {
			return nil, err
		}
		res_pastes[i].Paste_hits = p.Paste_hits
	}
	return res_pastes, nil
}
