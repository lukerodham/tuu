package tuu

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func newContext(r Route, res http.ResponseWriter, req *http.Request) *DefaultContext {
	data := make(map[string]interface{})
	data["path"] = r.Path
	data["env"] = r.Env

	params := req.URL.Query()
	vars := mux.Vars(req)
	for k, v := range vars {
		params.Set(k, v)
	}

	sessionStore, _ := r.Session.Get(req, r.SessionName)
	session := &Session{
		Session: sessionStore,
		req:     req,
		res:     res,
	}

	return &DefaultContext{
		response: res,
		request:  req,
		params:   params,
		data:     data,
		env:      r.Env,
		logger:   r.Logger,
		session:  session,
		flash:    newFlash(session),
	}
}

type DefaultContext struct {
	context.Context
	response    http.ResponseWriter
	request     *http.Request
	params      url.Values
	contentType string
	data        map[string]interface{}
	env         string
	logger      *logrus.Logger
	session     *Session
	flash       *Flash
}

// Response returns the original Response for the request.
func (d *DefaultContext) Response() http.ResponseWriter {
	return d.response
}

// Request returns the original Request.
func (d *DefaultContext) Request() *http.Request {
	return d.request
}

// Params returns all of the parameters for the request,
// including both named params and query string parameters.
func (d *DefaultContext) Params() url.Values {
	return d.params
}

// Param returns a param, either named or query string,
// based on the key.
func (d *DefaultContext) Param(key string) string {
	return d.Params().Get(key)
}

// Set a value onto the Context. Any value set onto the Context
// will be automatically available in templates.
func (d *DefaultContext) Set(key string, value interface{}) {
	d.data[key] = value
}

// Value that has previously stored on the context.
func (d *DefaultContext) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		if v, ok := d.data[k]; ok {
			return v
		}
	}
	return d.Context.Value(key)
}

func (d *DefaultContext) Render(status int, rr render.Renderer) error {
	if rr != nil {
		data := d.data
		pp := map[string]string{}
		for k, v := range d.params {
			pp[k] = v[0]
		}

		data["params"] = pp
		data["request"] = d.Request()
		data["session"] = d.Session()
		data["flash"] = d.Flash().data

		bb := &bytes.Buffer{}

		err := rr.Render(bb, data)
		if err != nil {
			return err
		}

		if d.Session() != nil {
			d.Flash().Clear()
			d.Flash().persist(d.Session())
		}

		d.Response().Header().Set("Content-Type", rr.ContentType())
		d.Response().WriteHeader(status)
		_, err = io.Copy(d.Response(), bb)
		if err != nil {
			return err
		}

		return nil
	}

	d.Response().WriteHeader(status)
	return nil
}

func (d *DefaultContext) Redirect(status int, url string) error {
	d.Flash().persist(d.Session())

	http.Redirect(d.Response(), d.Request(), url, status)
	return nil
}

func (d *DefaultContext) Env() string {
	return d.env
}

func (d *DefaultContext) Logger() *logrus.Logger {
	return d.logger
}

func (d *DefaultContext) Session() *Session {
	return d.session
}

func (d *DefaultContext) Flash() *Flash {
	return d.flash
}
