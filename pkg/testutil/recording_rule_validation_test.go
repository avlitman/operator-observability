package testutil

import (
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/avlitman/operator-observability/pkg/operatormetrics"
	"github.com/avlitman/operator-observability/pkg/operatorrules"
)

var _ = Describe("Validators", func() {
	var linter *Linter

	BeforeEach(func() {
		linter = New()
	})

	Context("RecordingRule Validation", func() {
		It("should validate recording rule with valid input", func() {
			recordingRule := &operatorrules.RecordingRule{
				MetricsOpts: operatormetrics.MetricOpts{
					Name: "ExampleRecordingRule",
				},
				Expr: intstr.FromString("sum(rate(http_requests_total[5m]))"),
			}
			problems := linter.LintRecordingRule(recordingRule)
			Expect(problems).To(BeEmpty())
		})

		It("should return error if recording rule name is empty", func() {
			recordingRule := &operatorrules.RecordingRule{
				MetricsOpts: operatormetrics.MetricOpts{},
				Expr:        intstr.FromString("sum(rate(http_requests_total[5m]))"),
			}
			problems := linter.LintRecordingRule(recordingRule)
			Expect(problems).To(HaveLen(1))
			Expect(problems[0].Description).To(ContainSubstring("recording rule must have a name"))
		})

		It("should return error if recording rule expression is empty", func() {
			recordingRule := &operatorrules.RecordingRule{
				MetricsOpts: operatormetrics.MetricOpts{
					Name: "ExampleRecordingRule",
				},
				Expr: intstr.FromString(""),
			}
			problems := linter.LintRecordingRule(recordingRule)
			Expect(problems).To(HaveLen(1))
			Expect(problems[0].Description).To(ContainSubstring("recording rule must have an expression"))
		})
	})
})
