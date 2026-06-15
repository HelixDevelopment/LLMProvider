package provider

import (
	"context"
	"errors"
	"testing"

	"digital.vasic.llmprovider/pkg/models"
)

// stubProvider is a unit-test test double implementing LLMProvider.
//
// Mocks/stubs are permitted in unit-test scope per CONST-050(A) / §11.4.27.
// Each field lets a table-driven case control the value returned by the
// corresponding interface method so the contract can be exercised
// deterministically without any network or credential dependency.
type stubProvider struct {
	completeResp *models.LLMResponse
	completeErr  error

	streamResp *models.LLMResponse
	streamErr  error

	healthErr error

	caps *models.ProviderCapabilities

	validateOK   bool
	validateMsgs []string

	// call-tracking proves the interface dispatch actually reached the impl.
	completeCalled       bool
	completeStreamCalled bool
	healthCalled         bool
	capsCalled           bool
	validateCalled       bool
	lastValidateConfig   map[string]interface{}
}

func (s *stubProvider) Complete(_ context.Context, _ *models.LLMRequest) (*models.LLMResponse, error) {
	s.completeCalled = true
	return s.completeResp, s.completeErr
}

func (s *stubProvider) CompleteStream(_ context.Context, _ *models.LLMRequest) (<-chan *models.LLMResponse, error) {
	s.completeStreamCalled = true
	if s.streamErr != nil {
		return nil, s.streamErr
	}
	ch := make(chan *models.LLMResponse, 1)
	if s.streamResp != nil {
		ch <- s.streamResp
	}
	close(ch)
	return ch, nil
}

func (s *stubProvider) HealthCheck() error {
	s.healthCalled = true
	return s.healthErr
}

func (s *stubProvider) GetCapabilities() *models.ProviderCapabilities {
	s.capsCalled = true
	return s.caps
}

func (s *stubProvider) ValidateConfig(config map[string]interface{}) (bool, []string) {
	s.validateCalled = true
	s.lastValidateConfig = config
	return s.validateOK, s.validateMsgs
}

// TestLLMProvider_InterfaceSatisfaction asserts a conforming implementation
// genuinely satisfies the LLMProvider interface AND is dispatchable through an
// interface-typed variable. If any method signature in the interface changes
// incompatibly, this fails to compile — which is the contract-stability guard.
func TestLLMProvider_InterfaceSatisfaction(t *testing.T) {
	var p LLMProvider = &stubProvider{}
	if p == nil {
		t.Fatal("expected non-nil LLMProvider")
	}
}

// TestLLMProvider_Complete_Success drives Complete through the interface and
// asserts the response flows back unchanged and the method was actually invoked.
func TestLLMProvider_Complete_Success(t *testing.T) {
	want := &models.LLMResponse{ID: "resp-1", Content: "hello", TokensUsed: 7}
	stub := &stubProvider{completeResp: want}
	var p LLMProvider = stub

	got, err := p.Complete(context.Background(), &models.LLMRequest{Prompt: "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stub.completeCalled {
		t.Fatal("Complete was not dispatched to the implementation")
	}
	if got == nil || got.ID != "resp-1" || got.Content != "hello" || got.TokensUsed != 7 {
		t.Fatalf("response did not flow through interface unchanged: %+v", got)
	}
}

// TestLLMProvider_Complete_ErrorPath asserts an error from the implementation
// propagates through the interface boundary and the response is nil.
func TestLLMProvider_Complete_ErrorPath(t *testing.T) {
	sentinel := errors.New("provider down")
	stub := &stubProvider{completeErr: sentinel}
	var p LLMProvider = stub

	got, err := p.Complete(context.Background(), &models.LLMRequest{Prompt: "hi"})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error to propagate, got %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil response on error, got %+v", got)
	}
}

// TestLLMProvider_CompleteStream_Success asserts the streaming channel contract:
// the channel yields the queued response and is closed (range terminates).
func TestLLMProvider_CompleteStream_Success(t *testing.T) {
	want := &models.LLMResponse{ID: "stream-1", Content: "chunk"}
	stub := &stubProvider{streamResp: want}
	var p LLMProvider = stub

	ch, err := p.CompleteStream(context.Background(), &models.LLMRequest{Prompt: "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stub.completeStreamCalled {
		t.Fatal("CompleteStream was not dispatched to the implementation")
	}
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}

	var received []*models.LLMResponse
	for resp := range ch { // proves the channel is closed; otherwise this blocks forever
		received = append(received, resp)
	}
	if len(received) != 1 {
		t.Fatalf("expected exactly 1 streamed response, got %d", len(received))
	}
	if received[0].ID != "stream-1" || received[0].Content != "chunk" {
		t.Fatalf("streamed response mismatch: %+v", received[0])
	}
}

