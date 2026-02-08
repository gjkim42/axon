package integration

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	axonv1alpha1 "github.com/gjkim42/axon/api/v1alpha1"
)

var _ = Describe("User", func() {
	var ns *corev1.Namespace

	BeforeEach(func() {
		ns = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "test-user-",
			},
		}
		Expect(k8sClient.Create(ctx, ns)).To(Succeed())
	})

	Context("Creating a User", func() {
		It("Should create a User with all fields", func() {
			By("Creating a Secret for the GitHub token")
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "gh-token",
					Namespace: ns.Name,
				},
				StringData: map[string]string{
					"GITHUB_TOKEN": "test-token",
				},
			}
			Expect(k8sClient.Create(ctx, secret)).To(Succeed())

			By("Creating a User with name, email, and githubToken")
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user-full",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name:  "Test User",
					Email: "test@example.com",
					GitHubToken: &axonv1alpha1.SecretReference{
						Name: "gh-token",
					},
				},
			}
			Expect(k8sClient.Create(ctx, user)).To(Succeed())

			By("Fetching the created User")
			fetched := &axonv1alpha1.User{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      "test-user-full",
				Namespace: ns.Name,
			}, fetched)).To(Succeed())
			Expect(fetched.Spec.Name).To(Equal("Test User"))
			Expect(fetched.Spec.Email).To(Equal("test@example.com"))
			Expect(fetched.Spec.GitHubToken).NotTo(BeNil())
			Expect(fetched.Spec.GitHubToken.Name).To(Equal("gh-token"))
		})

		It("Should create a User with only the required name field", func() {
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user-minimal",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name: "Minimal User",
				},
			}
			Expect(k8sClient.Create(ctx, user)).To(Succeed())

			fetched := &axonv1alpha1.User{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      "test-user-minimal",
				Namespace: ns.Name,
			}, fetched)).To(Succeed())
			Expect(fetched.Spec.Name).To(Equal("Minimal User"))
			Expect(fetched.Spec.Email).To(BeEmpty())
			Expect(fetched.Spec.GitHubToken).To(BeNil())
		})
	})

	Context("Updating a User", func() {
		It("Should update User fields", func() {
			By("Creating a User")
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user-update",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name: "Original Name",
				},
			}
			Expect(k8sClient.Create(ctx, user)).To(Succeed())

			By("Updating the User")
			fetched := &axonv1alpha1.User{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      "test-user-update",
				Namespace: ns.Name,
			}, fetched)).To(Succeed())
			fetched.Spec.Name = "Updated Name"
			fetched.Spec.Email = "updated@example.com"
			Expect(k8sClient.Update(ctx, fetched)).To(Succeed())

			By("Verifying the update")
			updated := &axonv1alpha1.User{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{
				Name:      "test-user-update",
				Namespace: ns.Name,
			}, updated)).To(Succeed())
			Expect(updated.Spec.Name).To(Equal("Updated Name"))
			Expect(updated.Spec.Email).To(Equal("updated@example.com"))
		})
	})

	Context("Deleting a User", func() {
		It("Should delete a User", func() {
			By("Creating a User")
			user := &axonv1alpha1.User{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-user-delete",
					Namespace: ns.Name,
				},
				Spec: axonv1alpha1.UserSpec{
					Name: "Delete Me",
				},
			}
			Expect(k8sClient.Create(ctx, user)).To(Succeed())

			By("Deleting the User")
			Expect(k8sClient.Delete(ctx, user)).To(Succeed())

			By("Verifying the User is deleted")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      "test-user-delete",
					Namespace: ns.Name,
				}, &axonv1alpha1.User{})
				return apierrors.IsNotFound(err)
			}, 10*time.Second, 100*time.Millisecond).Should(BeTrue())
		})
	})

	Context("Listing Users", func() {
		It("Should list all Users in a namespace", func() {
			By("Creating multiple Users")
			for _, name := range []string{"user-alpha", "user-beta", "user-gamma"} {
				user := &axonv1alpha1.User{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns.Name,
					},
					Spec: axonv1alpha1.UserSpec{
						Name: name,
					},
				}
				Expect(k8sClient.Create(ctx, user)).To(Succeed())
			}

			By("Listing Users")
			userList := &axonv1alpha1.UserList{}
			Expect(k8sClient.List(ctx, userList, client.InNamespace(ns.Name))).To(Succeed())
			Expect(userList.Items).To(HaveLen(3))
		})
	})
})
