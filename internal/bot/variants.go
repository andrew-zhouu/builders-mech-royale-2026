package bot

// Variants are named parameter tweaks for A/B runs in the local arena (roster names).
func Variants() map[string]func(p *Params) {
	return map[string]func(p *Params){
		// the published v0.5 behaviour (stealth in free-for-alls, hit rates as believed then)
		"v05": func(p *Params) { v05(p) },
		// v0.5 plus the visibility-map hideouts tried on Saturday morning
		"v05-hm": func(p *Params) { v05(p); p.UseHideMap = true; p.QuietHoldS = 3 },
		// never adapts to the lobby: always the squad (fighter) stance
		"fighter": func(p *Params) { v05(p); p.PracticeAware = false },
		// always the stealth stance, from the first tick
		"ghost": func(p *Params) { p.ForceFFA = true },
		// stealth but no looting at all
		"ghost-noloot": func(p *Params) { p.ForceFFA = true; p.FFALootRadius = 0 },
		// hideouts from the visibility map
		"nohm":     func(p *Params) { p.UseHideMap = false },
		"hm-near":  func(p *Params) { p.HideDistW = 0.30 },
		"hm-far":   func(p *Params) { p.HideDistW = 0.08 },
		"hm-hit":   func(p *Params) { p.HideSeenW, p.HideHitW = 0.2, 1.3 },
		"hm-seen":  func(p *Params) { p.HideSeenW, p.HideHitW = 1.2, 0.4 },
		"hm-noroof": func(p *Params) { p.UseRoofs = false },
		"sfr70":    func(p *Params) { p.StealthFightRange = 70 },
		"sfr100":   func(p *Params) { p.StealthFightRange = 100 },
		"sfr200":   func(p *Params) { p.StealthFightRange = 200 },
		// squad stance: farm first, kill later
		"patient60":  func(p *Params) { p.PatientSecs = 60 },
		"patient120": func(p *Params) { p.PatientSecs = 120 },
		"patient200": func(p *Params) { p.PatientSecs = 200; p.PatientCubes = 15 },
		// kiting several attackers
		"nokite":   func(p *Params) { p.KiteMin = 0 },
		"kite80":   func(p *Params) { p.KiteMin, p.KiteMax = 80, 130 },
		"kite140":  func(p *Params) { p.KiteMin, p.KiteMax = 140, 200 },
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

		// --- after the first real free-for-all matches (stealth scored 4.8 of 15) ---------------------
		// always fight, with the measured hit rates instead of the believed ones
		"f-hit": func(p *Params) { p.PracticeAware = false; p.HitModel = 1 },
		// ... and the fight plan that counts our own fire while changing range
		"f-plan": func(p *Params) { p.PracticeAware = false; p.HitModel = 1; p.PlanV2 = true },
		// ... and no backing away from several attackers without a laser
		"f-nk": func(p *Params) { p.PracticeAware = false; p.HitModel = 1; p.PlanV2 = true; p.KiteNeedsLaser = true },
		// ... and a laser is worth going for
		"f-las": func(p *Params) { brawl(p) },
		"f-las4": func(p *Params) { brawl(p); p.LaserValue = 4 },
		"f-las8": func(p *Params) { brawl(p); p.LaserValue = 8; p.ArmSwapMin = 0.5 },
		// ... and the running exchange corrects the estimate
		"f-ex":  func(p *Params) { brawl(p); p.ExchangeW = 0.6 },
		"f-ex1": func(p *Params) { brawl(p); p.ExchangeW = 1.0 },
		// leave a one-on-one sooner / later
		"f-d35": func(p *Params) { brawl(p); p.DuelDisengage = 0.35 },
		"f-d70": func(p *Params) { brawl(p); p.DuelDisengage = 0.70 },
		"f-d90": func(p *Params) { brawl(p); p.DuelDisengage = 0.90 },
		// pick fights more / less readily
		"f-e10": func(p *Params) { brawl(p); p.EngageRatio = 1.0 },
		"f-e16": func(p *Params) { brawl(p); p.EngageRatio = 1.6 },
		"f-nohunt": func(p *Params) { brawl(p); p.Hunt = false },
		// only the measured tables, old stance logic (does the model alone change anything?)
		"hit1": func(p *Params) { p.HitModel = 1; p.PlanV2 = true },

		// "brawl": practice path as published, free-for-all (and the first 24 s) as an armed farmer
		"brawl":    func(p *Params) { brawl2(p) },
		"b-stand":  func(p *Params) { brawl2(p); p.AvoidAdv = 0 },            // never slip: always fight back
		"b-av10":   func(p *Params) { brawl2(p); p.AvoidAdv = 1.0 },          // slip only from losing exchanges
		"b-av16":   func(p *Params) { brawl2(p); p.AvoidAdv = 1.6 },
		"b-avr60":  func(p *Params) { brawl2(p); p.AvoidRange = 60 },
		"b-avr120": func(p *Params) { brawl2(p); p.AvoidRange = 120 },
		"b-cov90":  func(p *Params) { brawl2(p); p.AvoidCover = 90 },
		"b-las3":   func(p *Params) { brawl2(p); p.LaserValue = 3 },
		"b-las10":  func(p *Params) { brawl2(p); p.LaserValue = 10; p.ArmSwapMin = 0.5 },
		"b-e10":    func(p *Params) { brawl2(p); p.BrawlEngage = 1.0 },
		"b-e18":    func(p *Params) { brawl2(p); p.BrawlEngage = 1.8 },
		"b-nohunt": func(p *Params) { brawl2(p); p.BrawlHunt = false },
		"b-hl200":  func(p *Params) { brawl2(p); p.HuntRangeLaser = 200 },
		"b-ex":     func(p *Params) { brawl2(p); p.ExchangeW = 0.7 },
		"b-c10":    func(p *Params) { brawl2(p); p.BrawlWantCubes = 10 },
		"b-c25":    func(p *Params) { brawl2(p); p.BrawlWantCubes = 25 },
		"b-lr120":  func(p *Params) { brawl2(p); p.BrawlLootRadius = 120 },
		"b-d35":    func(p *Params) { brawl2(p); p.DuelDisengage = 0.35 },
		"b-d75":    func(p *Params) { brawl2(p); p.DuelDisengage = 0.75 },
		"b-oldhit": func(p *Params) { brawl2(p); p.HitModel = 0 },
		// v0.6 exactly as published (with the first-goal path bug)
		"v06":       func(p *Params) { brawl2(p); p.PathFix = false },
		// experiments on top of the v0.6 brawl: sprint to lasers; leave mortars alone
		"b-sprint":   func(p *Params) { brawl2(p); p.LootBoost = 30 },
		"b-sprint10": func(p *Params) { brawl2(p); p.LootBoost = 30; p.LaserValue = 10; p.ArmSwapMin = 0.5 },
		"b-nomortar": func(p *Params) { brawl2(p); p.NoMortar = true },
		// experiments on top of the v0.6 brawl: hold fire until touched; keep farming under light fire
		"b-prov":      func(p *Params) { brawl2(p); p.ProvokeRule = true },
		"b-plink":     func(p *Params) { brawl2(p); p.PlinkMaxDPS = 60 },
		"b-provplink": func(p *Params) { brawl2(p); p.ProvokeRule = true; p.PlinkMaxDPS = 60 },
	}
}

// v05 restores the parameters of the published v0.5.
func v05(p *Params) {
	p.Brawl = false
	p.HitModel = 0
	p.PlanV2 = false
	p.KiteNeedsLaser = false
	p.LaserValue = 3.0
	p.AvoidAdv = 0
	p.UseHideMap = false
	p.QuietHoldS = 0
}

// brawl2 keeps the published practice behaviour (squad stance after the opening) and replaces the
// stealth stances with the armed-farmer one.
func brawl2(p *Params) {
	p.Brawl = true
	p.HitModel = 1
	p.PlanV2 = true
	p.KiteNeedsLaser = true
	p.LaserValue = 6
	p.AvoidAdv = 1.25
}

// brawl is the always-fight parameter set built after the first real matches.
func brawl(p *Params) {
	p.PracticeAware = false
	p.HitModel = 1
	p.PlanV2 = true
	p.KiteNeedsLaser = true
	p.LaserValue = 6
}
