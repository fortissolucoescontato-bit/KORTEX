package filemerge

import (
	"strings"
	"testing"
)

func TestInjectMarkdownSection_EmptyFile(t *testing.T) {
	result := InjectMarkdownSection("", "sdd", "## SDD Config\nSome content here.\n")

	want := "<!-- kortex:sdd -->\n## SDD Config\nSome content here.\n<!-- /kortex:sdd -->\n"
	if result != want {
		t.Fatalf("empty file inject:\ngot:  %q\nwant: %q", result, want)
	}
}

func TestInjectMarkdownSection_AppendToExistingContent(t *testing.T) {
	existing := "# My Config\n\nSome existing content.\n"
	result := InjectMarkdownSection(existing, "persona", "You are a senior architect.\n")

	want := "# My Config\n\nSome existing content.\n\n<!-- kortex:persona -->\nYou are a senior architect.\n<!-- /kortex:persona -->\n"
	if result != want {
		t.Fatalf("append to existing:\ngot:  %q\nwant: %q", result, want)
	}
}

func TestInjectMarkdownSection_UpdateExistingSection(t *testing.T) {
	existing := "# Config\n\n<!-- kortex:sdd -->\nOld SDD content.\n<!-- /kortex:sdd -->\n\nOther stuff.\n"
	result := InjectMarkdownSection(existing, "sdd", "New SDD content.\n")

	want := "# Config\n\n<!-- kortex:sdd -->\nNew SDD content.\n<!-- /kortex:sdd -->\n\nOther stuff.\n"
	if result != want {
		t.Fatalf("update existing section:\ngot:  %q\nwant: %q", result, want)
	}
}

func TestInjectMarkdownSection_MultipleSectionsOnlyTargetedOneUpdated(t *testing.T) {
	existing := "# Config\n\n<!-- kortex:persona -->\nPersona content.\n<!-- /kortex:persona -->\n\n<!-- kortex:sdd -->\nOld SDD.\n<!-- /kortex:sdd -->\n\n<!-- kortex:skills -->\nSkills content.\n<!-- /kortex:skills -->\n"

	result := InjectMarkdownSection(existing, "sdd", "Updated SDD.\n")

	// persona and skills should be unchanged
	want := "# Config\n\n<!-- kortex:persona -->\nPersona content.\n<!-- /kortex:persona -->\n\n<!-- kortex:sdd -->\nUpdated SDD.\n<!-- /kortex:sdd -->\n\n<!-- kortex:skills -->\nSkills content.\n<!-- /kortex:skills -->\n"
	if result != want {
		t.Fatalf("multiple sections:\ngot:  %q\nwant: %q", result, want)
	}
}

func TestInjectMarkdownSection_PreserveUserContentBeforeAndAfter(t *testing.T) {
	existing := "# User's custom intro\n\nHand-written notes.\n\n<!-- kortex:persona -->\nAuto persona.\n<!-- /kortex:persona -->\n\n# User's custom footer\n\nMore hand-written content.\n"

	result := InjectMarkdownSection(existing, "persona", "Updated persona.\n")

	want := "# User's custom intro\n\nHand-written notes.\n\n<!-- kortex:persona -->\nUpdated persona.\n<!-- /kortex:persona -->\n\n# User's custom footer\n\nMore hand-written content.\n"
	if result != want {
		t.Fatalf("preserve user content:\ngot:  %q\nwant: %q", result, want)
	}
}

func TestInjectMarkdownSection_MalformedMarkersTreatedAsNotFound(t *testing.T) {
	// Only opening marker, no closing marker — treat as not found, append.
	existing := "# Config\n\n<!-- kortex:sdd -->\nOrphaned content.\n"
	result := InjectMarkdownSection(existing, "sdd", "New SDD content.\n")

	// Should append since closing marker is missing.
	if result == existing {
		t.Fatalf("malformed markers: expected content to be appended, but got unchanged result")
	}

	// Result should contain the new properly-formed section.
	wantOpen := "<!-- kortex:sdd -->\nNew SDD content.\n<!-- /kortex:sdd -->\n"
	if !strings.Contains(result, wantOpen) {
		t.Fatalf("malformed markers: result should contain proper section:\ngot: %q", result)
	}
}

