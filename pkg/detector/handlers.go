package detector

import (
	"arbitrage-detector/pkg/models"
)

// RegisterOpportunityHandler registers a handler function that will be called
// whenever a new arbitrage opportunity is detected
func (a *ArbitrageDetector) RegisterOpportunityHandler(handler OpportunityHandler) {
	a.opportunityHandlers = append(a.opportunityHandlers, handler)
}

// notifyOpportunityHandlers calls all registered opportunity handlers
func (a *ArbitrageDetector) notifyOpportunityHandlers(opp models.ArbitrageOpportunity) {
	for _, handler := range a.opportunityHandlers {
		handler(opp)
	}
}