package ports

import "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"

func authDataToResponse(domainObj *domain.AuthenticationData) AuthenticationData {
	report := AuthenticationData{
		Username: domainObj.User(),
	}
	return report
}

func userToResponse(domainObj *domain.User) AuthenticationData {
	report := AuthenticationData{
		Id:           domainObj.ID(),
		Username:     domainObj.Username(),
		Token:        domainObj.Token(),
		RefreshToken: domainObj.RefreshToken(),
	}
	return report
}