func TestInjectMarkdownSection_CloseBeforeOpenTreatedAsNotFound(t *testing.T) {
	// Closing marker appears before opening — treat as not found.
	existing := "<!-- /kortex:sdd -->\nSome content.\n<!-- kortex:sdd -->\n"
	result := InjectMarkdownSection(existing, "sdd", "New content.\n")

	// Should append the section, not replace.
	wantSuffix := "<!-- kortex:sdd -->\nNew content.\n<!-- /kortex:sdd -->\n"
	if !strings.HasSuffix(result, wantSuffix) {
		t.Fatalf("close-before-open: expected appended section:\ngot: %q\nwant suffix: %q", result, wantSuffix)
	}
}

func TestInjectMarkdownSection_EmptyContentRemovesSection(t *testing.T) {
	existing := "# Config\n\n<!-- kortex:sdd -->\nSDD content here.\n<!-- /kortex:sdd -->\n\nOther stuff.\n"
	result := InjectMarkdownSection(existing, "sdd", "")

	want := "# Config\n\nOther stuff.\n"
	if result != want {
		t.Fatalf("empty content removes section:\ngot:  %q\nwant: %q", result, want)
	}
}

func TestInjectMarkdownSection_EmptyContentOnMissingSectionNoOp(t *testing.T) {
	existing := "# Config\n\nSome content.\n"
	result := InjectMarkdownSection(existing, "sdd", "")

	if result != existing {
		t.Fatalf("empty content on missing section should be no-op:\ngot:  %q\nwant: %q", result, existing)
	}
}

func TestInjectMarkdownSection_ContentWithoutTrailingNewline(t *testing.T) {
	result := InjectMarkdownSection("", "test", "no trailing newline")

	want := "<!-- kortex:test -->\nno trailing newline\n<!-- /kortex:test -->\n"
	if result != want {
		t.Fatalf("content without trailing newline:\ngot:  %q\nwant: %q", result, want)
	}
}

func TestInjectMarkdownSection_ExistingWithoutTrailingNewline(t *testing.T) {
	existing := "# Title"
	result := InjectMarkdownSection(existing, "test", "Content.\n")

	want := "# Title\n\n<!-- kortex:test -->\nContent.\n<!-- /kortex:test -->\n"
	if result != want {
		t.Fatalf("existing without trailing newline:\ngot:  %q\nwant: %q", result, want)
	}
}

// --- StripLegacyPersonaBlock tests ---

const legacyPersonaBlock = `## Rules

- NEVER add "Co-Authored-By" or any AI attribution to commits.

## Personality

Senior Architect, 15+ years experience, GDE & MVP.

## Language

- Spanish input → Rioplatense Spanish.

`

const kortexMarkerSection = `<!-- kortex:persona -->
## Personality

Senior Architect, 15+ years experience, GDE & MVP.
<!-- /kortex:persona -->
`

func TestStripLegacyPersonaBlock_NoFingerprintReturnsSame(t *testing.T) {
	input := "# My Config\n\nSome unrelated user content.\n"
	result := StripLegacyPersonaBlock(input)
	if result != input {
		t.Fatalf("no fingerprint: expected unchanged result:\ngot:  %q\nwant: %q", result, input)
	}
}

func TestStripLegacyPersonaBlock_FingerprintInsideMarkerReturnsSame(t *testing.T) {
	// Fingerprints only exist inside kortex markers — should NOT be stripped.
	input := "# My Config\n\n" + kortexMarkerSection
	result := StripLegacyPersonaBlock(input)
	if result != input {
		t.Fatalf("fingerprint inside marker: expected unchanged result:\ngot:  %q\nwant: %q", result, input)
	}
}

func TestStripLegacyPersonaBlock_LegacyBlockOnlyReturnsEmpty(t *testing.T) {
	// File contains only the legacy persona block with no markers.
	result := StripLegacyPersonaBlock(legacyPersonaBlock)
	if result != "" {
		t.Fatalf("legacy-only: expected empty string:\ngot: %q", result)
	}
}

