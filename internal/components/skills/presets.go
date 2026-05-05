package skills

import "github.com/EduardoVeraE/Gentle-QA/internal/model"

// sddSkills are the SDD orchestrator skills — always included.
var sddSkills = []model.SkillID{
	model.SkillSDDInit,
	model.SkillSDDExplore,
	model.SkillSDDPropose,
	model.SkillSDDSpec,
	model.SkillSDDDesign,
	model.SkillSDDTasks,
	model.SkillSDDApply,
	model.SkillSDDVerify,
	model.SkillSDDArchive,
	model.SkillSDDOnboard,
	model.SkillJudgmentDay,
}

// foundationSkills are baseline learning skills for the "recommended" tier.
var foundationSkills = []model.SkillID{
	model.SkillGoTesting,
	model.SkillCreator,
	model.SkillBranchPR,
	model.SkillIssueCreation,
	model.SkillSkillRegistry,
	model.SkillChainedPR,
	model.SkillCognitiveDoc,
	model.SkillCommentWriter,
	model.SkillWorkUnitCommits,
}

// qeFrontSkills covers E2E frontend testing on the Playwright stack: BDD,
// CLI, full E2E suite, live MCP debugging, and regression strategy.
var qeFrontSkills = []model.SkillID{
	model.SkillPlaywrightBDD,
	model.SkillPlaywrightCLI,
	model.SkillPlaywrightE2E,
	model.SkillPlaywrightMCPInspect,
	model.SkillPlaywrightRegression,
}

// qePerfSkills covers performance testing with k6.
var qePerfSkills = []model.SkillID{
	model.SkillK6LoadTest,
}

// qeAPISkills covers API and contract testing: Karate DSL plus the generic
// Playwright/REST-Assured api-testing skill.
var qeAPISkills = []model.SkillID{
	model.SkillKarateDSL,
	model.SkillAPITesting,
}

// qaSecuritySkills covers OWASP security testing.
var qaSecuritySkills = []model.SkillID{
	model.SkillQAOWASPSecurity,
}

// qaMobileSkills covers mobile testing (Appium + Detox + native).
var qaMobileSkills = []model.SkillID{
	model.SkillQAMobileTesting,
}

// qaVisualSkills covers visual regression testing.
var qaVisualSkills = []model.SkillID{
	model.SkillQAVisualRegression,
}

// qaContractSkills covers consumer-driven contract testing with Pact.
var qaContractSkills = []model.SkillID{
	model.SkillQAContractPact,
}

// qaManualSkills covers the ISTQB Foundation Level toolkit — included in ALL
// qe-* presets because the manual/test-design fundamentals are stack-agnostic.
var qaManualSkills = []model.SkillID{
	model.SkillQAManualISTQB,
}

// qaA11yPlaywrightSkills covers accessibility testing on the Playwright stack.
var qaA11yPlaywrightSkills = []model.SkillID{
	model.SkillA11yPlaywright,
}

// qaA11ySeleniumSkills covers accessibility testing on the Selenium stack.
// qe-sdet only — qe-front stays Playwright-pure.
var qaA11ySeleniumSkills = []model.SkillID{
	model.SkillA11ySelenium,
}

// qaSeleniumSkills covers Selenium WebDriver E2E testing. qe-sdet only.
var qaSeleniumSkills = []model.SkillID{
	model.SkillSeleniumE2E,
}

