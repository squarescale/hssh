package libhssh

import (
	"fmt"
	"regexp"
	"testing"

	cr "github.com/squarescale/cloudresolver"
	"github.com/stretchr/testify/suite"
)

type FilterTestSuite struct {
	suite.Suite
}

func (s *FilterTestSuite) TestNewFilterFromString() {
	reg, err := regexp.Compile("foo")
	s.Nil(err)

	testCases := []struct {
		desc           string
		filter         string
		expectError    bool
		expectedFilter *Filter
	}{
		{
			desc:           "Missing ':'",
			filter:         "x",
			expectError:    true,
			expectedFilter: nil,
		},
		{
			desc:           "Too much ':'",
			filter:         "x:y:z",
			expectError:    true,
			expectedFilter: nil,
		},
		{
			desc:           "Unmatchable field name && valid value regexp",
			filter:         "xxx:valid",
			expectError:    true,
			expectedFilter: nil,
		},
		{
			desc:           "Matchable field name && valid value regexp",
			filter:         "id:[syntax error)",
			expectError:    true,
			expectedFilter: nil,
		},
		{
			desc:        "Valid field name && valid regexp",
			filter:      fmt.Sprintf("ip:%s", reg),
			expectError: false,
			expectedFilter: &Filter{
				matchableFields: []string{
					"PrivateIpv4",
					"PrivateIpv6",
					"PublicIpv4",
					"PublicIpv6",
				},
				pattern: reg,
			},
		},
	}

	for _, tc := range testCases {
		f, err := NewFilterFromString(
			tc.filter,
		)

		if tc.expectError {
			s.NotNil(err, tc.desc)
			s.Nil(f, tc.desc)
			continue
		}

		s.Nil(err, tc.desc)
		s.Equal(tc.expectedFilter, f, tc.desc)
	}
}

func (s *FilterTestSuite) TestHostMatch() {
	testCases := []struct {
		desc    string
		filter  string
		host    *cr.Host
		matches bool
	}{
		{
			desc:    "matches",
			filter:  "id:aaa",
			host:    &cr.Host{Id: "aaa"},
			matches: true,
		},
		{
			desc:    "does not match",
			filter:  "zone:euwest",
			host:    &cr.Host{Zone: "x"},
			matches: false,
		},
	}

	for _, tc := range testCases {
		f, err := NewFilterFromString(tc.filter)
		s.Nil(err)
		s.NotNil(f)

		s.Equal(
			tc.matches,
			f.HostMatch(tc.host),
			tc.desc,
		)
	}
}

func TestFilterTestSuite(t *testing.T) {
	suite.Run(t, new(FilterTestSuite))
}
