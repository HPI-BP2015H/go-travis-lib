package travis

type RepositoryService struct {
	client *Client
}

type Repository struct {
	ID             *int    `json:"id,omitempty"`
	Name           *string `json:"name,omitempty"`
	Slug           *string `json:"slug,omitempty"`
	Description    *string `json:"description,omitempty"`
	GithubLanguage *string `json:"github_language,omitempty"`
	Active         *bool   `json:"active,omitempty"`
	Private        *bool   `json:"private,omitempty"`
	//Owner          *string //*Owner
	//DefaultBranch  *string //*Branch
	//Starred        *string //Unknown
}

func (r *RepositoryService) List() []Repository {

	req, err := r.client.NewRequest("GET", "repos")
	if err != nil {
		println("Error during new request")
		return nil
	}
	repos := new([]Repository)
	resp, err := r.client.Do(req, repos)
	println(resp) // TODO: remove
	if err != nil {
		println("Error during travisclient.Do")
		return nil
	}
	return *repos
}
