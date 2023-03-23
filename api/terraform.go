package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/canonical/microcluster/rest"
	"github.com/canonical/microcluster/state"
	"github.com/lxc/lxd/lxd/response"
	"github.com/lxc/lxd/lxd/util"

	"github.com/openstack-snaps/sunbeam-microcluster/api/types"
	"github.com/openstack-snaps/sunbeam-microcluster/sunbeam"
)

// /1.0/terraformstate endpoint.
// The endpoints are basically to provide REST URLs to Terraform http
// backend configuration to maintain Terraform state centrally with
// locking mechanism.
// Terraform 1.3.x doesnot support passing certs to the REST URL for
// authentications and so the endpoints are exposed as AllowUntrusted.
// TODO: Newer version yet to release 1.4.x supports TLS authentication
// to http backend. Once sunbeam moves to use 1.4.x, change the
// endpoints not to allow untrusted.
// https://github.com/hashicorp/terraform/commit/75e5ae27a258122fe6bf122beb943324c69de5b1
var terraformStateCmd = rest.Endpoint{
	Path: "terraformstate",

	Get:  rest.EndpointAction{Handler: cmdStateGet, AllowUntrusted: true},
	Post: rest.EndpointAction{Handler: cmdStatePost, AllowUntrusted: true},
}

// /1.0/terraformlock endpoint.
var terraformLockCmd = rest.Endpoint{
	Path: "terraformlock",

	Get:  rest.EndpointAction{Handler: cmdLockGet, AllowUntrusted: true},
	Post: rest.EndpointAction{Handler: cmdLockPost, AllowUntrusted: true},
}

// /1.0/terraformunlock endpoint.
var terraformUnlockCmd = rest.Endpoint{
	Path: "terraformunlock",

	Post: rest.EndpointAction{Handler: cmdUnlockPost, AllowUntrusted: true},
}

func cmdStateGet(s *state.State, _ *http.Request) response.Response {
	state, err := sunbeam.GetTerraformState(s)
	if err != nil {
		return response.InternalError(err)
	}

	// Just send state data instead of SyncResponse Json object as
	// terraform expects just state data.
	return response.ManualResponse(func(w http.ResponseWriter) error {
		return util.WriteJSON(w, state, nil)
	})
}

func cmdStatePost(s *state.State, r *http.Request) response.Response {
	b, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
	defer r.Body.Close()
	if err != nil {
		return response.InternalError(err)
	}

	err = sunbeam.UpdateTerraformState(s, string(b))
	if err != nil {
		return response.InternalError(err)
	}

	return response.EmptySyncResponse
}

func cmdLockGet(s *state.State, _ *http.Request) response.Response {
	lock, err := sunbeam.GetTerraformLock(s)
	if err != nil {
		return response.InternalError(err)
	}

	// Just send state data instead of SyncResponse Json object as
	// terraform expects just state data.
	return response.ManualResponse(func(w http.ResponseWriter) error {
		return util.WriteJSON(w, lock, nil)
	})
}

func cmdLockPost(s *state.State, r *http.Request) response.Response {
	var reqLock types.Lock
	var dbLock types.Lock

	err := json.NewDecoder(r.Body).Decode(&reqLock)
	if err != nil {
		return response.InternalError(err)
	}

	lockInDb, err := sunbeam.GetTerraformLock(s)
	if err != nil {
		return response.InternalError(err)
	}

	err = json.Unmarshal([]byte(lockInDb), &dbLock)
	if err != nil {
		return response.InternalError(err)
	}

	if dbLock.ID == "" {
		j, err := json.Marshal(reqLock)
		if err != nil {
			return response.InternalError(err)
		}

		err = sunbeam.UpdateTerraformLock(s, string(j))
		if err != nil {
			return response.InternalError(err)
		}

		return response.EmptySyncResponse
	}

	// If the lock from DB and request are same, send http 423
	if dbLock.ID == reqLock.ID && dbLock.Operation == reqLock.Operation && dbLock.Who == reqLock.Who {
		return response.ManualResponse(func(w http.ResponseWriter) error {
			w.WriteHeader(http.StatusLocked)
			return util.WriteJSON(w, nil, nil)
		})
	}

	// Already locked and request has different lockid, send http 409
	return response.ManualResponse(func(w http.ResponseWriter) error {
		w.WriteHeader(http.StatusConflict)
		return util.WriteJSON(w, nil, nil)
	})
}

func cmdUnlockPost(s *state.State, r *http.Request) response.Response {
	var reqLock types.Lock
	var dbLock types.Lock

	err := json.NewDecoder(r.Body).Decode(&reqLock)
	if err != nil {
		return response.InternalError(err)
	}

	lockInDb, err := sunbeam.GetTerraformLock(s)
	if err != nil {
		return response.InternalError(err)
	}

	err = json.Unmarshal([]byte(lockInDb), &dbLock)
	if err != nil {
		return response.InternalError(err)
	}

	// If the lock from DB and request are same, clear the lock from DB
	if dbLock.ID == reqLock.ID && dbLock.Operation == reqLock.Operation && dbLock.Who == reqLock.Who {
		err = sunbeam.UpdateTerraformLock(s, "{}")
		if err != nil {
			return response.InternalError(err)
		}

		return response.EmptySyncResponse
	}

	// Request has different lock id than in database, send http 409
	return response.ManualResponse(func(w http.ResponseWriter) error {
		w.WriteHeader(http.StatusConflict)
		return util.WriteJSON(w, nil, nil)
	})
}
