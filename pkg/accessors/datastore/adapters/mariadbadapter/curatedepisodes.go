package mariadbadapter

import "github.com/jecolasurdo/tbtlarchivist/pkg/contracts"

// UpsertEpisodeInfo inserts or updates clip info.
func (m *MariaDbConnection) UpsertEpisodeInfo(clipInfo contracts.EpisodeInfo) error {
	panic("not implemented")

	/*

	   TODO:
	   There are far more clips on the marsupial gurgle site than I'd anticipated.
	   (~57000). In order to be able to complete a meaningful amount of research, it
	   will make sense to prioritize the clips. To handle this, I will create an
	   integer field in ClipInfo that represents a clips "priority". The archivist
	   will assign lower priority clips after all high priority clips have been
	   researched for all episodes. For example: let's say we have three clip
	   curators.  One gets "commonly used tbtl drops" and assigns them a priority of
	   100.  The second gets "grab bag of drops" and assigns them a priority of 200.
	   The third gets all drops from the search engine and assigns them a priority of
	   300.  When the clip curator pushes these to the data store, it is likely that a
	   clip will be found in multiple priority categories. When this happens, the
	   higher priority clip is retained. Thus, if the same clip is found in the search
	   page and in the "commonly used drops" page, it will retain the higher priority
	   of 100 regardless of which source was curated first.  When an archivist assigns
	   work, it will only assign the 100 level priority clips until all 100 level
	   priority clips have been studied. Once all 100 level clips have been studied,
	   the archivist will begin to assign work for the 200 level clips.  Once all 200s
	   have been studied, work will begin on the 300s, etc...

	   I will create separate curators, one for each priority level. This will reduce
	   the scraping complexity of any individual curator.

	*/

}
