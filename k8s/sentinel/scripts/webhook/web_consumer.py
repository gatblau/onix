#
# Sentinel - Copyright (c) 2019 by www.gatblau.org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software distributed under
# the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
# either express or implied.
# See the License for the specific language governing permissions and limitations under the License.
#
# Contributors to this project, hereby assign copyright in this code to the project,
# to be licensed under the same terms as the rest of the code.
#

try:
    # try the import for python 3
    from http.server import HTTPServer, BaseHTTPRequestHandler
except ImportError:
    # if failed then go for the python 2 compatible import
    from BaseHTTPServer import HTTPServer, BaseHTTPRequestHandler

# This is an example of a very simple web scripts for the Sentinel Webhook publisher
# This scripts outputs the payload sent by Sentinel to the std out
# Use for testing purposes only
class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):

    # if you try it on your browser then it tells you what it is about
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'Web scripts is waiting for Sentinel post requests!')

    # process post requests writing the body to the std output
    def do_POST(self):
        # retrieves the body in the request
        content_length = int(self.headers['Content-Length'])
        body = self.rfile.read(content_length)

        # sends an OK response
        self.send_response(200)
        self.end_headers()

        # prints the body in the request to stdout
        print body

# sets the port the web server listens to
port = 8000

# creates an instance of the web server
httpd = HTTPServer(('localhost', port), SimpleHTTPRequestHandler)

# run the server
print 'Sentinel web comsumer listening on port: ' + str(port)
httpd.serve_forever()