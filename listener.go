package main

type IListener interface {
	GetSensuChan() chan *Alert
	Start()
}

type SensuResult struct {
	Client string     `json:"client"` // the client the check came from
	Check  SensuCheck `json:"check"`  // the check info
}

type SensuCheck struct {
	Name         string  `json:"name"`
	CapsuleName  string  `json:"capsule_name"`
	Output       string  `json:"output"`
	Status       float64 `json:"status"`
	CapsuleId    string  `json:"capsule_id,omitempty"`
	DeploymentId string  `json:"deployment_id,omitempty"`
	AccountSlug  string  `json:"account,omitempty"`
}
