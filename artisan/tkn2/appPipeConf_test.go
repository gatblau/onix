package tkn2

import (
	"testing"
)

func TestSurvey(t *testing.T) {
	// collects information to assemble the pipeline
	c := NewAppPipelineConfig(".", "", false)
	// assembles the pipeline
	_ = MergeArtPipe(c, false)
}