// TestLLMProvider_CompleteStream_ErrorPath asserts a stream-setup error
// propagates and returns a nil channel.
func TestLLMProvider_CompleteStream_ErrorPath(t *testing.T) {
	sentinel := errors.New("stream init failed")
	stub := &stubProvider{streamErr: sentinel}
	var p LLMProvider = stub

	ch, err := p.CompleteStream(context.Background(), &models.LLMRequest{Prompt: "hi"})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error to propagate, got %v", err)
	}
	if ch != nil {
		t.Fatal("expected nil channel on error")
	}
}

// TestLLMProvider_HealthCheck covers both the healthy and unhealthy contract.
func TestLLMProvider_HealthCheck(t *testing.T) {
	tests := []struct {
		name      string
		healthErr error
		wantErr   bool
	}{
		{name: "healthy", healthErr: nil, wantErr: false},
		{name: "unhealthy", healthErr: errors.New("503"), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubProvider{healthErr: tt.healthErr}
			var p LLMProvider = stub

			err := p.HealthCheck()
			if !stub.healthCalled {
				t.Fatal("HealthCheck was not dispatched to the implementation")
			}
			if (err != nil) != tt.wantErr {
				t.Fatalf("HealthCheck err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// TestLLMProvider_GetCapabilities asserts the capabilities object flows back
// through the interface unchanged.
func TestLLMProvider_GetCapabilities(t *testing.T) {
	want := &models.ProviderCapabilities{
		SupportedModels:   []string{"model-x"},
		SupportsStreaming: true,
		SupportsTools:     true,
		Limits:            models.ModelLimits{MaxTokens: 4096},
	}
	stub := &stubProvider{caps: want}
	var p LLMProvider = stub

	got := p.GetCapabilities()
	if !stub.capsCalled {
		t.Fatal("GetCapabilities was not dispatched to the implementation")
	}
	if got == nil {
		t.Fatal("expected non-nil capabilities")
	}
	if !got.SupportsStreaming || !got.SupportsTools {
		t.Fatalf("capability flags did not flow through: %+v", got)
	}
	if got.Limits.MaxTokens != 4096 {
		t.Fatalf("expected MaxTokens 4096, got %d", got.Limits.MaxTokens)
	}
	if len(got.SupportedModels) != 1 || got.SupportedModels[0] != "model-x" {
		t.Fatalf("supported models did not flow through: %v", got.SupportedModels)
	}
}

// TestLLMProvider_ValidateConfig covers the (bool, []string) contract for both
// a valid config (true, no messages) and an invalid one (false, with reasons),
// and asserts the config argument actually reaches the implementation.
func TestLLMProvider_ValidateConfig(t *testing.T) {
	tests := []struct {
		name       string
		ok         bool
		msgs       []string
		config     map[string]interface{}
		wantOK     bool
		wantMsgLen int
	}{
		{
			name:       "valid config",
			ok:         true,
			msgs:       nil,
			config:     map[string]interface{}{"api_key": "set"},
			wantOK:     true,
			wantMsgLen: 0,
		},
		{
			name:       "invalid config returns reasons",
			ok:         false,
			msgs:       []string{"missing api_key", "missing base_url"},
			config:     map[string]interface{}{},
			wantOK:     false,
			wantMsgLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := &stubProvider{validateOK: tt.ok, validateMsgs: tt.msgs}
			var p LLMProvider = stub

			ok, msgs := p.ValidateConfig(tt.config)
			if !stub.validateCalled {
				t.Fatal("ValidateConfig was not dispatched to the implementation")
			}
			if ok != tt.wantOK {
				t.Fatalf("ValidateConfig ok = %v, want %v", ok, tt.wantOK)
			}
			if len(msgs) != tt.wantMsgLen {
				t.Fatalf("ValidateConfig msgs len = %d, want %d", len(msgs), tt.wantMsgLen)
			}
			if len(stub.lastValidateConfig) != len(tt.config) {
				t.Fatalf("config argument did not reach implementation intact: got %v", stub.lastValidateConfig)
			}
		})
	}
}
