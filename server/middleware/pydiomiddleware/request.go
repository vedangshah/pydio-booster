// Package pydiomiddleware contains the logic for a middleware directive (repetitive task done for a Pydio request)
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
package pydiomiddleware

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
	pydhttp "github.com/pydio/pydio-booster/http"
	"github.com/pydio/pydio-booster/log"
	pydioworker "github.com/pydio/pydio-booster/worker"
)

var client = pydhttp.NewClient()

// RequestJob definition for the uploader
type RequestJob struct {
	Request    http.Request
	HandleFunc func(string, io.Reader, http.Header) error
	ErrorFunc  func()
}

// Do the job
func (j *RequestJob) Do() (err error) {

	resp, err := client.Do(&j.Request)
	if err != nil {
		j.ErrorFunc()
		return
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil && err != io.EOF {
		j.ErrorFunc()
		return
	}

	if resp.StatusCode != http.StatusOK || err != nil {
		logger.Errorf("Not authorized : %v %v", j.Request, resp)
		j.ErrorFunc()
		return err
	}

	// For each header, looping through the handling func
	for header, values := range resp.Header {
		for _, value := range values {
			if err := j.HandleFunc(header, strings.NewReader(value), nil); err != nil {
				return err
			}
		}
	}

	// And finishing with the body
	return j.HandleFunc("body", resp.Body, resp.Header)
}

// NewRequestJob prepares the job for the middleware request
// based on the rules
func NewRequestJob(
	ctx context.Context,
	u url.URL,
	headers [][2]string,
	cookies []*http.Cookie,
	out Out,
	replacer httpserver.Replacer,
	encoder Encoder,
	writer http.ResponseWriter,
	close func() error,
	cancel func(),
) (pydioworker.Job, error) {

	queryArgs := u.Query()
	logger.Infoln("Request Job Start", u, headers, queryArgs, out)

	node, err := getContextNode(ctx)
	if err == nil {
		logger.Debugln("Request : Retrieved node ", node.Repo)
		var repo *url.URL
		var dir *url.URL

		repo = &url.URL{Opaque: node.Repo.String()}
		dir = &url.URL{Opaque: strings.TrimPrefix(node.Dir.String(), "/")}

		replacer.Set("repo", repo.String())
		replacer.Set("nodedir", dir.String())
		replacer.Set("nodename", node.Basename)

		logger.Debugln("Request : Replacer is set ", replacer)
	} else {
		logger.Errorln("Request : Could not read node")
	}

	logger.Debugln("We have a node ", node)

	// Replacing any potential placeholder
	u.Path = replacer.Replace(u.Path)
	u.RawQuery = replacer.Replace(u.RawQuery)

	values := url.Values{}
	for arg, vals := range u.Query() {
		for _, val := range vals {
			log.Debugln(replacer.Replace(val))
			values.Add(arg, replacer.Replace(val))
		}
	}

	if len(cookies) == 0 {
		// If we don't read from a cookie, then we have the auth details already set
		if auth, errAuth := getContextAuthParams(ctx, u.Path); errAuth == nil {
			values.Set("auth_hash", auth.Hash)
			values.Set("auth_token", auth.Token)
			if auth.Key != "" {
				values.Set("key", auth.Key)
			}
		}
	}
	u.RawQuery = values.Encode()

	request, _ := http.NewRequest("GET", u.Path, nil)
	logger.Debugln("Doing headers")
	for _, header := range headers {
		request.Header.Add(header[0], replacer.Replace(header[1]))
	}

	logger.Debugln("Doing cookies")
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	request.URL = &u

	logger.Debugf("URL is %s - headers - %v cookies - %v", u.String(), headers, request.Cookies())

	job := &RequestJob{
		Request:   *request,
		ErrorFunc: cancel,
		HandleFunc: func(key string, body io.Reader, header http.Header) error {

			if key == "body" {
				// Always finishing by the body
				defer close()
			}

			if key != out.Key {
				// Not interested
				return nil
			}

			logger.Debugln("Out to ", out)

			switch out.Name {
			case "body":

				defer close()

				// Write the response body
				logger.Debugln("Starting to write")

				// Write response header
				copyHeader(writer.Header(), header)

				_, err = io.Copy(writer, body)
				if err != nil {
					logger.Errorf("Could not write body %v", err)
					return err
				}
				logger.Debugln("Ended body write")

				return nil
			}

			data, _ := ioutil.ReadAll(body)
			logger.Debugln("We have data ", string(data))
			encoder.Encode(string(data))

			return nil
		},
	}

	return pydioworker.Job(job), nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		if _, ok := dst[k]; ok {
			// otherwise, overwrite
			dst.Del(k)
		}
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

type decoder interface {
	Decode(v interface{}) error
}
