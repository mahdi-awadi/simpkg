package forms

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-per/simpkg/cache"
	"github.com/go-per/simpkg/format"
	"github.com/go-per/simpkg/helpers"
	"github.com/go-per/simpkg/i18n"
	"github.com/go-per/simpkg/parse"
	"github.com/go-per/simpkg/std"
	"github.com/imroc/req/v3"
)

// FormExecutor struct
type FormExecutor struct {
	rf               *Forms
	form             *Form
	formName         string
	urlParams        map[string]string
	headerParams     map[string]string
	bodyParams       map[string]any
	onBeforePrepare  func(*FormExecutor) error
	onBeforeSend     func(*FormExecutor) error
	onAfterSent      func(*FormExecutor)
	request          *req.Request
	resp             *req.Response
	cache            cache.ICache
	cachePolicy      []func(*FormExecutor) error
	cacheFileName    string
	cacheAsync       bool
	cacheAppendTs    bool
	cacheSkipIfError bool
	restoreIfExists  bool
	target           any
	successStatuses  []int
	checkStatusCode  bool
	err              error
	prepared         bool
	restored         bool
}

// Executor return new executor for form
func (rf *Forms) Executor(formName string) *FormExecutor {
	return &FormExecutor{
		rf:              rf,
		formName:        formName,
		successStatuses: []int{http.StatusOK},
		checkStatusCode: true,
		cachePolicy:     make([]func(*FormExecutor) error, 0),
	}
}

// Form return form
func (e *FormExecutor) Form() *Form {
	return e.form
}

// Request set request
func (e *FormExecutor) Request(request *req.Request) *FormExecutor {
	e.request = request
	return e
}

// GetRequest returns current request instance
func (e *FormExecutor) GetRequest() *req.Request {
	return e.request
}

// Response set target response
func (e *FormExecutor) Response(target any) *FormExecutor {
	e.target = target
	return e
}

// CacheInterface set cache interface
func (e *FormExecutor) CacheInterface(c cache.ICache, skipOnError ...bool) *FormExecutor {
	e.cache = c
	e.cacheAsync = true
	e.cacheSkipIfError = true
	if skipOnError != nil && len(skipOnError) > 0 && !skipOnError[0] {
		e.cacheSkipIfError = false
	}

	return e
}

// CacheAsync set cache async value
func (e *FormExecutor) CacheAsync(v ...bool) {
	e.cacheAsync = true
	if v != nil && len(v) > 0 && !v[0] {
		e.cacheAsync = false
	}
}

// CacheFileName set cache file name
func (e *FormExecutor) CacheFileName(name string, withTimestamp ...bool) *FormExecutor {
	e.cacheFileName = name
	if withTimestamp != nil && len(withTimestamp) > 0 && withTimestamp[0] {
		e.cacheAppendTs = true
	}

	return e
}

// CachePolicy set cache policy
func (e *FormExecutor) CachePolicy(policy func(*FormExecutor) error) *FormExecutor {
	e.cachePolicy = append(e.cachePolicy, policy)
	return e
}

// RestoreIfExists restore from cache if exists
func (e *FormExecutor) RestoreIfExists(v ...bool) *FormExecutor {
	if v != nil && len(v) > 0 {
		e.restoreIfExists = v[0]
	} else {
		e.restoreIfExists = true
	}

	return e
}

// IsRestored returns restore status
func (e *FormExecutor) IsRestored() bool {
	return e.restored
}

// OnBeforePrepare set on before prepare callback
func (e *FormExecutor) OnBeforePrepare(fn func(*FormExecutor) error) *FormExecutor {
	e.onBeforePrepare = fn
	return e
}

// OnBeforeSend set on before send callback
func (e *FormExecutor) OnBeforeSend(fn func(*FormExecutor) error) *FormExecutor {
	e.onBeforeSend = fn
	return e
}

