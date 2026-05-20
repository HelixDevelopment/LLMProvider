package ollama

import (
	"context"
	"strings"
	"testing"

	"digital.vasic.llmprovider/pkg/i18n"
)

// vcLocaleTranslator is a unit-test-only Translator returning fixed
// non-English strings for the two ValidateConfig message IDs the
// round-425 OllamaProvider migration routes through the i18n seam.
// Mocks are permitted in unit tests per CONST-050(A).
type vcLocaleTranslator struct{}

func (vcLocaleTranslator) T(_ context.Context, id string, _ map[string]any) (string, error) {
	switch id {
	case "llmprovider_validate_base_url_required":
		return "Bazni URL je obavezan", nil
	case "llmprovider_validate_model_required":
		return "Model je obavezan", nil
	}
	return id, nil
}

// TestOllamaValidateConfig_I18nSeam_Localized is the POSITIVE half of the
// round-425 CONST-046 paired mutation: with a real Translator wired both
// OllamaProvider validation errors are localized. Reverting either
// migrated literal to hardcoded English makes the wired translator inert
// for that case and this FAILS.
func TestOllamaValidateConfig_I18nSeam_Localized(t *testing.T) {
	defer i18n.SetTranslator(nil)
	i18n.SetTranslator(vcLocaleTranslator{})

	// The constructor substitutes localhost/llama2 defaults for empty
	// args; force the fields empty so both validation branches fire.
	p := NewOllamaProvider("", "")
	p.baseURL = ""
	p.model = ""
	ok, errs := p.ValidateConfig(nil)
	if ok || len(errs) != 2 {
		t.Fatalf("ValidateConfig with empty config should yield 2 errors, got %d", len(errs))
	}
	joined := strings.Join(errs, "|")
	for _, eng := range []string{"base URL is required", "model is required"} {
		if strings.Contains(joined, eng) {
			t.Fatalf("ValidateConfig emitted hardcoded English literal %q — CONST-046 round-425 migration regressed", eng)
		}
	}
	for _, loc := range []string{"Bazni URL", "Model je"} {
		if !strings.Contains(joined, loc) {
			t.Fatalf("ValidateConfig errors %q missing localized fragment %q — i18n seam not exercised", joined, loc)
		}
	}
}

// TestOllamaValidateConfig_I18nSeam_NoopFallback is the NEGATIVE half:
// with no Translator wired the NoopTranslator echoes each message ID
// verbatim — a loud, visible fallback, never a silent empty string.
func TestOllamaValidateConfig_I18nSeam_NoopFallback(t *testing.T) {
	i18n.SetTranslator(nil) // reset to NoopTranslator
	p := NewOllamaProvider("", "")
	p.baseURL = ""
	p.model = ""
	ok, errs := p.ValidateConfig(nil)
	if ok || len(errs) != 2 {
		t.Fatalf("ValidateConfig with empty config should yield 2 errors, got %d", len(errs))
	}
	want := map[string]bool{
		"llmprovider_validate_base_url_required": true,
		"llmprovider_validate_model_required":    true,
	}
	for _, e := range errs {
		if !want[e] {
			t.Fatalf("NoopTranslator fallback = %q, want a verbatim message ID echo", e)
		}
	}
}
