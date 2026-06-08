package projecthttp

import (
	projectusecases "macabi-back/internal/project/application/usecases"
)

type ProjectHandler struct {
	createUC       *projectusecases.CreateProject
	listUC         *projectusecases.ListProjects
	getUC          *projectusecases.GetProject
	updateUC       *projectusecases.UpdateProject
	deleteUC       *projectusecases.DeleteProject
	addMemberUC    *projectusecases.AddProjectMember
	removeMemberUC *projectusecases.RemoveProjectMember
	listMembersUC  *projectusecases.ListProjectMembers
}

func NewProjectHandler(
	createUC *projectusecases.CreateProject,
	listUC *projectusecases.ListProjects,
	getUC *projectusecases.GetProject,
	updateUC *projectusecases.UpdateProject,
	deleteUC *projectusecases.DeleteProject,
	addMemberUC *projectusecases.AddProjectMember,
	removeMemberUC *projectusecases.RemoveProjectMember,
	listMembersUC *projectusecases.ListProjectMembers,
) *ProjectHandler {
	return &ProjectHandler{
		createUC:       createUC,
		listUC:         listUC,
		getUC:          getUC,
		updateUC:       updateUC,
		deleteUC:       deleteUC,
		addMemberUC:    addMemberUC,
		removeMemberUC: removeMemberUC,
		listMembersUC:  listMembersUC,
	}
}
