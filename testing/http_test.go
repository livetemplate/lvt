package testing

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTPResponse_String(t *testing.T) {
	resp := &HTTPResponse{
		Response: &http.Response{},
		Body:     []byte("Hello, World!"),
	}

	if got := resp.String(); got != "Hello, World!" {
		t.Errorf("String() = %q, want %q", got, "Hello, World!")
	}
}

func TestHTTPResponse_HasTemplateErrors(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected bool
	}{
		{
			name:     "no errors",
			body:     "<html><body>Hello World</body></html>",
			expected: false,
		},
		{
			name:     "has dot expression",
			body:     "<html><body>{{.Name}}</body></html>",
			expected: true,
		},
		{
			name:     "has if block",
			body:     "<html><body>{{if .Show}}content{{end}}</body></html>",
			expected: true,
		},
		{
			name:     "has range block",
			body:     "<html><body>{{range .Items}}item{{end}}</body></html>",
			expected: true,
		},
		{
			name:     "has end block",
			body:     "<html><body>{{end}}</body></html>",
			expected: true,
		},
		{
			name:     "has else block",
			body:     "<html><body>{{else}}</body></html>",
			expected: true,
		},
		{
			name:     "has template block",
			body:     "<html><body>{{template \"header\"}}</body></html>",
			expected: true,
		},
		{
			name:     "has alternate delimiters",
			body:     "<html><body>[[.Name]]</body></html>",
			expected: true,
		},
		{
			name:     "normal braces ok",
			body:     "<html><body>{name: value}</body></html>",
			expected: false,
		},
		{
			name:     "javascript ok",
			body:     "<script>const x = {name: 'value'};</script>",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &HTTPResponse{
				Response: &http.Response{},
				Body:     []byte(tt.body),
			}
			if got := resp.HasTemplateErrors(); got != tt.expected {
				t.Errorf("HasTemplateErrors() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHTTPResponse_FindTemplateErrors(t *testing.T) {
	resp := &HTTPResponse{
		Response: &http.Response{},
		Body:     []byte("<html>{{.Name}} and {{if .Show}}content{{end}}</html>"),
	}

	errors := resp.FindTemplateErrors()
	if len(errors) < 2 {
		t.Errorf("FindTemplateErrors() found %d errors, expected at least 2", len(errors))
	}
}

func TestParseSelector(t *testing.T) {
	tests := []struct {
		selector      string
		expectedTag   string
		expectedID    string
		expectedClass string
	}{
		{"div", "div", "", ""},
		{"#main", "", "main", ""},
		{".container", "", "", "container"},
		{"div#main", "div", "main", ""},
		{"div.container", "div", "", "container"},
		{"tbody tr", "tr", "", ""}, // Descendant selector - takes last part
	}

	for _, tt := range tests {
		t.Run(tt.selector, func(t *testing.T) {
			tag, id, class := parseSelector(tt.selector)
			if tag != tt.expectedTag {
				t.Errorf("tag = %q, want %q", tag, tt.expectedTag)
			}
			if id != tt.expectedID {
				t.Errorf("id = %q, want %q", id, tt.expectedID)
			}
			if class != tt.expectedClass {
				t.Errorf("class = %q, want %q", class, tt.expectedClass)
			}
		})
	}
}

// Integration tests with real HTTP server
func TestHTTPAssert_StatusCodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/redirect":
			http.Redirect(w, r, "/ok", http.StatusFound)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		case "/unauthorized":
			w.WriteHeader(http.StatusUnauthorized)
		case "/forbidden":
			w.WriteHeader(http.StatusForbidden)
		case "/badrequest":
			w.WriteHeader(http.StatusBadRequest)
		case "/servererror":
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	t.Run("StatusOK", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/ok")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()
		assert := NewHTTPAssert(httpResp)
		assert.StatusOK(t)
	})

	t.Run("StatusRedirect", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/redirect")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()
		assert := NewHTTPAssert(httpResp)
		assert.StatusRedirect(t)
	})

	t.Run("StatusNotFound", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/notfound")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()
		assert := NewHTTPAssert(httpResp)
		assert.StatusNotFound(t)
	})

	t.Run("StatusUnauthorized", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/unauthorized")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()
		assert := NewHTTPAssert(httpResp)
		assert.StatusUnauthorized(t)
	})

	t.Run("StatusForbidden", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/forbidden")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()
		assert := NewHTTPAssert(httpResp)
		assert.StatusForbidden(t)
	})

	t.Run("StatusBadRequest", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/badrequest")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()
		assert := NewHTTPAssert(httpResp)
		assert.StatusBadRequest(t)
	})

	t.Run("StatusServerError", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/servererror")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()
		assert := NewHTTPAssert(httpResp)
		assert.StatusServerError(t)
	})
}

func TestHTTPAssert_Contains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Hello World</body></html>"))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.Contains(t, "Hello")
	assert.Contains(t, "World")
	assert.NotContains(t, "Goodbye")
}

