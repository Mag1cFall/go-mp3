//go:build arm64 && !noasm

package enc

// quantizeGranule quantizes two lines per NEON iteration (NEON is baseline on
// arm64, so no runtime feature check). The kernel stores zeros for a group of
// lines whose scaled magnitudes are all below quantZeroMax, which the full
// sqrt chain rounds to 0 as well, and otherwise computes sqrt(t*sqrt(t)) for
// every line of the group, so the result is bit-identical
// (TestQuantizeGranuleParity). The skip keeps the low-bitrate zero tail cheap
// on cores whose vector square root is slow.
func quantizeGranule(xr *[576]float64, gg int, sf *scfState, lay *bandLayout, ix *[576]int32) {
	var steps [39]float64
	quantizeSteps(gg, sf, lay, &steps)
	quantizeNEON(xr, ix, &steps, &lay.width, lay.nBands, 0.4054, maxQuant, quantZeroMax)
}

// quantizeNEON quantizes band by band in quantizeGranuleGo's operation order.
// Every band width is even.
//
//go:noescape
func quantizeNEON(xr *[576]float64, ix *[576]int32, steps *[39]float64, widths *[39]int, nBands int, rounding, limit, zeroMax float64)
