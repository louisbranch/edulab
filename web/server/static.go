package server

import (
	"html/template"
	"net/http"

	"github.com/louisbranch/edulab/web/presenter"
)

func (srv *Server) about(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)
	title := printer.Sprintf("About")
	page.Title = title
	page.Partials = []string{"about"}
	page.Content = struct {
		Breadcrumbs   template.HTML
		Title         string
		References    string
		Context       string
		Contributions string
		Source        string
	}{
		Breadcrumbs:   presenter.HomeBreadcrumbs(printer),
		Title:         title,
		References:    printer.Sprintf("References"),
		Context:       printer.Sprintf("This project was created as part of the course, Physical Science in Contemporary Society, at the University of Toronto with the intention of being a free resource for educators."),
		Contributions: printer.Sprintf("If you would like to contribute to the project, for example, adding more translations, get in touch:"),
		Source:        printer.Sprintf("Source Code"),
	}

	srv.render(w, page)
}

func (srv *Server) guide(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Educator's Guide")
	page.Title = title
	page.Partials = []string{"guide"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Title       string
		Texts       interface{}
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Title:       title,
		Texts: struct {
			Guide string
		}{
			Guide: printer.Sprintf(`
### Introduction
EduLab is designed to help educators incorporate scientific methods into their teaching strategies. This guide provides step-by-step instructions on using the platform to evaluate and refine your teaching methods with evidence-based insights.

---

### Step 1: Set Up an Experiment
1. **Define Your Teaching Interventions**  
   Identify the different teaching methods or approaches you want to compare (e.g., traditional lecture vs. interactive workshops).
   
2. **Create Cohorts**  
   Use EduLab's cohort feature to group students who will experience specific teaching interventions. For example:
   - **Control**: Traditional lecture method.
   - **Intervention**: Interactive workshop approach.

3. **Develop Assessments**  
   Design a set of pre- and post-assessment questions to measure the effectiveness of each teaching method. Ensure these questions align with the learning objectives.

---

### Step 2: Conduct Pre-Assessment
- Share the pre-assessment link with your cohorts before introducing any teaching intervention. 
- Encourage students to complete the assessment to establish a baseline for their knowledge.

---

### Step 3: Implement Your Teaching Interventions
- Conduct your planned teaching methods for each cohort.
- Ensure that the interventions are distinct and well-documented for accurate comparisons.

---

### Step 4: Conduct Post-Assessment
- After completing the intervention, share the post-assessment link with the same cohorts.
- Collect responses to measure the knowledge gained through each teaching method.

---

### Step 5: Analyze Results
- Use EduLab's **Learning Gain Analysis** to compare pre- and post-assessment scores within and across cohorts. This allows you to:
  - Identify which teaching method led to higher learning gains.
  - Understand how different demographic groups responded to the interventions.
  
- Utilize the demographic data to tailor future teaching methods to meet the diverse needs of your students.

---

### Step 6: Iterate and Refine
- Based on the results, refine your teaching strategies to optimize learning outcomes. Repeat the process to continually improve your methods.`),
		},
	}

	srv.render(w, page)
}

func (srv *Server) faq(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Frequently Asked Questions")
	page.Title = title
	page.Partials = []string{"faq"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Title       string
		Texts       interface{}
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Title:       title,
		Texts: struct {
			FAQ string
		}{
			FAQ: printer.Sprintf(`### How is data privacy ensured on EduLab?  
EduLab anonymizes all student data, ensuring no personally identifiable information is stored or shared. The platform also complies with data protection standards.

---

### Can I customize the assessments?  
Yes, you can create and edit multiple-choice questions to align with your specific learning objectives.

---

### What types of demographic data can I collect?  
EduLab allows you to collect data on gender, age group, year of study, and major, helping you understand how different factors influence learning outcomes.

---

### How do I interpret the learning gain analysis?  
Learning gains are calculated as the difference between pre- and post-assessment scores, normalized to account for the initial baseline. Higher gains indicate more effective teaching methods.

---

### Is the platform open-source?  
Yes, EduLab provides access to its open-source code, allowing you to customize the platform to fit your needs.

---

### Can I use EduLab for non-science subjects?  
Absolutely! While EduLab is designed with science education in mind, its features are applicable across disciplines.`),
		},
	}

	srv.render(w, page)
}

func (srv *Server) tos(w http.ResponseWriter, r *http.Request) {
	printer, page := srv.i18n(w, r)

	title := printer.Sprintf("Terms of Service")
	page.Title = title
	page.Partials = []string{"tos"}
	page.Content = struct {
		Breadcrumbs template.HTML
		Title       string
		Texts       interface{}
	}{
		Breadcrumbs: presenter.HomeBreadcrumbs(printer),
		Title:       title,
		Texts: struct {
			TOS string
		}{
			TOS: printer.Sprintf(`### 1. Purpose

EduLab is a prototype platform designed for educational purposes only. It is not intended for commercial use. By using this platform, you agree to these Terms of Service.

### 2. User-Generated Content

* You retain ownership of any content you create or upload to EduLab.
* EduLab does not claim ownership of user-generated content and acts solely as a tool to facilitate educational activities.
* By using the platform, you grant EduLab the right to store and process your content as part of its educational functionality.

### 3. Content Guidelines

* You agree not to upload or create content that:
  * Violates copyright, trademark, or other intellectual property rights.
  * Contains offensive, harmful, or inappropriate material.
  * Violates any applicable laws or regulations.
* EduLab reserves the right to remove content that violates these guidelines without prior notice.

### 4. Disclaimer of Liability

* EduLab is provided "as is," without warranties of any kind, expressed or implied.
* EduLab is not responsible for the accuracy, reliability, or legality of user-generated content.
* The platform is not moderated, and EduLab is not liable for any damages resulting from the use of the platform or the content hosted on it.

### 5. No Accounts or Personal Data

* EduLab does not require user accounts or collect personal data.
* Any data submitted is stored temporarily and used solely for educational purposes.

### 6. Indemnification

By using EduLab, you agree to indemnify and hold harmless the developers of EduLab from any claims or liabilities arising from your use of the platform or content you create.

### 7. Updates to Terms

These Terms of Service may be updated periodically. Continued use of the platform constitutes agreement to the updated terms.`),
		},
	}

	srv.render(w, page)
}
