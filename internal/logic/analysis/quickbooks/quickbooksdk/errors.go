// Copyright (c) 2020, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooksdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Failure is the outermost struct that holds an error response.
type Failure struct {
	Fault struct {
		Error []struct {
			Message string
			Detail  string
			Code    string `json:"code"`
			Element string `json:"element"`
		}
		Type string `json:"type"`
	}
	Time Date `json:"time"`
}

// Error implements the error interface.
func (f Failure) Error() string {
	text, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("unexpected error while marshalling error: %v", err)
	}

	return string(text)
}

func parseFailure(resp *http.Response) error {
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("When reading response body:" + err.Error())
	}

	var errStruct Failure

	if err = json.Unmarshal(msg, &errStruct); err != nil {
		return errors.New(strconv.Itoa(resp.StatusCode) + " " + string(msg))
	} else if len(string(msg)) > 0 {
		return errors.New(strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}

	return errStruct
}
