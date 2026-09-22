package client

import (
	"strconv"
	"strings"

	"github.com/ao-data/albiondata-client/lib"
	"github.com/ao-data/albiondata-client/log"
	uuid "github.com/nu7hatch/gouuid"
)

// parseFame extracts the numeric fame value safely from strings that might be enclosed
// in nested brackets (e.g. "[12345]", "[[12345]]") or plain numbers.
func parseFame(raw string) int {
	cleaned := strings.Trim(raw, "[] \t\r\n\"'")
	if cleaned == "" {
		return 0
	}
	val, err := strconv.Atoi(cleaned)
	if err != nil {
		log.Debugf("Could not parse fame value %q: %v", raw, err)
		return 0
	}
	return val
}

type eventSkillData struct {
	SkillIds    []int     `mapstructure:"1"`
	Levels      []int     `mapstructure:"2"`
	Percentages []float64 `mapstructure:"3"`
	Fame        []string  `mapstructure:"4"`
}

func (event eventSkillData) Process(state *albionState) {
	log.Debug("Got skill data event...")

	skills := []*lib.Skill{}

	for k := range event.SkillIds {
		skill := &lib.Skill{}
		skill.ID = event.SkillIds[k]
		skill.Level = event.Levels[k]
		skill.PercentNextLevel = event.Percentages[k]
		fame := parseFame(event.Fame[k])
		skill.Fame = fame

		skills = append(skills, skill)
	}

	if len(skills) < 1 {
		return
	}

	upload := lib.SkillsUpload{
		Skills: skills,
	}

	identifier, _ := uuid.NewV4()
	log.Infof("Sending %d skills of %v to ingest", len(skills), state.CharacterName)
	for _, sk := range skills {
		if sk.Level > 0 || sk.ID == 683 || sk.ID == 24 || sk.ID == 680 || sk.ID == 681 || sk.ID == 682 {
			log.Infof("[Skill Debug] ID: %d | Level: %d | Fame: %d", sk.ID, sk.Level, sk.Fame)
		}
	}
	sendMsgToPrivateUploaders(&upload, lib.NatsSkillData, state, identifier.String(), len(skills))
}