// SkillsForPreset returns which skills should be installed for a given preset.
//
//   - "minimal" / PresetMinimal:       SDD skills only
//   - "ecosystem-only" / PresetEcosystemOnly: SDD + common framework skills
//   - "full-gentleman" / PresetFullGentleman: all user-facing skills
//   - "custom" / PresetCustom:         empty (caller should provide explicit list)
//   - "qe-front" / PresetQEFront:      SDD + Playwright stack (BDD/CLI/E2E/MCP/regression)
//     + ISTQB foundation + Playwright a11y
//   - "qe-perf"  / PresetQEPerf:       SDD + k6 + ISTQB foundation
//   - "qe-api"   / PresetQEAPI:        SDD + Karate DSL + api-testing + ISTQB foundation
//   - "qe-sdet"  / PresetQESDET:       SDD + every QE/QA skill (front + perf + api +
//     security + mobile + visual + contract + manual + a11y Playwright + a11y Selenium
//     + Selenium E2E)
func SkillsForPreset(preset model.PresetID) []model.SkillID {
	switch preset {
	case model.PresetMinimal:
		return copySkills(sddSkills)
	case model.PresetEcosystemOnly:
		return copySkills(append(sddSkills, foundationSkills...))
	case model.PresetFullGentleman:
		return AllSkillIDs()
	case model.PresetQEFront:
		all := make([]model.SkillID, 0,
			len(sddSkills)+len(qeFrontSkills)+len(qaManualSkills)+len(qaA11yPlaywrightSkills))
		all = append(all, sddSkills...)
		all = append(all, qeFrontSkills...)
		all = append(all, qaManualSkills...)
		all = append(all, qaA11yPlaywrightSkills...)
		return all
	case model.PresetQEPerf:
		all := make([]model.SkillID, 0, len(sddSkills)+len(qePerfSkills)+len(qaManualSkills))
		all = append(all, sddSkills...)
		all = append(all, qePerfSkills...)
		all = append(all, qaManualSkills...)
		return all
	case model.PresetQEAPI:
		all := make([]model.SkillID, 0, len(sddSkills)+len(qeAPISkills)+len(qaManualSkills))
		all = append(all, sddSkills...)
		all = append(all, qeAPISkills...)
		all = append(all, qaManualSkills...)
		return all
	case model.PresetQESDET:
		all := make([]model.SkillID, 0,
			len(sddSkills)+len(qeFrontSkills)+len(qePerfSkills)+len(qeAPISkills)+
				len(qaSecuritySkills)+len(qaMobileSkills)+len(qaVisualSkills)+len(qaContractSkills)+
				len(qaManualSkills)+len(qaA11yPlaywrightSkills)+len(qaA11ySeleniumSkills)+len(qaSeleniumSkills))
		all = append(all, sddSkills...)
		all = append(all, qeFrontSkills...)
		all = append(all, qePerfSkills...)
		all = append(all, qeAPISkills...)
		all = append(all, qaSecuritySkills...)
		all = append(all, qaMobileSkills...)
		all = append(all, qaVisualSkills...)
		all = append(all, qaContractSkills...)
		all = append(all, qaManualSkills...)
		all = append(all, qaA11yPlaywrightSkills...)
		all = append(all, qaA11ySeleniumSkills...)
		all = append(all, qaSeleniumSkills...)
		return all
	case model.PresetCustom:
		return nil
	default:
		// Unknown preset — default to full.
		all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills))
		all = append(all, sddSkills...)
		all = append(all, foundationSkills...)
		return all
	}
}

// AllSkillIDs returns every user-facing skill ID. Maintainer-only skills
// (e.g. SkillUpstreamSync) are deliberately excluded — they ship in the binary
// but must be requested explicitly via `--skill <id>`.
func AllSkillIDs() []model.SkillID {
	all := make([]model.SkillID, 0,
		len(sddSkills)+len(foundationSkills)+len(qeFrontSkills)+len(qePerfSkills)+len(qeAPISkills)+
			len(qaSecuritySkills)+len(qaMobileSkills)+len(qaVisualSkills)+len(qaContractSkills)+
			len(qaManualSkills)+len(qaA11yPlaywrightSkills)+len(qaA11ySeleniumSkills)+len(qaSeleniumSkills))
	all = append(all, sddSkills...)
	all = append(all, foundationSkills...)
	all = append(all, qeFrontSkills...)
	all = append(all, qePerfSkills...)
	all = append(all, qeAPISkills...)
	all = append(all, qaSecuritySkills...)
	all = append(all, qaMobileSkills...)
	all = append(all, qaVisualSkills...)
	all = append(all, qaContractSkills...)
	all = append(all, qaManualSkills...)
	all = append(all, qaA11yPlaywrightSkills...)
	all = append(all, qaA11ySeleniumSkills...)
	all = append(all, qaSeleniumSkills...)
	return all
}

func copySkills(src []model.SkillID) []model.SkillID {
	dst := make([]model.SkillID, len(src))
	copy(dst, src)
	return dst
}
