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
}

// qeFrontSkills covers E2E frontend testing: Playwright BDD + CLI.
var qeFrontSkills = []model.SkillID{
	model.SkillPlaywrightBDD,
	model.SkillPlaywrightCLI,
}

// qePerfSkills covers performance testing with k6.
var qePerfSkills = []model.SkillID{
	model.SkillK6LoadTest,
}

// qeAPISkills covers API and contract testing with Karate DSL.
var qeAPISkills = []model.SkillID{
	model.SkillKarateDSL,
}

// SkillsForPreset returns which skills should be installed for a given preset.
//
//   - "minimal" / PresetMinimal:       SDD skills only
//   - "ecosystem-only" / PresetEcosystemOnly: SDD + common framework skills
//   - "full-gentleman" / PresetFullGentleman: all available skills
//   - "custom" / PresetCustom:         empty (caller should provide explicit list)
//   - "qe-front" / PresetQEFront:      SDD + Playwright BDD + Playwright CLI
//   - "qe-perf"  / PresetQEPerf:       SDD + k6 load testing
//   - "qe-api"   / PresetQEAPI:        SDD + Karate DSL
//   - "qe-sdet"  / PresetQESDET:       SDD + all QE skills
func SkillsForPreset(preset model.PresetID) []model.SkillID {
	switch preset {
	case model.PresetMinimal:
		return copySkills(sddSkills)
	case model.PresetEcosystemOnly:
		return copySkills(append(sddSkills, foundationSkills...))
	case model.PresetFullGentleman:
		all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills))
		all = append(all, sddSkills...)
		all = append(all, foundationSkills...)
		return all
	case model.PresetQEFront:
		all := make([]model.SkillID, 0, len(sddSkills)+len(qeFrontSkills))
		all = append(all, sddSkills...)
		all = append(all, qeFrontSkills...)
		return all
	case model.PresetQEPerf:
		all := make([]model.SkillID, 0, len(sddSkills)+len(qePerfSkills))
		all = append(all, sddSkills...)
		all = append(all, qePerfSkills...)
		return all
	case model.PresetQEAPI:
		all := make([]model.SkillID, 0, len(sddSkills)+len(qeAPISkills))
		all = append(all, sddSkills...)
		all = append(all, qeAPISkills...)
		return all
	case model.PresetQESDET:
		all := make([]model.SkillID, 0, len(sddSkills)+len(qeFrontSkills)+len(qePerfSkills)+len(qeAPISkills))
		all = append(all, sddSkills...)
		all = append(all, qeFrontSkills...)
		all = append(all, qePerfSkills...)
		all = append(all, qeAPISkills...)
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

// AllSkillIDs returns every known skill ID.
func AllSkillIDs() []model.SkillID {
	all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills)+len(qeFrontSkills)+len(qePerfSkills)+len(qeAPISkills))
	all = append(all, sddSkills...)
	all = append(all, foundationSkills...)
	all = append(all, qeFrontSkills...)
	all = append(all, qePerfSkills...)
	all = append(all, qeAPISkills...)
	return all
}

func copySkills(src []model.SkillID) []model.SkillID {
	dst := make([]model.SkillID, len(src))
	copy(dst, src)
	return dst
}
