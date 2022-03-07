package cats

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock_cats "github.com/itsHabib/go-kt-1/cats/mocks"
)

func TestClient_GetBreeds(t *testing.T) {
	defaultBody := func() []byte {
		return []byte(`[{"url": "url.com"}]`)
	}
	for _, tc := range []struct {
		desc    string
		wantErr bool
		errStr string
		breedId string
		doer func(t *testing.T, ctrl *gomock.Controller) HttpDo
	}{
		{
			desc: "should return and error if breed id is empty",
			wantErr: true,
			errStr: EmptyBreedID.Error(),
		},
		{
			desc:    "get breeds requests failure should result in an error",
			wantErr: true,
			breedId: "mcoo",
			doer: func(t *testing.T, ctrl *gomock.Controller) HttpDo {
				m := mock_cats.NewMockHttpDo(ctrl)
				m.
					EXPECT().
					Do(gomock.Any()).
					Return(nil, errors.New("random"))

				return m
			},
		},
		{
			desc:    "get breeds requests should throw an error on a non-200 response",
			wantErr: true,
			breedId: "mcoo",
			doer: func(t *testing.T, ctrl *gomock.Controller) HttpDo {
				m := mock_cats.NewMockHttpDo(ctrl)
				m.
					EXPECT().
					Do(gomock.Any()).
					DoAndReturn(func(req *http.Request) (*http.Response, error) {
						require.NotNil(t, req)
						url := req.URL
						query := url.Query()
						require.NotNil(t, query)

						assert.Equal(t, "mcoo", query.Get(breedIdQueryKey))
						assert.NotEmpty(t, req.Header.Get(apiKeyHeader))

						body := defaultBody()
						resp := http.Response{
							StatusCode:       400,
							Body:           io.NopCloser(bytes.NewReader(body)),
							ContentLength:    int64(len(body)),
						}

						return &resp, nil
					})

				return m
			},
		},
		{
			desc: "happy path",
			breedId: "mcoo",
			doer: func(t *testing.T, ctrl *gomock.Controller) HttpDo {
				m := mock_cats.NewMockHttpDo(ctrl)
				m.
					EXPECT().
					Do(gomock.Any()).
					DoAndReturn(func(req *http.Request) (*http.Response, error) {
						require.NotNil(t, req)
						url := req.URL
						query := url.Query()
						require.NotNil(t, query)

						assert.Equal(t, "mcoo", query.Get(breedIdQueryKey))
						assert.NotEmpty(t, req.Header.Get(apiKeyHeader))

						body := defaultBody()
						resp := http.Response{
							StatusCode:       200,
							Body:           io.NopCloser(bytes.NewReader(body)),
							ContentLength:    int64(len(body)),
						}

						return &resp, nil
					})

				return m
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			if tc.doer == nil {
				tc.doer = func(_ *testing.T, _ *gomock.Controller) HttpDo { return mock_cats.NewMockHttpDo(ctrl)}
			}

			c, err := NewClient("key")
			require.NoError(t, err)
			c.c = tc.doer(t, ctrl)

			resp, err := c.GetBreeds(tc.breedId)
			if tc.wantErr {
				assert.Error(t, err)
				if tc.errStr != "" {
					assert.EqualError(t, err, tc.errStr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
