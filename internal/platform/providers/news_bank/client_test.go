package newsbank

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchNews(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[
			{"link":"https://example.com/1","title":{"rendered":"One"}},
			{"link":"https://example.com/2","title":{"rendered":"Two"}}
		]`))
	}))
	defer server.Close()

	titles, err := FetchNews(server.URL)
	require.NoError(t, err)
	assert.Equal(t, []string{"One", "Two"}, titles)
}

func TestFetchNewsReturnsErrorOnBadStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	titles, err := FetchNews(server.URL)
	require.Error(t, err)
	assert.Nil(t, titles)
	assert.Contains(t, err.Error(), "api responded with status: 500")
}
