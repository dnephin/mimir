// SPDX-License-Identifier: AGPL-3.0-only

package storegateway

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/prometheus/prometheus/model/labels"

	"github.com/grafana/mimir/pkg/storegateway/indexcache"
)

func TestSnappyGobSeriesCacheEntryCodec(t *testing.T) {
	type testType struct {
		LabelSets   []labels.Labels
		MatchersKey indexcache.LabelMatchersKey
	}

	entry := testType{
		LabelSets: []labels.Labels{
			labels.FromStrings("foo", "bar"),
			labels.FromStrings("baz", "boo"),
		},
		MatchersKey: indexcache.CanonicalLabelMatchersKey([]*labels.Matcher{labels.MustNewMatcher(labels.MatchRegexp, "foo", "bar")}),
	}

	t.Run("happy case roundtrip", func(t *testing.T) {
		data, err := encodeSnappyGob(entry)
		require.NoError(t, err)

		var decoded testType
		err = decodeSnappyGob(data, &decoded)
		require.NoError(t, err)
		require.Equal(t, entry, decoded)
	})

	t.Run("can't decode wrong codec", func(t *testing.T) {
		data, err := encodeSnappyGob(entry)
		require.NoError(t, err)

		data[0] = 'x'

		var decoded testType
		err = decodeSnappyGob(data, &decoded)
		require.Error(t, err)
	})

	t.Run("can't decode wrong data", func(t *testing.T) {
		data, err := encodeSnappyGob(entry)
		require.NoError(t, err)

		data = data[:len(gobCodecPrefix)+1]

		var decoded testType
		err = decodeSnappyGob(data, &decoded)
		require.Error(t, err)
	})
}
