package all

import (
	// add the momory db
	_ "github.com/fvdveen/mu2/db/memory"
	// add the postgres db
	_ "github.com/fvdveen/mu2/db/postgres"
)
