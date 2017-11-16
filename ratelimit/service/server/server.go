// Package server defines a rate limit server. The server has the following
// HTTP methods:
//
//     POST /acquire/:API_KEY/:REGION
//     	 Returns a unique string token that can be used to finalize or cancel
//     	 the quota request. If this method returns HTTP OK, then the token must
//     	 be marked either done or cancelled within one minute, or it is
//     	 considered timed out. The method supports the following form fields:
//
//         method: relative HTTP path to the Riot method. If omitted, then
//         	 the request refers to the application-level quota.
//         uniquifier: token that, if provided, signifies a distinct quota bucket,
//           even if the method is the same.
//         noappquota: if set to T or t, indicates that the request should count
//           towards (possibly uniquified) method-level quota, but not application
//           quota.
//
//     POST /done/:TOKEN
//       Marks the request with the given token as complete, so that all
//       relevant quota can be returned after a delay. This request may
//       optionally include HTTP headers returned by the Riot API. If
//       available, the server will parse the headers and update internally
//       tracked quota availability.
//
//     POST /cancel/:TOKEN
//     	 Marks the request with the given token as cancelled, so that all
//     	 relevant quota can be returned immediately.
package server

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/nu7hatch/gouuid"
	"github.com/yuhanfang/riot/ratelimit"
)

const timeout = time.Minute

// callbacksForToken contains the function callbacks that can be invoked for a
// quota acquisition.
type callbacksForToken struct {
	done   ratelimit.Done
	cancel ratelimit.Cancel
}

type server struct {
	// tokens stores callbacks corresponding to unique tokens.
	tokens     map[string]*callbacksForToken
	tokensLock sync.Mutex

	limiter ratelimit.Limiter
}

func (s *server) HandleAcquire(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	key := vars["key"]
	region := vars["region"]
	method := r.Form.Get("method")
	uniquifier := r.Form.Get("uniquifier")
	noAppQuota := r.Form.Get("noappquota")

	inv := ratelimit.Invocation{
		ApplicationKey: key,
		Region:         strings.ToUpper(region),
		Method:         strings.ToLower(method),
		Uniquifier:     uniquifier,
		NoAppQuota:     noAppQuota == "t" || noAppQuota == "T",
	}

	done, cancel, err := s.limiter.Acquire(r.Context(), inv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Defined later in the same thread.
	var timer *time.Timer

	callbacks := callbacksForToken{
		done: func(res *http.Response) error {
			timer.Stop()
			return done(res)
		},
		cancel: func() error {
			timer.Stop()
			return cancel()
		},
	}

	var k string

	for {
		u, err := uuid.NewV4()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		k = u.String()

		s.tokensLock.Lock()
		if _, ok := s.tokens[k]; ok {
			s.tokensLock.Unlock()
			continue
		}
		s.tokens[k] = &callbacks
		fmt.Fprintf(w, "%s", k)
		s.tokensLock.Unlock()
		break
	}

	// Schedule automatic closing out.
	timer = time.AfterFunc(timeout, func() {
		s.tokensLock.Lock()
		defer s.tokensLock.Unlock()

		if got, ok := s.tokens[k]; ok {
			got.done(nil)
			delete(s.tokens, k)
		}
	})
}

func (s *server) HandleDone(w http.ResponseWriter, r *http.Request) {
	s.tokensLock.Lock()
	defer s.tokensLock.Unlock()
	vars := mux.Vars(r)
	token := vars["token"]
	got, ok := s.tokens[token]
	if !ok {
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	got.done(&http.Response{
		Header: r.Header,
	})
	delete(s.tokens, token)
}

func (s *server) HandleCancel(w http.ResponseWriter, r *http.Request) {
	s.tokensLock.Lock()
	defer s.tokensLock.Unlock()
	vars := mux.Vars(r)
	token := vars["token"]
	got, ok := s.tokens[token]
	if !ok {
		http.Error(w, "bad token", http.StatusBadRequest)
		return
	}
	got.cancel()
	delete(s.tokens, token)
}

// New returns an HTTP handler that implements the rate limit service. The
// return value can be used via code like:
// 		r := New()
//    http.Handle("/", r)
func New() http.Handler {
	s := server{
		tokens:  make(map[string]*callbacksForToken),
		limiter: ratelimit.NewLimiter(),
	}
	r := mux.NewRouter()
	r.HandleFunc("/acquire/{key}/{region}", s.HandleAcquire).Methods("POST")
	r.HandleFunc("/done/{token}", s.HandleDone).Methods("POST")
	r.HandleFunc("/cancel/{token}", s.HandleCancel).Methods("POST")
	return r
}
