package health

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestHealthMonitor_RecordSuccess_RequiresConsecutiveSuccessesToRecover is a
// reproduce-first RED test (§11.4.115) for the consecutive-success-threshold
// defect. The documented contract (CLAUDE.md "Consecutive failure/success
// thresholds") and the symmetry with RecordFailure (which uses ConsecutiveFails)
// require that an Unhealthy provider recover to Healthy only after
// HealthyThreshold *consecutive* successes. The defect: recovery is gated on the
// cumulative, monotonically-increasing SuccessCount, so an Unhealthy provider
// with prior lifetime successes flips back to Healthy after a SINGLE success,
// defeating the flap-suppression purpose of the threshold.
func TestHealthMonitor_RecordSuccess_RequiresConsecutiveSuccessesToRecover(t *testing.T) {
	config := HealthMonitorConfig{HealthyThreshold: 2, UnhealthyThreshold: 2, Enabled: false}
	hm := NewHealthMonitor(config)
	hm.RegisterProvider("test", &mockProvider{})

	// Drive provider Healthy (2 consecutive successes).
	hm.RecordSuccess("test")
	hm.RecordSuccess("test")
	health, _ := hm.GetHealth("test")
	assert.Equal(t, HealthStatusHealthy, health.Status, "precondition: provider should be Healthy")

	// Drive provider Unhealthy (2 consecutive failures).
	hm.RecordFailure("test", errors.New("err 1"))
	hm.RecordFailure("test", errors.New("err 2"))
	health, _ = hm.GetHealth("test")
	assert.Equal(t, HealthStatusUnhealthy, health.Status, "precondition: provider should be Unhealthy")

	// ONE success must NOT be enough to recover when HealthyThreshold == 2.
	hm.RecordSuccess("test")
	health, _ = hm.GetHealth("test")
	assert.NotEqual(t, HealthStatusHealthy, health.Status,
		"single success must NOT flip an Unhealthy provider to Healthy when HealthyThreshold=2")

	// A SECOND consecutive success completes the threshold and recovers.
	hm.RecordSuccess("test")
	health, _ = hm.GetHealth("test")
	assert.Equal(t, HealthStatusHealthy, health.Status,
		"two consecutive successes must recover the provider")
}

// TestHealthMonitor_CheckProvider_RequiresConsecutiveSuccessesToRecover is the
// monitor-loop analogue, exercising the same defect through checkProvider's
// success path (the code path the live monitor actually drives).
func TestHealthMonitor_CheckProvider_RequiresConsecutiveSuccessesToRecover(t *testing.T) {
	config := HealthMonitorConfig{HealthyThreshold: 2, UnhealthyThreshold: 2, Timeout: 5 * time.Second, Enabled: false}
	hm := NewHealthMonitor(config)
	hm.ctx = context.Background()
	mp := &mockProvider{}
	hm.RegisterProvider("test", mp)

	// Bootstrap to Healthy via two successful checks.
	hm.checkProvider("test", mp)
	hm.checkProvider("test", mp)
	health, _ := hm.GetHealth("test")
	assert.Equal(t, HealthStatusHealthy, health.Status, "precondition: Healthy")

	// Two failing checks -> Unhealthy.
	mp.SetHealthError(errors.New("down"))
	hm.checkProvider("test", mp)
	hm.checkProvider("test", mp)
	health, _ = hm.GetHealth("test")
	assert.Equal(t, HealthStatusUnhealthy, health.Status, "precondition: Unhealthy")

	// One recovering check must NOT flip to Healthy with HealthyThreshold=2.
	mp.SetHealthError(nil)
	hm.checkProvider("test", mp)
	health, _ = hm.GetHealth("test")
	assert.NotEqual(t, HealthStatusHealthy, health.Status,
		"single recovering check must NOT flip Unhealthy->Healthy when HealthyThreshold=2")
}
