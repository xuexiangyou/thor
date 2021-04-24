package bq

type (
	Beanstalk struct {
		Endpoint string
		Tube     string
	}

	Redis struct {
		Host string
    }

	BqConf struct {
		Beanstalk []Beanstalk
		Redis Redis
	}
)
