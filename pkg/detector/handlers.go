package detector

import (
	"apex-arbitrage/pkg/models"
)

// RegisterOpportunityHandler registers a handler function that will be called
// whenever a new arbitrage opportunity is detected
func (a *APEX) RegisterOpportunityHandler(handler OpportunityHandler) {
	a.opportunityHandlers = append(a.opportunityHandlers, handler)
}

// notifyOpportunityHandlers calls all registered opportunity handlers
func (a *APEX) notifyOpportunityHandlers(opp models.ArbitrageOpportunity) {
	for _, handler := range a.opportunityHandlers {
		handler(opp)
	}
}
