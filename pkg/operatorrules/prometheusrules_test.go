package operatorrules

import (
	"k8s.io/apimachinery/pkg/util/intstr"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"

	"github.com/avlitman/operator-observability/pkg/operatormetrics"
)

var _ = Describe("PrometheusRules", func() {
	Context("Building resource", func() {
		var recordingRules = []RecordingRule{
			{
				MetricsOpts: operatormetrics.MetricOpts{
					Name:        "number_of_pods",
					Help:        "Number of guestbook operator pods in the cluster",
					ConstLabels: map[string]string{"controller": "guestbook"},
				},
				MetricType: operatormetrics.GaugeType,
				Expr:       intstr.FromString("sum(up{namespace='default', pod=~'guestbook-operator-.*'}) or vector(0)"),
			},
		}

		var alerts = []promv1.Rule{
			{
				Alert: "GuestbookOperatorDown",
				Expr:  intstr.FromString("number_of_pods == 0"),
				Annotations: map[string]string{
					"summary":     "Guestbook operator is down",
					"description": "Guestbook operator is down for more than 5 minutes.",
				},
				Labels: map[string]string{
					"severity": "critical",
				},
			},
		}

		BeforeEach(func() {
			operatorRegistry = newRegistry()

			err := RegisterRecordingRules(recordingRules)
			Expect(err).To(Not(HaveOccurred()))

			err = RegisterAlerts(alerts)
			Expect(err).To(Not(HaveOccurred()))
		})

		It("should build PrometheusRule with valid input", func() {
			rules, err := BuildPrometheusRule(
				"guestbook-operator-prometheus-rules",
				"default",
				map[string]string{"app": "guestbook-operator"},
			)

			Expect(err).To(BeNil())
			Expect(rules).NotTo(BeNil())
			Expect(rules.Name).To(Equal("guestbook-operator-prometheus-rules"))
			Expect(rules.Namespace).To(Equal("default"))

			Expect(rules.Spec.Groups).To(HaveLen(2))

			Expect(rules.Spec.Groups[0].Name).To(Equal("recordingRules.rules"))
			Expect(rules.Spec.Groups[0].Rules).To(HaveLen(1))
			Expect(rules.Spec.Groups[0].Rules[0].Record).To(Equal("number_of_pods"))
			Expect(rules.Spec.Groups[0].Rules[0].Expr).To(Equal(intstr.FromString("sum(up{namespace='default', pod=~'guestbook-operator-.*'}) or vector(0)")))

			Expect(rules.Spec.Groups[1].Name).To(Equal("alerts.rules"))
			Expect(rules.Spec.Groups[1].Rules).To(HaveLen(1))
			Expect(rules.Spec.Groups[1].Rules[0].Alert).To(Equal("GuestbookOperatorDown"))
			Expect(rules.Spec.Groups[1].Rules[0].Expr).To(Equal(intstr.FromString("number_of_pods == 0")))
		})
	})
})
