package main

import (
	"github.com/chrnin/ganboard"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func getOrCreateProject(c *gin.Context) (ganboard.Project, error) {
	kb := c.Keys["KBCLIENT"].(*ganboard.Client)

	project, err := kb.GetProjectByIdentifier("signauxfaibles")

	if project.Identifier == "signauxfaibles" {
		return project, nil
	}

	test, err := kb.GetAllProjects()
	spew.Dump(test)

	params := ganboard.ProjectParams{
		Name:        "Signaux-Faibles™",
		Identifier:  "signauxfaibles",
		Description: "Test project for SignauxFaibles",
		Email:       "christophe.ninucci@direccte.gouv.fr",
	}
	idProject, err := kb.CreateProject(params)

	if err != nil {
		return ganboard.Project{}, err
	}
	project, err = kb.GetProjectByID(idProject)
	return project, err

}
func createTaskFromSiret(c *gin.Context) {

	project, err := getOrCreateProject(c)
	spew.Dump(err)
	spew.Dump(project)

}
