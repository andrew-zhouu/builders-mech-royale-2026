package bot

// Variants are named parameter tweaks for A/B runs in the local arena (roster names).
func Variants() map[string]func(p *Params) {
	return map[string]func(p *Params){
		// never adapts to the lobby: always the squad (fighter) stance
		"fighter": func(p *Params) { p.PracticeAware = false },
		// always the stealth stance, from the first tick
		"ghost": func(p *Params) { p.ForceFFA = true },
		// stealth but no looting at all
		"ghost-noloot": func(p *Params) { p.ForceFFA = true; p.FFALootRadius = 0 },
		// opening: roof before weapons
		"rooffirst":   func(p *Params) { p.RoofFirstSecs = 40; p.RoofFirstDist = 90 },
		"rooffirst60": func(p *Params) { p.RoofFirstSecs = 60; p.RoofFirstDist = 150 },
		// multi-attacker leave threshold
		"ttd6":  func(p *Params) { p.LeaveTTD = 6 },
		"ttd16": func(p *Params) { p.LeaveTTD = 16 },
		// never leave a one-on-one
		"duel0": func(p *Params) { p.UseRoofs = false; p.DuelDisengage = 0 },
		"duel25": func(p *Params) { p.UseRoofs = false; p.DuelDisengage = 0.25 },
		"duel75": func(p *Params) { p.UseRoofs = false; p.DuelDisengage = 0.75 },
		// no rooftop hiding
		"noroof": func(p *Params) { p.UseRoofs = false },
		// rooftops in the squad stance as well
		"roofsquad": func(p *Params) { p.RoofInSquad = true },
		// fighter stance without roofs (the pre-roof baseline for the squad stance)
		"fighter-noroof": func(p *Params) { p.PracticeAware = false; p.UseRoofs = false },
		// stealth, greedier looting
		"ghost-greedy": func(p *Params) { p.ForceFFA = true; p.FFALootRadius = 260; p.FFAWantCubes = 20 },
	}
}
