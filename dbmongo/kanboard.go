package main

import (
	"fmt"

	"github.com/chrnin/ganboard"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

func getOrCreateProject(c *gin.Context) (ganboard.Project, error) {
	kb := c.Keys["KBCLIENT"].(*ganboard.Client)

	project, err := kb.GetProjectByID(1)

	if project.Identifier == "SIGNAUXFAIBLES" {
		return project, nil
	}
	return project, err
	// test, err := kb.GetAllProjects()
	// spew.Dump(test)

	// params := ganboard.ProjectParams{
	// 	Name:        "Signaux-Faibles™",
	// 	Identifier:  "signauxfaibles",
	// 	Description: "Test project for SignauxFaibles",
	// 	Email:       "christophe.ninucci@direccte.gouv.fr",
	// }
	// idProject, err := kb.CreateProject(params)
	// fmt.Println(idProject)
	// if err != nil {
	// 	return ganboard.Project{}, err
	// }
	// project, err = kb.GetProjectByIdentifier("SIGNAUXFAIBLES")
	// return project, err

}
func getKBProject(c *gin.Context) {
	kb := c.Keys["KBCLIENT"].(*ganboard.Client)

	params := ganboard.ProjectParams{
		Name:        "Signaux-Faibles™",
		Identifier:  "signauxfaibles",
		Description: "Test project for SignauxFaibles",
		Email:       "christophe.ninucci@direccte.gouv.fr",
	}
	idProject, err := kb.CreateProject(params)
	fmt.Println(idProject)
	fmt.Println(err)
	project, err := kb.GetAllProjects()
	spew.Dump(err)
	spew.Dump(project)

}