func TestHTTPAssert_NoTemplateErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Clean HTML</body></html>"))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.NoTemplateErrors(t)
}

func TestHTTPAssert_HasElement(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
				<body>
					<div id="main" class="container">
						<h1>Title</h1>
						<p class="intro">Intro text</p>
						<table>
							<tbody>
								<tr><td>Row 1</td></tr>
								<tr><td>Row 2</td></tr>
							</tbody>
						</table>
					</div>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.HasElement(t, "div")
	assert.HasElement(t, "h1")
	assert.HasElement(t, "#main")
	assert.HasElement(t, ".container")
	assert.HasElement(t, "div.container")
	assert.HasElement(t, "div#main")
	assert.HasElement(t, "p.intro")
	assert.HasNoElement(t, "span")
	assert.HasNoElement(t, "#nonexistent")
}

func TestHTTPAssert_ElementCount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<table>
				<tbody>
					<tr><td>Row 1</td></tr>
					<tr><td>Row 2</td></tr>
					<tr><td>Row 3</td></tr>
				</tbody>
			</table>
		`))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.ElementCount(t, "tr", 3)
	assert.ElementCount(t, "td", 3)
}

func TestHTTPAssert_TableRowCount(t *testing.T) {
	// Note: Simple selector implementation doesn't support true descendant selectors
	// "tbody tr" is parsed as just "tr", counting all tr elements
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<table>
				<tbody>
					<tr><td>Row 1</td></tr>
					<tr><td>Row 2</td></tr>
				</tbody>
			</table>
		`))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.TableRowCount(t, 2)
}

func TestHTTPAssert_ContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<html></html>"))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.ContentTypeHTML(t)
}

func TestHTTPAssert_Header(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom", "test-value")
		w.Write([]byte(""))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.Header(t, "X-Custom", "test-value")
	assert.HasHeader(t, "X-Custom")
}

func TestHTTPAssert_FormField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<form>
				<input type="text" name="username" value="john">
				<input type="hidden" name="csrf_token" value="abc123">
				<textarea name="bio">Hello</textarea>
			</form>
		`))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.HasFormField(t, "username")
	assert.HasFormField(t, "csrf_token")
	assert.FormFieldValue(t, "username", "john")
	assert.HasCSRFToken(t)
}

func TestHTTPAssert_RedirectTo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	}))
	defer server.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, _ := client.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.StatusRedirect(t)
	assert.RedirectTo(t, "/dashboard")
}

func TestHTTPAssert_ContainsAll(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Hello World Welcome</body></html>"))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.ContainsAll(t, "Hello", "World", "Welcome")
}

func TestHTTPAssert_Matches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>User ID: 12345</body></html>"))
	}))
	defer server.Close()

	resp, _ := http.Get(server.URL)
	httpResp := &HTTPResponse{Response: resp}
	httpResp.readBody()
	assert := NewHTTPAssert(httpResp)

	assert.Matches(t, `User ID: \d+`)
}

// Integration test - full workflow
func TestHTTPTest_FullWorkflow(t *testing.T) {
	// Create a simple test server
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<html>
				<body>
					<h1>Welcome</h1>
					<form action="/submit" method="POST">
						<input type="hidden" name="csrf_token" value="test123">
						<input type="text" name="name" value="">
						<button type="submit">Submit</button>
					</form>
				</body>
			</html>
		`))
	})
	mux.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		r.ParseForm()
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "Name required", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/success", http.StatusSeeOther)
	})
	mux.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body><h1>Success!</h1></body></html>"))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Test GET request
	t.Run("GET home page", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/")
		if err != nil {
			t.Fatalf("GET failed: %v", err)
		}
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()

		assert := NewHTTPAssert(httpResp)
		assert.StatusOK(t)
		assert.Contains(t, "Welcome")
		assert.HasElement(t, "h1")
		assert.HasElement(t, "form")
		assert.HasCSRFToken(t)
		assert.NoTemplateErrors(t)
	})

	// Test POST with missing field
	t.Run("POST without required field", func(t *testing.T) {
		resp, err := client.PostForm(server.URL+"/submit", nil)
		if err != nil {
			t.Fatalf("POST failed: %v", err)
		}
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()

		assert := NewHTTPAssert(httpResp)
		assert.StatusBadRequest(t)
	})

	// Test successful POST
	t.Run("POST with valid data", func(t *testing.T) {
		resp, err := client.PostForm(server.URL+"/submit", url.Values{
			"name": {"John"},
		})
		if err != nil {
			t.Fatalf("POST failed: %v", err)
		}
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()

		assert := NewHTTPAssert(httpResp)
		assert.StatusRedirect(t)
		assert.RedirectTo(t, "/success")
	})

	// Follow redirect
	t.Run("Follow redirect to success", func(t *testing.T) {
		resp, _ := client.Get(server.URL + "/success")
		httpResp := &HTTPResponse{Response: resp}
		httpResp.readBody()

		assert := NewHTTPAssert(httpResp)
		assert.StatusOK(t)
		assert.Contains(t, "Success!")
	})
}
