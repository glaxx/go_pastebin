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
// go_pastebin offers functions to interact with the Pastebin API, for
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
	loginURL = "https://pastebin.com/api/api_login.php"
	postURL  = "https://pastebin.com/api/api_post.php"

	expireNever     = "N"
	expire10Minutes = "10M"
	expire1Hour     = "1H"
	expire1Day      = "1D"
	expire1Week     = "1W"
	expire2Weeks    = "2W"
	expire1Month    = "1M"

	exposurePublic   = "0"
	exposureUnlisted = "1"
	exposurePrivate  = "2"
)

// Paste represents output structure from pastebin api
type Paste struct {
	PasteKey         string
	PasteDate        time.Time
	PasteTitle       string
	PasteSize        int
	PasteExpireDate  time.Time
	PastePrivate     int
	PasteFormatLong  string
	PasteFormatShort string
	PasteURL         *url.URL
	PasteHits        int
}

// Pastebin represents pastebin session
type Pastebin struct {
	apiUserKey string
	apiDevKey  string
}

// NewPastebin creates a new session and hold api keys
func NewPastebin(apiDevKey string) Pastebin {
	return Pastebin{apiDevKey: apiDevKey}
}

// PasteAnonymous posts a text (paste) to Pastebin. If successfull this
// function will return a URL reference to the past and nil as error.
// format, expire and private should be set according to
// pastebin.com/api .
func (pb *Pastebin) PasteAnonymous(paste, title, format, expire, private string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}

	pasteOptions.Set("api_dev_key", pb.apiDevKey)
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", paste)
	pasteOptions.Set("api_paste_name", title)
	pasteOptions.Set("api_paste_format", format)
	pasteOptions.Set("api_paste_expire_date", expire)
	pasteOptions.Set("api_paste_private", private)

	return pasteRequest(postURL, pasteOptions)
}

// PasteAnonymousSimple posts a text (paste) to Pastebin. If successfull this
// function will return a URL reference to the past and nil as error.
// title, format, expire and private are defaulted by Pastebin.
func (pb *Pastebin) PasteAnonymousSimple(apiPasteCode string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}
	pasteOptions.Set("api_dev_key", pb.apiDevKey)
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", apiPasteCode)

	return pasteRequest(postURL, pasteOptions)
}

// ListTrendingPastes request (and returns) an array of trending pastes.
func (pb *Pastebin) ListTrendingPastes() (pastes []Paste, err error) {
	listOptions := url.Values{}
	listOptions.Set("api_option", "trends")
	listOptions.Set("api_dev_key", pb.apiDevKey)

	p, err := listRequest(postURL, listOptions)
	if err != nil {
		return nil, err
	}
	return p, err
}

// GenerateUserSession request (and returns) a Session key object, which
// is necessary for all user based tasks.
func (pb *Pastebin) GenerateUserSession(apiUsername, apiPassword string) (p *Pastebin, err error) {
	userOptions := url.Values{}
	userOptions.Set("api_user_name", apiUsername)
	userOptions.Set("api_user_password", apiPassword)
	userOptions.Set("api_dev_key", pb.apiDevKey)

	resp, err := http.PostForm(loginURL, userOptions)
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

	pb.apiUserKey = string(body)

	return pb, nil
}

// ListPastes request (and returns) an array of pastes, which are referenced
// by the user.
func (pb *Pastebin) ListPastes(resultLimit int) (pastes []Paste, err error) {
	if resultLimit < 0 || resultLimit > 1000 {
		return nil, errors.New("resultLimit is out of range")
	}
	listOptions := url.Values{}
	listOptions.Set("api_dev_key", pb.apiDevKey)
	listOptions.Set("api_user_key", pb.apiUserKey)
	listOptions.Set("api_option", "list")
	listOptions.Set("api_result_limit", strconv.Itoa(resultLimit))
	p, err := listRequest(postURL, listOptions)
	if err != nil {
		return nil, err
	}

	return p, err
}

// DeletePaste deletes a paste, which is referenced by the paste_key.
func (pb *Pastebin) DeletePaste(apiPasteKey string) (err error) {
	qryOptions := url.Values{}
	qryOptions.Set("api_paste_key", apiPasteKey)
	qryOptions.Set("api_user_key", pb.apiUserKey)
	qryOptions.Set("api_dev_key", pb.apiDevKey)
	qryOptions.Set("api_option", "delete")

	resp, err := http.PostForm(postURL, qryOptions)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) != "Paste Removed" {
		return errors.New(string(body))
	}
	return nil
}

