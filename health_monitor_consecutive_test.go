package llmprovider

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHealthMonitor_RecordSuccess_RequiresConsecutiveSuccessesToRecover is the
// reproduce-first RED test (§11.4.115) for the consecutive-success-threshold
// defect in the root-package copy of the health monitor. The HealthyThreshold
// field is documented "Consecutive successes to mark healthy" and RecordFailure
// uses ConsecutiveFails — but recovery was gated on the cumulative,
// monotonically-increasing SuccessCount, so an Unhealthy provider with prior
// lifetime successes flipped back to Healthy after a SINGLE success.
func TestHealthMonitor_RecordSuccess_RequiresConsecutiveSuccessesToRecover(t *testing.T) {
	config := HealthMonitorConfig{HealthyThreshold: 2, UnhealthyThreshold: 2, Enabled: false}
	hm := NewHealthMonitor(config)
	hm.RegisterProvider("test", &mockProvider{})

	// Drive Healthy (2 consecutive successes).
	hm.RecordSuccess("test")
	hm.RecordSuccess("test")
	health, _ := hm.GetHealth("test")
	assert.Equal(t, HealthStatusHealthy, health.Status, "precondition: Healthy")

	// Drive Unhealthy (2 consecutive failures).
	hm.RecordFailure("test", errors.New("err 1"))
	hm.RecordFailure("test", errors.New("err 2"))
	health, _ = hm.GetHealth("test")
	assert.Equal(t, HealthStatusUnhealthy, health.Status, "precondition: Unhealthy")

	// ONE success must NOT recover when HealthyThreshold == 2.
	hm.RecordSuccess("test")
	health, _ = hm.GetHealth("test")
	assert.NotEqual(t, HealthStatusHealthy, health.Status,
		"single success must NOT flip an Unhealthy provider to Healthy when HealthyThreshold=2")

	// A SECOND consecutive success completes the threshold.
	hm.RecordSuccess("test")
	health, _ = hm.GetHealth("test")
	assert.Equal(t, HealthStatusHealthy, health.Status,
		"two consecutive successes must recover the provider")
}
