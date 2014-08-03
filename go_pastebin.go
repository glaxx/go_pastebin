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
	"net/http"
	"net/url"
	"io/ioutil"
)

const (
	api_dev_key = ""
	login_url = "http://pastebin.com/api/api_login.php"
	paste_url = "http://pastebin.com/api/api_post.php"

	expire_never = "N"
	expire_10Minutes = "10M"
	expire_1Hour = "1H"
	expire_1Day = "1D"
	expire_1Week = "1W"
	expire_2Weeks = "2W"
	expire_1Month = "1M"

	private_public = "0"
	private_unlisted = "1"
	private_private = "2"
)


type session struct {
	api_user_key string
}

func PasteAnonymous(paste, title, format, expire, private string) (pasteURL string, err error){
	pasteOptions := url.Values{}

	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)
	pasteOptions.Set("api_paste_name", title)
	pasteOptions.Set("api_paste_format", format)
	pasteOptions.Set("api_paste_expire_date", expire)
	pasteOptions.Set("api_paste_private", private)

	return request(paste_url, pasteOptions)
}

func PasteAnonymousSimple(paste string) (pasteURL string, err error) {
	pasteOptions := url.Values{}
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)

	return request(paste_url, pasteOptions)
}

func GenerateUserSession(username, password string) (se *session, err error) {
	var s session
	userOptions := url.Values{}
	userOptions.Set("api_user_name", username)
	userOptions.Set("api_user_password", password)

	s.api_user_key, err = request(login_url, userOptions)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *session) Paste(paste, title, format, expire, private string) (pasteURL string, err error) {
	pasteOptions := url.Values{}

	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)
	pasteOptions.Set("api_paste_name", title)
	pasteOptions.Set("api_paste_format", format)
	pasteOptions.Set("api_paste_expire_date", expire)
	pasteOptions.Set("api_paste_private", private)

	pasteOptions.Set("api_user_key", s.api_user_key)

	return request(paste_url, pasteOptions)
}

func (s *session) PasteSimple(paste string) (pasteURL string, err error) {
	pasteOptions := url.Values{}
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)

	return request(paste_url, pasteOptions)
}

func request(req_url string, options url.Values) (result string, err error) {
	options.Set("api_dev_key", api_dev_key)

	resp, err := http.PostForm(req_url, options)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}


