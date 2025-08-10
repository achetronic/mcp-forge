package handlers

import (
	"io"
	"net/http"
)

// HandleOauthAuthorizationServer process requests for endpoint: /.well-known/oauth-authorization-server
func (h *HandlersManager) HandleOauthAuthorizationServer(response http.ResponseWriter, request *http.Request) {

	remoteUrl := h.dependencies.AppCtx.Config.OAuthAuthorizationServer.IssuerUri + "/.well-known/openid-configuration"
	remoteResponse, err := http.Get(remoteUrl)
	if err != nil {
		h.dependencies.AppCtx.Logger.Error("error getting content from /.well-known/openid-configuration", "error", err.Error())
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//
	remoteResponseBytes, err := io.ReadAll(remoteResponse.Body)
	if err != nil {
		h.dependencies.AppCtx.Logger.Error("error reading bytes from remote response", "error", err.Error())
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Cache-Control", "max-age=3600")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET")          // FIXME: TOO STRICT
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type") // FIXME: TOO STRICT

	_, err = response.Write(remoteResponseBytes)
	if err != nil {
		h.dependencies.AppCtx.Logger.Error("error sending response to client", "error", err.Error())
		return
	}
}