func TestStripLegacyPersonaBlock_LegacyBlockBeforeMarkersStripped(t *testing.T) {
	// Stale free-text persona block sits before a properly-marked section.
	input := legacyPersonaBlock + "\n" + kortexMarkerSection
	result := StripLegacyPersonaBlock(input)

	// The legacy block should be gone.
	if strings.Contains(result, "## Rules") {
		t.Fatal("stripped result should not contain legacy '## Rules' header")
	}
	// The marked section must survive.
	if !strings.Contains(result, "<!-- kortex:persona -->") {
		t.Fatal("stripped result missing kortex marker section")
	}
}

func TestStripLegacyPersonaBlock_MarkerSectionContentPreserved(t *testing.T) {
	// Markers and their content must be fully preserved after stripping.
	input := legacyPersonaBlock + "\n" + kortexMarkerSection + "\n# User Notes\n\nSome user text.\n"
	result := StripLegacyPersonaBlock(input)

	if !strings.Contains(result, "<!-- kortex:persona -->") {
		t.Fatal("marker open not preserved")
	}
	if !strings.Contains(result, "<!-- /kortex:persona -->") {
		t.Fatal("marker close not preserved")
	}
	if !strings.Contains(result, "# User Notes") {
		t.Fatal("user content after markers not preserved")
	}
}

func TestStripLegacyPersonaBlock_OnlyTwoOfThreeFingerprints(t *testing.T) {
	// File has "## Personality" and "Senior Architect" but NOT "## Rules" —
	// only two of three fingerprints, so it should NOT be stripped.
	input := "## Personality\n\nSenior Architect, 15+ years experience.\n\n" + kortexMarkerSection
	result := StripLegacyPersonaBlock(input)
	// With only 2/3 fingerprints, stripping should NOT occur.
	if result != input {
		t.Fatalf("partial fingerprint: expected unchanged result:\ngot:  %q\nwant: %q", result, input)
	}
}

func TestStripLegacyPersonaBlock_MixedZone_OnlyOneFingerprint_PreMarker(t *testing.T) {
	// Edge case: "## Rules" appears in user content before the first marker,
	// but the other two fingerprints ("## Personality" and "Senior Architect")
	// exist only inside a kortex marker block.
	//
	// Old behaviour (bug): one fingerprint in the pre-marker zone was enough to
	// trigger stripping, destroying the user's "## Rules" section.
	// New behaviour (fixed): ALL fingerprints must appear in the pre-marker zone;
	// since only one does, the file is returned unchanged.
	userRulesSection := "## Rules\n\n- Never do X.\n- Always do Y.\n\n"
	markerWithOtherFingerprints := "<!-- kortex:persona -->\n## Personality\n\nSenior Architect, 15+ years experience.\n<!-- /kortex:persona -->\n"

	input := userRulesSection + markerWithOtherFingerprints
	result := StripLegacyPersonaBlock(input)

	if result != input {
		t.Fatalf(
			"mixed-zone edge case: only one fingerprint in pre-marker zone, expected unchanged result:\ngot:  %q\nwant: %q",
			result, input,
		)
	}
}

func TestStripLegacyPersonaBlock_MixedZone_TwoFingerprints_PreMarker(t *testing.T) {
	// Two of the three fingerprints appear before the first marker, but only the
	// third ("## Rules") exists inside the marker block. Stripping must NOT fire
	// because not all fingerprints are in the pre-marker zone.
	preMarker := "## Personality\n\nSenior Architect, 15+ years experience.\n\n"
	markerWithRule := "<!-- kortex:persona -->\n## Rules\n\n- Rule inside marker.\n<!-- /kortex:persona -->\n"

	input := preMarker + markerWithRule
	result := StripLegacyPersonaBlock(input)

	if result != input {
		t.Fatalf(
			"mixed-zone (2 of 3 in pre-marker): expected unchanged result:\ngot:  %q\nwant: %q",
			result, input,
		)
	}
}

