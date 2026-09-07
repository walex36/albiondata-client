package client

import (
	"strconv"

	"github.com/ao-data/albiondata-client/lib"
	"github.com/ao-data/albiondata-client/log"
	uuid "github.com/nu7hatch/gouuid"
)

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
		// for some reason, the value is enclosed in [[]]. trying to get rid of them
		fame, err := strconv.Atoi(event.Fame[k][2 : len(event.Fame[k])-2])
		if err != nil {
			log.Error("Could not parse fame value. ", err)
			continue
		}
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
	sendMsgToPrivateUploaders(&upload, lib.NatsSkillData, state, identifier.String())
}
