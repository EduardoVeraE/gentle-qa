package catalog

import "github.com/EduardoVeraE/Gentle-QA/internal/model"

type Skill struct {
	ID       model.SkillID
	Name     string
	Category string
	Priority string
}

var mvpSkills = []Skill{
	// SDD skills
	{ID: model.SkillSDDInit, Name: "sdd-init", Category: "sdd", Priority: "p0"},

	{ID: model.SkillSDDApply, Name: "sdd-apply", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDVerify, Name: "sdd-verify", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDExplore, Name: "sdd-explore", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDPropose, Name: "sdd-propose", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDSpec, Name: "sdd-spec", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDDesign, Name: "sdd-design", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDTasks, Name: "sdd-tasks", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDArchive, Name: "sdd-archive", Category: "sdd", Priority: "p0"},
	{ID: model.SkillSDDOnboard, Name: "sdd-onboard", Category: "sdd", Priority: "p0"},
	// Foundation skills
	{ID: model.SkillGoTesting, Name: "go-testing", Category: "testing", Priority: "p0"},
	{ID: model.SkillCreator, Name: "skill-creator", Category: "workflow", Priority: "p0"},
	{ID: model.SkillJudgmentDay, Name: "judgment-day", Category: "workflow", Priority: "p0"},
	{ID: model.SkillBranchPR, Name: "branch-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillIssueCreation, Name: "issue-creation", Category: "workflow", Priority: "p0"},
	{ID: model.SkillSkillRegistry, Name: "skill-registry", Category: "workflow", Priority: "p0"},
	// Sustainable review skills (upstream)
	{ID: model.SkillChainedPR, Name: "chained-pr", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCognitiveDoc, Name: "cognitive-doc-design", Category: "workflow", Priority: "p0"},
	{ID: model.SkillCommentWriter, Name: "comment-writer", Category: "workflow", Priority: "p0"},
	{ID: model.SkillWorkUnitCommits, Name: "work-unit-commits", Category: "workflow", Priority: "p0"},
	// QE / SDET skills
	{ID: model.SkillPlaywrightBDD, Name: "playwright-bdd", Category: "qe-e2e", Priority: "p0"},
	{ID: model.SkillPlaywrightCLI, Name: "playwright-cli", Category: "qe-e2e", Priority: "p0"},
	{ID: model.SkillK6LoadTest, Name: "k6-load-test", Category: "qe-performance", Priority: "p0"},
	{ID: model.SkillKarateDSL, Name: "karate-dsl", Category: "qe-api", Priority: "p0"},
	{ID: model.SkillQAOWASPSecurity, Name: "qa-owasp-security", Category: "qe-security", Priority: "p0"},
	{ID: model.SkillQAMobileTesting, Name: "qa-mobile-testing", Category: "qe-mobile", Priority: "p0"},
	{ID: model.SkillQAVisualRegression, Name: "qa-visual-regression", Category: "qe-visual", Priority: "p0"},
	{ID: model.SkillQAContractPact, Name: "qa-contract-pact", Category: "qe-contract", Priority: "p0"},
}

func MVPSkills() []Skill {
	skills := make([]Skill, len(mvpSkills))
	copy(skills, mvpSkills)
	return skills
}
