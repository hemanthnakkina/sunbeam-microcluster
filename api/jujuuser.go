package api

import (
        "encoding/json"
        "net/http"
	"net/url"

	"github.com/gorilla/mux"
        "github.com/canonical/microcluster/rest"
        "github.com/canonical/microcluster/state"
        "github.com/lxc/lxd/lxd/response"

        "github.com/openstack-snaps/sunbeam-microcluster/api/types"
        "github.com/openstack-snaps/sunbeam-microcluster/sunbeam"
)

// /1.0/jujuusers endpoint.
var jujuusersCmd = rest.Endpoint{
        Path: "jujuusers",

        Get: rest.EndpointAction{Handler: cmdJujuUsersGet, ProxyTarget: true},
        Post: rest.EndpointAction{Handler: cmdJujuUsersPost, ProxyTarget: true},
}

// /1.0/jujuusers/<name> endpoint.
var jujuuserCmd = rest.Endpoint{
	Path: "jujuusers/{name}",

	Delete: rest.EndpointAction{Handler: cmdJujuUsersDelete, ProxyTarget: true},
}

func cmdJujuUsersGet(s *state.State, r *http.Request) response.Response {
        users, err := sunbeam.ListJujuUsers(s)
        if err != nil {
                return response.InternalError(err)
        }

        return response.SyncResponse(true, users)
}

func cmdJujuUsersPost(s *state.State, r *http.Request) response.Response {
        var req types.JujuUser

        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
                return response.InternalError(err)
        }

        err = sunbeam.AddJujuUser(s, req.Username, req.Token)
        if err != nil {
                return response.InternalError(err)
        }

        return response.EmptySyncResponse
}

func cmdJujuUsersDelete(s *state.State, r *http.Request) response.Response {
	name, err := url.PathUnescape(mux.Vars(r)["username"])
        if err != nil {
                return response.SmartError(err)
        }
        err = sunbeam.DeleteJujuUser(s, name)
        if err != nil {
                return response.InternalError(err)
        }

        return response.EmptySyncResponse
}
