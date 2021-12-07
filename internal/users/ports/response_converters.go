package ports

import "dev.azure.com/filimonovga/our-expenses/our-expenses-server/internal/users/domain"

func userToResponse(domainObj *domain.User) AuthenticationData {
	report := AuthenticationData{
		Id:           domainObj.ID(),
		Username:     domainObj.Username(),
		Token:        domainObj.Token(),
		RefreshToken: domainObj.RefreshToken(),
	}
	return report
}
