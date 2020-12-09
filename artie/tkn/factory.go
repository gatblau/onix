package tkn

const ApiVersion = "tekton.dev/v1alpha1"

func NewTask(builderImage string) *Task {
	t := new(Task)
	t.APIVersion = ApiVersion
	t.Kind = "Task"
	t.Spec = Spec{
		Inputs: Inputs{
			Resources: []Resources{
				{
					Name: "source",
					Type: "git",
				},
			},
		},
		Steps: []Steps{
			{
				Name:  "apply",
				Image: builderImage,
			},
		},
	}
	return t
}