func TestStripLegacyPersonaBlock_AllFingerprintsPreMarker_Strips(t *testing.T) {
	// Positive case: ALL three fingerprints appear before the first marker.
	// Stripping MUST fire, removing the pre-marker legacy block.
	preMarker := "## Rules\n\n- Some rule.\n\n## Personality\n\nSenior Architect, veteran.\n\n"
	markerSection := "<!-- kortex:persona -->\nUpdated persona.\n<!-- /kortex:persona -->\n"

	input := preMarker + markerSection
	result := StripLegacyPersonaBlock(input)

	if result == input {
		t.Fatal("all-fingerprints-pre-marker: expected stripping to occur, but got unchanged result")
	}
	if strings.Contains(result, "## Rules") {
		t.Fatal("all-fingerprints-pre-marker: legacy '## Rules' should have been stripped")
	}
	if !strings.Contains(result, "<!-- kortex:persona -->") {
		t.Fatal("all-fingerprints-pre-marker: marker section must be preserved")
	}
}

func TestStripLegacyPersonaBlock_EmptyFileReturnsSame(t *testing.T) {
	result := StripLegacyPersonaBlock("")
	if result != "" {
		t.Fatalf("empty file: expected empty result, got %q", result)
	}
}

func TestStripLegacyPersonaBlock_UserContentBeforeAndAfterMarkersPreserved(t *testing.T) {
	// User has hand-written notes before the legacy block — these should survive
	// IF they are not part of the legacy block.  Since the legacy detection works
	// by looking for fingerprints before the first marker, user content that
	// predates the legacy block would also be stripped.  This is an accepted
	// tradeoff documented in the function comment.
	input := legacyPersonaBlock + "\n" + kortexMarkerSection + "\n# Custom section\n\nUser stuff.\n"
	result := StripLegacyPersonaBlock(input)

	if !strings.Contains(result, "# Custom section") {
		t.Fatal("content after kortex markers must be preserved")
	}
}

// --- StripLegacyATLBlock tests ---

const legacyATLBlock = `<!-- BEGIN:kortex -->
## Agent Teams Orchestrator

You are a COORDINATOR, not an executor.

### Delegation Rules (ALWAYS ACTIVE)

| Rule | Instruction |
|------|------------|
| No inline work | Reading/writing code → delegate to sub-agent |
<!-- END:kortex -->`

func TestStripLegacyATLBlock_OnlyATLBlock_ReturnsEmpty(t *testing.T) {
	result := StripLegacyATLBlock(legacyATLBlock)
	if result != "" {
		t.Fatalf("only ATL block: expected empty string, got %q", result)
	}
}

func TestStripLegacyATLBlock_ATLBlockThenMarkers_StripsATLKeepsMarkers(t *testing.T) {
	sddSection := "<!-- kortex:sdd-orchestrator -->\nSome orchestrator content.\n<!-- /kortex:sdd-orchestrator -->\n"
	input := legacyATLBlock + "\n\n" + sddSection

	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("ATL open marker should have been stripped")
	}
	if strings.Contains(result, "<!-- END:kortex -->") {
		t.Fatal("ATL close marker should have been stripped")
	}
	if !strings.Contains(result, "<!-- kortex:sdd-orchestrator -->") {
		t.Fatal("sdd-orchestrator marker section must be preserved")
	}
	if !strings.Contains(result, "<!-- /kortex:sdd-orchestrator -->") {
		t.Fatal("sdd-orchestrator close marker must be preserved")
	}
}

func TestStripLegacyATLBlock_ContentBeforeATL_StripsOnlyATL(t *testing.T) {
	before := "# My Config\n\nSome user content.\n"
	sddSection := "<!-- kortex:sdd-orchestrator -->\nOrchestrator stuff.\n<!-- /kortex:sdd-orchestrator -->\n"
	input := before + "\n" + legacyATLBlock + "\n\n" + sddSection

	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("ATL open marker should have been stripped")
	}
	if !strings.Contains(result, "# My Config") {
		t.Fatal("user content before ATL block must be preserved")
	}
	if !strings.Contains(result, "<!-- kortex:sdd-orchestrator -->") {
		t.Fatal("sdd-orchestrator section must be preserved")
	}
}

func TestStripLegacyATLBlock_NoATLBlock_ReturnsUnchanged(t *testing.T) {
	input := "# My Config\n\nSome content without ATL block.\n"
	result := StripLegacyATLBlock(input)
	if result != input {
		t.Fatalf("no ATL block: expected unchanged result:\ngot:  %q\nwant: %q", result, input)
	}
}

