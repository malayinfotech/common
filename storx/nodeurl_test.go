// Copyright (C) 2019 Storx Labs, Inc.
// See LICENSE for copying information.

package storx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"common/pb"
	"common/storx"
)

func TestNodeURL(t *testing.T) {
	id, err := storx.NodeIDFromString("12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7")
	require.NoError(t, err)

	t.Run("Valid", func(t *testing.T) {
		type Test struct {
			String   string
			Expected storx.NodeURL
		}

		for _, testcase := range []Test{
			{"", storx.NodeURL{}},
			// host
			{"33.20.0.1:7777", storx.NodeURL{Address: "33.20.0.1:7777"}},
			{"[2001:db8:1f70::999:de8:7648:6e8]:7777", storx.NodeURL{Address: "[2001:db8:1f70::999:de8:7648:6e8]:7777"}},
			{"example.com:7777", storx.NodeURL{Address: "example.com:7777"}},
			// node id + host
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777", storx.NodeURL{ID: id, Address: "33.20.0.1:7777"}},
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@[2001:db8:1f70::999:de8:7648:6e8]:7777", storx.NodeURL{ID: id, Address: "[2001:db8:1f70::999:de8:7648:6e8]:7777"}},
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com:7777", storx.NodeURL{ID: id, Address: "example.com:7777"}},
			// node id
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@", storx.NodeURL{ID: id}},
			// debounce_limit
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@?debounce=3", storx.NodeURL{ID: id, DebounceLimit: 3}},
			// noise
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777?noise_proto=2",
				storx.NodeURL{
					ID:      id,
					Address: "33.20.0.1:7777",
					NoiseInfo: storx.NoiseInfo{
						Proto: storx.NoiseProto_IK_25519_AESGCM_BLAKE2b,
					},
				},
			},
			{"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777?noise_pub=12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7",
				storx.NodeURL{
					ID:      id,
					Address: "33.20.0.1:7777",
					NoiseInfo: storx.NoiseInfo{
						PublicKey: string(id.Bytes()),
					},
				},
			},
		} {
			url, err := storx.ParseNodeURL(testcase.String)
			require.NoError(t, err, testcase.String)

			assert.Equal(t, testcase.Expected, url)
			assert.Equal(t, testcase.String, url.String())

			copy := pb.NodeFromNodeURL(url).NodeURL()
			assert.Equal(t, testcase.Expected, copy)
			assert.Equal(t, testcase.String, copy.String())
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		for _, testcase := range []string{
			// invalid host
			"exampl e.com:7777",
			// invalid node id
			"12vha9oTFnerxgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@33.20.0.1:7777",
			"12vha9oTFnerx YRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@[2001:db8:1f70::999:de8:7648:6e8]:7777",
			"12vha9oTFnerxYRgeQ2BZqoFrLrn_5UWTCY2jA77dF3YvWew7@example.com:7777",
			// invalid node id
			"1112vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@",
		} {
			_, err := storx.ParseNodeURL(testcase)
			assert.Error(t, err, testcase)
		}
	})
}

func TestNodeURLs(t *testing.T) {
	id, err := storx.NodeIDFromString("12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7")
	require.NoError(t, err)

	s := "33.20.0.1:7777," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@[2001:db8:1f70::999:de8:7648:6e8]:7777," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com?noise_proto=2," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@example.com?noise_proto=1," +
		"12vha9oTFnerxYRgeQ2BZqoFrLrnmmf5UWTCY2jA77dF3YvWew7@"
	urls, err := storx.ParseNodeURLs(s)
	require.NoError(t, err)
	require.Equal(t, storx.NodeURLs{
		storx.NodeURL{Address: "33.20.0.1:7777"},
		storx.NodeURL{ID: id, Address: "[2001:db8:1f70::999:de8:7648:6e8]:7777"},
		storx.NodeURL{ID: id, Address: "example.com"},
		storx.NodeURL{ID: id, Address: "example.com", NoiseInfo: storx.NoiseInfo{Proto: storx.NoiseProto_IK_25519_AESGCM_BLAKE2b}},
		pb.NodeFromNodeURL(storx.NodeURL{ID: id, Address: "example.com", NoiseInfo: storx.NoiseInfo{Proto: storx.NoiseProto_IK_25519_ChaChaPoly_BLAKE2b}}).NodeURL(),
		storx.NodeURL{ID: id},
	}, urls)

	require.Equal(t, s, urls.String())
}
