// Copyright (c) 2015-2021 MinIO, Inc.; 2026 MaxIO
// SPDX-License-Identifier: AGPL-3.0-or-later
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.

package cmd

import "net/http"

// Standard cross domain policy information located at https://s3.amazonaws.com/crossdomain.xml
const crossDomainXML = `<?xml version="1.0"?><!DOCTYPE cross-domain-policy SYSTEM "http://www.adobe.com/xml/dtds/cross-domain-policy.dtd"><cross-domain-policy><allow-access-from domain="*" secure="false" /></cross-domain-policy>`

// Standard path where an app would find cross domain policy information.
const crossDomainXMLEntity = "/crossdomain.xml"

// A cross-domain policy file is an XML document that grants a web client, such as Adobe Flash Player
// or Adobe Acrobat (though not necessarily limited to these), permission to handle data across domains.
// When clients request content hosted on a particular source domain and that content make requests
// directed towards a domain other than its own, the remote domain needs to host a cross-domain
// policy file that grants access to the source domain, allowing the client to continue the transaction.
func setCrossDomainPolicyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cxml := crossDomainXML
		if globalServerCtxt.CrossDomainXML != "" {
			cxml = globalServerCtxt.CrossDomainXML
		}
		// Look for 'crossdomain.xml' in the incoming request.
		if r.URL.Path == crossDomainXMLEntity {
			// Write the standard cross domain policy xml.
			w.Write([]byte(cxml))
			// Request completed, no need to serve to other handlers.
			return
		}
		h.ServeHTTP(w, r)
	})
}
