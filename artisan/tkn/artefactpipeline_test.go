package tkn

import (
	"testing"
)

func TestSurvey(t *testing.T) {
	// collects information to assemble the pipeline
	c := NewArtPipelineConfig(".", "", false)
	// assembles the pipeline
	_ = MergeArtPipe(c, false)
}
