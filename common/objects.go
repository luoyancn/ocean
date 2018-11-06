package common

type base struct {
	Name string `json:"name"`
	Id   string `json:"id,omitempty"`
}

type domain struct {
	base
}

type project struct {
	Domain domain `json:"domain"`
	base
}

type user struct {
	Domain domain `json:"domain"`
	base
	PassWord string `json:"password,omitempty"`
	Expires  string `json:"password_expires_at,omitempty"`
}

type scope struct {
	Project project `json:"project"`
}

type password struct {
	User user `json:"user"`
}

type identity struct {
	Methods  []string `json:"methods"`
	Password password `json:"password"`
}

type auth struct {
	Identity identity `json:"identity"`
	Scope    scope    `json:"scope"`
}

type credAuth struct {
	Auth auth `json:"auth"`
}

func NewAuth(username, passwd, projectname, userdomain,
	projectdomain string) *credAuth {

	userDomain := domain{}
	userDomain.Name = userdomain

	projectDomain := domain{}
	projectDomain.Name = projectdomain

	usr := user{PassWord: passwd, Domain: userDomain}
	usr.Name = username

	proj := project{Domain: projectDomain}
	proj.Name = projectname

	return &credAuth{
		Auth: auth{
			Identity: identity{
				Methods: []string{"password"},
				Password: password{
					User: usr,
				},
			},
			Scope: scope{
				Project: proj,
			},
		},
	}
}

type role struct {
	base
}

type endpoint struct {
	Id        string `json:"id"`
	Interface string `json:"interface"`
	Region    string `json:"region"`
	RegionId  string `json:"region_id"`
	Url       string `json:"url"`
}

type catalog struct {
	Endpoints []endpoint `json:"endpoints"`
	base
	Type string `json:"type"`
}

type token struct {
	AuditIds []string  `json:"audit_ids"`
	Catalog  []catalog `json:"catalog"`
	Expires  string    `json:"expires_at"`
	IsDomain bool      `json:"is_domain"`
	IssuedAt string    `json:"issued_at"`
	Methods  []string  `json:"methods"`
	Project  project   `json:"project"`
	Role     []role    `json:"roles"`
	User     user      `json:"user"`
}

type RespToken struct {
	Token token `json:"token"`
}
