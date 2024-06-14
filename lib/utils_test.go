package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCapitalizeName(t *testing.T) {
	require.Equal(t, CapitalizeName("hello, world"), "helloWorld")
	require.Equal(t, CapitalizeName("k8s cluster overview"), "k8sClusterOverview")
	require.Equal(t, CapitalizeName("hamsa-api overview"), "hamsaApiOverview")
}
