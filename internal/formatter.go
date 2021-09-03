package internal

type Formatter interface {
	Format(envs []EnvVar) []string
}

type DelimiterFormatter struct {
	delimiter string
}

func (d *DelimiterFormatter) Format(envs []EnvVar) []string {
	output := make([]string, 0)
	for _, env := range envs {
		output = append(output, env.Name+d.delimiter+env.Value)
	}
	return output
}

func NewDelimiterFormatter(delimiter string) Formatter {
	return &DelimiterFormatter{delimiter: delimiter}
}
