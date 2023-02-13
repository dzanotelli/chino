package custodia

type Repository struct {
	repository_id string
	description string
}

func (r *Repository) getId() string {
	return r.repository_id 
}

func (r *Repository) getDescription() string {
	return r.description
}
