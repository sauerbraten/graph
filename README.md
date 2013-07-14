# Graph <a href="http://goci.me/project/github.com/sauerbraten/graph"><img src="http://goci.me/project/image/github.com/sauerbraten/graph" /></a>

A thread-safe implementation of a graph data structure in Go. See https://en.wikipedia.org/wiki/Graph_(abstract_data_type) for more information. This implementation is weighted, but undirected.

There is also a version of this package that supports storing values in the graph vertexes so that the graph can be used as a data store: https://github.com/sauerbraten/graph-store

## Usage

Get the package:

	$ go get -d github.com/sauerbraten/graph
	$ cd $GOROOT/src/github.com/sauerbraten/graph
	$ git checkout no-storage
	$ go install

Import the package:

	import (
		"github.com/sauerbraten/graph"
	)


## Documentation

For full package documentation, visit http://godoc.org/github.com/sauerbraten/graph.


## License

This code is licensed under a BSD License:

Copyright (c) 2013 Alexander Willing. All rights reserved.

- Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
- Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
