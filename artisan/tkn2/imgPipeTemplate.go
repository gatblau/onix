package tkn2

type imgTypeTemplate struct {
}

func newImgPipeTask(c *imgPipeConf) *Task {
	args, _ := c.getArgs()
	t := new(Task)
	t.APIVersion = ApiVersionTekton
	t.Kind = "Task"
	t.Metadata = &Metadata{
		Name: c.buildTaskName(),
	}
	t.Spec = &Spec{
		Steps: []*Steps{
			{
				Name:       "build-image",
				Image:      c.BuilderImage,
				WorkingDir: "/workspace/source",
				Command:    []string{"/bin/bash", "-c"},
				Args:       args,
			},
		},
	}
	return t
}