func TestStripLegacyATLBlock_OnlyOpenMarkerNoClose_StripsOrphanMarker(t *testing.T) {
	input := "<!-- BEGIN:kortex -->\nSome content without close marker.\n"
	result := StripLegacyATLBlock(input)
	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("orphan BEGIN marker should be stripped by post-loop cleanup")
	}
	if !strings.Contains(result, "Some content without close marker.") {
		t.Fatal("content around orphan BEGIN marker should be preserved")
	}
}

func TestStripLegacyATLBlock_ATLBlockAndSDDOrchestrator_StripsOnlyATL(t *testing.T) {
	sddSection := "<!-- kortex:sdd-orchestrator -->\nYou are a COORDINATOR.\n<!-- /kortex:sdd-orchestrator -->\n"
	input := legacyATLBlock + "\n\n" + sddSection

	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("ATL block should have been stripped")
	}
	if !strings.Contains(result, "<!-- kortex:sdd-orchestrator -->") {
		t.Fatal("sdd-orchestrator section must be preserved after ATL strip")
	}
	if !strings.Contains(result, "You are a COORDINATOR.") {
		t.Fatal("sdd-orchestrator content must be preserved")
	}
}

func TestStripLegacyATLBlock_EmptyFile_ReturnsEmpty(t *testing.T) {
	result := StripLegacyATLBlock("")
	if result != "" {
		t.Fatalf("empty file: expected empty result, got %q", result)
	}
}

func TestStripLegacyATLBlock_Idempotent(t *testing.T) {
	// Calling twice should produce the same result as calling once.
	sddSection := "<!-- kortex:sdd-orchestrator -->\nOrchestrator.\n<!-- /kortex:sdd-orchestrator -->\n"
	input := legacyATLBlock + "\n\n" + sddSection

	once := StripLegacyATLBlock(input)
	twice := StripLegacyATLBlock(once)

	if once != twice {
		t.Fatalf("idempotent: second call changed result:\nfirst:  %q\nsecond: %q", once, twice)
	}
}

func TestStripLegacyATLBlock_EmptyBetweenMarkers(t *testing.T) {
	// An ATL block with nothing between the markers should strip to empty.
	input := "<!-- BEGIN:kortex -->\n<!-- END:kortex -->"
	result := StripLegacyATLBlock(input)
	if result != "" {
		t.Fatalf("empty between markers: expected empty string, got %q", result)
	}
}

func TestStripLegacyATLBlock_DuplicateBlocks(t *testing.T) {
	// A file with two ATL blocks (e.g. pasted twice) — both must be stripped.
	block := "<!-- BEGIN:kortex -->\nsome content\n<!-- END:kortex -->"
	input := block + "\n\n" + block
	result := StripLegacyATLBlock(input)
	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("duplicate blocks: first ATL open marker should have been stripped")
	}
	if strings.Contains(result, "<!-- END:kortex -->") {
		t.Fatal("duplicate blocks: ATL close marker should have been stripped")
	}
	if result != "" {
		t.Fatalf("duplicate blocks: expected empty string after stripping both, got %q", result)
	}
}

func TestStripLegacyATLBlock_EndBeforeBeginWithValidPairAfter(t *testing.T) {
	// A stray END marker appears before a valid BEGIN...END pair.
	// The valid block must still be stripped.
	strayEnd := "<!-- END:kortex -->\n"
	validBlock := "<!-- BEGIN:kortex -->\nreal content\n<!-- END:kortex -->"
	after := "\n\nsome other content"
	input := strayEnd + validBlock + after

	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("end-before-begin: valid ATL open marker should have been stripped")
	}
	if strings.Contains(result, "real content") {
		t.Fatal("end-before-begin: valid ATL block content should have been stripped")
	}
	if !strings.Contains(result, "some other content") {
		t.Fatal("end-before-begin: content after valid ATL block must be preserved")
	}
	if strings.Contains(result, "<!-- END:kortex -->") {
		t.Fatal("end-before-begin: orphan END marker should have been removed from output")
	}
}

