package pipelines

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/hkoosha/giraffe"
)

const EkranPath = "/ekran"

//goland:noinspection GrazieInspection
type RequestCompensations struct {
	With     any     `json:"with"                 yaml:"with"`
	OnErrRe  *string `json:"on_err_re,omitempty"  yaml:"on_err_re,omitempty"`
	OnNameRe *string `json:"on_name_re,omitempty" yaml:"on_name_re,omitempty"`
	OnStep   *int    `json:"on_step,omitempty"    yaml:"on_step,omitempty"`
	WithFn   string  `json:"with_fn"              yaml:"with_fn"`
}

//nolint:lll
//goland:noinspection GrazieInspection
type Request struct {
	Init          any                     `json:"init,omitempty"          yaml:"init,omitempty"`
	Compensations *[]RequestCompensations `json:"compensations,omitempty" yaml:"compensations,omitempty"`
	Plan          string                  `json:"plan"                    yaml:"plan"`
}

func Remote(
	url string,
	plan string,
	hClient *http.Client,
) Fn {
	fn := remoteFn{
		hClient: hClient,
		plan:    plan,
		url:     url,
	}

	return fn.Ekran
}

type remoteFn struct {
	hClient *http.Client
	plan    string
	url     string
}

func (m *remoteFn) String() string {
	return reflect.TypeOf(m).Elem().String()
}

func mkPayload(
	plan string,
	dat giraffe.Datum,
) (*bytes.Buffer, error) {
	raw, err := dat.Raw()
	if err != nil {
		return nil, err
	}

	req := Request{
		Init:          raw,
		Plan:          plan,
		Compensations: nil,
	}

	payload := new(bytes.Buffer)
	if pErr := json.NewEncoder(payload).Encode(req); pErr != nil {
		return nil, pErr
	}

	return payload, nil
}

func mkRequest(
	ctx context.Context,
	url string,
	payload io.Reader,
) (*http.Request, error) {
	// TODO proper join.
	path := url + EkranPath

	hReq, err := http.NewRequestWithContext(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}

	// TODO headers && Content-Types.
	hReq.Header.Set("Content-Type", "application/json")

	// TODO headers.
	hReq.Header.Set("Accept", "application/json")

	return hReq, nil
}

func sendRequest(
	hClient *http.Client,
	hReq *http.Request,
) (*http.Response, error) {
	resp, err := hClient.Do(hReq)
	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Body == nil {
		return nil, newRemoteError(
			"empty response",
			nil,
		)
	}

	if resp.StatusCode != http.StatusOK {
		//goland:noinspection GoUnhandledErrorResult
		defer resp.Body.Close()

		body := "?"
		if b, err := io.ReadAll(resp.Body); err == nil {
			body = string(b)
		}

		return nil, newRemoteError(
			fmt.Sprintf(
				"unexpected status code: %d => %s",
				resp.StatusCode,
				body,
			),
			nil,
		)
	}

	return resp, nil
}

func decode(
	r io.Reader,
) (giraffe.Datum, error) {
	var res any
	dec := json.NewDecoder(r)
	dec.UseNumber()
	dec.DisallowUnknownFields()
	if err := dec.Decode(&res); err != nil {
		return giraffe.OfErr(), newRemoteError(
			"failed to decode response",
			err,
		)
	}

	dat, err := giraffe.Make(res)
	if err != nil {
		return giraffe.OfErr(), newRemoteError(
			"failed to decode response",
			err,
		)
	}

	return dat, nil
}

func (m *remoteFn) Ekran(
	ctx context.Context,
	dat giraffe.Datum,
) (giraffe.Datum, error) {
	payload, err := mkPayload(m.plan, dat)
	if err != nil {
		return giraffe.OfErr(), err
	}

	hReq, hErr := mkRequest(ctx, m.url, payload)
	if hErr != nil {
		return giraffe.OfErr(), hErr
	}

	resp, rErr := sendRequest(m.hClient, hReq)
	if rErr != nil {
		return giraffe.OfErr(), rErr
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	decoded, dErr := decode(resp.Body)
	if dErr != nil {
		return giraffe.OfErr(), dErr
	}

	return decoded, nil
}
