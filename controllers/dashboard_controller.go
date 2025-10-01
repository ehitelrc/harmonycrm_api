package controllers

import (
	"harmony_api/repository"
	"harmony_api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// Get campaign summary by company
func (dc *DashboardController) GetCampaignSummary(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))

	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid company id", nil, err)
		return
	}

	data, err := repository.NewDashboardCampaignRepository().GetCampaignSummaryByCompany(uint(companyID))

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error getting campaign summary", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Resument de campaña", data, nil)

}

// Get general dashboard by company
func (dc *DashboardController) GetGeneralDashboardByCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("company_id"))

	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid company id", nil, err)
		return
	}

	data, err := repository.NewDashboardCampaignRepository().GetGeneralDashboardByCompany(uint(companyID))

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error getting general dashboard", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "General dashboard", data, nil)

}

// Get campaign funnel
func (dc *DashboardController) GetCampaignFunnelSummary(c *gin.Context) {
	campaignID, err := strconv.Atoi(c.Param("campaign_id"))

	if err != nil {
		utils.Respond(c, http.StatusBadRequest, false, "Invalid campaign id", nil, err)
		return
	}

	data, err := repository.NewDashboardCampaignRepository().GetCampaignFunnelSummary(campaignID)

	if err != nil {
		utils.Respond(c, http.StatusInternalServerError, false, "Error getting campaign funnel summary", nil, err)
		return
	}

	utils.Respond(c, http.StatusOK, true, "Campaign funnel summary", data, nil)

}