func TestStripLegacyATLBlock_CRLFLineEndings(t *testing.T) {
	// CRLF line endings should be trimmed cleanly without stray \r characters.
	input := "before\r\n\r\n<!-- BEGIN:kortex -->\r\ncontent\r\n<!-- END:kortex -->\r\n\r\nafter\r\n"
	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("ATL block should be stripped")
	}
	if !strings.Contains(result, "before") {
		t.Fatal("content before block must be preserved")
	}
	if !strings.Contains(result, "after") {
		t.Fatal("content after block must be preserved")
	}
	// No stray \r should remain at the join point
	if strings.Contains(result, "\r\n\r\n\n") || strings.Contains(result, "\n\r\n\r") {
		t.Fatalf("CRLF: stray carriage returns at join point:\n%q", result)
	}
}

func TestStripLegacyPersonaBlock_CRLFLineEndings(t *testing.T) {
	// CRLF line endings in legacy block + markers should be handled cleanly.
	legacy := "## Rules\r\n\r\n- Some rule.\r\n\r\n## Personality\r\n\r\nSenior Architect, veteran.\r\n\r\n"
	marker := "<!-- kortex:persona -->\r\nUpdated persona.\r\n<!-- /kortex:persona -->\r\n"
	input := legacy + marker

	result := StripLegacyPersonaBlock(input)

	if strings.Contains(result, "## Rules") {
		t.Fatal("legacy block should be stripped")
	}
	if !strings.Contains(result, "<!-- kortex:persona -->") {
		t.Fatal("marker section must be preserved")
	}
	// The marker section should not have leading \r artifacts
	if strings.HasPrefix(result, "\r") {
		t.Fatal("result should not start with stray \\r")
	}
}

func TestStripLegacyATLBlock_InlineMarkerNotStripped(t *testing.T) {
	// ATL markers appearing inline (not at the start of a line) should NOT be stripped.
	input := "See <!-- BEGIN:kortex --> for reference.\nAnd <!-- END:kortex --> too.\n"
	result := StripLegacyATLBlock(input)
	if result != input {
		t.Fatalf("inline markers should not be stripped:\ngot:  %q\nwant: %q", result, input)
	}
}

func TestStripLegacyATLBlock_OrphanMarkersCRLF(t *testing.T) {
	// Orphan END marker with CRLF line endings — must be stripped without leaving stray \r.
	input := "before\r\n<!-- END:kortex -->\r\nafter"
	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- END:kortex -->") {
		t.Fatal("orphan END marker should be stripped")
	}
	if !strings.Contains(result, "before") {
		t.Fatal("content before orphan must be preserved")
	}
	if !strings.Contains(result, "after") {
		t.Fatal("content after orphan must be preserved")
	}
	// No stray \r between "before" and "after" — the marker line should be cleanly removed
	if strings.Contains(result, "\r\n\r\r") || strings.Contains(result, "\r\r") {
		t.Fatalf("orphan CRLF: stray \\r in output:\n%q", result)
	}
}

func TestStripLegacyATLBlock_OrphanBeginCRLF(t *testing.T) {
	// Orphan BEGIN marker with CRLF — must be stripped without stray \r.
	input := "before\r\n<!-- BEGIN:kortex -->\r\nsome content\r\n"
	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("orphan BEGIN marker should be stripped")
	}
	if !strings.Contains(result, "before") {
		t.Fatal("content before orphan must be preserved")
	}
	if !strings.Contains(result, "some content") {
		t.Fatal("content after orphan BEGIN must be preserved")
	}
}

func TestStripLegacyATLBlock_MultiBlocksWithContentBetween(t *testing.T) {
	// Two ATL blocks with user content between them — both blocks stripped,
	// user content preserved.
	block := "<!-- BEGIN:kortex -->\nATL stuff\n<!-- END:kortex -->"
	input := block + "\n\nuser content here\n\n" + block
	result := StripLegacyATLBlock(input)

	if strings.Contains(result, "<!-- BEGIN:kortex -->") {
		t.Fatal("both ATL blocks should be stripped")
	}
	if strings.Contains(result, "ATL stuff") {
		t.Fatal("ATL content should be stripped")
	}
	if !strings.Contains(result, "user content here") {
		t.Fatal("user content between blocks must be preserved")
	}
}
