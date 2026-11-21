package bot

// Cover that lasts.
//
// findCover returns the nearest cell out of a threat's sight, which by construction lies on the very edge
// of the shadow: in the real matches our mech, having got out of an attacker's sight, was seen again
// within 10 s in 79% of cases (the most of any team), typically after a second or three in which the
// attacker had moved two or three metres and we had not moved at all. The best survivor on the board
// was seen again half as often: it keeps walking. deeperCover is the second step: from the edge of the
// shadow to a place that stays hidden when the attacker shifts sideways, and, with a hurt mech and the
// attacker close, a place clearly farther from it as well (health comes back just as fast on the move).

// deeperCover returns the nearest cell that is hidden from every recent threat's last known position and
// from CoverMargin metres to either side of it, with the straight way there out of sight too. With
// farther set and the nearest threat inside RunOnRange, the cell must also be RunOnGain metres farther
// from that threat than we are. It returns -1 when where we stand already qualifies, or nothing does.
func (b *Bot) deeperCover(ths []threat, maxCost float64, farther bool) int {
	type eyeT struct{ x, y, z float64 }
	var eyes []eyeT
	d0 := 1e9
	for i := range ths {
		tr := ths[i].tr
		if (b.t-tr.LastSeen > 6 && !ths[i].phantom) || ths[i].d > 260 {
			continue
		}
		sc := tr.scale(b.sizePerCube)
		eyes = append(eyes, eyeT{tr.Pos[0], tr.Pos[1], tr.Pos[2] + 8*sc})
		if d := hyp(tr.Pos[0]-b.pos[0], tr.Pos[1]-b.pos[1]); d < d0 {
			d0 = d
		}
	}
	if len(eyes) == 0 {
		return -1
	}
	needFar := farther && b.P.RunOnRange > 0 && d0 <= b.P.RunOnRange
	margin := b.P.CoverMargin
	if margin <= 0 && !needFar {
		return -1
	}
	seen := func(ex, ey, ez, x, y, gz float64) bool {
		return b.W.Visible(ex, ey, ez, x, y, gz+b.height*0.5) || b.W.Visible(ex, ey, ez, x, y, gz+b.height*0.95)
	}
	hidden := func(x, y, gz float64) bool {
		for _, e := range eyes {
			if seen(e.x, e.y, e.z, x, y, gz) {
				return false
			}
		}
		return true
	}
	// hidden with room to spare: also from where the threat would be after a step to either side
	robust := func(x, y, gz float64) bool {
		for _, e := range eyes {
			if seen(e.x, e.y, e.z, x, y, gz) {
				return false
			}
			if margin <= 0 {
				continue
			}
			dx, dy := x-e.x, y-e.y
			l := hyp(dx, dy)
			if l < 1 {
				return false
			}
			px, py := -dy/l*margin, dx/l*margin
			if seen(e.x+px, e.y+py, e.z, x, y, gz) || seen(e.x-px, e.y-py, e.z, x, y, gz) {
				return false
			}
		}
		return true
	}
	if !needFar && robust(b.pos[0], b.pos[1], b.pos[2]) {
		return -1
	}
	z := &b.o.Zone
	found := -1
	checked := 0
	b.Nav.Flood(b.pos[0], b.pos[1], b.radius+0.15, b.height, maxCost, nil, func(k int, cost float32) bool {
		if (k%b.Nav.NX+k/b.Nav.NX)%2 != 0 {
			return true
		}
		if float64(b.Nav.Clear[k]) < b.radius+0.6 {
			return true
		}
		x, y := b.Nav.Pos(k)
		if hyp(x-z.Center[0], y-z.Center[1]) > z.Radius-6 {
			return true
		}
		if z.ShrinkEnd-b.t < 25 && hyp(x-z.NextCenter[0], y-z.NextCenter[1]) > z.NextRadius-4 {
			return true
		}
		if needFar {
			de := 1e9
			for _, e := range eyes {
				if d := hyp(x-e.x, y-e.y); d < de {
					de = d
				}
			}
			if de < d0+b.P.RunOnGain {
				return true
			}
		}
		checked++
		if checked > 300 {
			return false
		}
		if !robust(x, y, float64(b.Nav.Z[k])) {
			return true
		}
		// the straight way there: walkable points a quarter, half and three quarters along must be hidden too
		for _, f := range [3]float64{0.25, 0.5, 0.75} {
			mx, my := b.pos[0]+(x-b.pos[0])*f, b.pos[1]+(y-b.pos[1])*f
			mk := b.Nav.Index(mx, my)
			if b.Nav.Solid[mk] {
				continue // the walk goes round this
			}
			if !hidden(mx, my, float64(b.Nav.Z[mk])) {
				return true
			}
		}
		found = k
		return false
	})
	return found
}
