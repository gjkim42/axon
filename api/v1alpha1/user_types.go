package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UserSpec defines the desired state of User.
type UserSpec struct {
	// Name is the display name of the user.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Email is the email address of the user.
	// +optional
	Email string `json:"email,omitempty"`

	// GitHubToken references a Secret containing a GITHUB_TOKEN key for
	// GitHub authentication.
	// +optional
	GitHubToken *SecretReference `json:"githubToken,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Display Name",type=string,JSONPath=`.spec.name`
// +kubebuilder:printcolumn:name="Email",type=string,JSONPath=`.spec.email`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// User is the Schema for the users API.
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec UserSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// UserList contains a list of User.
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}