// OnAfterSent set on after send callback
func (e *FormExecutor) OnAfterSent(fn func(*FormExecutor)) *FormExecutor {
	e.onAfterSent = fn
	return e
}

// BodyParams set body params which replace keys with values in payload
// keys are dynamic keys in payload, format: {KEY}
func (e *FormExecutor) BodyParams(params map[string]any) *FormExecutor {
	e.bodyParams = params
	return e
}

// HeaderParams set params which replaces headers dynamic values
func (e *FormExecutor) HeaderParams(params map[string]string) *FormExecutor {
	e.headerParams = params
	return e
}

// UrlParams set url params which replace keys with values in url
func (e *FormExecutor) UrlParams(params map[string]string) *FormExecutor {
	e.urlParams = params
	return e
}

// ErrorIfStatusNoIn error if response status code not in given statuses
// default is http.StatusOK
func (e *FormExecutor) ErrorIfStatusNoIn(v ...int) *FormExecutor {
	e.successStatuses = v
	return e
}

// CheckStatusCode set check status code value
// check status by default is true
func (e *FormExecutor) CheckStatusCode(v ...bool) *FormExecutor {
	check := true
	if len(v) > 0 && !v[0] {
		check = false
	}

	e.checkStatusCode = check
	return e
}

// GetRawResponse returns response
func (e *FormExecutor) GetRawResponse() *req.Response {
	return e.resp
}

// GetResponse returns result
func (e *FormExecutor) GetResponse() any {
	return e.target
}

// Error returns error
func (e *FormExecutor) Error() error {
	return e.err
}

// Prepare form
func (e *FormExecutor) Prepare() *FormExecutor {
	if e.restoreIfExists && e.restore() {
		e.restored = true
	}

	// find form and process
	formItem, ok := e.rf.Get(e.formName)
	if !ok {
		e.err = errors.New(i18n.Translate("form_not_exists"))
		return e
	}
	if formItem.Endpoint == "" {
		e.err = errors.New(i18n.Translate("invalid_endpoint"))
		return e
	}
	if e.request == nil {
		e.err = errors.New(i18n.Translate("invalid_request"))
		return e
	}

	// clone new form
	e.form = &Form{}
	e.form.IsFormData = formItem.IsFormData
	e.form.Endpoint = formItem.Endpoint
	e.form.Method = formItem.Method
	e.form.Timeout = formItem.Timeout
	e.form.WithoutBody = formItem.WithoutBody
	e.form.bodyString = formItem.bodyString
	e.form.Headers = make(map[string]any, 0)
	e.form.Body = make(map[string]any, 0)
	e.form.Data = make(map[string]any, 0)

	if formItem.Headers != nil && len(formItem.Headers) > 0 {
		for key, value := range formItem.Headers {
			e.form.Headers[key] = value
		}
	}
	if formItem.Body != nil && len(formItem.Body) > 0 {
		for key, value := range formItem.Body {
			e.form.Body[key] = value
		}
	}
	if formItem.Data != nil && len(formItem.Data) > 0 {
		for key, value := range formItem.Data {
			e.form.Data[key] = value
		}
	}

	// discard process if restored
	if e.restored {
		return e
	}

	// on before prepare
	if e.onBeforePrepare != nil {
		err := e.onBeforePrepare(e)
		if err != nil {
			e.err = err
			return e
		}
	}

	// json replacer
	e.form.bodyString = format.Replace(e.form.bodyString, e.bodyParams)

	// set form data/payload
	if !e.form.WithoutBody && helpers.Includes([]string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}, e.form.Method) {
		if !e.form.IsFormData && e.form.bodyString != "" {
			e.request.SetBodyJsonString(e.form.bodyString)
		}

		if e.form.IsFormData && e.form.bodyString != "" {
			var payload map[string]any
			e.err = parse.ToStruct([]byte(e.form.bodyString), &payload)
			if e.err == nil {
				body := make(map[string]string)
				if payload != nil {
					for k, v := range payload {
						body[k] = format.Format("%v", v)
					}
				}
				e.request.SetFormData(body)
			}
		}
	}

	// update request headers
	if e.headerParams != nil && len(e.headerParams) > 0 {
		if e.form.Headers == nil {
			e.form.Headers = make(map[string]any)
		}

		for paramKey, paramValue := range e.headerParams {
			for key, value := range e.form.Headers {
				e.form.Headers[key] = strings.ReplaceAll(format.String(value), format.Format("{%v}", paramKey), paramValue)
			}
		}
	}

	// set request headers
	for key, value := range e.form.Headers {
		e.request.SetHeader(key, format.String(value))
	}

	// replace url params
	if e.urlParams != nil && len(e.urlParams) > 0 {
		for key, value := range e.urlParams {
			e.form.Endpoint = strings.ReplaceAll(e.form.Endpoint, format.Format("{%s}", key), value)
		}
	}

	e.prepared = true
	return e
}

