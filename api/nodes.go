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

// /1.0/nodes endpoint.
var nodesCmd = rest.Endpoint{
        Path: "nodes",

        Get: rest.EndpointAction{Handler: cmdNodesGet, ProxyTarget: true},
        Post: rest.EndpointAction{Handler: cmdNodesPost, ProxyTarget: true},
}

// /1.0/nodes/<name> endpoint.
var nodeCmd = rest.Endpoint{
	Path: "nodes/{name}",

	Delete: rest.EndpointAction{Handler: cmdNodesDelete, ProxyTarget: true},
}

func cmdNodesGet(s *state.State, r *http.Request) response.Response {
        nodes, err := sunbeam.ListNodes(s)
        if err != nil {
                return response.InternalError(err)
        }

        return response.SyncResponse(true, nodes)
}

func cmdNodesPost(s *state.State, r *http.Request) response.Response {
        var req types.Node

        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
                return response.InternalError(err)
        }

        err = sunbeam.AddNode(s, req.Name, req.Role)
        if err != nil {
                return response.InternalError(err)
        }

        return response.EmptySyncResponse
}

func cmdNodesDelete(s *state.State, r *http.Request) response.Response {
	name, err := url.PathUnescape(mux.Vars(r)["name"])
        if err != nil {
                return response.SmartError(err)
        }
        err = sunbeam.DeleteNode(s, name)
        if err != nil {
                return response.InternalError(err)
        }

        return response.EmptySyncResponse
}