// Paste posts a text (paste) to Pastebin. If successfull this
// function will return a URL reference to the past and nil as error.
// format, expire and private should be set according to
// pastebin.com/api .
func (pb *Pastebin) Paste(apiPasteCode, apiPasteName, apiPasteFormat, apiPasteExpire, apiPastePrivate string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}

	pasteOptions.Set("api_dev_key", pb.apiDevKey)
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", apiPasteCode)
	pasteOptions.Set("api_paste_name", apiPasteName)
	pasteOptions.Set("api_paste_format", apiPasteFormat)
	pasteOptions.Set("api_paste_expire_date", apiPasteExpire)
	pasteOptions.Set("api_paste_private", apiPastePrivate)
	pasteOptions.Set("api_user_key", pb.apiUserKey)

	return pasteRequest(postURL, pasteOptions)
}

// PasteSimple posts a text (paste) to Pastebin. If successfull this
// function will return a URL reference to the past and nil as error.
// title, format, expire and private are defaulted by Pastebin.
func (pb *Pastebin) PasteSimple(apiPasteCode string) (pasteURL *url.URL, err error) {
	pasteOptions := url.Values{}
	pasteOptions.Set("api_dev_key", pb.apiDevKey)
	pasteOptions.Set("api_option", "paste")
	pasteOptions.Set("api_paste_code", apiPasteCode)
	pasteOptions.Set("api_user_key", pb.apiUserKey)

	return pasteRequest(postURL, pasteOptions)
}

func pasteRequest(requestURL string, options url.Values) (pasteURL *url.URL, err error) {
	resp, err := http.PostForm(requestURL, options)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Inline function to process the rare case of "Post Limit" error
	// Looks like: Post%20limit,%20maximum%20pastes%20per%2024h%20reached
	urlFilter := func(url string) string {
		return strings.Join(strings.Split(url, "%20"), " ")
	}

	if strings.Contains(string(body), "Bad API request") {
		return nil, errors.New(string(body))
	} else if strings.Contains(urlFilter(string(body)), "Post Limit") {
		return nil, errors.New(urlFilter(string(body)))
	}

	return url.Parse(string(body))
}

func listRequest(requestURL string, options url.Values) (pastes []Paste, err error) {
	resp, err := http.PostForm(requestURL, options)
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

	result := "<fckn_invalid_xml>" + string(body) + "</fckn_invalid_xml>"

	type internalPaste struct {
		PasteKey         string `xml:"paste_key"`
		PasteDate        int64  `xml:"paste_date"`
		PasteTitle       string `xml:"paste_title"`
		PasteSize        int    `xml:"paste_size"`
		PasteExpireDate  int64  `xml:"paste_expire_date"`
		PastePrivate     int    `xml:"paste_private"`
		PasteFormatLong  string `xml:"paste_format_long"`
		PasteFormatShort string `xml:"paste_format_short"`
		PasteURL         string `xml:"paste_url"`
		PasteHits        int    `xml:"paste_hits"`
	}
	type wrapperPaste struct {
		XMLName xml.Name        `xml:"fckn_invalid_xml"`
		Pastes  []internalPaste `xml:"paste"`
	}
	p := wrapperPaste{}

	err = xml.Unmarshal([]byte(result), &p)
	if err != nil {
		return nil, err
	}

	resPastes := make([]Paste, len(p.Pastes))
	for i, p := range p.Pastes {
		resPastes[i].PasteKey = p.PasteKey
		resPastes[i].PasteDate = time.Unix(p.PasteDate, 0)
		resPastes[i].PasteTitle = p.PasteTitle
		resPastes[i].PasteSize = p.PasteSize
		if p.PasteExpireDate == 0 {
			resPastes[i].PasteExpireDate = time.Date(9999, time.Month(5), 23, 5, 23, 2, 3, time.UTC) // the law of fives is never wrong ;)
		} else {
			resPastes[i].PasteExpireDate = time.Unix(p.PasteExpireDate, 0)
		}
		resPastes[i].PastePrivate = p.PastePrivate
		resPastes[i].PasteFormatLong = p.PasteFormatLong
		resPastes[i].PasteFormatShort = p.PasteFormatShort
		resPastes[i].PasteURL, err = url.Parse(p.PasteURL)
		if err != nil {
			return nil, err
		}
		resPastes[i].PasteHits = p.PasteHits
	}

	return resPastes, nil
}