// Do call http request
func (e *FormExecutor) Do() *FormExecutor {
	// prepare
	if !e.prepared {
		e.Prepare()
	}

	// if restored, discard process
	if e.restored {
		return e
	}

	// check for errors
	if e.err != nil {
		return e
	}

	// on before send callback
	if e.onBeforeSend != nil {
		err := e.onBeforeSend(e)
		if err != nil {
			e.err = err
		}
		if e.err != nil {
			return e
		}
	}

	// execute and return response
	resp, err := e.request.Send(e.form.Method, e.form.Endpoint)
	if err != nil {
		e.err = err
	}

	// set response
	e.resp = resp

	// on after send
	if e.onAfterSent != nil {
		e.onAfterSent(e)
	}

	// if form has error
	if e.err != nil {
		return e
	}

	// check response
	if resp == nil {
		e.err = errors.New(i18n.Translate("response_is_empty"))
		return e
	}

	// check response status
	if e.checkStatusCode && !helpers.Includes(e.successStatuses, resp.GetStatusCode()) {
		e.err = errors.New(i18n.Translate("invalid_status", resp.GetStatusCode()))
		return e
	}

	// parse response
	if e.target != nil {
		body, err := resp.ToBytes()
		if err == nil {
			err = parse.ToStruct(body, e.target)
		}
		if err != nil {
			e.err = err
		}
		if e.target == nil {
			e.err = errors.New(i18n.Translate("response_is_empty"))
		}
	}

	// write cache
	e.writeCacheResponse()

	return e
}

// writeCacheResponse cache response
func (e *FormExecutor) writeCacheResponse() {
	if e.cache == nil || e.GetRawResponse() == nil {
		return
	}
	if e.cacheSkipIfError && e.Error() != nil {
		return
	}

	// check for cache policies
	if e.cachePolicy != nil {
		for _, policy := range e.cachePolicy {
			if err := policy(e); err != nil {
				return
			}
		}
	}

	// check for cache file name
	cacheFileName := e.cacheFileName
	if cacheFileName == "" {
		cacheFileName = e.formName
		e.cacheAppendTs = !e.restoreIfExists
	}

	body, err := e.GetRawResponse().ToBytes()
	if err == nil {
		if e.cacheAsync {
			e.cache.WriteAsync(cacheFileName, body, e.cacheAppendTs)
		} else {
			err = e.cache.Write(cacheFileName, body, e.cacheAppendTs)
			if err != nil {
				std.Error("CacheInterface error file name:%s error: %v", cacheFileName, err)
			}
		}
	}
}

// restore cached version
func (e *FormExecutor) restore() bool {
	if e.cache == nil || e.target == nil {
		return false
	}

	content, err := e.cache.Get(e.cacheFileName)
	if err != nil {
		return false
	}

	err = parse.ToStruct(content, e.target)
	if err != nil {
		return false
	}
	if e.target == nil {
		return false
	}

	return true
}
