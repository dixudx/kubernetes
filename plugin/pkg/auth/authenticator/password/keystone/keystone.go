/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package keystone

import (
	"errors"
	"strings"
	"crypto/tls"
	"net/http"

	"github.com/golang/glog"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"k8s.io/kubernetes/pkg/auth/user"
)

// KeystoneAuthenticator contacts openstack keystone to validate user's credentials passed in the request.
// The keystone endpoint is passed during apiserver startup
type KeystoneAuthenticator struct {
	authURL string
	insecureTLS bool
}

// AuthenticatePassword checks the username, password via keystone call
func (keystoneAuthenticator *KeystoneAuthenticator) AuthenticatePassword(username string, password string) (user.Info, bool, error) {
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: keystoneAuthenticator.authURL,
		Username:         username,
		Password:         password,
	}

	_, err := AuthenticatedClient(opts, keystoneAuthenticator.insecureTLS)
	if err != nil {
		glog.Info("Failed: Starting openstack authenticate client: "+ err.Error())
		return nil, false, errors.New("Failed to authenticate")
	}

	return &user.DefaultInfo{Name: username}, true, nil
}

// AuthenticatedClient logs in to an OpenStack cloud found at the identity endpoint specified by options, acquires a
// token, and returns a Client instance that's ready to operate.
func AuthenticatedClient(options gophercloud.AuthOptions, insecureTLS bool) (*gophercloud.ProviderClient, error) {

	client, err := openstack.NewClient(options.IdentityEndpoint)
	if err != nil {
		return nil, err
	}

	client.HTTPClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureTLS}, }

	err = openstack.Authenticate(client, options)
	if err != nil {
		return nil, err
	}
	return client, nil

}

// NewKeystoneAuthenticator returns a password authenticator that validates credentials using openstack keystone
func NewKeystoneAuthenticator(authURL string, insecureTLS bool) (*KeystoneAuthenticator, error) {
	if !strings.HasPrefix(authURL, "https") {
		return nil, errors.New("Auth URL should be secure and start with https")
	}
	if authURL == "" {
		return nil, errors.New("Auth URL is empty")
	}

	return &KeystoneAuthenticator{authURL, insecureTLS}, nil
}
